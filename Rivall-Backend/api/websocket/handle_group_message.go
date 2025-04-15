package websocket

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"

	"Rivall-Backend/db"
)

type SendGroupMessageEvent struct {
	MessageData string `json:"message_data"`
	ReceiverID  string `json:"receiver_id"`
	Timestamp   string `json:"timestamp"`
	MessageType string `json:"message_type"`
}

type NewGroupMessageEvent struct {
	SendGroupMessageEvent
	Sent   string   `json:"sent"`
	SeenBy []string `json:"seen_by"`
}

func SendGroupMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent SendGroupMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		log.Error().Err(err).Msg("bad payload in request")
		return nil
	}

	// Validate Event
	if exists := db.GroupExists(event.GroupID); !exists {
		log.Error().Msg("group does not exist")
		return nil
	}
	if exists := db.UserInGroup(event.GroupID, event.UserID); !exists {
		log.Error().Msgf("Sender user not in group: %s", event.UserID)
		return nil
	}

	// Save message to Group in Database
	bsonUserID, err := bson.ObjectIDFromHex(event.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert group ID")
		return err
	}

	var message = db.Message{
		ID:          bson.NewObjectID(),
		UserID:      bsonUserID,
		MessageData: chatevent.MessageData,
		Timestamp:   chatevent.Timestamp,
		MessageType: chatevent.MessageType,
		SeenBy:      []bson.ObjectID{bsonUserID},
	}
	if err := db.InsertGroupMessage(event.GroupID, message); err != nil {
		log.Error().Err(err).Msg("failed to insert message")
		return err
	}

	// Prepare an Outgoing Message to others
	var broadMessage NewGroupMessageEvent
	broadMessage.SendGroupMessageEvent = chatevent
	broadMessage.Sent = time.Now().Format(time.RFC3339)
	broadMessage.SeenBy = []string{event.UserID}
	broadMessageData, err := json.Marshal(broadMessage)

	var outgoingEvent Event
	outgoingEvent.Payload = broadMessageData
	outgoingEvent.Type = EventNewGroupMessage
	outgoingEvent.GroupID = event.GroupID
	outgoingEvent.UserID = event.UserID
	outgoingEvent.DirectMessageID = event.DirectMessageID

	if err != nil {
		log.Error().Err(err).Msg("failed to marshal message")
		return err
	}

	// Get all Group Members, Send them the message
	groupMembers, err := db.GetGroupMembers(event.GroupID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get group members")
		return err
	}
	for _, UserID := range groupMembers {
		if _, ok := c.Manager().Clients()[UserID]; !ok {
			continue
		}
		if UserID == event.UserID {
			continue
		}
		c.Manager().Clients()[UserID].Egress <- outgoingEvent
	}
	return nil
}
