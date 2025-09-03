package services

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebSocketService struct {
	connections   map[primitive.ObjectID]*models.WebSocketConnection
	subscriptions map[string][]primitive.ObjectID // channel -> user_ids
	mutex         sync.RWMutex
	logger        *logger.Logger
	upgrader      websocket.Upgrader
}

func NewWebSocketService(logger *logger.Logger) *WebSocketService {
	service := &WebSocketService{
		connections:   make(map[primitive.ObjectID]*models.WebSocketConnection),
		subscriptions: make(map[string][]primitive.ObjectID),
		logger:        logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, restrict this
			},
		},
	}

	// Start ping/pong handler
	go service.handlePings()

	return service
}

// GetUpgrader returns the WebSocket upgrader
func (s *WebSocketService) GetUpgrader() *websocket.Upgrader {
	return &s.upgrader
}

func (s *WebSocketService) HandleConnection(conn *websocket.Conn, userID primitive.ObjectID) {
	wsConn := &models.WebSocketConnection{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Conn:      conn,
		Channels:  []string{},
		Connected: true,
		LastPing:  time.Now(),
	}

	s.mutex.Lock()
	s.connections[userID] = wsConn
	s.mutex.Unlock()

	s.logger.Info("WebSocket connection established for user: %s", userID.Hex())

	// Handle messages
	defer s.disconnectUser(userID)
	for {
		var message models.WebSocketMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			s.logger.Error("WebSocket read error for user %s: %v", userID.Hex(), err)
			break
		}

		s.handleMessage(wsConn, message)
	}
}

func (s *WebSocketService) handleMessage(conn *models.WebSocketConnection, message models.WebSocketMessage) {
	switch message.Type {
	case "subscribe":
		s.subscribeToChannel(conn, message)
	case "unsubscribe":
		s.unsubscribeFromChannel(conn, message)
	case "ping":
		s.handlePing(conn)
	default:
		s.sendError(conn, "Unknown message type")
	}
}

func (s *WebSocketService) subscribeToChannel(conn *models.WebSocketConnection, message models.WebSocketMessage) {
	channel, ok := message.Data["channel"].(string)
	if !ok {
		s.sendError(conn, "Invalid channel")
		return
	}

	// Add channel to user's subscriptions
	s.mutex.Lock()
	conn.Channels = append(conn.Channels, channel)

	// Add user to channel subscribers
	subscribers := s.subscriptions[channel]
	subscribers = append(subscribers, conn.UserID)
	s.subscriptions[channel] = subscribers
	s.mutex.Unlock()

	s.sendSuccess(conn, "Subscribed to channel: "+channel)
}

func (s *WebSocketService) unsubscribeFromChannel(conn *models.WebSocketConnection, message models.WebSocketMessage) {
	channel, ok := message.Data["channel"].(string)
	if !ok {
		s.sendError(conn, "Invalid channel")
		return
	}

	s.mutex.Lock()
	// Remove channel from user's subscriptions
	for i, ch := range conn.Channels {
		if ch == channel {
			conn.Channels = append(conn.Channels[:i], conn.Channels[i+1:]...)
			break
		}
	}

	// Remove user from channel subscribers
	if subscribers, exists := s.subscriptions[channel]; exists {
		for i, userID := range subscribers {
			if userID == conn.UserID {
				s.subscriptions[channel] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}
	s.mutex.Unlock()

	s.sendSuccess(conn, "Unsubscribed from channel: "+channel)
}

func (s *WebSocketService) handlePing(conn *models.WebSocketConnection) {
	conn.LastPing = time.Now()
	response := models.WebSocketMessage{
		Type:      "pong",
		Timestamp: time.Now(),
	}
	conn.Conn.WriteJSON(response)
}

func (s *WebSocketService) disconnectUser(userID primitive.ObjectID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if conn, exists := s.connections[userID]; exists {
		conn.Conn.Close()
		conn.Connected = false

		// Remove user from all channel subscriptions
		for _, channel := range conn.Channels {
			if subscribers, exists := s.subscriptions[channel]; exists {
				for i, subUserID := range subscribers {
					if subUserID == userID {
						s.subscriptions[channel] = append(subscribers[:i], subscribers[i+1:]...)
						break
					}
				}
			}
		}

		delete(s.connections, userID)
		s.logger.Info("WebSocket connection closed for user: %s", userID.Hex())
	}
}

func (s *WebSocketService) BroadcastToChannel(channel string, message models.WebSocketMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if subscribers, exists := s.subscriptions[channel]; exists {
		for _, userID := range subscribers {
			if conn, exists := s.connections[userID]; exists && conn.Connected {
				conn.Conn.WriteJSON(message)
			}
		}
	}
}

func (s *WebSocketService) BroadcastToUser(userID primitive.ObjectID, message models.WebSocketMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if conn, exists := s.connections[userID]; exists && conn.Connected {
		conn.Conn.WriteJSON(message)
	}
}

func (s *WebSocketService) sendSuccess(conn *models.WebSocketConnection, message string) {
	response := models.WebSocketMessage{
		Type:      "success",
		Data:      map[string]interface{}{"message": message},
		Timestamp: time.Now(),
	}
	conn.Conn.WriteJSON(response)
}

func (s *WebSocketService) sendError(conn *models.WebSocketConnection, message string) {
	response := models.WebSocketMessage{
		Type:      "error",
		Data:      map[string]interface{}{"message": message},
		Timestamp: time.Now(),
	}
	conn.Conn.WriteJSON(response)
}

func (s *WebSocketService) handlePings() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.RLock()
		connections := make([]*models.WebSocketConnection, 0, len(s.connections))
		for _, conn := range s.connections {
			connections = append(connections, conn)
		}
		s.mutex.RUnlock()

		now := time.Now()
		for _, conn := range connections {
			if now.Sub(conn.LastPing) > time.Minute {
				s.disconnectUser(conn.UserID)
			}
		}
	}
}

func (s *WebSocketService) NotifyResourceChange(channel, resource, action string, data map[string]interface{}) {
	message := models.WebSocketMessage{
		Type:      "notification",
		Event:     "resource_change",
		Data:      data,
		Timestamp: time.Now(),
	}

	// Add resource and action info
	if message.Data == nil {
		message.Data = make(map[string]interface{})
	}
	message.Data["resource"] = resource
	message.Data["action"] = action

	s.BroadcastToChannel(channel, message)
}
