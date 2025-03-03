package session_manager

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Sessions struct {
	// sessions is a map of user_id to session_id
	Sessions map[string]Session

	// JWTSecretKey is the secret key used to sign JWT tokens
	JWTSecretKey string

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex
}

type Session struct {
	// session is a struct that holds the user_id and session_id
	UserID         string
	SessionID      string
	Token          string
	TokenExpiresAt time.Time
	Type           string
}

const ACCESS_TOKEN_TIMEOUT = time.Minute * 30

// const ACCESS_TOKEN_TIMEOUT = time.Second * 10
const REFRESH_TOKEN_TIMEOUT = time.Hour * 24

func NewSessionsManager(ctx context.Context, JWTSecretKey string) *Sessions {
	s := Sessions{
		Sessions:     make(map[string]Session),
		JWTSecretKey: JWTSecretKey,
	}

	go s.ExpireSessions(ctx)

	return &s
}

func (s *Sessions) GetJWTSecretKey() string {
	return s.JWTSecretKey
}

func (s *Sessions) CreateAccessToken(userID string) (string, time.Time) {
	exp := s.GetAccessTokenTimeout()
	return s.CreateJWTToken(userID, exp), exp
}

func (s *Sessions) CreateRefreshToken(userID string) (string, time.Time) {
	exp := s.GetRefreshTokenTimeout()
	return s.CreateJWTToken(userID, exp), exp
}

func (s *Sessions) GetAccessTokenTimeout() time.Time {
	return time.Now().Add(ACCESS_TOKEN_TIMEOUT)
}

func (s *Sessions) GetRefreshTokenTimeout() time.Time {
	return time.Now().Add(REFRESH_TOKEN_TIMEOUT)
}

func (s *Sessions) CreateJWTToken(userID string, exp time.Time) string {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.JWTSecretKey))
	if err != nil {
		log.Println(err)
	}
	return tokenString
}

func (s *Sessions) ValidateJWTToken(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecretKey), nil
	})

	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false
	}

	return claims, true
}

func (s *Sessions) NewAccessSession(userID string) Session {
	s.Lock()
	defer s.Unlock()

	accessToken, accessTokenExpiresAt := s.CreateAccessToken(userID)

	new_session := Session{
		UserID:         userID,
		Token:          accessToken,
		TokenExpiresAt: accessTokenExpiresAt,
		Type:           "access",
	}

	s.Sessions[accessToken] = new_session

	return new_session
}

func (s *Sessions) NewRefreshSession(userID string) Session {
	s.Lock()
	defer s.Unlock()

	refreshToken, refreshTokenExpiresAt := s.CreateRefreshToken(userID)

	new_session := Session{
		UserID:         userID,
		Token:          refreshToken,
		TokenExpiresAt: refreshTokenExpiresAt,
		Type:           "refresh",
	}

	s.Sessions[refreshToken] = new_session

	return new_session
}

func (s *Sessions) GetSession(token string) (Session, bool) {
	s.Lock()
	defer s.Unlock()

	session, ok := s.Sessions[token]
	return session, ok
}

func (s *Sessions) DeleteSession(token string) {
	s.Lock()
	defer s.Unlock()

	delete(s.Sessions, token)
}

func (s *Sessions) ExpireSessions(ctx context.Context) {
	ticker := time.NewTicker(400 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			s.Lock()
			for key, session := range s.Sessions {
				if session.TokenExpiresAt.Before(time.Now()) {
					delete(s.Sessions, key)
				}
			}
			s.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
