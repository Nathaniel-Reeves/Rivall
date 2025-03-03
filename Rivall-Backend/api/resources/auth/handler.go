package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"

	globals "Rivall-Backend/api/global"
	users "Rivall-Backend/api/resources/users"
)

type UserRes struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	AvatarImage string `json:"avatar_image"`
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Error().Msg("User ID not found in context")
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

type RecoveryCodeReq struct {
	Email string `json:"email"`
}

func SendAccountRecoveryEmail(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST account recovery email")

	// get user data
	emailReq := RecoveryCodeReq{}
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode user"))
		return
	}

	// check email is valid
	if users.ReadByUserEmail(emailReq.Email).ID == bson.NilObjectID {
		log.Error().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// create recovery code
	recoveryCode := globals.PasswordRecoveryMap.NewRecoveryOTP(emailReq.Email)

	// send email with code
	// TODO: make email service, log recovery code to console for now
	log.Info().Str("recovery_code", recoveryCode.Code).Msg("Recovery code sent")

	// return success
	w.WriteHeader(http.StatusOK)
}

type RecoveryCodeValidationReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func ValidateAccountRecoveryCode(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST validate account recovery code")

	// get body data
	req := RecoveryCodeValidationReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode request"))
		return
	}

	// get code
	code := req.Code
	if code == "" {
		log.Error().Msg("Code must be provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Code must be provided"))
		return
	}

	// get email
	email := req.Email
	if email == "" {
		log.Error().Msg("Email must be provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email must be provided"))
		return
	}

	// check code is valid for that email
	if !globals.PasswordRecoveryMap.VerifyRecoveryOTP(code, email) {
		log.Error().Msg("Invalid code")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid code"))
		return
	}

	// login user
	// get user from db
	user := users.ReadByUserEmail(email)
	if user.ID == bson.NilObjectID {
		log.Warn().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// Create Access Session
	accessSession := globals.SessionManager.NewAccessSession(user.ID.Hex())

	// Create Refresh Session
	refreshSession := globals.SessionManager.NewRefreshSession(user.ID.Hex())

	// Create OTP for websocket
	user.OTP = globals.WSManager.CreateOTP()

	// Clean Response Data
	su := UserRes{
		ID:          user.ID.Hex(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		AvatarImage: user.AvatarImage,
	}

	res := LoginUserRes{
		OTP:                   user.OTP,
		AccessToken:           accessSession.Token,
		RefreshToken:          refreshSession.Token,
		AccessTokenExpiresAt:  accessSession.TokenExpiresAt,
		RefreshTokenExpiresAt: refreshSession.TokenExpiresAt,
		User:                  su,
	}

	// return user
	json.NewEncoder(w).Encode(res)
}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST update user password")

	// get user id from path
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		log.Error().Msg("User ID must be provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User ID must be provided"))
		return
	}

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
	if user.ID == bson.NilObjectID || user.Password == "" {
		log.Error().Msg("User must have all required fields")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User must have all required fields."))
		return
	}

	// check user id matches the id in the path
	if user.ID.Hex() != userID {
		log.Error().Msg("User ID does not match path")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User ID does not match path"))
		return
	}

	// check user exists
	userFromDB := users.ReadByUserId(user.ID.Hex())
	if userFromDB.ID == bson.NilObjectID {
		log.Error().Msg("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// update password
	err = users.UpdateUserPassword(user.ID.Hex(), user.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update password")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update password"))
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
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

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	OTP                   string    `json:"otp"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  UserRes   `json:"user"`
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

	// Create Access Session
	accessSession := globals.SessionManager.NewAccessSession(user.ID.Hex())

	// Create Refresh Session
	refreshSession := globals.SessionManager.NewRefreshSession(user.ID.Hex())

	// Create OTP for websocket
	user.OTP = globals.WSManager.CreateOTP()

	// Clean Response Data
	su := UserRes{
		ID:          user.ID.Hex(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		AvatarImage: user.AvatarImage,
	}

	res := LoginUserRes{
		OTP:                   user.OTP,
		AccessToken:           accessSession.Token,
		RefreshToken:          refreshSession.Token,
		AccessTokenExpiresAt:  accessSession.TokenExpiresAt,
		RefreshTokenExpiresAt: refreshSession.TokenExpiresAt,
		User:                  su,
	}

	json.NewEncoder(w).Encode(res)
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST renew access token")

	// Get user id from path
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		log.Error().Msg("User ID must be provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User ID must be provided"))
		return
	}

	// get user data
	req := RenewAccessTokenReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode request"))
		return
	}

	// check user has all required fields
	if req.RefreshToken == "" {
		log.Error().Msg("Refresh token must be provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Refresh token must be provided"))
		return
	}

	// check user id matches the id in the path
	if userID != userID {
		log.Error().Msg("User ID does not match path")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User ID does not match path"))
		return
	}

	// get session from refresh token
	session, ok := globals.SessionManager.GetSession(req.RefreshToken)

	// check if session exists
	if !ok {
		log.Error().Msg("Session not found")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session not found"))
		return
	}

	// check session is a refresh token
	if session.Type != "refresh" {
		log.Error().Msg("Session is not a refresh token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session is not a refresh token"))
		return
	}

	// check session matches the user id
	if session.UserID != session.UserID {
		log.Error().Msg("Session does not match user ID")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session does not match user ID"))
		return
	}

	// Create Access Session
	accessSession := globals.SessionManager.NewAccessSession(session.UserID)

	// Clean Response Data
	res := RenewAccessTokenRes{
		AccessToken:          accessSession.Token,
		AccessTokenExpiresAt: accessSession.TokenExpiresAt,
	}

	json.NewEncoder(w).Encode(res)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("DELETE user login session")
	// Remember to delete the authorization token from the client side after this request

	// Get user id from path
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		log.Error().Msg("User ID must be provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User ID must be provided"))
		return
	}

	// Get Token from Header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		log.Error().Msg("Token not found in header")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Token not found in header"))
		return
	}

	// Delete Session
	globals.SessionManager.DeleteSession(tokenString)

	// remove token from header, the token will eventually expire
	w.Header().Set("Authorization", "")
	w.WriteHeader(http.StatusOK)

	// Close Websocket Connection
	globals.WSManager.RemoveClientByUserID(userID)
}
