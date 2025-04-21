package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"plefi/internal/config"
	"plefi/internal/db"
	"plefi/internal/server"
	"plefi/internal/services"
	"syscall"
	"time"

	"github.com/stripe/stripe-go/v82"
)

func main() {
	// Parse command line flags
	environment := flag.String("e", "development", "Environment to run the application (development, production)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [-e environment]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Initialize and run application components
	if err := runApp(*environment); err != nil {
		slog.Error("Failed to run application", "error", err)
		os.Exit(1)
	}
}

// initApp initializes all application components
func initApp(environment string) (*server.Server, error) {
	// Initialize configuration
	if err := config.Init(environment); err != nil {
		return nil, fmt.Errorf("config initialization error: %w", err)
	}

	// Create HTTP client with reasonable timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	if config.C.Proxy.Enabled {
		slog.Info("Proxy enabled, setting up HTTP client with proxy")
		proxyURL, _ := url.Parse(config.C.Proxy.Url)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	svcs := services.NewServices(httpClient)
	if config.C.Plex.AdminUserID == 0 {
		plexUser, err := svcs.Plex.GetUserDetails(context.Background(), config.C.Plex.Token.Value())
		if err != nil {
			return nil, fmt.Errorf("failed to get Plex admin user details: %w", err)
		}
		config.C.Plex.AdminUserID = plexUser.ID
		slog.Info("Plex admin user ID set in config",
			"plex_admin_user_id", config.C.Plex.AdminUserID,
			"plex_username", plexUser.Username)
	}
	if config.C.Plex.MachineIdentifier == "" {
		machineID, err := svcs.Plex.GetMachineIdentity(context.Background(), config.C.Plex.Url, config.C.Plex.Token.Value())
		if err != nil {
			return nil, fmt.Errorf("failed to get Plex machine identifier: %w", err)
		}
		config.C.Plex.MachineIdentifier = machineID
		slog.Info("Plex machine identifier set in config",
			"plex_machine_identifier", config.C.Plex.MachineIdentifier)
	}

	// Set Stripe API key
	stripe.Key = config.C.Stripe.SecretKey.Value()
	if stripe.Key == "" {
		return nil, fmt.Errorf("stripe API key not configured")
	}
	stripe.SetHTTPClient(httpClient)

	slog.Info("Initializing database connection",
		"driver", config.C.Database.Driver,
		"dsn", config.C.Database.Dsn)
	// Initialize database connection
	if err := db.Init(config.C.Database.Driver, config.C.Database.Dsn.Value()); err != nil {
		slog.Error("db failed to open", "error", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.DB.Migrate(context.Background()); err != nil {
		slog.Error("db failed to migrate", "error", err)
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize server components
	srv, err := server.Init(svcs, httpClient)
	if err != nil {
		return nil, fmt.Errorf("server initialization error: %w", err)
	}

	return srv, nil
}

// runApp initializes the application and starts the server with graceful shutdown
func runApp(environment string) error {
	// Initialize application
	srv, err := initApp(environment)
	if err != nil {
		return err
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("Server error", "error", err)
		}
	}()
	slog.Info("Server started")

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("Server exiting")
	return nil
}
