package websocket

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var (
	/**
	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		// Apply the Origin Checker
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// checkOrigin will check origin and return true if its allowed
func checkOrigin(r *http.Request) bool {

	// Grab the request origin
	origin := r.Header.Get("Origin")

	log.Debug().Msgf("Checking Origin: %s", origin)

	switch origin {
	// Update this to HTTPS
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clients ClientMap

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler
	// otps is a map of allowed OTP to accept connections from
	otps RetentionMap
}

func (m *Manager) Clients() ClientMap {
	return m.clients
}

func NewManager(ctx context.Context) *Manager {
	log.Info().Msg("Creating new Websocket Manager")
	m := &Manager{
		clients:  make(ClientMap),
		handlers: make(map[string]EventHandler),
		// Create a new retentionMap that removes Otps older than 5 seconds
		otps: NewRetentionMap(ctx, 5*time.Second),
	}
	m.setupEventHandlers()
	log.Info().Msg("Websocket Manager Created")
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventCreateGroup] = CreateGroupHandler
	m.handlers[EventAcceptGroupRequest] = AcceptGroupRequestHandler
	m.handlers[EventRejectGroupRequest] = RejectGroupRequestHandler
	m.handlers[EventSendGroupMessage] = SendGroupMessageHandler
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) CreateOTP() string {
	otp := m.otps.NewOTP()
	return otp.Key
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	log.Debug().Msg("Serving Websocket Connection")

	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Debug().Msgf("Verifying OTP: %s", otp)
	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Info().Msg("Unauthorized Connection")
		return
	}

	log.Info().Msg("Upgrading to Websocket Connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("Failed to upgrade connection to websocket")
		return
	}

	client := NewClient(conn, m, userID)
	m.addClient(client)
	log.Debug().Msg("Client Added to Manager")

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client.userID] = client

	log.Debug().Msg("Client Added")
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[client.userID]; ok {
		client.connection.Close()
		delete(m.clients, client.userID)
		log.Debug().Msg("Client Removed")
	} else {
		log.Warn().Msg("Client not found")
	}
}

// remove client by user id
func (m *Manager) RemoveClientByUserID(userID string) {
	m.Lock()
	defer m.Unlock()

	if client, ok := m.clients[userID]; ok {
		client.connection.Close()
		delete(m.clients, userID)
		log.Debug().Msg("Client Removed")
	} else {
		log.Warn().Msg("Client not found")
	}
}

var WSManager = NewManager(context.Background())
