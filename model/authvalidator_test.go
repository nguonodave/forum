package model

import (
	"errors"
	"testing"

	"forum/xerrors"
)

func TestIsValidPassword(t *testing.T) {
	password := "Test@1234"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if !IsValidPassword(password, hashedPassword) {
		t.Error("Expected password to be valid but got invalid")
	}

	if IsValidPassword("wrongPassword", hashedPassword) {
		t.Error("Expected password to be invalid but got valid")
	}
}

func TestHashPassword(t *testing.T) {
	password := "Test@1234"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Unexpected error hashing password: %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected hashed password to be non-empty")
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		password  string
		expectErr error
	}{
		{"Short1!", xerrors.ErrPasswordTooShort},
		{"nouppercase1!", errors.New("password must contain at least one uppercase letter")},
		{"NOLOWERCASE1!", errors.New("password must contain at least one lowercase letter")},
		{"NoNumber!", errors.New("password must contain at least one number")},
		{"NoSpecial1", errors.New("password must contain at least one special character")},
		{"Valid1!", nil},
	}

	for _, tc := range cases {
		err := ValidatePassword(tc.password)
		if (err != nil) != (tc.expectErr != nil) {
			t.Errorf("For password %q, expected error %v, got %v", tc.password, tc.expectErr, err)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		email     string
		expectErr bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"user@domain", true},
		{"user@.com", true},
	}

	for _, tc := range cases {
		err := ValidateEmail(tc.email)
		if (err != nil) != tc.expectErr {
			t.Errorf("For email %q, expected error %v, got %v", tc.email, tc.expectErr, err)
		}
	}
}
