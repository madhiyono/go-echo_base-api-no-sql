package models

import (
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebSocketMessage struct {
	Type      string         `json:"type"`
	Event     string         `json:"event"`
	Data      map[string]any `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}

type WebSocketConnection struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	Conn      *websocket.Conn    `json:"-"`
	Channels  []string           `json:"channels"`
	Connected bool               `json:"connected"`
	LastPing  time.Time          `json:"last_ping"`
}

type Subscription struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Channel  string             `json:"channel"`
	Resource string             `json:"resource"`
	Actions  []string           `json:"actions"`
}
