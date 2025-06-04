package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHandlerReadiness(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerReadiness)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandlerMetrics(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}
	cfg.fileserverHits.Store(5)

	req, err := http.NewRequest("GET", "/admin/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cfg.handlerMetrics)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("5 times")) {
		t.Errorf("handler should contain hit count in response body")
	}
}

func TestHandlerMetricsWrongMethod(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	req, err := http.NewRequest("POST", "/admin/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cfg.handlerMetrics)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestMiddlewareMetricsInc(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with middleware
	wrappedHandler := cfg.middlewareMetricsInc(testHandler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if cfg.fileserverHits.Load() != 1 {
		t.Errorf("middleware should increment hit counter: got %v want %v", cfg.fileserverHits.Load(), 1)
	}
}

func TestCleanProfanity(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"This is a normal message", "This is a normal message"},
		{"This contains kerfuffle word", "This contains **** word"},
		{"Multiple sharbert and fornax words", "Multiple **** and **** words"},
		{"KERFUFFLE in uppercase", "**** in uppercase"}, // Case insensitive matching
		{"", ""},
	}

	for _, tt := range tests {
		result := cleanProfanity(tt.input)
		if result != tt.expected {
			t.Errorf("cleanProfanity(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestUserStructJSON(t *testing.T) {
	userID := uuid.New()
	now := time.Now()
	
	user := User{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
		Email:     "test@example.com",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	if unmarshaled.ID != user.ID {
		t.Errorf("ID mismatch: got %v want %v", unmarshaled.ID, user.ID)
	}

	if unmarshaled.Email != user.Email {
		t.Errorf("Email mismatch: got %v want %v", unmarshaled.Email, user.Email)
	}
}

func TestChirpStructJSON(t *testing.T) {
	chirpID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	
	chirp := Chirp{
		ID:        chirpID,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      "Test chirp message",
		UserID:    userID,
	}

	jsonData, err := json.Marshal(chirp)
	if err != nil {
		t.Fatalf("Failed to marshal chirp: %v", err)
	}

	var unmarshaled Chirp
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal chirp: %v", err)
	}

	if unmarshaled.ID != chirp.ID {
		t.Errorf("ID mismatch: got %v want %v", unmarshaled.ID, chirp.ID)
	}

	if unmarshaled.Body != chirp.Body {
		t.Errorf("Body mismatch: got %v want %v", unmarshaled.Body, chirp.Body)
	}

	if unmarshaled.UserID != chirp.UserID {
		t.Errorf("UserID mismatch: got %v want %v", unmarshaled.UserID, chirp.UserID)
	}
}

func TestErrorResponseJSON(t *testing.T) {
	errResp := ErrorResponse{
		Error: "Test error message",
	}

	jsonData, err := json.Marshal(errResp)
	if err != nil {
		t.Fatalf("Failed to marshal error response: %v", err)
	}

	var unmarshaled ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if unmarshaled.Error != errResp.Error {
		t.Errorf("Error mismatch: got %v want %v", unmarshaled.Error, errResp.Error)
	}
}