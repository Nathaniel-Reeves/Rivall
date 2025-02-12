package global

import (
	"Rivall-Backend/api/resources/websocket"

	"github.com/go-playground/validator/v10"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Logger *zerolog.Logger
var Validator *validator.Validate
var MongoClient *mongo.Client
var WSManager *websocket.Manager
var JWTSecretKey string
