package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func createToken(user User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Error().Msg("User ID not found in context")
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// Register New User
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
//	@router			/auth/register [post]
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST new user")

	// get user data
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user, invalid JSON request."))
		return
	}

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

// Login User
//
//	@summary		Login User
//	@description	Log in a User
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/auth/login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET user by username")

	// get user data
	userLogin := User{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user"))
		return
	}

	// get user from db
	user := ReadUsername(userLogin.Username)
	if user.ID == bson.NilObjectID {
		log.Warn().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// check password
	if !ComparePasswords(user.Password, userLogin.Password) {
		log.Warn().Msg("Invalid password")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid password"))
		return
	}

	token, err := createToken(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create token"))
		return
	}

	w.Header().Set("Authorization", token)

	// set user online
	user.Online = true
	UpdateUserOnline(user)

	// return user
	json.NewEncoder(w).Encode(user)
}

// Logout User
//
//	@summary		Logout User
//	@description	Log out a User
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/auth/logout [delete]
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("DELETE user login session")

	// get user id from context
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user ID from context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user ID from context"))
		return
	}

	// get user from db
	user := ReadId(userID)
	if user.ID == bson.NilObjectID {
		log.Warn().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// remove token from header
	w.Header().Set("Authorization", "")

	// set user offline
	user.Online = false
	UpdateUserOnline(user)

	w.WriteHeader(http.StatusOK)
}

// Reading data from url parameters
// func LoginUser(w http.ResponseWriter, r *http.Request) {
// 	log.Info().Msg("GET user by username")

// 	// get user data
// 	vars := mux.Vars(r)
// 	username := vars["username"]
// 	password := vars["password"]

// 	// get user from db
// 	user := ReadUsername(username)
// 	if user.ID == bson.NilObjectID {
// 		log.Warn().Msg("User not found")
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("User not found"))
// 		return
// 	}

// 	// check password
// 	if !ComparePasswords(user.Password, password) {
// 		log.Warn().Msg("Invalid password")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte("Invalid password"))
// 		return
// 	}

// 	// return user
// 	json.NewEncoder(w).Encode(user)
// }
