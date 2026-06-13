package trisnake

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// HighScoreEntry represents a single high score record.
type HighScoreEntry struct {
	Date       string
	Score      int
	Speed      float64
	Difficulty string
}

// HighScoresFile is the path to the high scores markdown file.
const HighScoresFile = "HIGHSCORES.md"

// LoadHighScores reads and parses the HIGHSCORES.md file.
// Returns an empty slice if the file doesn't exist or is malformed (never panics).
func LoadHighScores() []HighScoreEntry {
	data, err := ioutil.ReadFile(HighScoresFile)
	if err != nil {
		return []HighScoreEntry{}
	}
	return ParseHighScores(string(data))
}

// ParseHighScores parses markdown table content into HighScoreEntry slice.
// Malformed rows are silently skipped.
func ParseHighScores(content string) []HighScoreEntry {
	var entries []HighScoreEntry
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "|--") || strings.HasPrefix(line, "|Date") {
			continue
		}
		if !strings.HasPrefix(line, "|") {
			continue
		}
		fields := strings.Split(line, "|")
		// Expected: ["", "date", "score", "speed", "difficulty", ""]
		if len(fields) < 5 {
			continue
		}
		score, err := strconv.Atoi(strings.TrimSpace(fields[2]))
		if err != nil {
			continue
		}
		speed, err := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
		if err != nil {
			continue
		}
		entries = append(entries, HighScoreEntry{
			Date:       strings.TrimSpace(fields[1]),
			Score:      score,
			Speed:      speed,
			Difficulty: strings.TrimSpace(fields[4]),
		})
	}
	// Sort by score descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	return entries
}

// SaveHighScore appends a new score row to HIGHSCORES.md.
// Creates the file with headers if it doesn't exist.
// Returns an error instead of crashing.
func SaveHighScore(score int, speed float64, difficulty string) error {
	datetime := time.Now()
	newRow := fmt.Sprintf("|%s|%d|%.0f|%s|\n",
		datetime.Format("01-02-2006 15:04:05"),
		score, speed, difficulty)

	// If file doesn't exist, create with header
	if _, err := os.Stat(HighScoresFile); os.IsNotExist(err) {
		header := "|Date|Score|Speed|Difficulty|\n|----|-----|-----|----------|\n"
		if err := ioutil.WriteFile(HighScoresFile, []byte(header+newRow), 0644); err != nil {
			return err
		}
		return nil
	}

	f, err := os.OpenFile(HighScoresFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(newRow)
	return err
}
