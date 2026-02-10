package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// ─── Message types ─────────────────────────────────────────

type IncomingMessage struct {
	Type       string `json:"type"`       // "chat" | "typing" | "read"
	ReceiverID int    `json:"receiverId"` // who we're talking to
	Body       string `json:"body,omitempty"`
	MessageID  int    `json:"messageId,omitempty"` // for read receipts
}

type OutgoingMessage struct {
	Type       string `json:"type"` // "chat" | "typing" | "read" | "online" | "offline"
	SenderID   int    `json:"senderId"`
	ReceiverID int    `json:"receiverId,omitempty"`
	Body       string `json:"body,omitempty"`
	MessageID  int    `json:"messageId,omitempty"`
	SentAt     string `json:"sentAt,omitempty"`
	ReadAt     string `json:"readAt,omitempty"`
}

// ─── Hub ────────────────────────────────────────────────────

type Hub struct {
	mu      sync.RWMutex
	clients map[int]*websocket.Conn // userID → ws connection
	db      *sqlx.DB
	rdb     *redis.Client
	ctx     context.Context
}

func NewHub(db *sqlx.DB, redisURL string) (*Hub, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	h := &Hub{
		clients: make(map[int]*websocket.Conn),
		db:      db,
		rdb:     rdb,
		ctx:     context.Background(),
	}

	// Subscribe to Redis pub/sub channel for cross-instance delivery
	go h.subscribeRedis()

	return h, nil
}

// Register a user's websocket connection
func (h *Hub) Register(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	// Close previous connection if exists
	if old, ok := h.clients[userID]; ok {
		old.Close()
	}
	h.clients[userID] = conn
	h.mu.Unlock()

	log.Printf("[chat] user %d connected", userID)

	// Broadcast online status to friends
	h.broadcastPresence(userID, "online")
}

// Unregister a user
func (h *Hub) Unregister(userID int) {
	h.mu.Lock()
	delete(h.clients, userID)
	h.mu.Unlock()

	log.Printf("[chat] user %d disconnected", userID)
	h.broadcastPresence(userID, "offline")
}

// HandleMessage processes an incoming websocket message
func (h *Hub) HandleMessage(senderID int, raw []byte) {
	var msg IncomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Printf("[chat] bad message from user %d: %v", senderID, err)
		return
	}

	switch msg.Type {
	case "chat":
		h.handleChat(senderID, msg)
	case "typing":
		h.handleTyping(senderID, msg)
	case "read":
		h.handleRead(senderID, msg)
	default:
		log.Printf("[chat] unknown message type %q from user %d", msg.Type, senderID)
	}
}

// ─── Chat message ──────────────────────────────────────────

func (h *Hub) handleChat(senderID int, msg IncomingMessage) {
	if msg.Body == "" || msg.ReceiverID == 0 {
		return
	}

	// Persist to MySQL
	now := time.Now().UTC()
	res, err := h.db.ExecContext(h.ctx,
		`INSERT INTO ChatMessage (SenderID, ReceiverID, Body, SentAt) VALUES (?, ?, ?, ?)`,
		senderID, msg.ReceiverID, msg.Body, now,
	)
	if err != nil {
		log.Printf("[chat] db insert error: %v", err)
		return
	}
	messageID64, _ := res.LastInsertId()
	messageID := int(messageID64)

	out := OutgoingMessage{
		Type:       "chat",
		SenderID:   senderID,
		ReceiverID: msg.ReceiverID,
		Body:       msg.Body,
		MessageID:  messageID,
		SentAt:     now.Format(time.RFC3339),
	}

	// Publish to Redis so other server instances can deliver
	data, _ := json.Marshal(out)
	h.rdb.Publish(h.ctx, "chat:messages", data)

	// Also deliver locally
	h.deliverToUser(msg.ReceiverID, out)
	// Echo back to sender with the messageID
	h.deliverToUser(senderID, out)
}

// ─── Typing indicator ──────────────────────────────────────

func (h *Hub) handleTyping(senderID int, msg IncomingMessage) {
	out := OutgoingMessage{
		Type:       "typing",
		SenderID:   senderID,
		ReceiverID: msg.ReceiverID,
	}
	h.deliverToUser(msg.ReceiverID, out)
}

// ─── Read receipt ──────────────────────────────────────────

