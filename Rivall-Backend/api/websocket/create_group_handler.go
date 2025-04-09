package websocket

import (
	"encoding/json"
	"time"

	db "Rivall-Backend/db"

	"github.com/rs/zerolog/log"
)

type CreateGroupPayload struct {
	GroupName string   `json:"group_name"`
	UserIDs   []string `json:"user_ids"`
	Message   string   `json:"message"`
}

type JoinGroupRequest struct {
	GroupID string    `json:"group_id"`
	UserID  string    `json:"user_id"`
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

func CreateGroupHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent CreateGroupPayload
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		log.Error().Err(err).Msg("bad payload in request")
		return err
	}

	// Get Admin UserID, Admin is the user creating the group
	AdminUserID := event.UserID

	// Confirm Users are legitimate
	for _, userID := range chatevent.UserIDs {
		if exists := db.UserExists(userID); !exists {
			log.Error().Msg("user does not exist")
			return nil
		}
	}

	// Create New Group in Database
	groupID, err := db.CreateGroup(chatevent.GroupName, AdminUserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create group")
		return err
	}

	log.Info().Msgf(`Created group with ID: %v`, groupID)

	// Add a Request to all users requested to be added to the group
	for _, userID := range chatevent.UserIDs {
		if err := db.CreateUserMessageRequest(AdminUserID, userID, groupID, chatevent.Message); err != nil {
			log.Error().Err(err).Msg("failed to send group request")
			return err
		}
	}

	// Prepare an Outgoing Message to others
	var broadMessage JoinGroupRequest

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.GroupID = groupID

	data, err := json.Marshal(broadMessage)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal broadcast message")
		return err
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage
	outgoingEvent.GroupID = groupID
	outgoingEvent.UserID = AdminUserID

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
