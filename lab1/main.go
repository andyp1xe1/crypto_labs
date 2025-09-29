package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type CipherOp string

const (
	Encrypt CipherOp = "encrypt"
	Decrypt CipherOp = "decrypt"
)

// main is the entry point of the application.
// It presents a menu to the user to choose the desired cipher task.
func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- Caesar Cipher Menu ---")
		fmt.Println("1. Standard Caesar Cipher (Task 1.1)")
		fmt.Println("2. Caesar Cipher with Permutation Key (Task 1.2)")
		fmt.Println("3. Exit")
		fmt.Print("Select an option: ")

		choiceStr, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
		if err != nil {
			fmt.Println("Invalid input. Please enter a number (1-3).")
			continue
		}

		switch choice {
		case 1:
			runStandardCaesar(reader)
		case 2:
			runPermutationCaesar(reader)
		case 3:
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid option. Please choose 1, 2, or 3.")
		}
	}
}

// runStandardCaesar handles the logic for Task 1.1.
func runStandardCaesar(reader *bufio.Reader) {
	fmt.Println("\n--- Standard Caesar Cipher ---")
	op := getOperation(reader)
	key := getShiftKey(reader)
	text := getText(reader)

	processedText, err := processText(text, key, alphabet, op)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\nResult: %s\n", processedText)
}

// runPermutationCaesar handles the logic for Task 1.2.
func runPermutationCaesar(reader *bufio.Reader) {
	fmt.Println("\n--- Caesar Cipher with Permutation Key ---")
	op := getOperation(reader)
	key1 := getShiftKey(reader)
	key2 := getPermutationKey(reader)
	text := getText(reader)

	permutedAlphabet := generatePermutedAlphabet(key2)
	fmt.Printf("Generated Permuted Alphabet: %s\n", permutedAlphabet)

	processedText, err := processText(text, key1, permutedAlphabet, op)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\nResult: %s\n", processedText)
}

// processText is the core cipher function that encrypts or decrypts text.
// It uses a provided alphabet and shift key.
// The operation is determined by the `op` parameter ("encrypt" or "decrypt").
func processText(inputText string, shiftKey int, currentAlphabet string, op CipherOp) (string, error) {
	charToIndex := make(map[rune]int)
	for i, char := range currentAlphabet {
		charToIndex[char] = i
	}

	// Sanitize and prepare the input text according to the requirements.
	sanitizedText := sanitizeText(inputText)
	var result strings.Builder
	alphabetSize := len(currentAlphabet)

	for _, char := range sanitizedText {
		index, ok := charToIndex[char]
		if !ok {
			// err := fmt.Errorf("character '%c' not in the current alphabet, must be in [a-zA-Z]* space", char)
			// return "", err
			continue
		}

		var newIndex int
		switch op {
		case Encrypt:
			newIndex = (index + shiftKey) % alphabetSize
		case Decrypt:
			newIndex = (index - shiftKey + alphabetSize) % alphabetSize
		}

		// Find the character at the new index in the current alphabet string.
		// We convert the string to a rune slice to handle UTF-8 characters correctly,
		// although our alphabet is simple ASCII.
		result.WriteRune([]rune(currentAlphabet)[newIndex])
	}

	return result.String(), nil
}

// generatePermutedAlphabet creates a new alphabet order based on a keyword.
// Duplicates in the keyword are removed, and the remaining standard alphabet
// letters are appended in their natural order.
func generatePermutedAlphabet(keyword string) string {
	var builder strings.Builder
	seen := make(map[rune]bool)

	// Add unique characters from the keyword first
	for _, char := range strings.ToUpper(keyword) {
		if !seen[char] {
			builder.WriteRune(char)
			seen[char] = true
		}
	}

	// Add the rest of the alphabet
	for _, char := range alphabet {
		if !seen[char] {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

// sanitizeText converts text to uppercase and removes any non-letter characters.
func sanitizeText(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if unicode.IsLetter(char) {
			builder.WriteRune(unicode.ToUpper(char))
		}
	}
	return builder.String()
}

// --- User Input Helper Functions ---

// getOperation prompts the user to choose between encryption and decryption.
func getOperation(reader *bufio.Reader) CipherOp {
	for {
		fmt.Print("Enter operation (encrypt/decrypt): ")
		op, _ := reader.ReadString('\n')
		op = strings.ToLower(strings.TrimSpace(op))
		if strings.Compare(op, "encrypt") == 0 || strings.Compare(op, "decrypt") == 0 {
			return CipherOp(op)
		}
		fmt.Println("Invalid operation. Please enter 'encrypt' or 'decrypt'.")
	}
}

// getShiftKey prompts for the integer shift key (k1) and validates it.
func getShiftKey(reader *bufio.Reader) int {
	for {
		fmt.Print("Enter the shift key (an integer between 1 and 25): ")
		keyStr, _ := reader.ReadString('\n')
		key, err := strconv.Atoi(strings.TrimSpace(keyStr))
		if err == nil && key >= 1 && key <= 25 {
			return key
		}
		fmt.Println("Invalid key. It must be an integer between 1 and 25.")
	}
}

// getPermutationKey prompts for the permutation keyword (k2) and validates it.
func getPermutationKey(reader *bufio.Reader) string {
	for {
		fmt.Print("Enter the permutation keyword (at least 7 letters long, no numbers/symbols): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		// Sanitize to count only valid characters for length check
		cleanKey := sanitizeText(key)

		if len(cleanKey) < 7 {
			fmt.Println("Invalid keyword. It must contain at least 7 letters.")
			continue
		}
		if len(cleanKey) != len(key) {
			fmt.Println("Invalid keyword. It must contain only letters ('A'-'Z', 'a'-'z').")
			continue
		}
		return key
	}
}

// getText prompts for the message or cryptogram.
func getText(reader *bufio.Reader) string {
	for {
		fmt.Print("Enter the text to process: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text != "" {
			return text
		}
		fmt.Println("Input cannot be empty.")
	}
}
