package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Post User Contact
//
//	@summary		Create User Contact
//	@description	Create a User Contact in the database
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/users/{user_id}/contacts/{contact_id} [post]
func PostUserContact(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST user contact")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Check the user exists
	if ReadId(userID).ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// get contact id from url parameters
	contactID := vars["contact_id"]

	// check the contact exists
	if ReadId(contactID).ID == bson.NilObjectID {
		log.Error().Msg("Contact does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Contact does not exist."))
		return
	}

	// create user contact with new contact
	err := CreateUserContact(userID, contactID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set user contact")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set user contact."))
		return
	}
}

// Get User
//
//	@summary		Get User
//	@description	Get a User from the database
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/users/{user_id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET user")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["user_id"]
	log.Debug().Msgf("User ID: %s", userID)

	// check user exists
	user := ReadId(userID)
	if user.ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// return user
	json.NewEncoder(w).Encode(user)
}

// Get User Contacts
//
//	@summary		Get User Contacts
//	@description	Get a User's contacts from the database
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/users/{user_id}/contacts [get]
func GetUserPopulateContacts(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET user contacts")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["id"]

	// check user exists
	user := ReadIdPopulateContacts(userID)
	if user.ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// return user contacts
	json.NewEncoder(w).Encode(user.Contacts)
}
