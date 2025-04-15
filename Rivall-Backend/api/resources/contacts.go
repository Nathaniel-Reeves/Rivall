package resources

import (
	db "Rivall-Backend/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FullContact struct {
	ID          bson.ObjectID `json:"_id" bson:"_id"`
	FirstName   string        `json:"first_name" bson:"first_name"`
	LastName    string        `json:"last_name" bson:"last_name"`
	Email       string        `json:"email" bson:"email"`
	AvatarImage string        `json:"avatar_image" bson:"avatar_image"`
}

func GetContact(w http.ResponseWriter, r *http.Request) {
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

	var contact FullContact
	contact.ID = user.ID
	contact.FirstName = user.FirstName
	contact.LastName = user.LastName
	contact.Email = user.Email
	contact.AvatarImage = user.AvatarImage

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}

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
	user := db.ReadByUserIdWithPopulatedFields(userID)
	if user.ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}
	contacts := user.Contacts
	for _, contact := range contacts {
		if contact.ContactID.Hex() == contactID {
			log.Error().Msg("User already has this contact")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User already has this contact."))
			return
		}
	}

	// create user contact with new contact
	err = db.CreateContact(userID, contactID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set user contact")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to set user contact."))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Contact added successfully."))
	log.Info().Msg("Contact added successfully.")
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET chat")

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

	// get contact id from url parameters
	chatID := vars["chat_id"]

	// check the contact exists
	dm, err := db.ReadDirectMessages(chatID)
	if err != nil {
		log.Error().Msg("Contact does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Contact does not exist."))
		return
	}

	// check if user is in the direct message group
	if !db.UserInDirectMessage(chatID, userID) {
		log.Error().Msg("User is not in the direct message group")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User is not in the direct message group."))
		return
	}

	type ChatUser struct {
		ID          bson.ObjectID `json:"_id" bson:"_id"`
		FirstName   string        `json:"first_name" bson:"first_name"`
		LastName    string        `json:"last_name" bson:"last_name"`
		Email       string        `json:"email" bson:"email"`
		AvatarImage string        `json:"avatar_image" bson:"avatar_image"`
	}

	type messageData struct {
		GroupMembers map[string]ChatUser `json:"group_members"`
		Messages     []db.Message        `json:"messages"`
	}

	var userA = db.ReadByUserId(dm.UserAID.Hex())
	var userB = db.ReadByUserId(dm.UserBID.Hex())

	var data messageData
	data.GroupMembers = make(map[string]ChatUser)
	data.GroupMembers[userA.ID.Hex()] = ChatUser{
		ID:          userA.ID,
		FirstName:   userA.FirstName,
		LastName:    userA.LastName,
		Email:       userA.Email,
		AvatarImage: userA.AvatarImage,
	}
	data.GroupMembers[userB.ID.Hex()] = ChatUser{
		ID:          userB.ID,
		FirstName:   userB.FirstName,
		LastName:    userB.LastName,
		Email:       userB.Email,
		AvatarImage: userB.AvatarImage,
	}
	data.Messages = dm.Messages
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	w.WriteHeader(http.StatusOK)
}
