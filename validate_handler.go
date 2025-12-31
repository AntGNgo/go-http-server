package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body string `json:"body"`
	}

	type resParams struct {
		Error       string `json:"error"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	const maxChripLength = 140
	if len(params.Body) > maxChripLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	} else {
		// fix
		bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
		cleaned := []string{}
		bodySplit := strings.Split(params.Body, " ")
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
