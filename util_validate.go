package main

import (
	"fmt"
	"slices"
	"strings"
)

const MAX_CHIRP_LENGTH = 140
var BAD_WORDS = []string{"kerfuffle", "sharbert", "fornax"}

func filterBadWords(chirp string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if slices.Contains(BAD_WORDS, strings.ToLower(word)) {
			words[i] = "****"
		} 

	}
	return strings.Join(words, " ")
}

func validateChirp(chirp string) (string, error) {
	if len(chirp) > MAX_CHIRP_LENGTH {
		return "", fmt.Errorf("Chirp is too long")
	}

	filteredBody := filterBadWords(chirp)
	return filteredBody, nil
}


