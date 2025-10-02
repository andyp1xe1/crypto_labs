package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

//go:embed message.txt
var message string

var englishFreq = map[rune]float64{
	'A': 8.17, 'B': 1.49, 'C': 2.78, 'D': 4.25, 'E': 12.70, 'F': 2.23,
	'G': 2.01, 'H': 6.09, 'I': 6.97, 'J': 0.15, 'K': 0.77, 'L': 4.03,
	'M': 2.41, 'N': 6.75, 'O': 7.51, 'P': 1.93, 'Q': 0.09, 'R': 5.99,
	'S': 6.33, 'T': 9.06, 'U': 2.76, 'V': 0.98, 'W': 2.36, 'X': 0.15,
	'Y': 1.97, 'Z': 0.07,
}

type freqEntry struct {
	char  rune
	count int
	freq  float64
}

func findCommonPatterns() {
	cleanMessage := regexp.MustCompile(`[^A-Z]`).ReplaceAllString(strings.ToUpper(message), "")

	// Find double letters
	doubles := make(map[string]int)
	for i := 0; i < len(cleanMessage)-1; i++ {
		pattern := cleanMessage[i : i+2]
		if pattern[0] == pattern[1] {
			doubles[pattern]++
		}
	}

	fmt.Println("Double letters found:")
	for pattern, count := range doubles {
		fmt.Printf("%s: %d times\n", pattern, count)
	}

	// Find common digraphs
	digraphs := make(map[string]int)
	for i := 0; i < len(cleanMessage)-1; i++ {
		pattern := cleanMessage[i : i+2]
		digraphs[pattern]++
	}

	// Sort digraphs by frequency
	type patternEntry struct {
		pattern string
		count   int
	}

	var sortedDigraphs []patternEntry
	for pattern, count := range digraphs {
		if count > 1 { // Only show patterns that appear more than once
			sortedDigraphs = append(sortedDigraphs, patternEntry{pattern, count})
		}
	}

	sort.Slice(sortedDigraphs, func(i, j int) bool {
		return sortedDigraphs[i].count > sortedDigraphs[j].count
	})

	fmt.Println("\nMost common digraphs:")
	for i, entry := range sortedDigraphs {
		if i >= 10 { // Show top 10
			break
		}
		fmt.Printf("%s: %d times\n", entry.pattern, entry.count)
	}

	// Find common trigraphs
	trigraphs := make(map[string]int)
	for i := 0; i < len(cleanMessage)-2; i++ {
		pattern := cleanMessage[i : i+3]
		trigraphs[pattern]++
	}

	var sortedTrigraphs []patternEntry
	for pattern, count := range trigraphs {
		if count > 1 {
			sortedTrigraphs = append(sortedTrigraphs, patternEntry{pattern, count})
		}
	}

	sort.Slice(sortedTrigraphs, func(i, j int) bool {
		return sortedTrigraphs[i].count > sortedTrigraphs[j].count
	})

	fmt.Println("\nMost common trigraphs:")
	for i, entry := range sortedTrigraphs {
		if i >= 10 {
			break
		}
		fmt.Printf("%s: %d times\n", entry.pattern, entry.count)
	}
}

func main() {
	fmt.Println("=== Frequency Analysis of Encrypted Message ===")

	// Count letter frequencies
	letterCounts := make(map[rune]int)
	totalLetters := 0

	for _, char := range strings.ToUpper(message) {
		if char >= 'A' && char <= 'Z' {
			letterCounts[char]++
			totalLetters++
		}
	}

	// Convert to frequency entries and sort
	var frequencies []freqEntry
	for char, count := range letterCounts {
		freq := (float64(count) / float64(totalLetters)) * 100
		frequencies = append(frequencies, freqEntry{char, count, freq})
	}

	// Sort by frequency (descending)
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].freq > frequencies[j].freq
	})

	fmt.Printf("Total letters: %d\n\n", totalLetters)
	fmt.Println("Letter Frequency Analysis:")
	fmt.Printf("%-6s %-7s %-10s\n", "Letter", "Count", "Frequency%")
	fmt.Println(strings.Repeat("-", 25))

	for _, entry := range frequencies {
		fmt.Printf("%-6c %-7d %.2f%%\n", entry.char, entry.count, entry.freq)
	}

	fmt.Println("\n=== Comparison with English Frequencies ===")
	fmt.Printf("%-6s %-10s %-12s %-10s\n", "Letter", "Message%", "English%", "Difference")
	fmt.Println(strings.Repeat("-", 40))

	for _, entry := range frequencies {
		englishFreq := englishFreq[entry.char]
		diff := entry.freq - englishFreq
		fmt.Printf("%-6c %-10.2f %-12.2f %+.2f\n", entry.char, entry.freq, englishFreq, diff)
	}

	// Find common patterns
	fmt.Println("\n=== Common Letter Patterns ===")
	findCommonPatterns()
}
