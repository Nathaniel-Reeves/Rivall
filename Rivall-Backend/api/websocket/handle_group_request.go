package websocket

import (
	db "Rivall-Backend/db"
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type AcceptGroupRequestEvent struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

type AcceptGroupRequestPayload struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

func AcceptGroupRequestHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent AcceptGroupRequestPayload
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		log.Error().Err(err).Msg("bad payload in request")
		return err
	}

	// Validate Event
	if exists := db.GroupExists(event.GroupID); !exists {
		log.Error().Msg("group does not exist")
		return nil
	}
	if exists := db.UserWasRequestedToJoinGroup(event.GroupID, event.UserID); !exists {
		log.Error().Msgf("User was not requested to join group: %s", event.UserID)
		return nil
	}

	// Accept Group Request in Database
	if err := db.AcceptGroupRequest(chatevent.GroupID, chatevent.UserID); err != nil {
		log.Error().Err(err).Msg("failed to accept group request")
		return err
	}

	var acceptEvent AcceptGroupRequestEvent
	acceptEvent.GroupID = chatevent.GroupID
	acceptEvent.UserID = chatevent.UserID

	data, err := json.Marshal(acceptEvent)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal accept group request event")
		return err
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventAcceptGroupRequest
	outgoingEvent.UserID = chatevent.UserID
	outgoingEvent.GroupID = chatevent.GroupID
	outgoingEvent.DirectMessageID = ""

	// Send event to admin user
	AdminID, err := db.GetGroupAdminID(chatevent.GroupID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get group admin")
		return err
	}

	if _, ok := c.Manager().Clients()[AdminID]; !ok {
		return nil
	}
	c.Manager().Clients()[AdminID].Egress <- outgoingEvent
	return nil
}

type RejectGroupRequestEvent struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

type RejectGroupRequestPayload struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

func RejectGroupRequestHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent RejectGroupRequestPayload
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		log.Error().Err(err).Msg("bad payload in request")
		return err
	}

	// Validate Event
	if exists := db.GroupExists(event.GroupID); !exists {
		log.Error().Msg("group does not exist")
		return nil
	}
	if exists := db.UserWasRequestedToJoinGroup(event.GroupID, event.UserID); !exists {
		log.Error().Msgf("User was not requested to join group: %s", event.UserID)
		return nil
	}

	// Reject Group Request in Database
	if err := db.RejectGroupRequest(chatevent.GroupID, chatevent.UserID); err != nil {
		log.Error().Err(err).Msg("failed to reject group request")
		return err
	}

	var rejectEvent RejectGroupRequestEvent
	rejectEvent.GroupID = chatevent.GroupID
	rejectEvent.UserID = chatevent.UserID

	data, err := json.Marshal(rejectEvent)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal reject group request event")
		return err
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventRejectGroupRequest
	outgoingEvent.UserID = chatevent.UserID
	outgoingEvent.GroupID = chatevent.GroupID
	outgoingEvent.DirectMessageID = ""

	// Send event to admin user
	AdminID, err := db.GetGroupAdminID(chatevent.GroupID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get group admin")
		return err
	}

	if _, ok := c.Manager().Clients()[AdminID]; !ok {
		return nil
	}
	c.Manager().Clients()[AdminID].Egress <- outgoingEvent

	return nil

}
