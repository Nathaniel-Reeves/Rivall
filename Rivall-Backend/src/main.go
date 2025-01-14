package main

import (
	"context"
	"flag"
	"log"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This gives global access to the mongodb collection
var collection_users *mongo.Collection
var client *mongo.Client

type User struct {
	ID int64 `bson: "_id,omitempty" json: "_id,omitempty"` 
	FullText string `bson: "full_text,omitempty" json: "full_text,omitempty"`
	Contacts struct {
		UserID int64 `bson: "user_id" json: "user_id"`
	} `bson: "contacts,omitempty" json: "contacts,omitempty"`
}

func getUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users []User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection_users.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	if err = cursor.All(ctx, &users); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(users)
}

func searchUsers(response http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	response.Header().Set("content-type", "application/json")
	var users []User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	searchStage := bson.D{
		{ "$search", bson.D{
			{ "index", "default" },
			{ "text", bson.D{
				{ "query", queryParams.Get("query") },
				{ "path", "full_text" },
			}},
		}},
	}
	cursor, err := collection_users.Aggregate(ctx, mongo.Pipeline{searchStage})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	if err = cursor.All(ctx, &users); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(users)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}

func main() {
	fmt.Println("Starting Rivall Backend...")

	// Connect MongoDB
	mongo_ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(mongo_ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(mongo_ctx); err != nil {
			panic(err)
		}
	}()
	
collection_users = client.Database("Rivall-DB").Collection("Users")

	// Start Server
	var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()

    r := mux.NewRouter()
    // Add your routes as needed
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/v1/users", getUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/search", searchUsers).Methods("GET")

    srv := &http.Server{
        Addr:         "0.0.0.0:8080",
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
    }

    // Run our server in a goroutine so that it doesn't block.
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()

    c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    srv.Shutdown(ctx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}
