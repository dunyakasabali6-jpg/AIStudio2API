package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Mag1cFall/AIStudio2API/internal/aistudio"
)

type bootstrapCookieAccount struct {
	Email   string          `json:"email"`
	Cookie  string          `json:"cookie"`
	Cookies json.RawMessage `json:"cookies"`
}

// bootstrapAccounts restores accounts from deployment secrets before the runtime pool is created.
// The environment variable is not copied into .env or logs; account credentials are persisted as normal storage-state files.
func bootstrapAccounts(store *aistudio.AccountStore, accounts []*aistudio.Account, requests *requestRegistry) ([]*aistudio.Account, error) {
	rawJSON := strings.TrimSpace(os.Getenv("BOOTSTRAP_COOKIES_JSON"))
	rawCookie := strings.TrimSpace(os.Getenv("BOOTSTRAP_COOKIE"))
	rawEmail := strings.TrimSpace(os.Getenv("BOOTSTRAP_EMAIL"))
	if rawJSON == "" && rawCookie == "" {
		if rawEmail != "" {
			return nil, fmt.Errorf("BOOTSTRAP_EMAIL requires BOOTSTRAP_COOKIE")
		}
		return accounts, nil
	}
	if rawJSON != "" && (rawCookie != "" || rawEmail != "") {
		return nil, fmt.Errorf("use BOOTSTRAP_COOKIES_JSON or BOOTSTRAP_EMAIL with BOOTSTRAP_COOKIE, not both")
	}

	var entries []bootstrapCookieAccount
	if rawJSON != "" {
		if err := json.Unmarshal([]byte(rawJSON), &entries); err != nil {
			return nil, fmt.Errorf("parse BOOTSTRAP_COOKIES_JSON: %w", err)
		}
	} else {
		entries = []bootstrapCookieAccount{{Email: rawEmail, Cookie: rawCookie}}
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("bootstrap cookie list is empty")
	}

	existing := make(map[string]struct{}, len(accounts))
	for _, account := range accounts {
		existing[strings.ToLower(account.ID)] = struct{}{}
	}
	for index, entry := range entries {
		email := strings.TrimSpace(entry.Email)
		cookieText := strings.TrimSpace(entry.Cookie)
		if len(entry.Cookies) != 0 && string(entry.Cookies) != "null" {
			var value string
			if err := json.Unmarshal(entry.Cookies, &value); err == nil {
				cookieText = value
			} else {
				cookieText = string(entry.Cookies)
			}
		}
		if email == "" || cookieText == "" {
			return nil, fmt.Errorf("bootstrap account %d must contain email and cookie/cookies", index+1)
		}
		email = strings.ToLower(email)
		if _, ok := existing[email]; ok {
			requests.log("auth", "INFO", "bootstrap account already exists; skipped", "account", email)
			continue
		}
		state, err := aistudio.ParseCookieText(cookieText)
		if err != nil {
			return nil, fmt.Errorf("bootstrap account %s has invalid cookie data or is missing required authentication cookies", email)
		}
		if _, err := aistudio.NewSigner().Authorization(state); err != nil {
			return nil, fmt.Errorf("bootstrap account %s cannot sign AI Studio requests: %w", email, err)
		}
		accountConfig := aistudio.DefaultAccountConfig(email)
		if err := accountConfig.Validate(); err != nil {
			return nil, fmt.Errorf("bootstrap account %s: %w", email, err)
		}
		account, publish, err := store.Create(accountConfig, state)
		if errors.Is(err, aistudio.ErrAccountExists) {
			requests.log("auth", "INFO", "bootstrap account directory already exists; skipped", "account", email)
			existing[email] = struct{}{}
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("persist bootstrap account %s: %w", email, err)
		}
		if err := publish.Release(); err != nil {
			return nil, fmt.Errorf("release bootstrap account %s: %w", email, err)
		}
		accounts = append(accounts, account)
		existing[email] = struct{}{}
		requests.log("auth", "INFO", "bootstrap account restored", "account", email)
	}
	return accounts, nil
}
