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
	"strconv"
	"syscall"
	"time"

	"plefi/config"
	"plefi/server"
	"plefi/services"

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
	if config.Config.GetBool("proxy.enabled") {
		proxyURL, _ := url.Parse(config.Config.GetString("proxy.url"))
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	svcs := services.NewServices(httpClient)
	if config.Config.GetString("plex.admin_user_id") == "" {
		plexUser, err := svcs.Plex.GetUserDetails(context.Background(), config.Config.GetString("plex.token"))
		if err != nil {
			return nil, fmt.Errorf("failed to get Plex admin user details: %w", err)
		}
		config.Config.Set("plex.admin_user_id", strconv.Itoa(plexUser.ID))
		slog.Info("Plex admin user ID set in config",
			"plex_admin_user_id", config.Config.GetString("plex.admin_user_id"),
			"plex_username", plexUser.Username)
	}

	// Set Stripe API key
	stripe.Key = config.Config.GetString("stripe.secret_key")
	if stripe.Key == "" {
		return nil, fmt.Errorf("stripe API key not configured")
	}
	stripe.SetHTTPClient(httpClient)

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
