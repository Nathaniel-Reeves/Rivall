package router

import (
	"net/http"

	"github.com/gorilla/mux"

	auth "Rivall-Backend/api/resources/auth"
	"Rivall-Backend/api/resources/health"
	messagegroup "Rivall-Backend/api/resources/message-group"
	users "Rivall-Backend/api/resources/users"
	"Rivall-Backend/api/router/middleware"

	globals "Rivall-Backend/api/global"
)

func New() *mux.Router {
	r := mux.NewRouter()

	// Add health routes
	r.HandleFunc("/health", health.Read).Methods(http.MethodGet)

	// Add v1 routes
	publicRouter := r.PathPrefix("/api/v1").Subrouter()
	publicRouter.HandleFunc("/auth/register", auth.RegisterNewUser).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/login", auth.LoginUser).Methods(http.MethodPost)

	privateRouter := r.PathPrefix("/api/v1/ws").Subrouter()
	privateRouter.Use(middleware.AuthenticationMiddleware)
	privateRouter.HandleFunc("/", globals.WSManager.ServeWS)
	privateRouter.HandleFunc("/auth/logout", auth.LogoutUser).Methods(http.MethodDelete)
	privateRouter.HandleFunc("/users/{user_id}", users.GetUser).Methods(http.MethodGet)
	privateRouter.HandleFunc("/users/{user_id}/contacts", users.PostUserContact).Methods(http.MethodPost)
	privateRouter.HandleFunc("/users/{user_id}/contacts", users.DeleteUserContact).Methods(http.MethodDelete)

	privateRouter.HandleFunc("/message-group", messagegroup.WriteNewMessageGroup).Methods(http.MethodPost)
	privateRouter.HandleFunc("/message-group/{group_id}/user-request", messagegroup.AcceptGroupRequest).Methods(http.MethodPut)
	privateRouter.HandleFunc("/message-group/{group_id}/user-request", messagegroup.RejectGroupRequest).Methods(http.MethodDelete)

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.ContentTypeJSON)
	r.Use(middleware.RequestLogging)

	return r
}
