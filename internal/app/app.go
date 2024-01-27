package app

import (
	"github.com/EvgeniyBudaev/love-server/internal/config"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

type App struct {
	Logger logger.Logger
	config *config.Config
	db     *Database
	fiber  *fiber.App
}

func NewApp() *App {
	// Default logger
	defaultLogger, err := logger.NewLogger(logger.GetDefaultLevel())
	if err != nil {
		log.Fatal("error func NewApp, method NewLogger by path internal/app/app.go", err)
	}
	// Config
	cfg, err := config.Load(defaultLogger)
	if err != nil {
		log.Fatal("error func NewApp, method Load by path internal/app/app.go", err)
	}
	// Logger level
	loggerLevel, err := logger.NewLogger(cfg.LoggerLevel)
	if err != nil {
		log.Fatal("error func NewApp, method NewLogger by path internal/app/app.go", err)
	}
	// Database connection
	postgresConnection, err := newPostgresConnection(cfg)
	if err != nil {
		log.Fatal("error func NewApp, method newPostgresConnection by path internal/app/app.go", err)
	}
	database := NewDatabase(loggerLevel, postgresConnection)
	err = postgresConnection.Ping()
	if err != nil {
		log.Fatal("error func NewApp, method NewDatabase by path internal/app/app.go", err)
	}
	// Fiber
	f := fiber.New(fiber.Config{
		ReadBufferSize: 16384,
	})
	// CORS
	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, X-Requested-With, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	return &App{
		config: cfg,
		db:     database,
		Logger: loggerLevel,
		fiber:  f,
	}
}
