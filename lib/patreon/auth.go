package patreon

import (
	"errors"
	"net/http"
	"strings"

	"github.com/adamgoose/abots/lib/structure"
	"github.com/charmbracelet/log"
	"github.com/dghubble/sling"
	"github.com/spf13/viper"
)

type (
	// Authenticator is an interface for getting a session ID.
	Authenticator interface {
		// GetSessionID returns the session ID.
		GetSessionID() (string, error)
	}

	// SessionIDAuthenticator is an Authenticator that uses a static session ID.
	SessionIDAuthenticator struct {
		SessionID string
	}
	// CredentialAuthenticator is an Authenticator that uses credentials to get a session ID.
	CredentialAuthenticator struct {
		Email    string
		Password string
	}
	// CachedAuthenticator is an Authenticator that caches the session ID.
	CachedAuthenticator struct {
		Authenticator Authenticator
		Cache         *structure.DB
		Log           *log.Logger
	}
)

// NewAuthenticator returns a new Authenticator.
func NewAuthenticator(log *log.Logger, db *structure.DB) (Authenticator, error) {
	sida, err := NewSessionIDAuthenticator()
	if err == nil {
		log.Debug("using session_id authenticator")
		return sida, nil
	}

	ca, err := NewCredentialAuthenticator()
	if err == nil {
		log.Debug("using cached credential authenticator")
		return NewCachedAuthenticator(log, ca, db), nil
	}

	return nil, err
}

// NewSessionIDAuthenticator returns a new SessionIDAuthenticator.
func NewSessionIDAuthenticator() (*SessionIDAuthenticator, error) {
	sid := viper.GetString("patreon_session_id")
	if sid == "" {
		return nil, errors.New("no session_id found")
	}

	return &SessionIDAuthenticator{
		SessionID: sid,
	}, nil
}

// GetSessionID returns the session ID.
func (a *SessionIDAuthenticator) GetSessionID() (string, error) {
	return a.SessionID, nil
}

// NewCredentialAuthenticator returns a new CredentialAuthenticator.
func NewCredentialAuthenticator() (*CredentialAuthenticator, error) {
	e := viper.GetString("patreon_email")
	p := viper.GetString("patreon_password")

	if e == "" || p == "" {
		return nil, errors.New("no email or password found")
	}

	return &CredentialAuthenticator{
		Email:    e,
		Password: p,
	}, nil
}

// GetSessionID returns the session ID.
func (a *CredentialAuthenticator) GetSessionID() (string, error) {
	login, err := sling.New().
		Post("https://www.patreon.com/api/auth").
		QueryStruct(&LoginQuery{
			Include:        "user.null",
			UserFields:     "[]",
			JsonAPIVersion: "1.0",
		}).
		BodyJSON(map[string]interface{}{
			"data": map[string]interface{}{
				"attributes": map[string]interface{}{
					"auth_context": "auth",
					"patreon_auth": map[string]interface{}{
						"allow_account_creation": false,
						"email":                  a.Email,
						"password":               a.Password,
						"redirect_target":        "https://www.patreon.com/home",
					},
				},
				"relationships": map[string]interface{}{},
				"type":          "genericPatreonApi",
			},
		}).
		Request()
	if err != nil {
		return "", err
	}

	r, err := http.DefaultClient.Do(login)
	if err != nil {
		return "", err
	}

	for _, h := range r.Header.Values("Set-Cookie") {
		if !strings.HasPrefix(h, "session_id") {
			continue
		}

		kv := strings.Split(h, ";")[0]
		v := strings.Split(kv, "=")[1]

		return v, nil
	}

	return "", errors.New("failed to retrieve session")
}

// NewCachedAuthenticator returns a new CachedAuthenticator.
func NewCachedAuthenticator(log *log.Logger, a Authenticator, c *structure.DB) *CachedAuthenticator {
	return &CachedAuthenticator{
		Log:           log.With("component", "cached_authenticator"),
		Authenticator: a,
		Cache:         c,
	}
}

// GetSessionID returns the session ID.
func (a *CachedAuthenticator) GetSessionID() (string, error) {
	var sid string

	// Check for a cached value
	if err := a.Cache.View(func(tx *structure.Tx) error {
		s, err := tx.Get(PatreonScraperBucket, []byte("session_id"))
		if err == nil {
			sid = string(s.Value)
		}
		return nil
	}); err != nil {
		return "", err
	}

	// If we have a cached value, return it
	if sid != "" {
		a.Log.Debug("using cached session_id")
		return sid, nil
	}

	// Otherwise, get a new session ID
	a.Log.Debug("getting new session_id")
	sid, err := a.Authenticator.GetSessionID()
	if err != nil {
		return "", err
	}

	// Cache the session ID
	a.Log.Debug("caching session_id")
	if err := a.Cache.Update(func(tx *structure.Tx) error {
		return tx.Put(PatreonScraperBucket, []byte("session_id"), []byte(sid), 0)
	}); err != nil {
		return "", err
	}

	return sid, nil
}

type LoginQuery struct {
	Include        string `url:"include"`
	UserFields     string `url:"fields[user]"`
	JsonAPIVersion string `url:"json-api-version"`
}
