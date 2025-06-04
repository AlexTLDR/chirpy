package auth

import (
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