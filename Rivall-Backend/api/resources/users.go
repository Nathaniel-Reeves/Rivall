package resources

import (
	"encoding/json"
	"net/http"

	db "Rivall-Backend/db"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PostUserContact(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST user contact")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Check the user exists
	if db.ReadByUserId(userID).ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// get contact id from content body
	vars = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&vars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode contact ID, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode contact ID, invalid JSON request."))
		return
	}
	contactID := vars["contact_id"]

	// check the contact exists
	if db.ReadByUserId(contactID).ID == bson.NilObjectID {
		log.Error().Msg("Contact does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Contact does not exist."))
		return
	}

	// check if user already has this contact
	contacts := db.ReadByUserId(userID).ContactIDs
	for _, contact := range contacts {
		id, err := bson.ObjectIDFromHex(contactID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert contact ID")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to convert contact ID."))
			return
		}
		if contact == id {
			log.Error().Msg("User already has this contact")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User already has this contact."))
			return
		}
	}

	// create user contact with new contact
	err = db.CreateUserContact(userID, contactID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set user contact")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set user contact."))
		return
	}
}

func DeleteUserContact(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("DELETE user contact")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["user_id"]
	log.Debug().Msgf("User ID: %s", userID)

	// get contact id from content body
	vars = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&vars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode contact ID, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode contact ID, invalid JSON request."))
		return
	}

	// delete user contact
	err = db.RemoveUserContact(userID, vars["contact_id"])
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user contact")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete user contact."))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET user")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["user_id"]
	log.Debug().Msgf("User ID: %s", userID)

	// check user exists
	user := db.ReadByUserId(userID)
	if user.ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// return user
	json.NewEncoder(w).Encode(user)
}
