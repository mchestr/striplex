package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/spf13/viper"
)

var C AppConfig

type Secret string

func (secret Secret) String() string {
	if secret == "" {
		return ""
	}
	return "*****"
}

func (secret Secret) Value() string {
	return string(secret)
}

func (secret Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(secret.String())
}

type AppConfig struct {
	Auth   AuthConfig
	Server ServerConfig
	Stripe StripeConfig
	Plex   PlexConfig
	Proxy  ProxyConfig
	Debug  bool
}

type AuthConfig struct {
	SessionSecret Secret
	SessionName   string
}

type ServerConfig struct {
	Address        string
	Hostname       string
	StaticPath     string
	Mode           string
	TrustedProxies *net.IPNet
}
type StripeConfig struct {
	PaymentMethodTypes  []string
	SecretKey           Secret
	WebhookSecret       Secret
	EntitlementName     string
	SubscriptionPriceID string
	DonationPriceID     string
}
type PlexConfig struct {
	ClientID        string
	AdminUserID     int
	ProductName     string
	SharedLibraries []string
	Token           Secret
	ServerID        string
}
type ProxyConfig struct {
	Enabled bool
	Url     string
}

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) error {
	var err error

	// Merge the environment configuration file
	defaultFileConfig := viper.New()
	defaultFileConfig.AddConfigPath("api/config/")
	defaultFileConfig.SetConfigName("default")
	if err = defaultFileConfig.ReadInConfig(); err != nil {
		slog.Warn("error on parsing default configuration file", "error", err)
	}

	// Merge the environment configuration file
	envFileConfig := viper.New()
	envFileConfig.AddConfigPath("api/config/")
	envFileConfig.SetConfigName(env)
	if err = envFileConfig.ReadInConfig(); err != nil {
		slog.Warn("error on parsing environment configuration file", "env", env, "error", err)
	}

	if err = defaultFileConfig.MergeConfigMap(envFileConfig.AllSettings()); err != nil {
		slog.Warn("error on merging configuration file", "env", env, "error", err)
	}

	config := viper.New()
	if err = config.MergeConfigMap(defaultFileConfig.AllSettings()); err != nil {
		slog.Warn("error on merging default configuration", "error", err)
		return err
	}
	config.SetConfigType("env")
	config.SetEnvPrefix("plefi")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()
	setDefaults(config)
	generateConfig(config)
	return nil
}

func setDefaults(config *viper.Viper) {
	config.SetDefault("server.address", ":8080")
	config.SetDefault("server.mode", "release")
	config.SetDefault("stripe.payment_method_types", []string{"card"})
	config.SetDefault("auth.session_secret", "changeme")
	config.SetDefault("auth.session_name", "plefi_session")
	config.SetDefault("debug", false)
}

func generateConfig(config *viper.Viper) {
	_, ipNet, err := net.ParseCIDR(config.GetString("server.trusted_proxies"))
	if err != nil {
		slog.Warn("error on parsing trusted proxies", "error", err)
	}
	C = AppConfig{
		Auth: AuthConfig{
			SessionSecret: Secret(config.GetString("auth.session_secret")),
			SessionName:   config.GetString("auth.session_name"),
		},
		Server: ServerConfig{
			Address:        config.GetString("server.address"),
			Hostname:       config.GetString("server.hostname"),
			StaticPath:     config.GetString("server.static_path"),
			Mode:           config.GetString("server.mode"),
			TrustedProxies: ipNet,
		},
		Stripe: StripeConfig{
			PaymentMethodTypes:  config.GetStringSlice("stripe.payment_method_types"),
			SecretKey:           Secret(config.GetString("stripe.secret_key")),
			WebhookSecret:       Secret(config.GetString("stripe.webhook_secret")),
			EntitlementName:     config.GetString("stripe.entitlement_name"),
			SubscriptionPriceID: config.GetString("stripe.subscription_price_id"),
			DonationPriceID:     config.GetString("stripe.donation_price_id"),
		},
		Plex: PlexConfig{
			ClientID:        config.GetString("plex.client_id"),
			AdminUserID:     config.GetInt("plex.admin_user_id"),
			ProductName:     config.GetString("plex.product_name"),
			SharedLibraries: strings.Split(config.GetString("plex.shared_libraries"), ","),
			Token:           Secret(config.GetString("plex.token")),
			ServerID:        config.GetString("plex.server_id"),
		},
		Proxy: ProxyConfig{
			Enabled: config.GetBool("proxy.enabled"),
			Url:     config.GetString("proxy.url"),
		},
	}
	if C.Debug {
		printJSON(C)
	}
}

func printJSON(obj interface{}) {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	fmt.Println(string(bytes))
}
