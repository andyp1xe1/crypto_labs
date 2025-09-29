package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSanitizeText checks the text sanitization logic.
// It ensures that input strings are correctly converted to uppercase
// and that all non-alphabetic characters are removed.
func TestSanitizeText(t *testing.T) {
	// A map of test cases: input string -> expected output string
	testCases := map[string]string{
		"hello world":       "HELLOWORLD",
		"Hello, World!":     "HELLOWORLD",
		"123 ABC xyz 456":   "ABCXYZ",
		"NoChanges":         "NOCHANGES",
		"!@#$%^&*()_+":      "",
		"   leading space":  "LEADINGSPACE",
		"trailing space   ": "TRAILINGSPACE",
	}

	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual := sanitizeText(input)
			assert.Equal(t, expected, actual, "Sanitization did not produce the expected result.")
		})
	}
}

func processTextNoErr(t *testing.T, inputText string, shiftKey int, currentAlphabet string, op CipherOp) string {
	txt, err := processText(inputText, shiftKey, currentAlphabet, op)
	if err != nil {
		t.Log(err)
	}
	return txt
}

// TestGeneratePermutedAlphabet checks the logic for creating
// a new alphabet based on a permutation keyword.
func TestGeneratePermutedAlphabet(t *testing.T) {
	testCases := []struct {
		name     string
		keyword  string
		expected string
	}{
		{
			name:     "Standard Example from PDF",
			keyword:  "cryptography",
			expected: "CRYPTOGAHBDEFIJKLMNQSUVWXZ",
		},
		{
			name:     "Keyword with repeated letters",
			keyword:  "hello",
			expected: "HELOABCDFGIJKMNPQRSTUVWXYZ",
		},
		{
			name:     "Full alphabet as keyword",
			keyword:  "ZYXWVUTSRQPONMLKJIHGFEDCBA",
			expected: "ZYXWVUTSRQPONMLKJIHGFEDCBA",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := generatePermutedAlphabet(tc.keyword)
			assert.Equal(t, tc.expected, actual, "Generated alphabet does not match the expected one.")
		})
	}
}

// TestProcessText covers the main encryption and decryption logic for both
// the standard and permuted Caesar ciphers. It uses table-driven tests
// to check multiple scenarios efficiently.
func TestProcessText(t *testing.T) {
	// Define alphabets for testing
	permutedAlphabet := generatePermutedAlphabet("cryptography") // CRYPTOGAHBDEFIJKMLNQSUVWXZ

	testCases := []struct {
		name      string
		text      string
		shiftKey  int
		alphabet  string
		operation CipherOp
		expected  string
	}{
		// --- Standard Caesar Cipher Tests ---
		{
			name:      "Standard Encrypt - No wrap",
			text:      "HELLO",
			shiftKey:  3,
			alphabet:  alphabet,
			operation: "encrypt",
			expected:  "KHOOR",
		},
		{
			name:      "Standard Encrypt - With wrap",
			text:      "XYZ",
			shiftKey:  3,
			alphabet:  alphabet,
			operation: "encrypt",
			expected:  "ABC",
		},
		{
			name:      "Standard Decrypt - No wrap",
			text:      "KHOOR",
			shiftKey:  3,
			alphabet:  alphabet,
			operation: "decrypt",
			expected:  "HELLO",
		},
		{
			name:      "Standard Decrypt - With wrap",
			text:      "ABC",
			shiftKey:  3,
			alphabet:  alphabet,
			operation: "decrypt",
			expected:  "XYZ",
		},
		// --- Permutation Caesar Cipher Tests ---
		{
			name:      "Permutation Encrypt - Verified Logic",
			text:      "cezar", // will be sanitized to CEZAR
			shiftKey:  3,
			alphabet:  permutedAlphabet,
			operation: "encrypt",
			expected:  "PJYDT", // C(0)->P(3), E(11)->J(14), Z(25)->Y(2), A(7)->D(10), R(1)->T(4)
		},
		{
			name:      "Permutation Decrypt - Reverse of above",
			text:      "PJYDT",
			shiftKey:  3,
			alphabet:  permutedAlphabet,
			operation: "decrypt",
			expected:  "CEZAR",
		},
		// --- Full Cycle Identity Tests (Encrypt then Decrypt) ---
		{
			name:      "Full Cycle - Standard",
			text:      processTextNoErr(t, sanitizeText("Attack at Dawn!"), 17, alphabet, "encrypt"), // Encrypt first
			shiftKey:  17,
			alphabet:  alphabet,
			operation: "decrypt",
			expected:  "ATTACKATDAWN", // Should return to original sanitized text
		},
		{
			name:      "Full Cycle - Permutation",
			text:      processTextNoErr(t, sanitizeText("Secret Message"), 8, permutedAlphabet, "encrypt"),
			shiftKey:  8,
			alphabet:  permutedAlphabet,
			operation: "decrypt",
			expected:  "SECRETMESSAGE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: processText expects already sanitized text.
			// The tests are structured to reflect this.
			sanitizedInput := sanitizeText(tc.text)
			actual := processTextNoErr(t, sanitizedInput, tc.shiftKey, tc.alphabet, tc.operation)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
