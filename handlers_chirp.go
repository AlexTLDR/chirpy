package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/AlexTLDR/chirpy/internal/auth"
	"github.com/AlexTLDR/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.handlerCreateChirp(w, r)
	case http.MethodGet:
		// Parse the path to see if we have an ID
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if len(pathParts) == 2 && pathParts[0] == "api" && pathParts[1] == "chirps" {
			// Get all chirps
			cfg.handlerGetChirps(w, r)
		} else if len(pathParts) == 3 && pathParts[0] == "api" && pathParts[1] == "chirps" {
			// Get specific chirp by ID
			chirpIDStr := pathParts[2]
			cfg.handlerGetChirpByID(w, r, chirpIDStr)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		// Parse the path to get chirp ID
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		
		if len(pathParts) == 3 && pathParts[0] == "api" && pathParts[1] == "chirps" {
			// Delete specific chirp by ID
			chirpIDStr := pathParts[2]
			cfg.handlerDeleteChirp(w, r, chirpIDStr)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	// Extract and validate JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err = decoder.Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong"})
		return
	}

	if reqBody.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Body is required"})
		return
	}

	if len(reqBody.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"})
		return
	}

	// Clean profane words
	cleanedBody := cleanProfanity(reqBody.Body)

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong"})
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check for author_id query parameter
	authorIDStr := r.URL.Query().Get("author_id")
	
	var dbChirps []database.Chirp
	var err error
	
	if authorIDStr != "" {
		// Parse the author ID
		authorID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid author ID"})
			return
		}
		
		// Get chirps by specific author
		dbChirps, err = cfg.dbQueries.GetChirpsByUserID(r.Context(), authorID)
	} else {
		// Get all chirps
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
	}
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong"})
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	// Check for sort query parameter
	sortParam := r.URL.Query().Get("sort")
	
	// Sort chirps based on the sort parameter (default is ascending)
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		// Default to ascending order (asc or no parameter)
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request, chirpIDStr string) {
	w.Header().Set("Content-Type", "application/json")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid chirp ID"})
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp not found"})
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request, chirpIDStr string) {
	w.Header().Set("Content-Type", "application/json")

	// Extract and validate JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Parse chirp ID
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid chirp ID"})
		return
	}

	// Get the chirp to check if it exists and if user owns it
	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp not found"})
		return
	}

	// Check if the user is the author of the chirp
	if dbChirp.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "You can only delete your own chirps"})
		return
	}

	// Delete the chirp
	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong"})
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func cleanProfanity(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(text)

	for i, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
