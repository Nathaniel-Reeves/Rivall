package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Rivall-Backend/api/router"
	"Rivall-Backend/api/utils"
	"Rivall-Backend/config"
	"Rivall-Backend/util/logger"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

//	@title			Rivall-Backend
//	@version		1.0
//	@description	This is a sample RESTful API with a CRUD

//	@contact.name	Nathaniel Reeves
//	@contact.url

//	@license.name	Apache License Version 2.0, January 2004
//	@license.url	http://www.apache.org/licenses/

// @host		localhost:8080
// @basePath	/v1

func ConnectMongoDB(l *zerolog.Logger, c *config.Conf) *mongo.Client {
	// Connect to MongoDB
	var uri string = c.DB.MongoURI
	if uri == "" {
		l.Fatal().Msg("MongoDB URI is empty")
		panic("MongoDB URI is empty")
	}
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		l.Fatal().Err(err).Msg("MongoDB connection failure")
	}
	l.Info().Msg("MongoDB connection initialized")

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := mongoClient.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	l.Info().Msg("Pinged your deployment. You successfully connected to MongoDB!")
	return mongoClient
}

func DisconnectMongoDB(l *zerolog.Logger, mongoClient *mongo.Client) {
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		l.Fatal().Err(err).Msg("MongoDB disconnection failure")
	}
	l.Info().Msg("MongoDB disconnected")
}

func main() {
	// Initialize logger, validator, and config
	c := config.New()
	logLevel := c.Server.Debug
	l := logger.New(logLevel)
	v := validator.New()
	l.Info().Msg("Building Rivall Backend API...")
	fmt.Println("Logging Level: ", logLevel)

	// Connect MongoDB
	mongoClient := ConnectMongoDB(l, c)

	// Initialize uility globals
	utils.New(l, v, mongoClient)

	// Initialize router
	r := router.New(l, v, mongoClient)

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

		l.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), c.Server.TimeoutIdle)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		DisconnectMongoDB(l, mongoClient)

		close(closed)
	}()

	// Start server
	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
