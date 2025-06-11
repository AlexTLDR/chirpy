package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	
	if hash == password {
		t.Fatal("HashPassword returned the original password instead of a hash")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectToken string
		expectError bool
	}{
		{
			name:        "valid bearer token",
			authHeader:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError: false,
		},
		{
			name:        "valid bearer token with extra spaces",
			authHeader:  "Bearer  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9  ",
			expectToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError: false,
		},
		{
			name:        "case insensitive bearer",
			authHeader:  "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError: false,
		},
		{
			name:        "mixed case bearer",
			authHeader:  "BeArEr eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError: false,
		},
		{
			name:        "empty authorization header",
			authHeader:  "",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "wrong format - no bearer prefix",
			authHeader:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "wrong format - basic auth",
			authHeader:  "Basic dXNlcjpwYXNzd29yZA==",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "bearer with empty token",
			authHeader:  "Bearer ",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "bearer with only spaces",
			authHeader:  "Bearer    ",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "too many parts",
			authHeader:  "Bearer token extra part",
			expectToken: "token extra part",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			token, err := GetBearerToken(headers)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if token != tt.expectToken {
					t.Errorf("Expected token %q but got %q", tt.expectToken, token)
				}
			}
		})
	}
}

func TestGetBearerTokenWithMissingHeader(t *testing.T) {
	headers := make(http.Header)
	// Don't set any Authorization header

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("GetBearerToken should have failed with missing Authorization header")
	}

	expectedError := "authorization header not found"
	if err.Error() != expectedError {
		t.Errorf("Expected error %q but got %q", expectedError, err.Error())
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("MakeRefreshToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("MakeRefreshToken returned empty token")
	}

	// Token should be 64 characters (32 bytes hex encoded)
	if len(token) != 64 {
		t.Fatalf("Expected token length 64, got %d", len(token))
	}

	// Token should be valid hex
	for _, char := range token {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			t.Fatalf("Token contains invalid hex character: %c", char)
		}
	}
}

func TestMakeRefreshTokenUniqueness(t *testing.T) {
	tokens := make(map[string]bool)
	numTokens := 100

	for i := 0; i < numTokens; i++ {
		token, err := MakeRefreshToken()
		if err != nil {
			t.Fatalf("MakeRefreshToken failed on iteration %d: %v", i, err)
		}

		if tokens[token] {
			t.Fatalf("MakeRefreshToken generated duplicate token: %s", token)
		}
		tokens[token] = true
	}

	if len(tokens) != numTokens {
		t.Fatalf("Expected %d unique tokens, got %d", numTokens, len(tokens))
	}
}

func TestMakeRefreshTokenMultipleCalls(t *testing.T) {
	token1, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("First MakeRefreshToken failed: %v", err)
	}

	token2, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("Second MakeRefreshToken failed: %v", err)
	}

	if token1 == token2 {
		t.Fatal("MakeRefreshToken returned the same token on consecutive calls")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	// Test correct password
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed for correct password: %v", err)
	}
	
	// Test incorrect password
	err = CheckPasswordHash(hash, "wrongpassword")
	if err == nil {
		t.Fatal("CheckPasswordHash should have failed for incorrect password")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour
	
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	
	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}
	
	// Test that different users get different tokens
	userID2 := uuid.New()
	token2, err := MakeJWT(userID2, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed for second user: %v", err)
	}
	
	if token == token2 {
		t.Fatal("MakeJWT returned the same token for different users")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour
	
	// Create a valid token
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	
	// Validate the token
	validatedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	
	if validatedUserID != userID {
		t.Fatalf("ValidateJWT returned wrong user ID. Expected %v, got %v", userID, validatedUserID)
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	expiresIn := time.Hour
	
	// Create a token with the correct secret
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	
	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed with wrong secret")
	}
}

func TestValidateJWTWithExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Millisecond * 1 // Very short expiration
	
	// Create a token that expires quickly
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	
	// Wait for the token to expire
	time.Sleep(time.Millisecond * 10)
	
	// Try to validate expired token
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed with expired token")
	}
}

func TestValidateJWTWithInvalidToken(t *testing.T) {
	secret := "test-secret"
	invalidToken := "invalid.token.here"
	
	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed with invalid token")
	}
}

func TestValidateJWTWithEmptyToken(t *testing.T) {
	secret := "test-secret"
	
	_, err := ValidateJWT("", secret)
	if err == nil {
		t.Fatal("ValidateJWT should have failed with empty token")
	}
}

func TestJWTRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		userID    uuid.UUID
		secret    string
		expiresIn time.Duration
	}{
		{
			name:      "standard case",
			userID:    uuid.New(),
			secret:    "my-secret-key",
			expiresIn: time.Hour,
		},
		{
			name:      "long expiration",
			userID:    uuid.New(),
			secret:    "another-secret",
			expiresIn: time.Hour * 24 * 7, // 1 week
		},
		{
			name:      "short expiration",
			userID:    uuid.New(),
			secret:    "short-secret",
			expiresIn: time.Minute,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create token
			token, err := MakeJWT(tt.userID, tt.secret, tt.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT failed: %v", err)
			}
			
			// Validate token
			validatedUserID, err := ValidateJWT(token, tt.secret)
			if err != nil {
				t.Fatalf("ValidateJWT failed: %v", err)
			}
			
			if validatedUserID != tt.userID {
				t.Fatalf("User ID mismatch. Expected %v, got %v", tt.userID, validatedUserID)
			}
		})
	}
}