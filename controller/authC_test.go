package controller

import (
	"database/sql"
	"forum/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables with strict constraints
	_, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            username TEXT NOT NULL,
            CONSTRAINT email_not_empty CHECK(email != ''),
            CONSTRAINT username_not_empty CHECK(username != '')
        );
        CREATE TABLE sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            token TEXT UNIQUE NOT NULL,
            expires_at DATETIME NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
    `)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	SetDB(db)
	return db
}

func TestRegisterUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name       string
		email      string
		password   string
		username   string
		setupFunc  func(*testing.T, *sql.DB)
		wantErr    bool
		errMessage string
	}{
		{
			name:     "Valid registration",
			email:    "test@example.com",
			password: "Password123!",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:       "Invalid email",
			email:      "invalid-email",
			password:   "Password123!",
			username:   "testuser",
			wantErr:    true,
			errMessage: "invalid email format",
		},
		{
			name:       "Invalid password",
			email:      "test@example.com",
			password:   "weak",
			username:   "testuser",
			wantErr:    true,
			errMessage: "password length too short, minimum of 8 characters required",
		},
		{
			name:     "Duplicate email",
			email:    "existing@example.com",
			password: "Password123!",
			username: "testuser2",
			setupFunc: func(t *testing.T, db *sql.DB) {
				_, err := db.Exec(
					"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
					"existing@example.com", "hashedpass", "existinguser",
				)
				if err != nil {
					t.Fatalf("Failed to setup duplicate user: %v", err)
				}
			},
			wantErr:    true,
			errMessage: "email is already taken",
		},
		{
			name:       "Database error",
			email:      "test2@example.com",
			password:   "Password123!",
			username:   "", // Violates NOT NULL constraint
			wantErr:    true,
			errMessage: "failed to create user",
		},
		{
			name:       "Password hashing error",
			email:      "test3@example.com",
			password:   "", // This should trigger internal error during hashing
			username:   "testuser3",
			wantErr:    true,
			errMessage: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset database for each test
			db.Exec("DELETE FROM users")
			db.Exec("DELETE FROM sessions")

			if tt.setupFunc != nil {
				tt.setupFunc(t, db)
			}

			err := RegisterUser(tt.email, tt.password, tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMessage != "" && err.Error() != tt.errMessage {
				t.Errorf("RegisterUser() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}

func TestVerifyLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	hashedPassword, err := model.HashPassword("Password123!")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Create a test user
	result, err := db.Exec(
		"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
		"test@example.com",
		hashedPassword,
		"testuser",
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	tests := []struct {
		name      string
		email     string
		password  string
		setupFunc func(*testing.T, *sql.DB)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "Valid login",
			email:    "test@example.com",
			password: "Password123!",
			wantErr:  false,
		},
		{
			name:     "Invalid email",
			email:    "nonexistent@example.com",
			password: "Password123!",
			wantErr:  true,
			errMsg:   "invalid credentials",
		},
		{
			name:     "Invalid password",
			email:    "test@example.com",
			password: "WrongPassword123!",
			wantErr:  true,
			errMsg:   "invalid credentials",
		},
		{
			name:     "Database error on session creation",
			email:    "test@example.com",
			password: "Password123!",
			setupFunc: func(t *testing.T, db *sql.DB) {
				// Drop sessions table to force database error
				db.Exec("DROP TABLE sessions")
			},
			wantErr: true,
			errMsg:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset sessions for each test
			db.Exec("DELETE FROM sessions")

			if tt.setupFunc != nil {
				tt.setupFunc(t, db)
			}

			token, expires, err := VerifyLogin(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("VerifyLogin() error message = %v, want %v", err.Error(), tt.errMsg)
			}
			if !tt.wantErr {
				if token == "" {
					t.Error("VerifyLogin() token is empty")
				}
				if expires == "" {
					t.Error("VerifyLogin() expires is empty")
				}

				// If we expect success, verify session was stored
				var storedToken string
				err = db.QueryRow("SELECT token FROM sessions WHERE token = ? AND user_id = ?", token, userID).Scan(&storedToken)
				if err != nil {
					t.Errorf("Session not stored in database: %v", err)
				}
			}
		})
	}
}

func TestValidateSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a test user and sessions
	result, err := db.Exec(
		"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
		"test@example.com", "hashedpass", "testuser",
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	userID, _ := result.LastInsertId()
	validToken := generateSessionToken()
	expiredToken := generateSessionToken()

	// Valid session
	_, err = db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		userID, validToken, time.Now().Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("Failed to create valid session: %v", err)
	}

	// Expired session
	_, err = db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		userID, expiredToken, time.Now().Add(-1*time.Hour),
	)
	if err != nil {
		t.Fatalf("Failed to create expired session: %v", err)
	}

	tests := []struct {
		name       string
		cookie     *http.Cookie
		setupFunc  func(*testing.T, *sql.DB)
		wantStatus int
	}{
		{
			name:       "Valid session",
			cookie:     &http.Cookie{Name: "session", Value: validToken},
			wantStatus: http.StatusOK,
		},
		{
			name:       "No cookie",
			cookie:     nil,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid session token",
			cookie:     &http.Cookie{Name: "session", Value: "invalid-token"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Expired session",
			cookie:     &http.Cookie{Name: "session", Value: expiredToken},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "Database error",
			cookie: &http.Cookie{Name: "session", Value: validToken},
			setupFunc: func(t *testing.T, db *sql.DB) {
				// Drop sessions table to force database error
				db.Exec("DROP TABLE sessions")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc(t, db)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()
			handler := ValidateSession(db, nextHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("ValidateSession() status = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestGenerateSessionToken(t *testing.T) {
	tokens := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		token := generateSessionToken()
		if token == "" {
			t.Error("generateSessionToken() generated empty token")
		}
		if tokens[token] {
			t.Error("generateSessionToken() generated duplicate token")
		}
		tokens[token] = true
	}
}

func TestSetDB(t *testing.T) {
	// Test with valid database
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer testDB.Close()

	SetDB(testDB)
	if db != testDB {
		t.Error("SetDB() failed to set the database connection")
	}

	// Test with nil database
	SetDB(nil)
	if db != nil {
		t.Error("SetDB() failed to set nil database")
	}
}
