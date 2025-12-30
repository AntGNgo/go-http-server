package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body string `json:"body"`
	}

	type resParams struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		response := resParams{
			Error: "Chirp is too long",
			Valid: false,
		}
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(dat)
	} else {
		response := resParams{
			Valid: true,
		}
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(dat)
	}

}
