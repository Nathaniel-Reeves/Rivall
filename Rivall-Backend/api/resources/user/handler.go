package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Create User
//
//	@summary		Create User
//	@description	Create a User in the database
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/users [post]
func PostUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Writing user")

	// get user data
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user"))
		return
	}

	// validate user data
	// err = router.Validator.Struct(user)
	// if err != nil {
	// 	router.Logger.Error().Err(err).Msg("Failed to validate user")
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// insert user data
	err = CreateUser(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert user"))
		return
	}

	// return success
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read user
//	@description	Read user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/users/{id} [get]
func GetUserById(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Reading user")

	// get user ID /users/{:id}
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok == false {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// get user data
	user := ReadId(id)

	// handle not found user
	if user.ID == bson.NilObjectID {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// return user data
	json.NewEncoder(w).Encode(user)
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Reading user")

	// get username /users?username={:username}
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// get user data
	user := ReadUsername(username)

	// handle not found user
	if user.ID == bson.NilObjectID {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// return user data
	json.NewEncoder(w).Encode(user)
}
