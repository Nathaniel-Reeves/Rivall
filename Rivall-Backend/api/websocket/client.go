package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager
	// egress is used to avoid concurrent writes on the WebSocket
	Egress chan Event

	// chatroom is used to know what room user is in
	chatroom string

	userID string
}

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager, userID string) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		Egress:     make(chan Event),
		userID:     userID,
	}
}

func (c *Client) Manager() *Manager {
	return c.manager
}

func (c *Client) Chatroom() string {
	return c.chatroom
}

func (c *Client) Send(message Event) {
	c.Egress <- message
}

// readMessages will start the client to read messages and handle them
// appropriatly.
// This is suppose to be ran as a goroutine
func (c *Client) readMessages() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.manager.removeClient(c)
	}()
	// Set Max Size of Messages in Bytes
	c.connection.SetReadLimit(512)
	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Err(err).Msg("error setting read deadline")
		return
	}
	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Err(err).Msg("error reading message")
			}
			break // Break the loop to close conn & Cleanup
		}
		// Marshal incoming data into a Event struct
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Err(err).Msg("error unmarshaling message")
			break // Breaking the connection here might be harsh xD
		}
		// Route the Event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Err(err).Msg("error routing event")
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Debug().Msg("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	// Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		// Graceful close if this triggers a closing
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.Egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Err(err).Msg("error writing close message")
				}
				// Return to close the goroutine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Err(err).Msg("error marshaling message")
				return // closes the connection, should we really
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Err(err).Msg("error writing message")
			}
			log.Debug().Msg("message sent")
		case <-ticker.C:
			log.Debug().Msg("ping")
			// Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Info().Msg("ping failed")
				return // return to break this goroutine triggeing cleanup
			}
		}

	}
}
