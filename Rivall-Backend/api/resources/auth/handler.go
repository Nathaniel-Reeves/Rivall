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

	"Rivall-Backend/api/global"
	users "Rivall-Backend/api/resources/users"
)

func createToken(user users.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.JWTSecretKey))
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Error().Msg("User ID not found in context")
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// Test Authorization
//
//	@summary		Test Authorization
//	@description	Test Authorization
//	@tags			users
//	@accept			json
//	@produce		json
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/auth/authorize [get]
func TestAuthorization(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("User is Authorized")
	w.Write([]byte("Users is Authorized"))
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
	user := users.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user, invalid JSON request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user, invalid JSON request."))
		return
	}

	// check user has all required fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		log.Error().Msg("User must have all required fields")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User must have all required fields."))
		return
	}

	// check user does not already exist
	if users.ReadEmail(user.Email).ID != bson.NilObjectID {
		log.Error().Msg("User already exists")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User already exists."))
		return
	}

	// insert user data
	err = users.CreateUser(user)
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
	userLogin := users.User{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user"))
		return
	}

	// get user from db
	user := users.ReadEmail(userLogin.Email)
	if user.ID == bson.NilObjectID {
		log.Warn().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// check password
	if !users.ComparePasswords(user.Password, userLogin.Password) {
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
	// Remember to delete the authorization token from the client side after this request

	// get user id from context
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user ID from context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user ID from context"))
		return
	}

	// get user from db
	user := users.ReadId(userID)
	if user.ID == bson.NilObjectID {
		log.Warn().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// remove token from header
	w.Header().Set("Authorization", "")
	w.WriteHeader(http.StatusOK)
}
