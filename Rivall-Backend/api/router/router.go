package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"Rivall-Backend/api/resources"
	"Rivall-Backend/api/router/middleware"
	"Rivall-Backend/api/websocket"
)

var WSManager *websocket.Manager

func New() *mux.Router {
	r := mux.NewRouter()

	// Add health routes
	r.HandleFunc("/health", resources.Read).Methods(http.MethodGet)

	// Add v1 routes
	publicRouter := r.PathPrefix("/api/v1").Subrouter()
	publicRouter.HandleFunc("/auth/register", resources.RegisterNewUser).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/login", resources.LoginUser).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/recovery/send-code", resources.SendAccountRecoveryEmail).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/recovery/validate-code", resources.ValidateAccountRecoveryCode).Methods(http.MethodPost)
	publicRouter.HandleFunc("/contacts/{user_id}", resources.GetContact).Methods(http.MethodGet)

	privateRouter := r.PathPrefix("/api/v1").Subrouter()
	privateRouter.Use(middleware.AuthMiddleware)

	privateRouter.HandleFunc("/users/{user_id}", resources.GetUser).Methods(http.MethodGet)
	privateRouter.HandleFunc("/auth/recovery/{user_id}/reset-password", resources.UpdateUserPassword).Methods(http.MethodPut)
	privateRouter.HandleFunc("/auth/{user_id}/refresh", resources.RenewAccessToken).Methods(http.MethodPost)
	privateRouter.HandleFunc("/auth/{user_id}/logout", resources.LogoutUser).Methods(http.MethodDelete)
	privateRouter.HandleFunc("/users/{user_id}", resources.GetUser).Methods(http.MethodGet)
	privateRouter.HandleFunc("/users/{user_id}/contacts", resources.PostUserContact).Methods(http.MethodPost)
	privateRouter.HandleFunc("/users/{user_id}/contacts/{chat_id}/chat", resources.GetChat).Methods(http.MethodGet)
	// privateRouter.HandleFunc("/contacts/{user_id}", resources.GetContact).Methods(http.MethodGet)

	privateWSRouter := r.PathPrefix("/api/v1/ws").Subrouter()
	privateWSRouter.HandleFunc("/connect/{user_id}", websocket.WSManager.ServeWS)

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.ContentTypeJSON)
	r.Use(middleware.RequestLogging)
	r.Use(middleware.SecureConnection)

	return r
}
