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
	Type       string `json:"type"`       // "chat" | "typing" | "read" | "game_*"
	ReceiverID int    `json:"receiverId"` // who we're talking to
	Body       string `json:"body,omitempty"`
	MessageID  int    `json:"messageId,omitempty"` // for read receipts

	// ── Game fields ──
	GameType   string          `json:"gameType,omitempty"`
	InviteID   string          `json:"inviteId,omitempty"`
	SessionID  string          `json:"sessionId,omitempty"`
	SenderName string          `json:"senderName,omitempty"`
	Move       json.RawMessage `json:"move,omitempty"`
	WinnerID   *int            `json:"winnerId,omitempty"`
}

type OutgoingMessage struct {
	Type       string `json:"type"` // "chat" | "typing" | "read" | "online" | "offline" | "game_*"
	SenderID   int    `json:"senderId"`
	ReceiverID int    `json:"receiverId,omitempty"`
	Body       string `json:"body,omitempty"`
	MessageID  int    `json:"messageId,omitempty"`
	SentAt     string `json:"sentAt,omitempty"`
	ReadAt     string `json:"readAt,omitempty"`

	// ── Game fields (relayed as-is) ──
	GameType   string          `json:"gameType,omitempty"`
	InviteID   string          `json:"inviteId,omitempty"`
	SessionID  string          `json:"sessionId,omitempty"`
	SenderName string          `json:"senderName,omitempty"`
	Move       json.RawMessage `json:"move,omitempty"`
	WinnerID   *int            `json:"winnerId,omitempty"`
}

