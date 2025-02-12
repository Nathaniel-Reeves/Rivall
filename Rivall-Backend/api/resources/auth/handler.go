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
	globals "Rivall-Backend/api/global"
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
	if users.ReadByUserEmail(user.Email).ID != bson.NilObjectID {
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
	user := users.ReadByUserEmail(userLogin.Email)
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

	// create token
	token, err := createToken(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create token"))
		return
	}
	w.Header().Set("Authorization", token)

	// Create OTP for websocket
	user.OTP = globals.WSManager.CreateOTP()

	// Clean Response Data
	type sendUser struct {
		ID        string `json:"id" bson:"_id,omitempty"`
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		Email     string `json:"email" bson:"email"`
		OTP       string `json:"otp" bson:"otp"`
	}

	su := sendUser{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		OTP:       user.OTP,
	}

	// return user
	json.NewEncoder(w).Encode(su)
}

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
	user := users.ReadByUserEmail(userID)
	if user.ID == bson.NilObjectID {
		log.Warn().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// remove token from header
	w.Header().Set("Authorization", "")
	w.WriteHeader(http.StatusOK)

	// The websocket connection will be closed via timeout
}
