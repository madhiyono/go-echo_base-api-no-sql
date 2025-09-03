# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### New Features

- (none)

### Changes

- (none)

## [1.0.1-beta] - 2025-09-03

### New Features

- Standardized API response template for all endpoints (success and error)
- Request validation for incoming data (with custom rules)
- Basic authentication using JWT
- Authorization based on user roles
- Role management stored in the database
- Role permission management (assign permissions to roles)

### Changes

- Improved security and flexibility with role-based access control
- Added endpoints for role and permission management

## [1.0.0-beta] - 2025-08-30

- First stable release
- REST API base with Echo
- MongoDB repository pattern
- Basic user CRUD endpoints
- Project documentation and setup instructions