// AdminMessage is sent from admin panel → player via WebSocket
type AdminMessage struct {
	Type    string          `json:"type"`    // "admin_money","admin_event","admin_stats","admin_sickness","admin_heal","admin_revive"
	Payload json.RawMessage `json:"payload"` // action-specific data
	Message string          `json:"message"` // human-readable summary
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

	// Flush any pending admin notifications from while user was offline
	go h.FlushPendingOnConnect(userID)
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
	case "game_invite", "game_accept", "game_decline", "game_move", "game_state", "game_end":
		h.handleGameRelay(senderID, msg)
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

// ─── Game relay ────────────────────────────────────────────

// handleGameRelay relays game messages between two players.
// The server does NOT validate moves — it acts as a thin relay
// so both clients handle game logic locally.
func (h *Hub) handleGameRelay(senderID int, msg IncomingMessage) {
	if msg.ReceiverID == 0 {
		log.Printf("[game] %s from user %d missing receiverId", msg.Type, senderID)
		return
	}

	out := OutgoingMessage{
		Type:       msg.Type,
		SenderID:   senderID,
		ReceiverID: msg.ReceiverID,
		GameType:   msg.GameType,
		InviteID:   msg.InviteID,
		SessionID:  msg.SessionID,
		SenderName: msg.SenderName,
		Move:       msg.Move,
		WinnerID:   msg.WinnerID,
	}

	log.Printf("[game] relay %s from user %d → user %d (session=%s)",
		msg.Type, senderID, msg.ReceiverID, msg.SessionID)

	// Deliver to receiver
	h.deliverToUser(msg.ReceiverID, out)

	// For accept messages, also echo back to sender with the session info
	if msg.Type == "game_accept" {
		h.deliverToUser(senderID, out)
	}
}

// ─── Presence (continued) ──────────────────────────────────

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
		log.Printf("[chat] user %d not connected, skipping delivery", userID)
		return // user not connected on this instance
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("[chat] write error to user %d: %v", userID, err)
	} else {
		log.Printf("[chat] delivered %s to user %d (msgId=%d)", msg.Type, userID, msg.MessageID)
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

// ─── Conversation list ─────────────────────────────────────

// ConversationSummary is one row in the conversations list.
type ConversationSummary struct {
	FriendID    int       `db:"friendId"    json:"friendId"`
	UserName    string    `db:"userName"    json:"userName"`
	Name        string    `db:"name"        json:"name"`
	LastMessage string    `db:"lastMessage" json:"lastMessage"`
	LastTime    time.Time `db:"lastTime"    json:"lastTime"`
	Unread      int       `db:"unread"      json:"unread"`
	Online      bool      `json:"online"`
}

// GetConversations returns all conversations for a user (friends + latest message + unread count).
func (h *Hub) GetConversations(userID int) ([]ConversationSummary, error) {
	// Step 1: Get all accepted friend IDs with user info
	type friendRow struct {
		FriendID int    `db:"friendId"`
		UserName string `db:"userName"`
		Name     string `db:"name"`
	}
	var friends []friendRow
	err := h.db.SelectContext(h.ctx, &friends,
		`SELECT
		   CASE WHEN f.SenderID = ? THEN f.ReceiverID ELSE f.SenderID END AS friendId,
		   u.UserName AS userName,
		   u.Name     AS name
		 FROM Friends f
		 JOIN Users u ON u.UserID = CASE WHEN f.SenderID = ? THEN f.ReceiverID ELSE f.SenderID END
		 WHERE f.Status = 'accepted' AND (f.SenderID = ? OR f.ReceiverID = ?)`,
		userID, userID, userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get friends: %w", err)
	}

	convos := make([]ConversationSummary, 0, len(friends))
	for _, fr := range friends {
		// Last message between these two users
		var lastMsg struct {
			Body   string    `db:"Body"`
			SentAt time.Time `db:"SentAt"`
		}
		hasMsg := true
		err := h.db.GetContext(h.ctx, &lastMsg,
			`SELECT Body, SentAt FROM ChatMessage
			 WHERE (SenderID = ? AND ReceiverID = ?) OR (SenderID = ? AND ReceiverID = ?)
			 ORDER BY SentAt DESC LIMIT 1`,
			userID, fr.FriendID, fr.FriendID, userID,
		)
		if err != nil {
			hasMsg = false
		}

		// Unread count from this friend
		var unread int
		_ = h.db.GetContext(h.ctx, &unread,
			`SELECT COUNT(*) FROM ChatMessage WHERE ReceiverID = ? AND SenderID = ? AND ReadAt IS NULL`,
			userID, fr.FriendID,
		)

		cs := ConversationSummary{
			FriendID: fr.FriendID,
			UserName: fr.UserName,
			Name:     fr.Name,
			Unread:   unread,
			Online:   h.IsOnline(fr.FriendID),
		}
		if hasMsg {
			cs.LastMessage = lastMsg.Body
			cs.LastTime = lastMsg.SentAt
		}
		convos = append(convos, cs)
	}

	return convos, nil
}

// GetTotalUnread returns the total unread message count across all conversations.
func (h *Hub) GetTotalUnread(userID int) (int, error) {
	var count int
	err := h.db.GetContext(h.ctx, &count,
		`SELECT COUNT(*) FROM ChatMessage WHERE ReceiverID = ? AND ReadAt IS NULL`,
		userID,
	)
	return count, err
}

// ─── Admin push ────────────────────────────────────────────

// SendAdminPush delivers an admin action to a user via WebSocket.
// If the user is offline, the notification is persisted to DB so it
// can be fetched on reconnection.
func (h *Hub) SendAdminPush(targetUserID int, msgType string, payload json.RawMessage, message string) (delivered bool, err error) {
	// 1. Try to deliver via WebSocket first
	adminMsg := AdminMessage{
		Type:    msgType,
		Payload: payload,
		Message: message,
	}

	h.mu.RLock()
	conn, online := h.clients[targetUserID]
	h.mu.RUnlock()

	if online {
		data, _ := json.Marshal(adminMsg)
		if writeErr := conn.WriteMessage(websocket.TextMessage, data); writeErr != nil {
			log.Printf("[admin-push] ws write error to user %d: %v", targetUserID, writeErr)
			// Fall through to persist for later
		} else {
			log.Printf("[admin-push] delivered %s to user %d", msgType, targetUserID)
			return true, nil
		}
	}

	// 2. User is offline (or WS write failed) — persist for later delivery on reconnect
	_, err = h.db.ExecContext(h.ctx,
		`INSERT INTO AdminNotification (TargetUserId, Type, Payload, Message) VALUES (?, ?, ?, ?)`,
		targetUserID, msgType, string(payload), message,
	)
	if err != nil {
		log.Printf("[admin-push] failed to persist notification for user %d: %v", targetUserID, err)
		return false, fmt.Errorf("persist notification: %w", err)
	}

	log.Printf("[admin-push] user %d offline, notification persisted for later", targetUserID)
	return false, nil
}

// SendAdminBroadcast delivers an admin message to ALL connected users and
// persists notifications for offline ones.
func (h *Hub) SendAdminBroadcast(msgType string, payload json.RawMessage, message string) (onlineCount int, err error) {
	// Get all user IDs
	var userIDs []int
	err = h.db.SelectContext(h.ctx, &userIDs, `SELECT UserId FROM Users`)
	if err != nil {
		return 0, fmt.Errorf("list users: %w", err)
	}

	onlineCount = 0
	for _, uid := range userIDs {
		delivered, pushErr := h.SendAdminPush(uid, msgType, payload, message)
		if pushErr != nil {
			log.Printf("[admin-push] broadcast error for user %d: %v", uid, pushErr)
			continue
		}
		if delivered {
			onlineCount++
		}
	}
	return onlineCount, nil
}

// GetPendingNotifications returns unread admin notifications for a user.
func (h *Hub) GetPendingNotifications(userID int) ([]AdminNotificationRow, error) {
	var rows []AdminNotificationRow
	err := h.db.SelectContext(h.ctx, &rows,
		`SELECT NotificationId, TargetUserId, Type, Payload, Message, CreatedAt, ReadAt
		 FROM AdminNotification
		 WHERE TargetUserId = ? AND ReadAt IS NULL
		 ORDER BY CreatedAt DESC`,
		userID,
	)
	return rows, err
}

// MarkNotificationsRead marks all unread notifications as read for a user.
func (h *Hub) MarkNotificationsRead(userID int) error {
	_, err := h.db.ExecContext(h.ctx,
		`UPDATE AdminNotification SET ReadAt = NOW() WHERE TargetUserId = ? AND ReadAt IS NULL`,
		userID,
	)
	return err
}

// AdminNotificationRow is a DB-mapped row from AdminNotification.
type AdminNotificationRow struct {
	NotificationID int        `db:"NotificationId" json:"notificationId"`
	TargetUserID   int        `db:"TargetUserId"   json:"targetUserId"`
	Type           string     `db:"Type"           json:"type"`
	Payload        string     `db:"Payload"        json:"payload"`
	Message        string     `db:"Message"        json:"message"`
	CreatedAt      time.Time  `db:"CreatedAt"      json:"createdAt"`
	ReadAt         *time.Time `db:"ReadAt"         json:"readAt"`
}

// FlushPendingOnConnect sends all pending admin notifications to a
// user who just connected, then marks them as read.
func (h *Hub) FlushPendingOnConnect(userID int) {
	rows, err := h.GetPendingNotifications(userID)
	if err != nil || len(rows) == 0 {
		return
	}

	log.Printf("[admin-push] flushing %d pending notifications to user %d", len(rows), userID)

	for _, row := range rows {
		adminMsg := AdminMessage{
			Type:    row.Type,
			Payload: json.RawMessage(row.Payload),
			Message: row.Message,
		}
		data, _ := json.Marshal(adminMsg)

		h.mu.RLock()
		conn, ok := h.clients[userID]
		h.mu.RUnlock()

		if ok {
			_ = conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	// Mark all as read after delivery
	_ = h.MarkNotificationsRead(userID)
}
