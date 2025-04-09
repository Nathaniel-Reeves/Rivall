package globals

import (
	"Rivall-Backend/util/password_recovery"
	"Rivall-Backend/util/session_manager"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Logger *zerolog.Logger
var Validator *validator.Validate
var MongoClient *mongo.Client
var JWTSecretKey string
var SessionManager *session_manager.Sessions
var PasswordRecoveryMap *password_recovery.RecoveryRetentionMap
