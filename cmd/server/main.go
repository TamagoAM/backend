package main

import (
	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/handler"

	"tamagoam/internal/config"
	"tamagoam/internal/db"
	"tamagoam/internal/graphql"
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
	}
	log.Println("db migrated")

	schema, err := graphql.NewSchema(dbConn)
	if err != nil {
		log.Fatalf("graphql schema failed: %v", err)
	}

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

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

	app.Post("/graphql", adaptor.HTTPHandler(gqlHandler))
	app.Get("/playground", adaptor.HTTPHandler(playgroundHandler))

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
