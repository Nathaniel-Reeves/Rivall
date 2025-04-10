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
	Payload         json.RawMessage `json:"payload"`
	GroupID         string          `json:"group_id"`
	DirectMessageID string          `json:"direct_message_id"`
	UserID          string          `json:"user_id"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, c *Client) error

const (
	// Action Events
	EventSendMessage        = "send_message"
	EventCreateGroup        = "create_group"
	EventAcceptGroupRequest = "accept_group_request"
	EventRejectGroupRequest = "reject_group_request"
	EventSendGroupMessage   = "send_group_message"

	// Don't forget to add new action events to the setupEventHandlers func in manager.go

	// Listen Events
	EventNewMessage           = "new_message"
	EventNewGroupRequest      = "new_group_request"
	EventGroupRequestAccepted = "group_request_accepted"
	EventGroupRequestRejected = "group_request_rejected"
	EventNewGroupMessage      = "new_group_message"
)
