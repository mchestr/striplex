[<img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" align="right" width="25%" padding-right="350">]()

# `PleFi`

#### Integrate Stripe payments with Plex servers for subscription-based access!

<p align="left">
	<img src="https://img.shields.io/github/license/mchestr/plefi?style=for-the-badge&logo=opensourceinitiative&logoColor=white&color=00ADD8" alt="license">
	<img src="https://img.shields.io/github/last-commit/mchestr/plefi?style=for-the-badge&logo=git&logoColor=white&color=00ADD8" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/mchestr/plefi?style=for-the-badge&color=00ADD8" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/mchestr/plefi?style=for-the-badge&color=00ADD8" alt="repo-language-count">
</p>
<p align="left">
		<em>_Built with:_</em>
</p>
<p align="left">
	<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white" alt="Go">
	<img src="https://img.shields.io/badge/Stripe-008CDD.svg?style=for-the-badge&logo=Stripe&logoColor=white" alt="Stripe">
	<img src="https://img.shields.io/badge/Plex-E5A00D.svg?style=for-the-badge&logo=Plex&logoColor=white" alt="Plex">
	<img src="https://img.shields.io/badge/Docker-2496ED.svg?style=for-the-badge&logo=Docker&logoColor=white" alt="Docker">
</p>
<br>

## 🔗 Table of Contents

I. [📍 Overview](#-overview)
II. [👾 Features](#-features)
III. [📁 Project Structure](#-project-structure)
IV. [🚀 Getting Started](#-getting-started)
V. [🔧 Configuration](#-configuration)
VI. [📡 API Endpoints](#-api-endpoints)
VII. [🔰 Contributing](#-contributing)
VIII. [🎗 License](#-license)

---

## 📍 Overview

PleFi is a service that integrates Stripe payment processing with Plex media server, enabling subscription-based access management for Plex servers. It provides user authentication, payment handling, and subscription management to automate Plex server access control.

---

## 👾 Features

|      | Feature         | Summary       |
| :--- | :---:           | :---          |
| 🔐  | **Authentication**  | <ul><li>Plex user authentication using OAuth</li><li>Secure session management with encryption</li></ul> |
| 💳 | **Payments**  | <ul><li>Stripe subscription management and payment processing</li><li>Configurable subscription plans</li></ul> |
| 🔄 | **Webhooks** | <ul><li>Webhook support for Stripe events</li><li>Automated user access management based on subscription status</li></ul> |
| 🛠️ | **Configuration**  | <ul><li>Hierarchical configuration system with environment variables</li><li>Development and production modes</li></ul> |
| 🔌 | **API**    | <ul><li>RESTful endpoints for authentication and payment flows</li><li>Health check endpoints for monitoring</li></ul> |

---

## 📁 Project Structure

```sh
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

### 📂 Project Index

<details open>
	<summary><b><code>PleFi/</code></b></summary>
	<details>
		<summary><b>__root__</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b>main.go</b></td>
				<td>Application entry point that initializes configuration, services, and the HTTP server</td>
			</tr>
			<tr>
				<td><b>Dockerfile</b></td>
				<td>Container definition for building and running the PleFi application</td>
			</tr>
			<tr>
				<td><b>go.mod</b></td>
				<td>Go module definition specifying dependencies</td>
			</tr>
			<tr>
				<td><b>go.sum</b></td>
				<td>Go dependency checksums ensuring consistent builds</td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details>
		<summary><b>config/</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b>config.go</b></td>
				<td>Configuration initialization and management from files and environment variables</td>
			</tr>
			<tr>
				<td><b>default.yaml</b></td>
				<td>Default configuration values for all environments</td>
			</tr>
			<tr>
				<td><b>development.yaml</b></td>
				<td>Development-specific configuration settings</td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details>
		<summary><b>controllers/</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b>controller.go</b></td>
				<td>Main application controller handling common functionality</td>
			</tr>
			</table>
			<details>
				<summary><b>api/</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b>controller.go</b></td>
						<td>API controller implementation for general endpoints</td>
					</tr>
					</table>
				</blockquote>
			</details>
			<details>
				<summary><b>plex/</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b>controller.go</b></td>
						<td>Plex controller handling OAuth authentication</td>
					</tr>
					</table>
				</blockquote>
			</details>
			<details>
				<summary><b>stripe/</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b>controller.go</b></td>
						<td>Stripe controller for payment processing and webhooks</td>
					</tr>
					</table>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details>
		<summary><b>services/</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b>services.go</b></td>
				<td>Service container for dependency injection</td>
			</tr>
			<tr>
				<td><b>plex.go</b></td>
				<td>Plex service for interacting with Plex API</td>
			</tr>
			</table>
		</blockquote>
	</details>
</details>

---

## 🚀 Getting Started

### ☑️ Prerequisites

Before getting started with PleFi, ensure your runtime environment meets the following requirements:

- **Go 1.24 or higher**
- **Stripe account with webhook setup**
- **Plex account and server** (optional for development)


### ⚙️ Installation

1. Clone the repository:

```sh
❯ git clone https://github.com/mchestr/plefi.git
❯ cd plefi
```

2. Install dependencies:

```sh
❯ go mod download
```

3. Create a `.env` file in the project root with your configuration:

```sh
❯ cat > .env << EOL
PLEFI_STRIPE__SECRET_KEY="sk_test_your_stripe_secret_key"
PLEFI_STRIPE__WEBHOOK_SECRET="whsec_your_stripe_webhook_secret"
PLEFI_STRIPE__DEFAULT_PRICE_ID="price_your_default_price_id"

PLEFI_PLEX__CLIENT_ID="your_plex_client_id"
PLEFI_PLEX__PRODUCT="Your Plex Server Name"

PLEFI_SERVER__HOSTNAME="your-server-hostname.com"
PLEFI_SERVER__SESSION_SECRET="generate_a_random_secret_key"
PLEFI_SERVER__MODE="development"
EOL
```

### 🤖 Usage

**Run locally with Go:**

```sh
❯ go run main.go
```

**Run in development mode:**

```sh
❯ go run main.go -e development
```

**Using Docker:**

```sh
❯ docker build -t plefi .
❯ docker run -p 8080:8080 --env-file .env plefi
```

## 🔧 Configuration

PleFi uses a hierarchical configuration system:

1. Default values
2. Configuration files in `config/` directory
3. Environment variables (prefixed with `PLEFI_`)

### Environment Variables Reference

<details open>
<summary><b>Server Configuration</b></summary>
<blockquote>

- `PLEFI_SERVER__MODE` - Server mode (debug, release)
- `PLEFI_SERVER__ADDRESS` - Server bind address (default: `:8080`)
- `PLEFI_SERVER__HOSTNAME` - Server hostname for callbacks
- `PLEFI_SERVER__SESSION_SECRET` - Secret for session encryption
- `PLEFI_SERVER__TRUSTED_PROXIES` - Comma-separated list of trusted proxy IPs

</blockquote>
</details>

<details>
<summary><b>Stripe Configuration</b></summary>
<blockquote>

- `PLEFI_STRIPE__SECRET_KEY` - Stripe API secret key
- `PLEFI_STRIPE__WEBHOOK_SECRET` - Stripe webhook signing secret
- `PLEFI_STRIPE__DEFAULT_PRICE_ID` - Default subscription price ID

</blockquote>
</details>

<details>
<summary><b>Plex Configuration</b></summary>
<blockquote>

- `PLEFI_PLEX__CLIENT_ID` - Plex client identifier
- `PLEFI_PLEX__PRODUCT` - Plex product name

</blockquote>
</details>

<details>
<summary><b>Logging Configuration</b></summary>
<blockquote>

- `PLEFI_LOG__LEVEL` - Logging level (debug, info, warn, error)
- `PLEFI_LOG__FORMAT` - Logging format (json, text)

</blockquote>
</details>

## 📡 API Endpoints

<details open>
<summary><b>Application</b></summary>
<blockquote>

- `GET /` - Landing page
- `GET /health` - Health check endpoint

</blockquote>
</details>

<details>
<summary><b>Plex Authentication</b></summary>
<blockquote>

- `GET /plex/auth` - Initiate Plex authentication
- `GET /plex/callback` - Plex authentication callback

</blockquote>
</details>

<details>
<summary><b>Stripe</b></summary>
<blockquote>

- `POST /stripe/webhook` - Stripe webhook endpoint for subscription events
- `GET /stripe/checkout` - Create a Stripe checkout session for subscription
- `GET /stripe/success` - Handle successful subscription checkout
- `GET /stripe/cancel` - Handle cancelled subscription checkout

</blockquote>
</details>

---

## 🔰 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

<details closed>
<summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your GitHub account.
2. **Clone Locally**: Clone the forked repository to your local machine.
   ```sh
   git clone https://github.com/yourusername/plefi.git
   ```
3. **Create a New Branch**: Always work on a new branch.
   ```sh
   git checkout -b feature/amazing-feature
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message.
   ```sh
   git commit -m 'Add amazing feature'
   ```
6. **Push to GitHub**: Push the changes to your forked repository.
   ```sh
   git push origin feature/amazing-feature
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository.
</details>

---

## 🎗 License

See the [LICENSE](LICENSE) file for details.

---
