package main

import (
	"context"
	"log"
	"net"
	"net/http"
	api "orkidslearning/src/api"
	"orkidslearning/src/config"
	"orkidslearning/src/database"
	"orkidslearning/src/services"
	"orkidslearning/src/telemetry"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var serviceName = "orkidslearning"

func main() {
	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Load environment variables
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Set up OpenTelemetry
	otelShutdown, err := telemetry.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatalf("Failed to set up telemetry: %v", err)
	}
	defer func() {
		if shutdownErr := otelShutdown(context.Background()); shutdownErr != nil {
			log.Printf("Failed to shut down telemetry: %v", shutdownErr)
		}
	}()

	// Connect to MongoDB
	db, err := database.NewDatabase(ctx, env.MongoURI, env.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if disconnectErr := db.Disconnect(ctx); disconnectErr != nil {
			log.Printf("Failed to disconnect from the database: %v", disconnectErr)
		}
	}()

	// connect to postgres
	port, err := strconv.ParseUint(env.PostgresPort, 10, 16)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	connConfig := pgx.ConnConfig{
		Host:     env.PostgresHost,
		Port:     uint16(port),
		User:     env.PostgresUser,
		Password: env.PostgresPassword,
		Database: env.PostgresDB,
	}
	conn, err := database.NewPostgresDatabase(ctx, connConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if disconnectErr := conn.Disconnect(); disconnectErr != nil {
			log.Printf("Failed to disconnect from the database: %v", disconnectErr)
		}
	}()

	// Initialize services
	jwtService := services.NewJWTService(env.JWTSecretKey, env.JWTExpirationTime)
	contextService := services.NewContextService(db, jwtService, conn)

	// Create a Gin router
	router := gin.New()
	router.Use(otelgin.Middleware(serviceName))

	// CORS setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Register routes
	api.InitializeRoutes(router, contextService)

	// HTTP server configuration
	srv := &http.Server{
		Addr:         ":" + env.Port,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  5 * time.Second, // Adjusted timeout values
		WriteTimeout: 15 * time.Second,
		Handler:      router,
	}

	// Start the HTTP server
	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Server running at http://localhost:%s", env.Port)
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption or server error
	select {
	case err = <-srvErr:
		log.Printf("Server error: %v", err)
	case <-ctx.Done():
		log.Println("Shutting down gracefully...")
		stop()
	}

	// Gracefully shut down the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if shutdownErr := srv.Shutdown(shutdownCtx); shutdownErr != nil {
		log.Printf("Error during server shutdown: %v", shutdownErr)
	}
}
