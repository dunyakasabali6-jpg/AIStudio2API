package config

import "testing"

// TestValidateProxyCredentials 验证代理 URL 接受账号密码并拒绝空用户名
func TestValidateProxyCredentials(t *testing.T) {
	for _, value := range []string{"http://user:secret@127.0.0.1:8080", "socks5://user:secret@127.0.0.1:1080"} {
		if err := ValidateProxy(value); err != nil {
			t.Errorf("ValidateProxy(%q) = %v", value, err)
		}
	}
	if err := ValidateProxy("http://:secret@127.0.0.1:8080"); err == nil {
		t.Error("空用户名被接受")
	}
}

func TestServerDefaultsUseSmallSpiderMonkeyProfile(t *testing.T) {
	cfg := Default()
	if cfg.WAABackend != WAABackendSpiderMonkey {
		t.Fatalf("WAA backend = %q, want %q", cfg.WAABackend, WAABackendSpiderMonkey)
	}
	if cfg.WarmWorkerLimit != 1 || cfg.MaxActiveWorkers != 2 || cfg.WarmStartupConcurrency != 1 || cfg.PerAccountConcurrency != 1 {
		t.Fatalf("worker defaults = warm:%d active:%d startup:%d per-account:%d, want 1/2/1/1",
			cfg.WarmWorkerLimit, cfg.MaxActiveWorkers, cfg.WarmStartupConcurrency, cfg.PerAccountConcurrency)
	}
	if cfg.LogFile {
		t.Fatal("LogFile default = true, want false for container profile")
	}
}

func TestCamoufoxBackendIsRejected(t *testing.T) {
	cfg := Default()
	if err := setWAABackend(&cfg, WAABackendCamoufox); err == nil {
		t.Fatal("setWAABackend(camoufox) error = nil, want unsupported backend error")
	}
	if cfg.WAABackend != WAABackendSpiderMonkey {
		t.Fatalf("WAA backend after rejected value = %q, want SpiderMonkey", cfg.WAABackend)
	}
}

func TestAdminAuthenticationRequiresProxyAPIKey(t *testing.T) {
	cfg := Default()
	cfg.AdminAuthEnabled = true
	cfg.AdminUsername = "admin"
	cfg.AdminPassword = "strong-password"
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want missing PROXY_API_KEY error")
	}
	cfg.ProxyAPIKey = "separate-api-key"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() with API key error = %v", err)
	}
}
