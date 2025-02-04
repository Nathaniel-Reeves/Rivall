package messagegroup

import (
	user "Rivall-Backend/api/resources/auth"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Create Message Group
//
//	@summary		Create Message Group
//	@description	Create a Message Group in the database
//	@tags			messagegroups
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/messagegroups [post]
func PostNewMessageGroup(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST new message group")

	// parse json data
	messageGroup := MessageGroup{}
	err := json.NewDecoder(r.Body).Decode(&messageGroup)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode message group, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode message group, invalid JSON request."))
		return
	}

	// check message group has at least one user
	if len(messageGroup.Users) == 0 {
		log.Error().Msg("Message group must have at least one user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message group must have at least one user."))
		return
	}

	// check at least one user in the message group is the logged in user
	userID := r.Context().Value("user_id").(string)
	found := false
	for _, u := range messageGroup.Users {
		if u.ID.Hex() == userID {
			found = true
			break
		}
	}
	if !found {
		log.Error().Msg("Message group must have the logged in user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message group must have the logged in user."))
		return
	}

	// create message group
	user := user.User{}
	user.ID, err = bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user ID")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user ID"))
		return
	}
	err = CreateMessageGroup(user, messageGroup)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create message group")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create message group"))
		return
	}

	// add message group to each user in the message group
	err = AddMessageGroupToUsers(user, messageGroup)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add message group to users")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to add message group to users"))
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
}
