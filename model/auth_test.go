package model

import (
	"errors"
	"testing"

	"forum/xerrors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIsValidPassword(t *testing.T) {
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), Cost)

	t.Run("ValidPassword", func(t *testing.T) {
		assert.True(t, IsValidPassword(password, string(hashedPassword)))
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		assert.False(t, IsValidPassword("wrongpassword", string(hashedPassword)))
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("ValidPassword", func(t *testing.T) {
		password := "ValidPass123!"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("PasswordTooShort", func(t *testing.T) {
		password := "short"
		hashedPassword, err := HashPassword(password)
		assert.Equal(t, xerrors.ErrPasswordTooShort, err)
		assert.Empty(t, hashedPassword)
	})

	t.Run("BcryptError", func(t *testing.T) {
		// Force bcrypt to fail by setting an invalid cost
		originalCost := Cost
		Cost = 100 // Invalid cost
		defer func() { Cost = originalCost }()

		password := "ValidPass123!"
		hashedPassword, err := HashPassword(password)
		assert.Error(t, err)
		assert.Empty(t, hashedPassword)
	})
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected error
	}{
		{"ValidPassword", "ValidPass123!", nil},
		{"TooShort", "short", xerrors.ErrPasswordTooShort},
		{"NoUppercase", "nopass123!", errors.New("password must contain at least one uppercase letter")},
		{"NoLowercase", "NOPASS123!", errors.New("password must contain at least one lowercase letter")},
		{"NoNumber", "NoPassWord!", errors.New("password must contain at least one number")},
		{"NoSpecialChar", "NoPass1234", errors.New("password must contain at least one special character")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"ValidEmail", "test@example.com", true},
		{"InvalidEmail", "invalid-email", false},
		{"InvalidFormat", "test@.com", false},
		{"NoDomain", "test@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestIsEmailTaken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("EmailTaken", func(t *testing.T) {
		email := "taken@example.com"
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\?\\)").
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		assert.True(t, IsEmailTaken(db, email))
	})

	t.Run("EmailNotTaken", func(t *testing.T) {
		email := "not_taken@example.com"
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\?\\)").
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		assert.False(t, IsEmailTaken(db, email))
	})

	t.Run("DatabaseError", func(t *testing.T) {
		email := "error@example.com"
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\?\\)").
			WithArgs(email).
			WillReturnError(errors.New("database error"))

		assert.False(t, IsEmailTaken(db, email))
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
