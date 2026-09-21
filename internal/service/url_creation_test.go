package service

import (
	"testing"
)

func TestCreateURL(t *testing.T) {
	shortCode, err := CreateURL()

	if err != nil {
		t.Fatalf("CreateURL() returned an error: %v", err)
	}

	if len(shortCode) != 10 {
		t.Fatalf(
			"CreateURL() returned %d characters, want 10",
			len(shortCode),
		)
		
	}
}

func TestCreateURL_ValidCharacters(t *testing.T) {
	shortCode, err := CreateURL()

	if err != nil {
		t.Fatalf("CreateURL() returned an error: %v", err)
	}

	for _, char := range shortCode {
		if !containsRune(charset, char) {
			t.Errorf(
				"CreateURL() returned invalid character %q",
				char,
			)
		}
	}
}

func TestCreateURL_Randomness(t *testing.T) {
	const iterations = 100

	codes := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		shortCode, err := CreateURL()

		if err != nil {
			t.Fatalf("CreateURL() returned an error: %v", err)
		}

		if codes[shortCode] {
			t.Errorf(
				"CreateURL() generated duplicate code: %s",
				shortCode,
			)
		}

		codes[shortCode] = true
	}
}

func containsRune(s string, target rune) bool {
	for _, char := range s {
		if char == target {
			return true
		}
	}

	return false
}