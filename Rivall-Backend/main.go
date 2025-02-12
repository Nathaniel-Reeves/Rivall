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

func main() {
	// Initialize logger, validator, and config
	c := config.New()
	logLevel := c.Server.Debug
	l := logger.New(logLevel, c)
	v := validator.New()
	log.Info().Msg("Building Rivall Backend API...")

	// Initialize context
	ctx := context.Background()

	// Connect MongoDB
	MongoClient := ConnectMongoDB(ctx, c)

	// Setup Websocket Management
	WSManager := websocket.NewManager(ctx)

	// Initialize global variables
	Logger = l
	Validator = v

	// Initialize router
	r := router.New(l, v, MongoClient)

	// Inject global variables
	global.Logger = l
	global.Validator = v
	global.MongoClient = MongoClient
	global.WSManager = WSManager
	global.JWTSecretKey = c.Server.JWTSecretKey

	// Initialize server
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
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
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
