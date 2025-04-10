package websocket

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"

	"Rivall-Backend/db"
)

type SendMessageEvent struct {
	MessageData string `json:"message_data"`
	ReceiverID  string `json:"receiver_id"`
	Timestamp   string `json:"timestamp"`
	MessageType string `json:"message_type"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent   time.Time `json:"sent"`
	SeenBy []string  `json:"seen_by"`
}

func SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		log.Error().Err(err).Msg("bad payload in request")
		return err
	}

	// Validate Event
	if exists := db.DirectMessageExists(event.DirectMessageID); !exists {
		log.Error().Msg("direct message does not exist")
		return nil
	}
	if exists := db.UserInDirectMessage(event.DirectMessageID, event.UserID); !exists {
		log.Error().Msgf("Sender user not in direct message: %s", event.UserID)
		return nil
	}
	if exists := db.UserInDirectMessage(event.DirectMessageID, chatevent.ReceiverID); !exists {
		log.Error().Msgf("Receiver user not in direct message: %s", chatevent.ReceiverID)
		return nil
	}

	// Save message to Group in Database
	bsonFromID, err := bson.ObjectIDFromHex(event.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert direct message ID")
	}

	var message = db.Message{
		ID:          bson.NewObjectID(),
		SenderID:    bsonFromID,
		MessageData: chatevent.MessageData,
		Timestamp:   chatevent.Timestamp,
		MessageType: chatevent.MessageType,
		SeenBy:      []bson.ObjectID{bsonFromID},
	}
	if err := db.InsertMessage(event.DirectMessageID, message); err != nil {
		log.Error().Err(err).Msg("failed to insert message")
		return err
	}

	// Prepare an Outgoing Message to others
	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.MessageData = chatevent.MessageData
	broadMessage.ReceiverID = chatevent.ReceiverID
	broadMessage.Timestamp = chatevent.Timestamp
	broadMessage.MessageType = chatevent.MessageType
	broadMessage.SeenBy = []string{event.UserID}

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
	outgoingEvent.DirectMessageID = event.DirectMessageID

	// Broadcast to the other Client if user is online
	if _, ok := c.Manager().Clients()[chatevent.ReceiverID]; !ok {
		return nil
	}
	c.Manager().Clients()[chatevent.ReceiverID].Egress <- outgoingEvent
	return nil
}
