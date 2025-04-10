package resources

import (
	"encoding/json"
	"net/http"

	db "Rivall-Backend/db"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type NewGroupReq struct {
	GroupName      string   `json:"group_name"`
	RequestMessage string   `json:"message"`
	UserIDs        []string `json:"user_ids"`
}

func WriteNewMessageGroup(w http.ResponseWriter, r *http.Request) {
	// Create new message group in database

	// Validate request
	// Check the admin user is a valid logged in user
	vars := mux.Vars(r)
	adminUserID := vars["user_id"]
	if db.ReadByUserId(adminUserID).ID == bson.NilObjectID {
		log.Error().Msg("Admin user does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Admin user does not exist."))
		return
	}

	// Check all users requested to be added to the group are valid users
	body := make(map[string][]string)
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user IDs, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user IDs, invalid JSON request."))
		return
	}
	for _, userID := range body["user_ids"] {
		if db.ReadByUserId(userID).ID == bson.NilObjectID {
			log.Error().Msg("User does not exist")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User does not exist."))
			return
		}
	}

	// Check the group name is valid
	if len(body["group_name"]) == 0 {
		log.Error().Msg("Group name is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Group name is empty."))
		return
	}

	// Check the message is valid
	if len(body["message"]) == 0 {
		log.Error().Msg("Message is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message is empty."))
		return
	}

	// Create new message group in database, Only add the Admin user to the group
	insertID, err2 := db.CreateGroup(body["group_name"][0], adminUserID)
	if err2 != nil { // check insertID is type string
		log.Error().Err(err).Msg("Failed to create new message group")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create new message group."))
		return
	}

	// Send Group Request to all users requested to be added to the group
	failedRequests := make([]map[string]interface{}, 0)
	for _, userID := range body["user_ids"] {
		err3 := db.CreateGroupRequest(adminUserID, userID, insertID, body["message"][0])
		if err3 != nil {
			log.Error().Err(err).Msg("Failed to send group request")
			r := map[string]interface{}{
				"userID":  userID,
				"error":   err3,
				"message": "Failed to send group request",
			}
			failedRequests = append(failedRequests, r)
		}
	}
	if len(failedRequests) > 0 {
		log.Error().Msg("Failed to send group request to some users")
		w.WriteHeader(http.StatusInternalServerError)
		response, err := json.Marshal(failedRequests)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal failed requests")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal failed requests."))
			return
		}
		w.Write(response)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message group created."))
}

func AcceptGroupRequest(w http.ResponseWriter, r *http.Request) {
	// Accept group request

	// Validate request, Check the user is a valid logged in user
	body := make(map[string][]string)
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user IDs, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user IDs, invalid JSON request."))
		return
	}
	// Check the group is a valid group
	groupID := body["group_id"][0]
	groupObj := db.ReadByGroupId(groupID)
	if groupObj.ID == bson.NilObjectID {
		log.Error().Msg("Group does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Group does not exist."))
		return
	}
	// Check the user is a valid user
	if db.ReadByUserId(body["user_id"][0]).ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// Add the group to the user's groups and update the message request status
	err2 := db.AcceptGroupRequest(body["user_id"][0], groupID)
	if err2 != nil {
		log.Error().Err(err2).Msg("Failed to accept group request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to accept group request."))
		return
	}

	// Add the user to the group
	err3 := db.AddUserToGroup(body["user_id"][0], groupID)
	if err3 != nil {
		log.Error().Err(err3).Msg("Failed to add user to group")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to add user to group."))
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}

func RejectGroupRequest(w http.ResponseWriter, r *http.Request) {
	// Reject group request

	// Validate request, Check the user is a valid logged in user
	// Check the group is a valid group
	body := make(map[string][]string)
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user IDs, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user IDs, invalid JSON request."))
		return
	}

	// Check the user is a valid user
	if db.ReadByUserId(body["user_id"][0]).ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// Update the group request status
	err2 := db.RejectGroupRequest(body["user_id"][0], body["group_id"][0])
	if err2 != nil {
		log.Error().Err(err2).Msg("Failed to reject group request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reject group request."))
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}
