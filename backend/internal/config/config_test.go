package config

import (
	"reflect"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("CORS_ORIGIN", "")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if !reflect.DeepEqual(cfg.AllowedOrigins, []string{"*"}) {
		t.Errorf("AllowedOrigins = %v, want [*]", cfg.AllowedOrigins)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("CORS_ORIGIN", "http://localhost:5173, https://example.com ,")

	cfg := Load()
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want 9090", cfg.Port)
	}
	want := []string{"http://localhost:5173", "https://example.com"}
	if !reflect.DeepEqual(cfg.AllowedOrigins, want) {
		t.Errorf("AllowedOrigins = %v, want %v", cfg.AllowedOrigins, want)
	}
}