func (h *Hub) handleRead(senderID int, msg IncomingMessage) {
	if msg.MessageID == 0 {
		return
	}

	now := time.Now().UTC()
	_, _ = h.db.ExecContext(h.ctx,
		`UPDATE ChatMessage SET ReadAt = ? WHERE MessageId = ? AND ReceiverID = ?`,
		now, msg.MessageID, senderID,
	)

	out := OutgoingMessage{
		Type:      "read",
		SenderID:  senderID,
		MessageID: msg.MessageID,
		ReadAt:    now.Format(time.RFC3339),
	}

	// Notify the original sender that their message was read
	// We need to look up who sent the message
	var originalSender int
	err := h.db.GetContext(h.ctx, &originalSender,
		`SELECT SenderID FROM ChatMessage WHERE MessageId = ?`, msg.MessageID)
	if err == nil {
		h.deliverToUser(originalSender, out)
	}
}

// ─── Presence ──────────────────────────────────────────────

func (h *Hub) broadcastPresence(userID int, status string) {
	// Get user's accepted friends
	friendIDs := h.getAcceptedFriendIDs(userID)

	out := OutgoingMessage{
		Type:     status, // "online" or "offline"
		SenderID: userID,
	}

	for _, fid := range friendIDs {
		h.deliverToUser(fid, out)
	}
}

func (h *Hub) getAcceptedFriendIDs(userID int) []int {
	var ids []int
	rows, err := h.db.QueryContext(h.ctx,
		`SELECT CASE WHEN SenderID = ? THEN ReceiverID ELSE SenderID END AS friendId
		 FROM Friends WHERE Status = 'accepted' AND (SenderID = ? OR ReceiverID = ?)`,
		userID, userID, userID,
	)
	if err != nil {
		return ids
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if rows.Scan(&id) == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// IsOnline checks if a user is connected
func (h *Hub) IsOnline(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// GetOnlineFriends returns which of the user's friends are currently online
func (h *Hub) GetOnlineFriends(userID int) []int {
	friendIDs := h.getAcceptedFriendIDs(userID)
	var online []int
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, fid := range friendIDs {
		if _, ok := h.clients[fid]; ok {
			online = append(online, fid)
		}
	}
	return online
}

// ─── Delivery ──────────────────────────────────────────────

func (h *Hub) deliverToUser(userID int, msg OutgoingMessage) {
	h.mu.RLock()
	conn, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return // user not connected on this instance
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[chat] write error to user %d: %v", userID, err)
	}
}

// ─── Redis pub/sub ─────────────────────────────────────────

func (h *Hub) subscribeRedis() {
	sub := h.rdb.Subscribe(h.ctx, "chat:messages")
	ch := sub.Channel()

	for msg := range ch {
		var out OutgoingMessage
		if err := json.Unmarshal([]byte(msg.Payload), &out); err != nil {
			continue
		}
		// Deliver to the receiver if connected on this instance
		// (avoid double-delivery: only deliver if we didn't already)
		h.mu.RLock()
		_, senderHere := h.clients[out.SenderID]
		_, receiverHere := h.clients[out.ReceiverID]
		h.mu.RUnlock()

		// Only deliver from Redis if we didn't already deliver locally
		if !senderHere && receiverHere {
			h.deliverToUser(out.ReceiverID, out)
		}
		if !receiverHere && senderHere {
			h.deliverToUser(out.SenderID, out)
		}
	}
}

// ─── Chat history queries ──────────────────────────────────

type ChatHistoryMessage struct {
	MessageID  int        `db:"MessageId" json:"messageId"`
	SenderID   int        `db:"SenderID" json:"senderId"`
	ReceiverID int        `db:"ReceiverID" json:"receiverId"`
	Body       string     `db:"Body" json:"body"`
	SentAt     time.Time  `db:"SentAt" json:"sentAt"`
	ReadAt     *time.Time `db:"ReadAt" json:"readAt"`
}

// GetConversation returns messages between two users, paginated
func (h *Hub) GetConversation(userA, userB, limit, offset int) ([]ChatHistoryMessage, error) {
	var msgs []ChatHistoryMessage
	err := h.db.SelectContext(h.ctx, &msgs,
		`SELECT MessageId, SenderID, ReceiverID, Body, SentAt, ReadAt
		 FROM ChatMessage
		 WHERE (SenderID = ? AND ReceiverID = ?) OR (SenderID = ? AND ReceiverID = ?)
		 ORDER BY SentAt DESC
		 LIMIT ? OFFSET ?`,
		userA, userB, userB, userA, limit, offset,
	)
	return msgs, err
}

// GetUnreadCount returns how many unread messages a user has from a specific sender
func (h *Hub) GetUnreadCount(receiverID, senderID int) (int, error) {
	var count int
	err := h.db.GetContext(h.ctx, &count,
		`SELECT COUNT(*) FROM ChatMessage WHERE ReceiverID = ? AND SenderID = ? AND ReadAt IS NULL`,
		receiverID, senderID,
	)
	return count, err
}
