package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Rivall-Backend/api/resources/websocket"
	"Rivall-Backend/api/router"
	"Rivall-Backend/config"
	"Rivall-Backend/util/logger"
	"Rivall-Backend/util/password_recovery"
	"Rivall-Backend/util/session_manager"

	"Rivall-Backend/api/global"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//	@title			Rivall-Backend
//	@version		1.0
//	@description	This is a sample RESTful API with a CRUD

//	@contact.name	Nathaniel Reeves
//	@contact.url

// @host		localhost:8080
// @basePath	/v1

func ConnectMongoDB(ctx context.Context, c *config.Conf) *mongo.Client {
	// Connect to MongoDB
	var uri string = c.DB.MongoURI
	if uri == "" {
		log.Fatal().Msg("MongoDB URI is empty")
		panic("MongoDB URI is empty")
	}
	MongoClient, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal().Err(err).Msg("MongoDB connection failure")
	}
	log.Info().Msg("MongoDB connection initialized")

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := MongoClient.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping MongoDB")
		panic(err)
	}
	log.Info().Msg("Pinged your deployment. You successfully connected to MongoDB!")
	return MongoClient
}

func DisconnectMongoDB(MongoClient *mongo.Client) {
	if err := MongoClient.Disconnect(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("MongoDB disconnection failure")
	}
	log.Info().Msg("MongoDB disconnected")
}

var Logger *zerolog.Logger
var Validator *validator.Validate
var MongoClient *mongo.Client
var JWTSecretKey string

// func sendEmail() {
// 	message := mail.NewMsg()
// 	if err := message.From("rivall@gmail.com"); err != nil {
// 		log.Error().Err(err).Msg("Failed to set sender")
// 		return
// 	}
// 	if err := message.To("nathaniel.jacob.reeves@gmail.com"); err != nil {
// 		log.Error().Err(err).Msg("Failed to set recipient")
// 		return
// 	}
// 	message.Subject("Test Email")
// 	message.SetBodyString(mail.TypeTextPlain, "This is a test email")
// 	client, err := mail.NewClient("smtp.example.com", mail.WithSMTPAuth(mail.SMTPAuthPlain),
// 		mail.WithUsername("my_username"), mail.WithPassword("extremely_secret_pass"))
// 	if err != nil {
// 		log.Error().Err(err).Msg("Failed to create mail client")
// 	}
// 	if err := client.DialAndSend(message); err != nil {
// 		log.Error().Err(err).Msg("Failed to send mail")
// 	}
// }

func main() {

	// Send Email
	// sendEmail()

	// Initialize logger, validator, and config
	c := config.New()
	logLevel := c.Server.Debug
	l := logger.New(logLevel, c)
	v := validator.New()
	pr := password_recovery.NewRecoveryRetentionMap(context.Background())
	log.Info().Msg("Building Rivall Backend API...")

	// Initialize context
	ctx := context.Background()

	// Connect MongoDB
	MongoClient := ConnectMongoDB(ctx, c)

	// Setup Logined User Management
	Sessions := session_manager.NewSessionsManager(ctx, c.Server.JWTSecretKey)

	// Setup Websocket Management
	WSManager := websocket.NewManager(ctx)

	// Inject global variables
	global.Logger = l
	global.Validator = v
	global.MongoClient = MongoClient
	global.WSManager = WSManager
	global.SessionManager = Sessions
	global.JWTSecretKey = c.Server.JWTSecretKey
	global.PasswordRecoveryMap = pr

	// Initialize router
	r := router.New()

	// Initialize server
	// cfg := &tls.Config{
	// 	MinVersion:               tls.VersionTLS12,
	// 	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	// 	PreferServerCipherSuites: true,
	// 	CipherSuites: []uint16{
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	// 	},
	// }
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
		// TLSConfig:    cfg,
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Graceful shutdown functionality
	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msgf("Shutting down server %v:%v", c.Server.Address, c.Server.Port)

		ctx, cancel := context.WithTimeout(context.Background(), c.Server.TimeoutIdle)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server shutdown failure")
		}

		DisconnectMongoDB(MongoClient)

		close(closed)
	}()

	// Start server
	log.Info().Msgf("Starting server %v:%v", c.Server.Address, c.Server.Port)
	// if err := s.ListenAndServeTLS("fullchain.pem", "cert-key.pem"); err != nil && err != http.ErrServerClosed {
	// 	log.Fatal().Err(err).Msg("Server startup failure")
	// }
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
