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
	ID            string    `json:"_id"`
	SendUserID    string    `json:"send_user_id"`
	RecieveUserID string    `json:"receive_user_id"`
	GroupID       string    `json:"group_id"`
	GroupName     string    `json:"group_name"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
	Status        int8      `json:"status"`
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
	for _, UserID := range chatevent.UserIDs {
		id, err := db.CreateGroupRequest(AdminUserID, UserID, groupID, chatevent.GroupName, chatevent.Message)
		if err != nil {
			log.Error().Err(err).Msg("failed to send group request")
			return err
		}

		// Prepare an Outgoing Message to others
		var broadMessage JoinGroupRequest
		broadMessage.ID = id
		broadMessage.SendUserID = AdminUserID
		broadMessage.RecieveUserID = UserID
		broadMessage.GroupID = groupID
		broadMessage.GroupName = chatevent.GroupName
		broadMessage.Message = chatevent.Message
		broadMessage.Timestamp = time.Now()
		broadMessage.Status = 0

		data, err := json.Marshal(broadMessage)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal broadcast message")
			return err
		}

		// Place payload into an Event
		var outgoingEvent Event
		outgoingEvent.Payload = data
		outgoingEvent.Type = EventNewGroupRequest
		outgoingEvent.GroupID = groupID
		outgoingEvent.UserID = AdminUserID
		outgoingEvent.DirectMessageID = event.DirectMessageID

		// Broadcast to all other Clients in the group if the users are online
		if _, ok := c.Manager().Clients()[UserID]; !ok {
			continue
		}
		c.Manager().Clients()[UserID].Egress <- outgoingEvent
	}
	return nil
}
