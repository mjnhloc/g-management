package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
	UserID   string
	UserType string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	redis      *RedisClient
	mu         sync.RWMutex
}

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func NewHub(redis *RedisClient) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      redis,
	}
}

func (h *Hub) Run(ctx context.Context) {
	// Subscribe to Redis channels for different event types
	eventTypes := []string{"class_updates", "member_updates", "trainer_updates"}
	for _, eventType := range eventTypes {
		go h.subscribeToRedisChannel(ctx, eventType)
	}

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) subscribeToRedisChannel(ctx context.Context, channel string) {
	messages, err := h.redis.SubscribeToChannel(ctx, channel)
	if err != nil {
		log.Printf("Failed to subscribe to Redis channel %s: %v", channel, err)
		return
	}

	for msg := range messages {
		h.broadcast <- []byte(msg.Payload)
	}
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	client := &Client{
		ID:       generateID(), // Implement this function
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      h,
		UserID:   c.GetString("user_id"),
		UserType: c.GetString("user_type"),
	}

	client.Hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages if needed
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Process message based on type
		switch msg.Type {
		case "subscribe":
			// Handle subscription requests
		case "unsubscribe":
			// Handle unsubscription requests
		}
	}
}

// generateID generates a unique ID for a client
func generateID() string {
	return uuid.New().String()
}
