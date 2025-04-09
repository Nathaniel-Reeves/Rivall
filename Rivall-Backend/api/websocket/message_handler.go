package websocket

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

func SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		log.Error().Err(err).Msg("bad payload in request")
		return err
	}

	// Get UserIDs in Group

	// Confirm User is in the Group

	// Save message to Group in Database

	// Prepare an Outgoing Message to others
	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal broadcast message")
		return err
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage
	outgoingEvent.GroupID = event.GroupID
	outgoingEvent.UserID = event.UserID

	// Broadcast to all other Clients in the Group
	for client := range c.Manager().Clients() {
		// Only send to clients inside the same Group
		// if client.Chatroom() == c.Chatroom() {
		// 	client.Egress <- outgoingEvent
		// }
		client.Egress <- outgoingEvent
	}
	return nil
}
