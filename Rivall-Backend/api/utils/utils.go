package utils

import (
	"github.com/go-playground/validator/v10"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Logger *zerolog.Logger
var Validator *validator.Validate
var MongoClient *mongo.Client

func New(l *zerolog.Logger, v *validator.Validate, mongoClient *mongo.Client) {
	// Set global variables
	Logger = l
	Validator = v
	MongoClient = mongoClient
}
