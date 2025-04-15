package resources

import (
	"encoding/json"
	"net/http"

	db "Rivall-Backend/db"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET user")

	// get user id from url parameters
	vars := mux.Vars(r)
	userID := vars["user_id"]
	log.Debug().Msgf("User ID: %s", userID)

	// check user exists
	user := db.ReadByUserIdWithPopulatedFields(userID)
	if user.ID == bson.NilObjectID {
		log.Error().Msg("User does not exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist."))
		return
	}

	// return user
	json.NewEncoder(w).Encode(user)
}
