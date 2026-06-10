package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLdapUnavailable    = errors.New("ldap server unavailable")
)

type ldapConfig struct {
	URL         string
	BaseDN      string
	BindDN      string
	BindPass    string
	UserFilter  string
	AdminGroup  string
	TLSInsecure bool
	StartTLS    bool
}

func loadLdapConfig() *ldapConfig {
	return &ldapConfig{
		URL:         os.Getenv("LDAP_URL"),
		BaseDN:      os.Getenv("LDAP_BASE_DN"),
		BindDN:      os.Getenv("LDAP_BIND_DN"),
		BindPass:    os.Getenv("LDAP_BIND_PASSWORD"),
		UserFilter:  os.Getenv("LDAP_USER_FILTER"),
		AdminGroup:  os.Getenv("LDAP_ADMIN_GROUP"),
		TLSInsecure: os.Getenv("LDAP_TLS_INSECURE") == "true",
		StartTLS:    os.Getenv("LDAP_STARTTLS") == "true",
	}
}

func LdapAuthenticate(username, password string) (*User, error) {
	cfg := loadLdapConfig()

	slog.Debug("Connecting to LDAP server", "url", cfg.URL)
	conn, err := ldap.DialURL(cfg.URL)
	if err != nil {
		slog.Error("Failed to connect to LDAP server", "error", err, "url", cfg.URL)
		return nil, ErrLdapUnavailable
	}
	defer conn.Close()

	if cfg.StartTLS {
		tlsCfg := &tls.Config{InsecureSkipVerify: cfg.TLSInsecure}
		if err := conn.StartTLS(tlsCfg); err != nil {
			slog.Error("LDAP StartTLS failed", "error", err)
			return nil, ErrLdapUnavailable
		}
	}

	if err := conn.Bind(cfg.BindDN, cfg.BindPass); err != nil {
		slog.Error("LDAP service account bind failed", "error", err)
		return nil, ErrLdapUnavailable
	}

	filter := fmt.Sprintf(cfg.UserFilter, ldap.EscapeFilter(username))
	searchReq := ldap.NewSearchRequest(
		cfg.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"dn", "memberOf"},
		nil,
	)

	slog.Debug("Searching LDAP for user", "username", username, "filter", filter)
	result, err := conn.Search(searchReq)
	if err != nil {
		slog.Error("LDAP search failed", "error", err)
		return nil, ErrLdapUnavailable
	}

	if len(result.Entries) == 0 {
		slog.Warn("User not found in LDAP", "username", username)
		return nil, ErrInvalidCredentials
	}

	userDN := result.Entries[0].DN

	if err := conn.Bind(userDN, password); err != nil {
		slog.Warn("LDAP user credential bind failed", "username", username)
		return nil, ErrInvalidCredentials
	}

	role := "viewer"
	adminGroup := strings.ToLower(strings.TrimSpace(cfg.AdminGroup))
	for _, entry := range result.Entries {
		for _, attr := range entry.Attributes {
			if attr.Name == "memberOf" {
				for _, val := range attr.Values {
					if strings.EqualFold(strings.TrimSpace(val), adminGroup) {
						role = "admin"
					}
				}
			}
		}
	}

	slog.Info("LDAP authentication successful", "username", username, "role", role)
	return &User{Username: username, Role: role}, nil
}
