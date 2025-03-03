package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"Rivall-Backend/api/global"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ====================
		// Authenticate User
		// ====================

		// get token from header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization header"))
			return
		}

		// remove "Bearer " from token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// get token claims
		claims, ok := global.SessionManager.ValidateJWTToken(tokenString)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
			return
		}

		// get user_id from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token: user_id"))
			return
		}

		// check if session exists
		session, ok := global.SessionManager.GetSession(tokenString)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token: token not found"))
			return
		}

		// check if session is valid
		if session.TokenExpiresAt.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token expired"))
			return
		}

		// ====================
		// Authorize User
		// ====================
		vars := mux.Vars(r)

		// if user_id in url, check if it matches the token user_id
		if vars["user_id"] != "" {
			if vars["user_id"] != userID {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized user"))
				return
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		log.Debug().Str("user_id", userID).Msg("Authenticated user")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
