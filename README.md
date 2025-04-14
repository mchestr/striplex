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
plefi
├── config/                # Configuration files and setup
│   ├── config.go          # Configuration initialization
│   ├── default.yaml       # Default configuration values
│   └── development.yaml   # Development environment settings
├── controllers/           # Request handlers
│   ├── app/               # Core application endpoints
│   │   └── controller.go  # App controller implementation
│   ├── plex/              # Plex authentication handlers
│   │   └── controller.go  # Plex controller implementation
│   └── stripe/            # Stripe payment handlers
│       └── controller.go  # Stripe controller implementation
├── db/                    # Database connection logic
│   └── db.go              # Database initialization
├── model/                 # Data models
│   ├── discord_token.go   # Discord authentication token model
│   ├── discord_user.go    # Discord user model
│   ├── plex_token.go      # Plex authentication token model
│   ├── user_info.go       # User information model
│   └── plex_user.go       # Plex user model
├── services/              # Service implementations
│   └── wizarr.go          # Wizarr service for Plex invitations
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
git clone https://github.com/yourusername/plefi.git
cd plefi
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file in the project root with your configuration:

```
PLEFI_DATABASE__DSN="host=localhost user=username password=password dbname=plefi port=5432 sslmode=disable TimeZone=UTC"
PLEFI_STRIPE__WEBHOOK_SECRET="your_stripe_webhook_secret"
PLEFI_PLEX__CLIENT_ID="your_plex_client_id"
PLEFI_SERVER__HOSTNAME="your-server-hostname.com"
PLEFI_PLEX__PRODUCT="Your Plex Server Name"
PLEFI_SERVER__SESSION_SECRET="generate_a_random_secret_key"
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
docker build -t plefi .
docker run -p 8080:8080 --env-file .env plefi
```

## API Endpoints

### Application

- `GET /` - Landing page
- `GET /health` - Health check endpoint
- `GET /whoami` - Get authenticated user information

### Plex Authentication

- `GET /plex/auth` - Initiate Plex authentication
- `GET /plex/callback` - Plex authentication callback

### Stripe

- `POST /stripe/webhook` - Stripe webhook endpoint for subscription events
- `GET /stripe/checkout` - Create a Stripe checkout session for subscription
- `GET /stripe/success` - Handle successful subscription checkout
- `GET /stripe/cancel` - Handle cancelled subscription checkout

### Wizarr Integration

- Automatic generation of Plex server invites through Wizarr after successful subscription
- Email delivery of invitation links to subscribers

## Configuration

Striplex uses a hierarchical configuration system:

1. Default values
2. Configuration files in `config/` directory
3. Environment variables (prefixed with `PLEFI_`)

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