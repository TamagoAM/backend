package main

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/graphql-go/handler"

	"tamagoam/internal/auth"
	"tamagoam/internal/config"
	"tamagoam/internal/db"
	gql "tamagoam/internal/graphql"
)

// contextKey is unexported to avoid collisions.
type contextKey string

const userClaimsKey contextKey = "userClaims"

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
		if err := db.Migrate(dbConn, "migrations/006_datetime_columns.sql"); err != nil {
			log.Fatalf("db migrate 006 failed: %v", err)
		}
	}
	log.Println("db migrated")

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
				ctx := context.WithValue(c.Context(), userClaimsKey, claims)
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

	app.Post("/graphql", jwtMiddleware, adaptor.HTTPHandler(gqlHandler))
	app.Get("/playground", adaptor.HTTPHandler(playgroundHandler))

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
