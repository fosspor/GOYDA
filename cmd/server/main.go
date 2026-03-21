package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fosspor/GOYDA/internal/config"
	"github.com/fosspor/GOYDA/internal/handlers"
	"github.com/fosspor/GOYDA/internal/migrate"
	"github.com/fosspor/GOYDA/internal/middleware"
	"github.com/fosspor/GOYDA/internal/spa"
	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	migDir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR"))
	if migDir == "" {
		migDir = "./migrations"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	if err := migrate.Up(ctx, pool, migDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	st := &store.Store{Pool: pool}
	if err := st.SeedDemo(ctx); err != nil {
		log.Fatalf("seed: %v", err)
	}

	api := handlers.NewAPI(cfg, st)

	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"detail": err.Error()})
		},
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Get("/health", api.Health)

	app.Post("/api/auth/register", api.Register)
	app.Post("/api/auth/login", api.Login)

	j := middleware.JWT([]byte(cfg.JWTSecret))
	app.Get("/api/me", j, api.Me)
	app.Patch("/api/me", j, api.PatchMe)

	app.Get("/api/locations", api.ListLocations)
	app.Get("/api/locations/:id", api.GetLocation)
	app.Post("/api/locations", j, api.CreateLocation)

	app.Get("/api/routes", j, api.ListMyRoutes)
	app.Post("/api/routes", j, api.CreateRoute)
	app.Get("/api/routes/:id", j, api.GetRoute)

	app.Post("/api/ai/generate-route", api.GenerateRoute)
	app.Get("/api/ai/recommendations", api.AIRecommendations)

	spa.Register(app)

	log.Printf("listening on %s", cfg.HTTPAddr)
	go func() {
		<-ctx.Done()
		_ = app.Shutdown()
	}()
	if err := app.Listen(cfg.HTTPAddr); err != nil {
		// При штатном Shutdown() не считаем это фатальной ошибкой запуска.
		if ctx.Err() != nil {
			return
		}
		log.Fatalf("listen: %v", err)
	}
}
