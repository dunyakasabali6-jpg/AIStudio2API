package app

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Mag1cFall/AIStudio2API/internal/aistudio"
)

const testBootstrapCookie = "__Secure-1PSID=session; SAPISID=sapi; __Secure-1PAPISID=pap; __Secure-3PAPISID=pap3"

func TestBootstrapAccountsCreatesAccountFromCookieEnvironment(t *testing.T) {
	t.Setenv("BOOTSTRAP_EMAIL", "bootstrap@example.com")
	t.Setenv("BOOTSTRAP_COOKIE", testBootstrapCookie)
	t.Setenv("BOOTSTRAP_COOKIES_JSON", "")

	store := aistudio.NewAccountStore(filepath.Join(t.TempDir(), "auth"))
	accounts, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	requests := newRequestRegistry(ctx)

	accounts, err = bootstrapAccounts(store, accounts, requests)
	if err != nil {
		t.Fatalf("bootstrapAccounts() error = %v", err)
	}
	if len(accounts) != 1 || accounts[0].ID != "bootstrap@example.com" {
		t.Fatalf("bootstrapAccounts() accounts = %#v, want one bootstrap account", accounts)
	}

	reloaded, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(reloaded) != 1 || reloaded[0].ID != "bootstrap@example.com" {
		t.Fatalf("persisted accounts = %#v, want one bootstrap account", reloaded)
	}
}

func TestBootstrapAccountsRejectsConflictingEnvironmentSources(t *testing.T) {
	t.Setenv("BOOTSTRAP_EMAIL", "bootstrap@example.com")
	t.Setenv("BOOTSTRAP_COOKIE", testBootstrapCookie)
	t.Setenv("BOOTSTRAP_COOKIES_JSON", "[]")

	store := aistudio.NewAccountStore(filepath.Join(t.TempDir(), "auth"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := bootstrapAccounts(store, nil, newRequestRegistry(ctx))
	if err == nil {
		t.Fatal("bootstrapAccounts() error = nil, want conflicting source error")
	}
}

func TestBootstrapAccountsRejectsCookieWithoutEmail(t *testing.T) {
	t.Setenv("BOOTSTRAP_EMAIL", "")
	t.Setenv("BOOTSTRAP_COOKIE", testBootstrapCookie)
	t.Setenv("BOOTSTRAP_COOKIES_JSON", "")

	store := aistudio.NewAccountStore(filepath.Join(t.TempDir(), "auth"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := bootstrapAccounts(store, nil, newRequestRegistry(ctx))
	if err == nil {
		t.Fatal("bootstrapAccounts() error = nil, want missing email error")
	}
}

func TestBootstrapAccountsRejectsEmailWithoutCookie(t *testing.T) {
	t.Setenv("BOOTSTRAP_EMAIL", "bootstrap@gmail.com")
	t.Setenv("BOOTSTRAP_COOKIE", "")
	t.Setenv("BOOTSTRAP_COOKIES_JSON", "")

	store := aistudio.NewAccountStore(filepath.Join(t.TempDir(), "auth"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := bootstrapAccounts(store, nil, newRequestRegistry(ctx))
	if err == nil {
		t.Fatal("bootstrapAccounts() error = nil, want BOOTSTRAP_COOKIE requirement")
	}
}
