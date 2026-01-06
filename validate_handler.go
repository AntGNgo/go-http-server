package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func createChirp(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type resParams struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParams{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

}

func validateChirpHandler(msg string) {
	const maxChripLength = 140
	if len(msg) > maxChripLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	} else {
		// fix
		bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
		cleaned := []string{}
		bodySplit := strings.Split(msg, " ")
		for _, word := range bodySplit {
			if slices.Contains(bannedWords, strings.ToLower(word)) {
				cleaned = append(cleaned, "****")
			} else {
				cleaned = append(cleaned, word)
			}
		}
		cleanedJoined := strings.Join(cleaned, " ")
		respondWithJSON(w, http.StatusOK, resParams{
			CleanedBody: cleanedJoined,
		})
	}

}
