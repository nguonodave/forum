package controller

import "testing"

// import (
// 	"database/sql"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"forum/model"
// )

// // setupTestDB creates an in-memory SQLite database for testing
// func setupTestDB(t *testing.T) *model.Database {
// 	db, err := sql.Open("sqlite3", ":memory:")
// 	if err != nil {
// 		t.Fatalf("Failed to open test database: %v", err)
// 	}

// 	// Create necessary tables
// 	_, err = db.Exec(`
// 		CREATE TABLE users (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			email TEXT UNIQUE,
// 			password TEXT,
// 			username TEXT
// 		);
// 		CREATE TABLE sessions (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			user_id INTEGER,
// 			token TEXT,
// 			expires_at DATETIME,
// 			FOREIGN KEY(user_id) REFERENCES users(id)
// 		);
// 	`)
// 	if err != nil {
// 		t.Fatalf("Failed to create test tables: %v", err)
// 	}

// 	return &model.Database{Db: db}
// }

// func TestGenerateSessionToken(t *testing.T) {
// 	token1 := generateSessionToken()
// 	token2 := generateSessionToken()

// 	if token1 == "" {
// 		t.Error("Generated token should not be empty")
// 	}
// 	if token1 == token2 {
// 		t.Error("Generated tokens should be unique")
// 	}
// 	if len(token1) != 36 {
// 		t.Errorf("Expected token length of 36, got %d", len(token1))
// 	}
// }

// func TestHandleRegister(t *testing.T) {
// 	db := setupTestDB(t)
// 	defer db.Db.Close()

// 	tests := []struct {
// 		name     string
// 		username string
// 		email    string
// 		password string
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "Valid registration",
// 			username: "testuser",
// 			email:    "test@example.com",
// 			password: "Password123!",
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "Invalid email",
// 			username: "testuser",
// 			email:    "invalid-email",
// 			password: "Password123!",
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid password - too short",
// 			username: "testuser",
// 			email:    "test2@example.com",
// 			password: "Pw1!",
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid password - no uppercase",
// 			username: "testuser",
// 			email:    "test3@example.com",
// 			password: "password123!",
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid password - no lowercase",
// 			username: "testuser",
// 			email:    "test4@example.com",
// 			password: "PASSWORD123!",
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid password - no number",
// 			username: "testuser",
// 			email:    "test5@example.com",
// 			password: "Password!!",
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid password - no special char",
// 			username: "testuser",
// 			email:    "test6@example.com",
// 			password: "Password123",
// 			wantErr:  true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := HandleRegister(db, tt.username, tt.email, tt.password)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("HandleRegister() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			if err == nil {
// 				// Verify user was stored in database
// 				var count int
// 				err := db.Db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", tt.email).Scan(&count)
// 				if err != nil {
// 					t.Errorf("Failed to verify user creation: %v", err)
// 				}
// 				if count != 1 {
// 					t.Error("User not stored in database")
// 				}
// 			}
// 		})
// 	}

// 	// Test duplicate email
// 	t.Run("Duplicate email", func(t *testing.T) {
// 		// First registration should succeed
// 		err := HandleRegister(db, "testuser", "duplicate@example.com", "Password123!")
// 		if err != nil {
// 			t.Fatalf("Failed to create initial user: %v", err)
// 		}

// 		// Second registration with same email should fail
// 		err = HandleRegister(db, "testuser2", "duplicate@example.com", "Password123!")
// 		if err == nil {
// 			t.Error("Expected error for duplicate email")
// 		}
// 	})
// }

// func TestHandleLogin(t *testing.T) {
// 	db := setupTestDB(t)
// 	defer db.Db.Close()

// 	// Create a test user
// 	err := HandleRegister(db, "testuser", "test@example.com", "Password123!")
// 	if err != nil {
// 		t.Fatalf("Failed to create test user: %v", err)
// 	}

