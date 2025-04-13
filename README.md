# Striplex

Striplex is a service that integrates Stripe payment processing with Plex media server, enabling subscription-based access management for Plex servers.

## Features

- Plex user authentication using OAuth
- Stripe subscription management and payment processing
- Discord user integration
- RESTful API for service interaction
- Webhook support for Stripe events
- PostgreSQL database for persistent storage

## Project Structure

```
striplex
├── config/                # Configuration files and setup
│   ├── config.go          # Configuration initialization
│   ├── default.yaml       # Default configuration values
│   └── development.yaml   # Development environment settings
├── controllers/           # Request handlers
│   ├── app.go             # Core application endpoints
│   ├── plex.go            # Plex authentication handlers
│   └── stripe.go          # Stripe webhook handlers
├── db/                    # Database connection logic
│   └── db.go              # Database initialization
├── model/                 # Data models
│   ├── discord_token.go   # Discord authentication token model
│   ├── discord_user.go    # Discord user model
│   ├── plex_token.go      # Plex authentication token model
│   └── plex_user.go       # Plex user model
├── server/                # HTTP server setup
│   ├── router.go          # Route definitions
│   └── server.go          # Server initialization
├── .env                   # Environment variables (not in git)
├── .gitignore             # Git ignore rules
├── .mise.toml             # Development environment configuration
├── Dockerfile             # Container definition
├── go.mod                 # Go module definition
├── go.sum                 # Go dependency checksums
├── main.go                # Application entry point
└── README.md              # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL database
- Stripe account with webhook setup
- Plex account and server (optional for development)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/striplex.git
cd striplex
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file in the project root with your configuration:

```
STRIPLEX_DATABASE__DSN="host=localhost user=username password=password dbname=striplex port=5432 sslmode=disable TimeZone=UTC"
STRIPLEX_STRIPE__WEBHOOK_SECRET="your_stripe_webhook_secret"
STRIPLEX_PLEX__CLIENT_ID="your_plex_client_id"
STRIPLEX_SERVER__HOSTNAME="your-server-hostname.com"
STRIPLEX_PLEX__PRODUCT="Your Plex Server Name"
STRIPLEX_SERVER__SESSION_SECRET="generate_a_random_secret_key"
```

### Running the Application

Start the application with:

```bash
go run main.go
```

For development mode:

```bash
go run main.go -e development
```

### Docker Deployment

Build and run the Docker container:

```bash
docker build -t striplex .
docker run -p 8080:8080 --env-file .env striplex
```

## API Endpoints

### Application

- `GET /` - Landing page
- `GET /health` - Health check endpoint
- `GET /whoami` - Get authenticated user information

### Plex Authentication

- `GET /api/v1/plex/auth` - Initiate Plex authentication
- `GET /api/v1/plex/callback` - Plex authentication callback

### Stripe

- `POST /api/v1/stripe/webhook` - Stripe webhook endpoint for subscription events

## Configuration

Striplex uses a hierarchical configuration system:

1. Default values
2. Configuration files in `config/` directory
3. Environment variables (prefixed with `STRIPLEX_`)

### Key Configuration Options

- `database.dsn` - PostgreSQL connection string
- `server.mode` - Server mode (debug, release)
- `server.address` - Server bind address (default: `:8080`)
- `stripe.webhook_secret` - Stripe webhook signing secret
- `plex.client_id` - Plex client identifier
- `plex.product` - Plex product name
- `server.hostname` - Server hostname for callbacks
- `server.session_secret` - Secret for session encryption

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.