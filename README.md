# PleFi

PleFi is a service that integrates Stripe payment processing with Plex media server, enabling subscription-based access management for Plex servers.

## Features

- Plex user authentication using OAuth
- Stripe subscription management and payment processing
- Webhook support for Stripe events

## Project Structure

```
plefi
├── config/                 # Configuration files and setup
│   ├── config.go           # Configuration initialization
│   ├── default.yaml        # Default configuration values
│   └── development.yaml    # Development environment settings
├── controllers/            # Request handlers
│   ├── api/                # API endpoints
│   │   └── controller.go   # API controller implementation
│   ├── controller.go       # Main application controller
│   ├── plex/               # Plex authentication handlers
│   │   └── controller.go   # Plex controller implementation
│   └── stripe/             # Stripe payment handlers
│       └── controller.go   # Stripe controller implementation
├── model/                  # Data models
│   └── user_info.go        # User information model
├── services/               # Service implementations
│   └── services.go         # Service container
│   └── plex.go             # Plex Service
├── server/                 # HTTP server setup
│   ├── router.go           # Route definitions
│   └── server.go           # Server initialization with graceful shutdown
├── views/                  # HTML templates
│   ├── index.tmpl          # Main landing page template
│   ├── stripe_success.tmpl # Subscription success page
│   └── stripe_cancel.tmpl  # Subscription cancellation page
├── .env                    # Environment variables (not in git)
├── .gitignore              # Git ignore rules
├── Dockerfile              # Container definition
├── go.mod                  # Go module definition
├── go.sum                  # Go dependency checksums
├── main.go                 # Application entry point
└── README.md               # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Stripe account with webhook setup
- Plex account and server (optional for development)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/mchestr/plefi.git
cd plefi
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file in the project root with your configuration:

```
PLEFI_STRIPE__SECRET_KEY="sk_test_your_stripe_secret_key"
PLEFI_STRIPE__WEBHOOK_SECRET="whsec_your_stripe_webhook_secret"
PLEFI_STRIPE__DEFAULT_PRICE_ID="price_your_default_price_id"

PLEFI_PLEX__CLIENT_ID="your_plex_client_id"
PLEFI_PLEX__PRODUCT="Your Plex Server Name"

PLEFI_SERVER__HOSTNAME="your-server-hostname.com"
PLEFI_SERVER__SESSION_SECRET="generate_a_random_secret_key"
PLEFI_SERVER__MODE="development"
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

PleFi uses a hierarchical configuration system:

1. Default values
2. Configuration files in `config/` directory
3. Environment variables (prefixed with `PLEFI_`)

### Environment Variables Reference

Below is a complete reference of all environment variables used in the application. All environment variables should be prefixed with `PLEFI_` and use double underscores to represent nesting.

#### Server Configuration
- `PLEFI_SERVER__MODE` - Server mode (debug, release)
- `PLEFI_SERVER__ADDRESS` - Server bind address (default: `:8080`)
- `PLEFI_SERVER__HOSTNAME` - Server hostname for callbacks
- `PLEFI_SERVER__SESSION_SECRET` - Secret for session encryption
- `PLEFI_SERVER__TRUSTED_PROXIES` - Comma-separated list of trusted proxy IPs

#### Stripe Configuration
- `PLEFI_STRIPE__SECRET_KEY` - Stripe API secret key
- `PLEFI_STRIPE__WEBHOOK_SECRET` - Stripe webhook signing secret
- `PLEFI_STRIPE__DEFAULT_PRICE_ID` - Default subscription price ID

#### Plex Configuration
- `PLEFI_PLEX__CLIENT_ID` - Plex client identifier
- `PLEFI_PLEX__PRODUCT` - Plex product name

#### Logging Configuration
- `PLEFI_LOG__LEVEL` - Logging level (debug, info, warn, error)
- `PLEFI_LOG__FORMAT` - Logging format (json, text)ret

### Key Configuration Optionstripe webhook signing secret
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
4. Push to the branch (`git push origin feature/amazing-feature`)details.

## License

This project is licensed under the MIT License - see the LICENSE file for details.This project is licensed under the MIT License - see the LICENSE file for details.5. Open a Pull Request## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
