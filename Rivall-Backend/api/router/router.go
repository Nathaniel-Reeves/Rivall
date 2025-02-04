package router

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"

	auth "Rivall-Backend/api/resources/auth"
	"Rivall-Backend/api/resources/health"
	message_group "Rivall-Backend/api/resources/message-group"
	"Rivall-Backend/api/router/middleware"
)

var Logger *zerolog.Logger
var Validator *validator.Validate
var MongoClient *mongo.Client

func New(l *zerolog.Logger, v *validator.Validate, mongoClient *mongo.Client) *mux.Router {
	r := mux.NewRouter()

	// Set global variables
	Logger = l
	Validator = v
	MongoClient = mongoClient

	// Add health routes
	r.HandleFunc("/health", health.Read).Methods(http.MethodGet)

	// Add v1 routes
	publicRouter := r.PathPrefix("/api/v1").Subrouter()
	publicRouter.HandleFunc("/auth/register", auth.RegisterNewUser).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/login", auth.LoginUser).Methods(http.MethodPost)

	privateRouter := r.PathPrefix("/api/v1").Subrouter()
	privateRouter.Use(middleware.AuthenticationMiddleware)
	privateRouter.HandleFunc("/auth/logout", auth.LogoutUser).Methods(http.MethodDelete)
	privateRouter.HandleFunc("/messagegroups", message_group.PostNewMessageGroup).Methods(http.MethodPost)

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.ContentTypeJSON)
	r.Use(middleware.RequestLogging)

	return r
}
