package database

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserModel(t *testing.T) {
	userID := uuid.New()
	now := time.Now().UTC()
	
	user := User{
		ID:             userID,
		CreatedAt:      now,
		UpdatedAt:      now,
		Email:          "test@example.com",
		HashedPassword: "hashed_password_123",
	}

	// Test that all fields are set correctly
	if user.ID != userID {
		t.Errorf("Expected ID %v, got %v", userID, user.ID)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", user.Email)
	}
	
	if user.HashedPassword != "hashed_password_123" {
		t.Errorf("Expected hashed password hashed_password_123, got %v", user.HashedPassword)
	}
}

func TestChirpModel(t *testing.T) {
	chirpID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()
	
	chirp := Chirp{
		ID:        chirpID,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      "This is a test chirp message",
		UserID:    userID,
	}

	// Test that all fields are set correctly
	if chirp.ID != chirpID {
		t.Errorf("Expected ID %v, got %v", chirpID, chirp.ID)
	}
	
	if chirp.Body != "This is a test chirp message" {
		t.Errorf("Expected body 'This is a test chirp message', got %v", chirp.Body)
	}
	
	if chirp.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, chirp.UserID)
	}
}

func TestUserModelJSONSerialization(t *testing.T) {
	userID := uuid.New()
	now := time.Now().UTC()
	
	user := User{
		ID:             userID,
		CreatedAt:      now,
		UpdatedAt:      now,
		Email:          "test@example.com",
		HashedPassword: "hashed_password_123",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal user from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaled.ID != user.ID {
		t.Errorf("ID mismatch after JSON round trip: got %v want %v", unmarshaled.ID, user.ID)
	}
	
	if unmarshaled.Email != user.Email {
		t.Errorf("Email mismatch after JSON round trip: got %v want %v", unmarshaled.Email, user.Email)
	}
	
	if unmarshaled.HashedPassword != user.HashedPassword {
		t.Errorf("HashedPassword mismatch after JSON round trip: got %v want %v", unmarshaled.HashedPassword, user.HashedPassword)
	}
}

func TestChirpModelJSONSerialization(t *testing.T) {
	chirpID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()
	
	chirp := Chirp{
		ID:        chirpID,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      "This is a test chirp message",
		UserID:    userID,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(chirp)
	if err != nil {
		t.Fatalf("Failed to marshal chirp to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled Chirp
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal chirp from JSON: %v", err)
	}

	// Verify all fields match
	if unmarshaled.ID != chirp.ID {
		t.Errorf("ID mismatch after JSON round trip: got %v want %v", unmarshaled.ID, chirp.ID)
	}
	
	if unmarshaled.Body != chirp.Body {
		t.Errorf("Body mismatch after JSON round trip: got %v want %v", unmarshaled.Body, chirp.Body)
	}
	
	if unmarshaled.UserID != chirp.UserID {
		t.Errorf("UserID mismatch after JSON round trip: got %v want %v", unmarshaled.UserID, chirp.UserID)
	}
}

func TestCreateUserParams(t *testing.T) {
	params := CreateUserParams{
		Email:          "test@example.com",
		HashedPassword: "hashed_password_123",
	}

	if params.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", params.Email)
	}
	
	if params.HashedPassword != "hashed_password_123" {
		t.Errorf("Expected hashed password hashed_password_123, got %v", params.HashedPassword)
	}
}

func TestCreateChirpParams(t *testing.T) {
	userID := uuid.New()
	params := CreateChirpParams{
		Body:   "Test chirp body",
		UserID: userID,
	}

	if params.Body != "Test chirp body" {
		t.Errorf("Expected body 'Test chirp body', got %v", params.Body)
	}
	
	if params.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, params.UserID)
	}
}

func TestModelFieldTypes(t *testing.T) {
	// Test that the model fields have the correct types
	var user User
	var chirp Chirp

	// These assignments should compile without error if types are correct
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Email = "string"
	user.HashedPassword = "string"

	chirp.ID = uuid.New()
	chirp.CreatedAt = time.Now()
	chirp.UpdatedAt = time.Now()
	chirp.Body = "string"
	chirp.UserID = uuid.New()

	// If we get here, all type assignments worked
	t.Log("All model field types are correct")
}