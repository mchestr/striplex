package config

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings" // new import
	"testing"

	"github.com/spf13/viper"
)

func TestSecretMethods(t *testing.T) {
	tests := []struct {
		name     string
		input    Secret
		wantStr  string
		wantVal  string
		wantJSON string
	}{
		{"empty", Secret(""), "", "", "\"\""},
		{"nonempty", Secret("abc"), "*****", "abc", "\"*****\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}
			if got := tt.input.Value(); got != tt.wantVal {
				t.Errorf("Value() = %q, want %q", got, tt.wantVal)
			}
			b, err := tt.input.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(b) != tt.wantJSON {
				t.Errorf("MarshalJSON() = %s, want %s", b, tt.wantJSON)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	v := viper.New()
	setDefaults(v)
	if got := v.GetString("server.address"); got != ":8080" {
		t.Errorf("default server.address = %q, want %q", got, ":8080")
	}
	if got := v.GetString("server.mode"); got != "release" {
		t.Errorf("default server.mode = %q, want %q", got, "release")
	}
	if got := v.GetString("auth.session_secret"); got != "changeme" {
		t.Errorf("default auth.session_secret = %q, want %q", got, "changeme")
	}
	if got := v.GetString("auth.session_name"); got != "plefi_session" {
		t.Errorf("default auth.session_name = %q, want %q", got, "plefi_session")
	}
	if mt := v.GetStringSlice("stripe.payment_method_types"); len(mt) != 1 || mt[0] != "card" {
		t.Errorf("default stripe.payment_method_types = %v, want [card]", mt)
	}
}

func TestGenerateConfig(t *testing.T) {
	v := viper.New()
	// ... set all needed keys ...
	v.Set("server.address", ":9999")
	v.Set("server.hostname", "host")
	v.Set("server.static_path", "/static")
	v.Set("server.mode", "debug")
	v.Set("server.trusted_proxies", "10.0.0.0/8")
	v.Set("auth.session_secret", "sec")
	v.Set("auth.session_name", "name")
	v.Set("stripe.payment_method_types", []string{"card", "ideal"})
	v.Set("stripe.secret_key", "sk")
	v.Set("stripe.webhook_secret", "wh")
	v.Set("stripe.entitlement_name", "ent")
	v.Set("stripe.subscription_price_id", "sub")
	v.Set("stripe.donation_price_id", "don")
	v.Set("plex.client_id", "cid")
	v.Set("plex.admin_user_id", 7)
	v.Set("plex.product_name", "prod")
	v.Set("plex.shared_libraries", "lib1,lib2")
	v.Set("plex.token", "tok")
	v.Set("plex.server_id", "sid")
	v.Set("proxy.enabled", true)
	v.Set("proxy.url", "http://p")

	// capture stdout from printJSON
	var buf bytes.Buffer
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	generateConfig(v)

	w.Close()
	os.Stdout = oldOut
	io.Copy(&buf, r)

	// assert the global C was populated correctly
	if C.Server.Address != ":9999" {
		t.Errorf("Server.Address = %q, want %q", C.Server.Address, ":9999")
	}
	if C.Server.Hostname != "host" {
		t.Errorf("Server.Hostname = %q, want %q", C.Server.Hostname, "host")
	}
	if C.Server.Mode != "debug" {
		t.Errorf("Server.Mode = %q, want %q", C.Server.Mode, "debug")
	}
	if C.Auth.SessionSecret.Value() != "sec" {
		t.Errorf("Auth.SessionSecret.Value() = %q, want %q", C.Auth.SessionSecret.Value(), "sec")
	}
	if len(C.Stripe.PaymentMethodTypes) != 2 || C.Stripe.PaymentMethodTypes[1] != "ideal" {
		t.Errorf("Stripe.PaymentMethodTypes = %v, want [card ideal]", C.Stripe.PaymentMethodTypes)
	}
	if C.Stripe.SecretKey.Value() != "sk" {
		t.Errorf("Stripe.SecretKey.Value() = %q, want %q", C.Stripe.SecretKey.Value(), "sk")
	}
	ipnet := C.Server.TrustedProxies
	if ipnet == nil || ipnet.IP.String() != "10.0.0.0" {
		t.Errorf("TrustedProxies = %v, want 10.0.0.0/8", C.Server.TrustedProxies)
	}
	if !C.Proxy.Enabled || C.Proxy.Url != "http://p" {
		t.Errorf("Proxy = %+v, want Enabled=true Url=http://p", C.Proxy)
	}
}

func TestInitIntegration(t *testing.T) {
	// create temp project layout
	tmp := t.TempDir()
	apiconf := filepath.Join(tmp, "internal", "config")
	if err := os.MkdirAll(apiconf, 0755); err != nil {
		t.Fatal(err)
	}
	// default config
	defaultY := `
server:
  address: ":1111"
  hostname: "h"
  trusted_proxies: "192.168.1.0/24"
auth:
  session_secret: "dsecret"
  session_name: "dname"
`
	// test env override
	testY := `
server:
  address: ":2222"
auth:
  session_secret: "tsecret"
stripe:
  secret_key: "tkey"
plex:
  admin_user_id: 99
`
	if err := os.WriteFile(filepath.Join(apiconf, "default.yaml"), []byte(defaultY), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(apiconf, "test.yaml"), []byte(testY), 0644); err != nil {
		t.Fatal(err)
	}

	// switch cwd so Init finds our files
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmp)

	// clear any PLEFI_* environment variables so AutomaticEnv won't override file settings
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "PLEFI_") {
			os.Unsetenv(strings.SplitN(e, "=", 2)[0])
		}
	}

	if err := Init("test"); err != nil {
		t.Errorf("Init returned error: %v", err)
	}

	// validate merged config
	if C.Server.Address != ":2222" {
		t.Errorf("Init Server.Address = %q, want %q", C.Server.Address, ":2222")
	}
	if C.Server.Hostname != "h" {
		t.Errorf("Init Server.Hostname = %q, want %q", C.Server.Hostname, "h")
	}
	if C.Auth.SessionSecret.Value() != "tsecret" {
		t.Errorf("Init SessionSecret = %q, want %q", C.Auth.SessionSecret.Value(), "tsecret")
	}
	if C.Stripe.SecretKey.Value() != "tkey" {
		t.Errorf("Init Stripe.SecretKey = %q, want %q", C.Stripe.SecretKey.Value(), "tkey")
	}
	if C.Plex.AdminUserID != 99 {
		t.Errorf("Init Plex.AdminUserID = %d, want %d", C.Plex.AdminUserID, 99)
	}
	// check that default trusted_proxies still applied
	if C.Server.TrustedProxies == nil || C.Server.TrustedProxies.IP.String() != "192.168.1.0" {
		t.Errorf("Init TrustedProxies = %v, want 192.168.1.0/24", C.Server.TrustedProxies)
	}
}
