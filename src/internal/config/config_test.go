package config

import (
	"reflect"
	"testing"
)

func Test_getEnv(t *testing.T) {
	t.Run("returns default when unset", func(t *testing.T) {
		t.Setenv("SOME_KEY", "")
		got := getEnv("SOME_KEY", "default-value")
		if got != "default-value" {
			t.Fatalf("expected default-value, got %q", got)
		}
	})

	t.Run("returns default when explicitly empty", func(t *testing.T) {
		t.Setenv("EMPTY_KEY", "")
		got := getEnv("EMPTY_KEY", "fallback")
		if got != "fallback" {
			t.Fatalf("expected fallback, got %q", got)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv("SET_KEY", "real-value")
		got := getEnv("SET_KEY", "ignored-default")
		if got != "real-value" {
			t.Fatalf("expected real-value, got %q", got)
		}
	})
}

func TestLoad(t *testing.T) {
	const defaultDB = "host=localhost port=5432 user=user password=password dbname=workoutpal sslmode=disable"

	t.Run("uses defaults when env not set or empty", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "")
		t.Setenv("GOOGLE_CLIENT_ID", "")
		t.Setenv("JWT_SECRET", "")
		t.Setenv("PORT", "")

		cfg := Load()
		if cfg == nil {
			t.Fatal("Load() returned nil")
		}

		want := &Config{
			DatabaseURL:    defaultDB,
			GoogleClientID: "",
			JWTSecret:      "your-secret-key",
			Port:           "8080",
		}
		if !reflect.DeepEqual(cfg, want) {
			t.Fatalf("defaults mismatch\n got:  %#v\n want: %#v", cfg, want)
		}
	})

	t.Run("uses env values when set", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/app?sslmode=disable")
		t.Setenv("GOOGLE_CLIENT_ID", "abc.apps.googleusercontent.com")
		t.Setenv("JWT_SECRET", "shh-its-a-secret")
		t.Setenv("PORT", "9090")

		cfg := Load()
		if cfg == nil {
			t.Fatal("Load() returned nil")
		}

		want := &Config{
			DatabaseURL:    "postgres://user:pass@db:5432/app?sslmode=disable",
			GoogleClientID: "abc.apps.googleusercontent.com",
			JWTSecret:      "shh-its-a-secret",
			Port:           "9090",
		}
		if !reflect.DeepEqual(cfg, want) {
			t.Fatalf("env override mismatch\n got:  %#v\n want: %#v", cfg, want)
		}
	})
}
