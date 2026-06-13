package trisnake

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const highScoresFile = "HIGHSCORES.md"

// LoadHighScores reads and parses the HIGHSCORES.md file.
// It returns a sorted slice (highest score first) of valid entries.
// If the file is missing or contains no valid entries, it returns an empty slice (no crash).
func LoadHighScores() []HighScoreEntry {
	data, err := os.ReadFile(highScoresFile)
	if err != nil {
		// File missing or unreadable — return empty, no crash.
		return nil
	}
	return ParseHighScores(string(data))
}

// ParseHighScores parses the markdown table content into HighScoreEntry structs.
// Malformed lines are silently skipped.
func ParseHighScores(content string) []HighScoreEntry {
	var entries []HighScoreEntry
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines, header, and separator rows.
		if line == "" || strings.HasPrefix(line, "|Date") || strings.HasPrefix(line, "|--") || strings.HasPrefix(line, "|-") {
			continue
		}
		// Expected format: |date|score|speed|difficulty|
		entry, ok := parseLine(line)
		if ok {
			entries = append(entries, entry)
		}
	}
	// Sort by score descending.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	return entries
}

// parseLine parses a single markdown table row into a HighScoreEntry.
func parseLine(line string) (HighScoreEntry, bool) {
	parts := strings.Split(line, "|")
	// Filter out empty parts (leading/trailing pipes).
	var fields []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			fields = append(fields, p)
		}
	}
	if len(fields) < 4 {
		return HighScoreEntry{}, false
	}

	score, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil {
		return HighScoreEntry{}, false
	}
	speed, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
	if err != nil {
		return HighScoreEntry{}, false
	}

	return HighScoreEntry{
		Date:       strings.TrimSpace(fields[0]),
		Score:      score,
		Speed:      speed,
		Difficulty: strings.TrimSpace(fields[3]),
	}, true
}

// SaveHighScoreEntry appends a new score row to HIGHSCORES.md.
// If the file does not exist, it creates it with a proper header.
func SaveHighScoreEntry(score int, speed float64, difficulty string) error {
	// Ensure file exists with header.
	if err := ensureHighScoresFile(); err != nil {
		return fmt.Errorf("cannot create high scores file: %w", err)
	}

	datetime := time.Now().Format("01-02-2006 15:04:05")
	row := fmt.Sprintf("\n|%s|%d|%.0f|%s|  ", datetime, score, speed, difficulty)

	f, err := os.OpenFile(highScoresFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(row); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

// ensureHighScoresFile creates HIGHSCORES.md with a header if it doesn't exist.
func ensureHighScoresFile() error {
	if _, err := os.Stat(highScoresFile); err == nil {
		return nil
	}
	header := "|Date|Score|Speed|Difficulty|\n|----|-----|-----|----------| \n"
	return os.WriteFile(highScoresFile, []byte(header), 0644)
}

// FormatHighScores builds display lines for the high scores screen.
// Returns up to maxEntries lines, formatted as a readable table.
func FormatHighScores(entries []HighScoreEntry, maxEntries int) []string {
	if len(entries) == 0 {
		return []string{"  No high scores yet! Play a game first."}
	}
	n := len(entries)
	if n > maxEntries {
		n = maxEntries
	}

	lines := make([]string, 0, n+1)
	lines = append(lines, fmt.Sprintf("  %-4s %-8s %-22s %-10s %-6s", "#", "Score", "Date", "Difficulty", "Speed"))
	lines = append(lines, "  "+strings.Repeat("-", 54))
	for i := 0; i < n; i++ {
		e := entries[i]
		lines = append(lines, fmt.Sprintf("  %-4d %-8d %-22s %-10s %-6.0f",
			i+1, e.Score, e.Date, e.Difficulty, e.Speed))
	}
	return lines
}
