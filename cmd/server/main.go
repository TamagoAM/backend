package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/graphql-go/handler"

	"time"

	"tamagoam/internal/auth"
	"tamagoam/internal/chat"
	"tamagoam/internal/config"
	"tamagoam/internal/db"
	gql "tamagoam/internal/graphql"
	"tamagoam/internal/ticker"
)

func main() {
	log.Println("starting server...")
	cfg := config.Load()
	log.Printf("config loaded: %+v", cfg)

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer dbConn.Close()
	log.Println("db connected")
	if cfg.MigrateOnStart {
		if err := db.Migrate(dbConn, "migrations/001_init.sql"); err != nil {
			log.Fatalf("db migrate failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/002_add_auth_fields.sql"); err != nil {
			log.Fatalf("db migrate 002 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/003_seed_admin.sql"); err != nil {
			log.Fatalf("db migrate 003 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/004_seed_races.sql"); err != nil {
			log.Fatalf("db migrate 004 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/005_seed_game_data.sql"); err != nil {
			log.Fatalf("db migrate 005 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/006_friend_requests_chat.sql"); err != nil {
			log.Fatalf("db migrate 006 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/007_admin_notifications.sql"); err != nil {
			log.Fatalf("db migrate 007 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/008_life_choice_options.sql"); err != nil {
			log.Fatalf("db migrate 008 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/009_add_happiness_to_tama_stats.sql"); err != nil {
			log.Fatalf("db migrate 009 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/010_stat_history.sql"); err != nil {
			log.Fatalf("db migrate 010 failed: %v", err)
		}
		if err := db.Migrate(dbConn, "migrations/011_add_last_tick_at.sql"); err != nil {
			log.Fatalf("db migrate 011 failed: %v", err)
		}
	}
	log.Println("db migrated")

	// ─── Background Ticker (decay engine) ────────────────────
	gameTicker := ticker.New(dbConn, 5*time.Minute)
	gameTicker.Start()
	defer gameTicker.Stop()
	log.Println("game ticker started (5m interval)")

	// ─── Chat Hub (WebSocket + Redis) ─────────────────────────
	chatHub, err := chat.NewHub(dbConn, cfg.RedisURL)
	if err != nil {
		log.Fatalf("chat hub init failed: %v", err)
	}
	log.Println("chat hub initialised")

	schema, err := gql.NewSchema(dbConn)
	if err != nil {
		log.Fatalf("graphql schema failed: %v", err)
	}

	store := gql.NewSQLStore(dbConn)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// ─── Auth endpoints ────────────────────────────────────────
	app.Post("/auth/register", func(c *fiber.Ctx) error {
		var body struct {
			Name     string `json:"name"`
			LastName string `json:"lastName"`
			UserName string `json:"userName"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
		}
		if body.Email == "" || body.Password == "" || body.UserName == "" {
			return c.Status(400).JSON(fiber.Map{"error": "email, userName and password are required"})
		}

		// Check for existing email
		if existing, _ := store.GetUserByEmail(c.Context(), body.Email); existing != nil {
			return c.Status(409).JSON(fiber.Map{"error": "email already registered"})
		}
		// Check for existing username
		if existing, _ := store.GetUserByUserName(c.Context(), body.UserName); existing != nil {
			return c.Status(409).JSON(fiber.Map{"error": "username already taken"})
		}

		hashed, err := auth.HashPassword(body.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
		}

		user, err := store.CreateUser(c.Context(), gql.CreateUserInput{
			Name:           body.Name,
			LastName:       body.LastName,
			UserName:       body.UserName,
			Email:          body.Email,
			PasswordHash:   hashed,
			ClearanceLevel: 0,
			Verified:       false,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to create user: " + err.Error()})
		}

		token, err := auth.GenerateToken(cfg.JWTSecret, user.UserID, user.UserName, user.ClearanceLevel)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":             user.UserID,
				"name":           user.Name,
				"lastName":       user.LastName,
				"userName":       user.UserName,
				"email":          user.Email,
				"clearanceLevel": user.ClearanceLevel,
				"verified":       user.Verified,
			},
		})
	})

	app.Post("/auth/login", func(c *fiber.Ctx) error {
		var body struct {
			Login    string `json:"login"` // email or username
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
		}
		if body.Login == "" || body.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "login and password are required"})
		}

		// Try by email first, then by username
		user, err := store.GetUserByEmail(c.Context(), body.Login)
		if err != nil {
			user, err = store.GetUserByUserName(c.Context(), body.Login)
		}
		if err != nil || user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		if !auth.CheckPassword(user.PasswordHash, body.Password) {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// Update last connection timestamp
		if err := store.UpdateLastConnection(c.Context(), user.UserID); err != nil {
			log.Printf("warning: failed to update last connection: %v", err)
		}

		token, err := auth.GenerateToken(cfg.JWTSecret, user.UserID, user.UserName, user.ClearanceLevel)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":             user.UserID,
				"name":           user.Name,
				"lastName":       user.LastName,
				"userName":       user.UserName,
				"email":          user.Email,
				"clearanceLevel": user.ClearanceLevel,
				"verified":       user.Verified,
			},
		})
	})

	app.Get("/auth/me", func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(fiber.Map{"error": "no token provided"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ValidateToken(cfg.JWTSecret, tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		user, err := store.GetUser(c.Context(), claims.UserID)
		if err != nil || user == nil {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(fiber.Map{
			"user": fiber.Map{
				"id":             user.UserID,
				"name":           user.Name,
				"lastName":       user.LastName,
				"userName":       user.UserName,
				"email":          user.Email,
				"clearanceLevel": user.ClearanceLevel,
				"verified":       user.Verified,
			},
		})
	})

	// ─── JWT middleware for GraphQL ─────────────────────────────
	// Optional auth: if a valid token is present, inject claims into context.
	// Requests without a token still pass through (for public queries).
	jwtMiddleware := func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header != "" {
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if claims, err := auth.ValidateToken(cfg.JWTSecret, tokenStr); err == nil {
				ctx := context.WithValue(c.Context(), auth.UserClaimsKey, claims)
				c.SetUserContext(ctx)
			}
		}
		return c.Next()
	}

	gqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: false,
	})
	playgroundHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Wrap the GraphQL handler so that the JWT claims stored in
	// c.UserContext() are forwarded into the http.Request.Context()
	// that the graphql-go resolver receives as p.Context.
	// adaptor.HTTPHandler does NOT do this automatically.
	gqlFiberHandler := func(c *fiber.Ctx) error {
		var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Merge the Fiber user-context (with JWT claims) into the http.Request
			if userCtx := c.UserContext(); userCtx != nil {
				if claims, ok := userCtx.Value(auth.UserClaimsKey).(*auth.Claims); ok && claims != nil {
					r = r.WithContext(context.WithValue(r.Context(), auth.UserClaimsKey, claims))
				}
			}
			gqlHandler.ServeHTTP(w, r)
		})
		return adaptor.HTTPHandler(httpHandler)(c)
	}

	app.Post("/graphql", jwtMiddleware, gqlFiberHandler)
	app.Get("/playground", adaptor.HTTPHandler(playgroundHandler))

	// ─── WebSocket for real-time chat ─────────────────────────
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Authenticate via query param ?token=<jwt>
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "no token"))
			c.Close()
			return
		}
		claims, err := auth.ValidateToken(cfg.JWTSecret, tokenStr)
		if err != nil {
			c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "invalid token"))
			c.Close()
			return
		}

		userID := claims.UserID
		chatHub.Register(userID, c)
		defer chatHub.Unregister(userID)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			chatHub.HandleMessage(userID, msg)
		}
	}))

	// ─── Chat REST endpoints ──────────────────────────────────
	app.Get("/chat/history", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		withID, err := strconv.Atoi(c.Query("with"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "'with' query param (user id) is required"})
		}
		limit, _ := strconv.Atoi(c.Query("limit", "50"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		msgs, err := chatHub.GetConversation(claims.UserID, withID, limit, offset)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(msgs)
	})

	app.Get("/chat/unread", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		senderID, _ := strconv.Atoi(c.Query("from", "0"))
		count, err := chatHub.GetUnreadCount(claims.UserID, senderID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"unread": count})
	})

	app.Get("/chat/conversations", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		convos, err := chatHub.GetConversations(claims.UserID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(convos)
	})

	app.Get("/chat/total-unread", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		count, err := chatHub.GetTotalUnread(claims.UserID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"unread": count})
	})

	// ─── Admin action endpoints (real-time push) ─────────────

	// Require admin clearance (level >= 2)
	adminMiddleware := func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		if claims.ClearanceLevel < 2 {
			return c.Status(403).JSON(fiber.Map{"error": "admin access required"})
		}
		return c.Next()
	}

	// POST /admin/give-money — atomically add money + push via WS
	app.Post("/admin/give-money", jwtMiddleware, adminMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			TamaStatsID  int `json:"tamaStatsId"`
			TargetUserID int `json:"targetUserId"`
			Amount       int `json:"amount"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.Amount <= 0 || body.TamaStatsID == 0 || body.TargetUserID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "amount, tamaStatsId, and targetUserId are required"})
		}

		// Atomic money update — no race condition
		_, err := dbConn.ExecContext(c.Context(),
			`UPDATE Tama_stats SET Money = Money + ? WHERE TamaStatId = ?`,
			body.Amount, body.TamaStatsID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db update failed: " + err.Error()})
		}

		// Read back new balance
		var newBalance int
		_ = dbConn.GetContext(c.Context(), &newBalance,
			`SELECT Money FROM Tama_stats WHERE TamaStatId = ?`, body.TamaStatsID)

		payload, _ := json.Marshal(fiber.Map{"amount": body.Amount, "newBalance": newBalance})
		msg := fmt.Sprintf("💰 An admin gave you %d coins!", body.Amount)

		delivered, _ := chatHub.SendAdminPush(body.TargetUserID, "admin_money", payload, msg)

		return c.JSON(fiber.Map{
			"success":    true,
			"newBalance": newBalance,
			"delivered":  delivered,
		})
	})

	// POST /admin/send-event — create active event + push via WS
	app.Post("/admin/send-event", jwtMiddleware, adminMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			EventID      int  `json:"eventId"`
			TargetUserID int  `json:"targetUserId"`
			IsGlobal     bool `json:"isGlobal"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.EventID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "eventId is required"})
		}

		// Fetch event details for the payload
		var ev struct {
			Name     string  `db:"Name"`
			Desc     *string `db:"Desc"`
			Severity string  `db:"Severity"`
			Bonus    *string `db:"Bonus"`
			Malus    *string `db:"Malus"`
		}
		if err := dbConn.GetContext(c.Context(), &ev,
			`SELECT Name, `+"`Desc`"+`, Severity, Bonus, Malus FROM Event WHERE EventId = ?`, body.EventID); err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "event not found"})
		}

		// Insert into ActiveEvent
		var targetUID *int
		if !body.IsGlobal {
			targetUID = &body.TargetUserID
		}
		_, err := dbConn.ExecContext(c.Context(),
			`INSERT INTO ActiveEvent (EventId, TargetUserId, IsGlobal) VALUES (?, ?, ?)`,
			body.EventID, targetUID, body.IsGlobal)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db insert failed: " + err.Error()})
		}

		descStr := ""
		if ev.Desc != nil {
			descStr = *ev.Desc
		}
		bonusStr := ""
		if ev.Bonus != nil {
			bonusStr = *ev.Bonus
		}
		malusStr := ""
		if ev.Malus != nil {
			malusStr = *ev.Malus
		}

		payload, _ := json.Marshal(fiber.Map{
			"eventId":  body.EventID,
			"name":     ev.Name,
			"desc":     descStr,
			"severity": ev.Severity,
			"bonus":    bonusStr,
			"malus":    malusStr,
		})
		msg := fmt.Sprintf("🎭 Event: %s — %s", ev.Name, descStr)

		var onlineCount int
		var delivered bool
		if body.IsGlobal {
			onlineCount, _ = chatHub.SendAdminBroadcast("admin_event", payload, msg)
		} else {
			delivered, _ = chatHub.SendAdminPush(body.TargetUserID, "admin_event", payload, msg)
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"isGlobal":    body.IsGlobal,
			"delivered":   delivered,
			"onlineCount": onlineCount,
		})
	})

	// POST /admin/adjust-stats — apply stat deltas + push via WS
	app.Post("/admin/adjust-stats", jwtMiddleware, adminMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			TamaStatsID  int            `json:"tamaStatsId"`
			TargetUserID int            `json:"targetUserId"`
			Deltas       map[string]int `json:"deltas"` // {"hunger": 20, "boredom": -10, ...}
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.TamaStatsID == 0 || body.TargetUserID == 0 || len(body.Deltas) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "tamaStatsId, targetUserId, and deltas are required"})
		}

		// Build SET clauses for each delta — use GREATEST/LEAST to clamp 0-100
		allowed := map[string]string{
			"hunger": "Hunger", "boredom": "Boredom", "hygiene": "Hygiene",
			"socialSatis": "SocialSatis", "workSatis": "WorkSatis", "personalSatis": "PersonalSatis",
		}
		sets := []string{}
		args := []interface{}{}
		for k, delta := range body.Deltas {
			col, ok := allowed[k]
			if !ok || delta == 0 {
				continue
			}
			sets = append(sets, fmt.Sprintf("%s = GREATEST(0, LEAST(100, %s + ?))", col, col))
			args = append(args, delta)
		}
		if len(sets) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "no valid deltas"})
		}

		query := "UPDATE Tama_stats SET " + strings.Join(sets, ", ") + " WHERE TamaStatId = ?"
		args = append(args, body.TamaStatsID)
		_, err := dbConn.ExecContext(c.Context(), query, args...)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db update failed: " + err.Error()})
		}

		payload, _ := json.Marshal(fiber.Map{"deltas": body.Deltas})
		msg := "📊 An admin adjusted your stats!"
		delivered, _ := chatHub.SendAdminPush(body.TargetUserID, "admin_stats", payload, msg)

		return c.JSON(fiber.Map{"success": true, "delivered": delivered})
	})

	// POST /admin/give-sickness — update tama sickness + push via WS
	app.Post("/admin/give-sickness", jwtMiddleware, adminMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			TamaID       int    `json:"tamaId"`
			TargetUserID int    `json:"targetUserId"`
			Sickness     string `json:"sickness"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.TamaID == 0 || body.TargetUserID == 0 || body.Sickness == "" {
			return c.Status(400).JSON(fiber.Map{"error": "tamaId, targetUserId, and sickness are required"})
		}

		_, err := dbConn.ExecContext(c.Context(),
			`UPDATE Tama SET Sickness = ? WHERE TamaId = ?`, body.Sickness, body.TamaID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db update failed: " + err.Error()})
		}

		payload, _ := json.Marshal(fiber.Map{"sickness": body.Sickness})
		msg := fmt.Sprintf("🦠 Your tama caught %s!", body.Sickness)
		delivered, _ := chatHub.SendAdminPush(body.TargetUserID, "admin_sickness", payload, msg)

		return c.JSON(fiber.Map{"success": true, "delivered": delivered})
	})

	// POST /admin/heal — remove sickness + push via WS
	app.Post("/admin/heal", jwtMiddleware, adminMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			TamaID       int    `json:"tamaId"`
			TargetUserID int    `json:"targetUserId"`
			OldSickness  string `json:"oldSickness"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.TamaID == 0 || body.TargetUserID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "tamaId and targetUserId are required"})
		}

		_, err := dbConn.ExecContext(c.Context(),
			`UPDATE Tama SET Sickness = NULL WHERE TamaId = ?`, body.TamaID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db update failed: " + err.Error()})
		}

		payload, _ := json.Marshal(fiber.Map{"healed": true, "oldSickness": body.OldSickness})
		msg := fmt.Sprintf("💊 An admin cured your %s!", body.OldSickness)
		delivered, _ := chatHub.SendAdminPush(body.TargetUserID, "admin_heal", payload, msg)

		return c.JSON(fiber.Map{"success": true, "delivered": delivered})
	})

	// POST /admin/revive — revive dead tama + reset stats + push via WS
	app.Post("/admin/revive", jwtMiddleware, adminMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			TamaID       int `json:"tamaId"`
			TamaStatsID  int `json:"tamaStatsId"`
			TargetUserID int `json:"targetUserId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.TamaID == 0 || body.TamaStatsID == 0 || body.TargetUserID == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "tamaId, tamaStatsId, and targetUserId are required"})
		}

		// Clear death info
		_, err := dbConn.ExecContext(c.Context(),
			`UPDATE Tama SET DeathDay = NULL, CauseOfDeath = NULL, Sickness = NULL WHERE TamaId = ?`, body.TamaID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "tama update failed: " + err.Error()})
		}

		// Reset stats to healthy values
		_, err = dbConn.ExecContext(c.Context(),
			`UPDATE Tama_stats SET Hunger = 70, Boredom = 30, Hygiene = 80, SocialSatis = 50, WorkSatis = 50, PersonalSatis = 50 WHERE TamaStatId = ?`,
			body.TamaStatsID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "stats update failed: " + err.Error()})
		}

		payload, _ := json.Marshal(fiber.Map{"revived": true})
		msg := "✨ An admin revived your tama!"
		delivered, _ := chatHub.SendAdminPush(body.TargetUserID, "admin_revive", payload, msg)

		return c.JSON(fiber.Map{"success": true, "delivered": delivered})
	})

	// GET /admin/notifications — get pending notifications for current user
	app.Get("/admin/notifications", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		rows, err := chatHub.GetPendingNotifications(claims.UserID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(rows)
	})

	// POST /admin/notifications/read — mark all notifications as read
	app.Post("/admin/notifications/read", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		if err := chatHub.MarkNotificationsRead(claims.UserID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	// ─── Widget data endpoint (lightweight REST) ──────────────
	// Returns pre-computed tama status for home-screen widgets.
	app.Get("/api/widget-data", jwtMiddleware, func(c *fiber.Ctx) error {
		claims, ok := c.UserContext().Value(auth.UserClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}

		tamas, err := store.TamasByUser(c.Context(), claims.UserID)
		if err != nil || len(tamas) == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "no tama found"})
		}
		tama := tamas[0]

		stats, err := store.TamaStatsByUser(c.Context(), claims.UserID)
		if err != nil || len(stats) == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "no stats found"})
		}
		stat := stats[0]

		// Happiness = average of the three satisfactions (mirrors frontend formula)
		happiness := (stat.SocialSatis + stat.WorkSatis + stat.PersonalSatis) / 3.0

		isAlive := tama.DeathDay == nil
		isSick := tama.Sickness != nil && *tama.Sickness != ""

		return c.JSON(fiber.Map{
			"tamaName":      tama.Name,
			"happiness":     happiness,
			"hunger":        stat.Hunger,
			"boredom":       stat.Boredom,
			"hygiene":       stat.Hygiene,
			"socialSatis":   stat.SocialSatis,
			"workSatis":     stat.WorkSatis,
			"personalSatis": stat.PersonalSatis,
			"money":         stat.Money,
			"isAlive":       isAlive,
			"isSick":        isSick,
		})
	})

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
