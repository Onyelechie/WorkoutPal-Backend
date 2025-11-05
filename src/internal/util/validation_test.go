package util

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty", "", "email is required"},
		{"invalid format", "not-an-email", "invalid email format"},
		{"missing domain", "test@", "invalid email format"},
		{"valid email", "user@example.com", ""},
		{"valid with dots", "first.last@domain.co.uk", ""},
		{"valid with plus", "user+alias@gmail.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr) {
				t.Fatalf("expected '%s', got '%v'", tt.wantErr, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty", "", "username is required"},
		{"too short", "ab", "username must be at least 3 characters"},
		{"too long", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", "username must be less than 50 characters"},
		{"valid", "valid_user123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.input)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr) {
				t.Fatalf("expected '%s', got '%v'", tt.wantErr, err)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty", "", "name is required"},
		{"only spaces", "   ", "name must be at least 2 characters"},
		{"too short", "a", "name must be at least 2 characters"},
		{"valid", "John", ""},
		{"valid with spaces", " Jane ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr) {
				t.Fatalf("expected '%s', got '%v'", tt.wantErr, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty", "", "password is required"},
		{"too short", "12345", "password must be at least 6 characters"},
		{"valid", "strongPass123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr) {
				t.Fatalf("expected '%s', got '%v'", tt.wantErr, err)
			}
		})
	}
}
