package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"tamagoam/internal/engine"
)

const expoPushURL = "https://exp.host/--/api/v2/push/send"

// Service handles push notification delivery and throttling.
type Service struct {
	db     *sqlx.DB
	client *http.Client
}

// NewService creates a new notification service.
func NewService(db *sqlx.DB) *Service {
	return &Service{
		db:     db,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// RegisterToken saves a push token for a user (upsert).
func (s *Service) RegisterToken(ctx context.Context, userID int, token string, platform string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO PushToken (UserId, Token, Platform)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE Platform = VALUES(Platform), CreatedAt = CURRENT_TIMESTAMP
	`, userID, token, platform)
	return err
}

// UnregisterToken removes a push token for a user.
func (s *Service) UnregisterToken(ctx context.Context, userID int, token string) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM PushToken WHERE UserId = ? AND Token = ?",
		userID, token)
	return err
}

// GetUserTokens returns all push tokens for a user.
func (s *Service) GetUserTokens(ctx context.Context, userID int) ([]PushToken, error) {
	var tokens []PushToken
	err := s.db.SelectContext(ctx, &tokens,
		"SELECT TokenId, UserId, Token, Platform FROM PushToken WHERE UserId = ?",
		userID)
	return tokens, err
}

// SendToUser sends a push notification to all devices of a user.
// Respects throttling and night-time suppression.
func (s *Service) SendToUser(ctx context.Context, userID int, notifType NotifType, title string, body string, data any) error {
	// Check if user is in night mode (suppress notifications during sleep)
	if s.isUserSleeping(ctx, userID) {
		log.Printf("[notif] suppressed %s for user %d (sleeping)", notifType, userID)
		return nil
	}

	// Check throttle
	if !s.shouldSend(ctx, userID, notifType) {
		log.Printf("[notif] throttled %s for user %d", notifType, userID)
		return nil
	}

	tokens, err := s.GetUserTokens(ctx, userID)
	if err != nil || len(tokens) == 0 {
		return err
	}

	// Build messages
	messages := make([]PushMessage, len(tokens))
	for i, t := range tokens {
		messages[i] = PushMessage{
			To:    t.Token,
			Title: title,
			Body:  body,
			Data:  data,
			Sound: "default",
		}
	}

	// Send via Expo Push API
	if err := s.sendExpoPush(messages); err != nil {
		log.Printf("[notif] expo push error for user %d: %v", userID, err)
		return err
	}

	// Log for throttling
	s.logNotification(ctx, userID, notifType)

	return nil
}

// CheckAndNotifyStats checks stat thresholds and sends notifications if needed.
func (s *Service) CheckAndNotifyStats(ctx context.Context, userID int, tamaName string, stats *engine.TamaStats, happiness float64, isSick bool, wasSick bool, isDead bool) {
	if isDead {
		_ = s.SendToUser(ctx, userID, NotifDeath,
			"💀 Your tama has died!",
			fmt.Sprintf("%s has passed away from neglect…", tamaName),
			map[string]string{"type": "death"})
		return
	}

	if stats.Hunger <= engine.NotifLowHungerThreshold {
		_ = s.SendToUser(ctx, userID, NotifLowHunger,
			"🍽️ Your tama is hungry!",
			fmt.Sprintf("%s's hunger is at %.0f%% — feed them soon!", tamaName, stats.Hunger),
			map[string]string{"type": "low_hunger"})
	}

	if happiness <= engine.NotifLowHappinessThreshold {
		_ = s.SendToUser(ctx, userID, NotifLowHappiness,
			"😢 Your tama is unhappy!",
			fmt.Sprintf("%s's happiness is very low (%.0f%%). Give them some attention!", tamaName, happiness),
			map[string]string{"type": "low_happiness"})
	}

	if stats.Boredom >= engine.NotifHighBoredomThreshold {
		_ = s.SendToUser(ctx, userID, NotifHighBoredom,
			"😴 Your tama is bored!",
			fmt.Sprintf("%s is bored out of their mind (%.0f%%). Play with them!", tamaName, stats.Boredom),
			map[string]string{"type": "high_boredom"})
	}

	if stats.Hygiene <= engine.NotifLowHygieneThreshold {
		_ = s.SendToUser(ctx, userID, NotifLowHygiene,
			"🧼 Your tama needs a bath!",
			fmt.Sprintf("%s's hygiene is at %.0f%% — time to clean up!", tamaName, stats.Hygiene),
			map[string]string{"type": "low_hygiene"})
	}

	if isSick && !wasSick {
		_ = s.SendToUser(ctx, userID, NotifSickness,
			"🤒 Your tama is sick!",
			fmt.Sprintf("%s has fallen ill. Visit the app to heal them!", tamaName),
			map[string]string{"type": "sickness"})
	}
}

// SendBedtimeReminder sends a reminder to put the tama to bed.
func (s *Service) SendBedtimeReminder(ctx context.Context, userID int, tamaName string) {
	_ = s.SendToUser(ctx, userID, NotifBedtime,
		"🌙 Bedtime!",
		fmt.Sprintf("Time to put %s to sleep! Turn off the lights to avoid happiness decay.", tamaName),
		map[string]string{"type": "bedtime"})
}

// SendWakeUpReminder sends a wake-up notification.
func (s *Service) SendWakeUpReminder(ctx context.Context, userID int, tamaName string) {
	_ = s.SendToUser(ctx, userID, NotifWakeUp,
		"☀️ Good morning!",
		fmt.Sprintf("%s is waking up! Turn on the lights and start the day.", tamaName),
		map[string]string{"type": "wake_up"})
}

// SendFriendRequestNotif notifies a user of a new friend request.
func (s *Service) SendFriendRequestNotif(ctx context.Context, userID int, fromName string) {
	_ = s.SendToUser(ctx, userID, NotifFriendRequest,
		"👋 Friend request!",
		fmt.Sprintf("%s wants to be your friend!", fromName),
		map[string]string{"type": "friend_request"})
}

// SendChatMessageNotif notifies a user of a new chat message.
func (s *Service) SendChatMessageNotif(ctx context.Context, userID int, fromName string, preview string) {
	_ = s.SendToUser(ctx, userID, NotifChatMessage,
		fmt.Sprintf("💬 Message from %s", fromName),
		preview,
		map[string]string{"type": "chat_message"})
}

// ── Internal helpers ──────────────────────────────

// isUserSleeping checks if the user's tama lights are off during nighttime.
func (s *Service) isUserSleeping(ctx context.Context, userID int) bool {
	var row struct {
		LightsOff bool   `db:"LightsOff"`
		Timezone  string `db:"Timezone"`
	}
	err := s.db.GetContext(ctx, &row, `
		SELECT COALESCE(ts.LightsOff, FALSE) AS LightsOff, COALESCE(u.Timezone, 'Europe/Paris') AS Timezone
		FROM Users u
		LEFT JOIN Tama t ON t.UserId = u.UserId AND t.DeathDay IS NULL
		LEFT JOIN Tama_stats ts ON t.TamaStatsID = ts.TamaStatId
		WHERE u.UserId = ?
		LIMIT 1
	`, userID)
	if err != nil {
		return false
	}

	loc, err := time.LoadLocation(row.Timezone)
	if err != nil {
		loc = time.UTC
	}
	hour := time.Now().In(loc).Hour()
	return row.LightsOff && engine.IsNightHour(hour)
}

// sendExpoPush sends messages via Expo's Push API.
func (s *Service) sendExpoPush(messages []PushMessage) error {
	body, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("marshal push: %w", err)
	}

	req, err := http.NewRequest("POST", expoPushURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create push request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send push: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("expo push returned status %d", resp.StatusCode)
	}

	return nil
}

// SendExpoPushPublic is a public wrapper around sendExpoPush for admin use.
func (s *Service) SendExpoPushPublic(messages []PushMessage) error {
	return s.sendExpoPush(messages)
}

// logNotification records that a notification was sent (for throttling).
func (s *Service) logNotification(ctx context.Context, userID int, notifType NotifType) {
	// Find current escalation level
	escalation := s.getEscalationLevel(ctx, userID, notifType) + 1
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO NotificationLog (UserId, NotifType, EscalationN)
		VALUES (?, ?, ?)
	`, userID, string(notifType), escalation)
	if err != nil {
		log.Printf("[notif] failed to log notification: %v", err)
	}
}

// getEscalationLevel returns how many times this notif type was sent recently.
func (s *Service) getEscalationLevel(ctx context.Context, userID int, notifType NotifType) int {
	var count int
	err := s.db.GetContext(ctx, &count, `
		SELECT COALESCE(MAX(EscalationN), 0)
		FROM NotificationLog
		WHERE UserId = ? AND NotifType = ? AND SentAt > DATE_SUB(NOW(), INTERVAL 24 HOUR)
	`, userID, string(notifType))
	if err != nil {
		return 0
	}
	return count
}

// shouldSend implements smart escalation: 1h → 4h → 12h between repeated notifs.
func (s *Service) shouldSend(ctx context.Context, userID int, notifType NotifType) bool {
	var lastSent struct {
		SentAt      time.Time `db:"SentAt"`
		EscalationN int       `db:"EscalationN"`
	}
	err := s.db.GetContext(ctx, &lastSent, `
		SELECT SentAt, EscalationN
		FROM NotificationLog
		WHERE UserId = ? AND NotifType = ?
		ORDER BY SentAt DESC
		LIMIT 1
	`, userID, string(notifType))
	if err != nil {
		// No previous notification — send it
		return true
	}

	// Smart escalation intervals
	var minInterval time.Duration
	switch {
	case lastSent.EscalationN <= 1:
		minInterval = 1 * time.Hour
	case lastSent.EscalationN <= 2:
		minInterval = 4 * time.Hour
	default:
		minInterval = 12 * time.Hour
	}

	elapsed := time.Since(lastSent.SentAt)
	return elapsed >= minInterval
}

// ResetEscalation resets the escalation counter for a notification type.
// Call this when the user takes corrective action (e.g., feeds tama → reset low_hunger).
func (s *Service) ResetEscalation(ctx context.Context, userID int, notifType NotifType) {
	_, _ = s.db.ExecContext(ctx, `
		DELETE FROM NotificationLog
		WHERE UserId = ? AND NotifType = ?
	`, userID, string(notifType))
}
