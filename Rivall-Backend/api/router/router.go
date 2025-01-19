package router

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"Rivall-Backend/api/resources/health"
	"Rivall-Backend/api/resources/user"
	"Rivall-Backend/api/router/middleware"
	"Rivall-Backend/api/router/middleware/requestlog"
	"Rivall-Backend/api/utils"
)

type APIRoute struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var v1Routes = []APIRoute{
	{
		Method:  "POST",
		Path:    "/users/",
		Handler: user.WriteOne,
	},
}

func AddRoutes(router *mux.Router, routes []APIRoute) {
	for _, route := range routes {
		router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Method)
			if r.Method == route.Method {
				requestlog.NewHandler(route.Handler, utils.Logger)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})
	}
}

func New(l *zerolog.Logger, v *validator.Validate, mongoClient *mongo.Client) *mux.Router {
	r := mux.NewRouter()
	// Add health routes
	r.HandleFunc("/health", health.Read)

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.ContentTypeJSON)

	// Add v1 routes
	AddRoutes(r, v1Routes)

	return r
}