// 	tests := []struct {
// 		name     string
// 		email    string
// 		password string
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "Valid login",
// 			email:    "test@example.com",
// 			password: "Password123!",
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "Invalid email",
// 			email:    "nonexistent@example.com",
// 			password: "Password123!",
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid password",
// 			email:    "test@example.com",
// 			password: "WrongPassword123!",
// 			wantErr:  true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			token, expiresAt, err := HandleLogin(db, tt.email, tt.password)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("HandleLogin() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			if !tt.wantErr {
// 				if token == "" {
// 					t.Error("Expected non-empty token")
// 				}
// 				if expiresAt.Before(time.Now()) {
// 					t.Error("Expected future expiration time")
// 				}

// 				// Verify session was stored
// 				var count int
// 				err := db.Db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = ?", token).Scan(&count)
// 				if err != nil {
// 					t.Errorf("Failed to verify session creation: %v", err)
// 				}
// 				if count != 1 {
// 					t.Error("Session not stored in database")
// 				}
// 			}
// 		})
// 	}
// }

// func TestValidateSession(t *testing.T) {
// 	db := setupTestDB(t)
// 	defer db.Db.Close()

// 	// Create a test user and get a valid session
// 	err := HandleRegister(db, "testuser", "test@example.com", "Password123!")
// 	if err != nil {
// 		t.Fatalf("Failed to create test user: %v", err)
// 	}

// 	token, _, err := HandleLogin(db, "test@example.com", "Password123!")
// 	if err != nil {
// 		t.Fatalf("Failed to create test session: %v", err)
// 	}

// 	// Create an expired session
// 	var userID int
// 	err = db.Db.QueryRow("SELECT id FROM users WHERE email = ?", "test@example.com").Scan(&userID)
// 	if err != nil {
// 		t.Fatalf("Failed to get user ID: %v", err)
// 	}

// 	expiredToken := generateSessionToken()
// 	_, err = db.Db.Exec(
// 		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
// 		userID,
// 		expiredToken,
// 		time.Now().Add(-24*time.Hour),
// 	)
// 	if err != nil {
// 		t.Fatalf("Failed to create expired session: %v", err)
// 	}

// 	tests := []struct {
// 		name           string
// 		cookie         *http.Cookie
// 		expectedStatus int
// 		checkUserID    bool
// 	}{
// 		{
// 			name: "Valid session",
// 			cookie: &http.Cookie{
// 				Name:  "session",
// 				Value: token,
// 			},
// 			expectedStatus: http.StatusOK,
// 			checkUserID:    true,
// 		},
// 		{
// 			name:           "No cookie",
// 			cookie:         nil,
// 			expectedStatus: http.StatusUnauthorized,
// 			checkUserID:    false,
// 		},
// 		{
// 			name: "Invalid session token",
// 			cookie: &http.Cookie{
// 				Name:  "session",
// 				Value: "invalid-token",
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			checkUserID:    false,
// 		},
// 		{
// 			name: "Expired session",
// 			cookie: &http.Cookie{
// 				Name:  "session",
// 				Value: expiredToken,
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			checkUserID:    false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				if tt.checkUserID {
// 					userID := r.Context().Value("userID")
// 					if userID == nil {
// 						t.Error("Expected userID in context")
// 					}
// 				}
// 				w.WriteHeader(http.StatusOK)
// 			})

// 			handler := ValidateSession(db.Db, nextHandler)
// 			req := httptest.NewRequest("GET", "/", nil)
// 			if tt.cookie != nil {
// 				req.AddCookie(tt.cookie)
// 			}
// 			rr := httptest.NewRecorder()

// 			handler.ServeHTTP(rr, req)

// 			if rr.Code != tt.expectedStatus {
// 				t.Errorf("Handler returned wrong status code: got %v want %v",
// 					rr.Code, tt.expectedStatus)
// 			}
// 		})
// 	}
// }

func TestGenerateSession(t *testing.T) {

}
