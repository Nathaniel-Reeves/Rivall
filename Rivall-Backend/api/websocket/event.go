package websocket

import (
	"encoding/json"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`

	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, c *Client) error

const (
	// Action Events
	EventSendMessage = "send_message"
	EventCreateGroup = "create_group"

	// Don't forget to add new action events to the setupEventHandlers func in manager.go

	// Listen Events
	EventNewMessage = "new_message"
)
