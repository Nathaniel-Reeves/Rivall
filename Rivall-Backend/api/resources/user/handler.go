package user

import (
	"encoding/json"
	"net/http"

	"Rivall-Backend/api/utils"
)

func WriteOne(w http.ResponseWriter, r *http.Request) {
	utils.Logger.Info().Msg("Writing user")

	// get user data
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to decode user")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate user data
	// err = utils.Validator.Struct(user)
	// if err != nil {
	// 	utils.Logger.Error().Err(err).Msg("Failed to validate user")
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// insert user data
	err = Insert(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return success
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read user
//	@description	Read user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"user ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/users/{id} [get]
