package router

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"Rivall-Backend/api/resources/health"
	"Rivall-Backend/api/resources/user"
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
	r.HandleFunc("/api/v1/users", user.GetUserByUsername).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", user.GetUserById).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", user.PostUser).Methods(http.MethodPost)

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.ContentTypeJSON)
	r.Use(middleware.RequestLogging)

	return r
}
