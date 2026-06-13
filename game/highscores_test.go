package trisnake

import (
	"os"
	"testing"
)

func TestParseHighScores(t *testing.T) {
	content := `|Date|Score|Speed|Difficulty|
|----|-----|-----|----------|
|10-28-2019 10:54:29|5|12|Normal|
|10-29-2019 11:00:00|10|15|Hard|
|10-30-2019 12:00:00|3|8|Easy|`

	entries := ParseHighScores(content)

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Should be sorted by score descending
	if entries[0].Score != 10 {
		t.Errorf("Expected first entry score 10, got %d", entries[0].Score)
	}
	if entries[0].Difficulty != "Hard" {
		t.Errorf("Expected first entry difficulty Hard, got %s", entries[0].Difficulty)
	}
	if entries[1].Score != 5 {
		t.Errorf("Expected second entry score 5, got %d", entries[1].Score)
	}
	if entries[2].Score != 3 {
		t.Errorf("Expected third entry score 3, got %d", entries[2].Score)
	}
}

func TestParseHighScoresMalformed(t *testing.T) {
	content := `|Date|Score|Speed|Difficulty|
|----|-----|-----|----------|
|10-28-2019 10:54:29|5|12|Normal|
|bad row without pipes
|10-29-2019 11:00:00|invalid|15|Hard|
|10-30-2019 12:00:00|3|8|Easy|`

	entries := ParseHighScores(content)

	if len(entries) != 2 {
		t.Errorf("Expected 2 valid entries, got %d", len(entries))
	}
}

func TestParseHighScoresEmpty(t *testing.T) {
	entries := ParseHighScores("")
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for empty content, got %d", len(entries))
	}

	entries = ParseHighScores("|Date|Score|Speed|Difficulty|\n|----|-----|-----|----------|\n")
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for header-only content, got %d", len(entries))
	}
}

func TestParseHighScoresTrailingSpaces(t *testing.T) {
	// The existing HIGHSCORES.md format has trailing spaces
	content := "|Date|Score|Speed|Difficulty|\n|----|-----|-----|----------| \n|10-28-2019 10:54:29|5|12|Normal|  \n"

	entries := ParseHighScores(content)
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Score != 5 || entries[0].Difficulty != "Normal" {
		t.Errorf("Unexpected entry values: %+v", entries[0])
	}
}

func TestSaveHighScore(t *testing.T) {
	// Back up the original file
	backupFile := "HIGHSCORES.md.testbak"
	os.Rename(HighScoresFile, backupFile)
	defer os.Rename(backupFile, HighScoresFile)

	// Save a score to a fresh file
	err := SaveHighScore(42, 15, "Hard")
	if err != nil {
		t.Fatalf("Failed to save high score: %v", err)
	}

	// Load and verify
	entries := LoadHighScores()
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Score != 42 {
		t.Errorf("Expected score 42, got %d", entries[0].Score)
	}
	if entries[0].Speed != 15 {
		t.Errorf("Expected speed 15, got %.0f", entries[0].Speed)
	}
	if entries[0].Difficulty != "Hard" {
		t.Errorf("Expected difficulty Hard, got %s", entries[0].Difficulty)
	}

	// Save another score
	err = SaveHighScore(100, 20, "Normal")
	if err != nil {
		t.Fatalf("Failed to save second high score: %v", err)
	}

	entries = LoadHighScores()
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}
	// Should be sorted descending by score
	if entries[0].Score != 100 {
		t.Errorf("Expected first entry score 100, got %d", entries[0].Score)
	}

	// Clean up test file
	os.Remove(HighScoresFile)
}

func TestLoadHighScoresFileNotExist(t *testing.T) {
	// Back up the original file if it exists
	backupFile := "HIGHSCORES.md.nobak"
	err := os.Rename(HighScoresFile, backupFile)
	if err != nil {
		// File might not exist, which is fine
		t.Skip("Cannot rename HIGHSCORES.md, skipping")
	}
	defer os.Rename(backupFile, HighScoresFile)

	entries := LoadHighScores()
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries when file doesn't exist, got %d", len(entries))
	}
}
