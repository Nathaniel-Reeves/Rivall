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
	clients ClientList

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler
	// otps is a map of allowed OTP to accept connections from
	otps RetentionMap
}

// NewManager is used to initalize all the values inside the manager
func NewManager(ctx context.Context) *Manager {
	log.Info().Msg("Creating new Websocket Manager")
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		// Create a new retentionMap that removes Otps older than 5 seconds
		otps: NewRetentionMap(ctx, 5*time.Second),
	}
	m.setupEventHandlers()
	log.Info().Msg("Websocket Manager Created")
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventChangeRoom] = ChatRoomHandler
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// loginHandler is used to verify an user authentication and return a one time password
func (m *Manager) CreateOTP() string {

	// add a new OTP
	otp := m.otps.NewOTP()

	return otp.Key
}

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	log.Debug().Msg("Serving Websocket Connection")

	// Grab the User ID
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		// Tell the user its not authorized
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check user ID is valid
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Grab the OTP in the Get param
	otp := r.URL.Query().Get("otp")
	if otp == "" {
		// Tell the user its not authorized
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify OTP is existing
	log.Debug().Msgf("Verifying OTP: %s", otp)
	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Info().Msg("Unauthorized Connection")
		return
	}

	// Begin by upgrading the HTTP request
	log.Info().Msg("Upgrading to Websocket Connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("Failed to upgrade connection to websocket")
		return
	}

	// Create New Client
	client := NewClient(conn, m, userID)
	// Add the newly created client to the manager
	m.addClient(client)
	log.Debug().Msg("Client Added to Manager")

	go client.readMessages()
	go client.writeMessages()
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true

	log.Debug().Msg("Client Added")
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clients, client)
		log.Debug().Msg("Client removed")
	}
}

// remove client by user id
func (m *Manager) RemoveClientByUserID(userID string) {
	m.Lock()
	defer m.Unlock()

	for client := range m.clients {
		if client.userID == userID {
			// close connection
			client.connection.Close()
			// remove
			delete(m.clients, client)
			log.Debug().Msg("Client removed")
			return
		}
	}

	log.Warn().Msg("Client not found")
}
