package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"striplex/config"
	"striplex/db"
	"striplex/server"

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

	// Initialize application components
	if err := initApp(*environment); err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
}

// initApp initializes all application components
func initApp(environment string) error {
	// Initialize configuration
	if err := config.Init(environment); err != nil {
		return fmt.Errorf("config initialization error: %w", err)
	}

	// Set Stripe API key
	stripe.Key = config.Config.GetString("stripe.secret_key")
	if stripe.Key == "" {
		return fmt.Errorf("stripe API key not configured")
	}

	// Initialize database connection
	if err := db.Connect(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Initialize server components
	if err := server.Init(); err != nil {
		return fmt.Errorf("server initialization error: %w", err)
	}

	return nil
}
