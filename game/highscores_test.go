package trisnake

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseHighScores_ValidContent(t *testing.T) {
	content := `|Date|Score|Speed|Difficulty|
|----|-----|-----|----------|
|10-28-2019 10:54:29|5|12|Normal|
|06-13-2026 15:30:00|42|20|Hard|
|01-01-2025 08:00:00|10|8|Easy|  `

	entries := ParseHighScores(content)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	// Should be sorted by score descending.
	if entries[0].Score != 42 {
		t.Errorf("expected top score 42, got %d", entries[0].Score)
	}
	if entries[0].Difficulty != "Hard" {
		t.Errorf("expected top entry difficulty Hard, got %s", entries[0].Difficulty)
	}
	if entries[1].Score != 10 {
		t.Errorf("expected second score 10, got %d", entries[1].Score)
	}
	if entries[2].Score != 5 {
		t.Errorf("expected third score 5, got %d", entries[2].Score)
	}
	if entries[0].Speed != 20 {
		t.Errorf("expected top speed 20, got %.0f", entries[0].Speed)
	}
}

func TestParseHighScores_EmptyContent(t *testing.T) {
	entries := ParseHighScores("")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for empty content, got %d", len(entries))
	}
}

func TestParseHighScores_MalformedLines(t *testing.T) {
	content := `|Date|Score|Speed|Difficulty|
|----|-----|-----|----------|
|bad line with no pipes
|10-28-2019 10:54:29|notanumber|12|Normal|
|06-13-2026 15:30:00|42|20|Hard|`

	entries := ParseHighScores(content)
	if len(entries) != 1 {
		t.Fatalf("expected 1 valid entry (skipping malformed), got %d", len(entries))
	}
	if entries[0].Score != 42 {
		t.Errorf("expected score 42, got %d", entries[0].Score)
	}
}

func TestParseHighScores_HeaderOnly(t *testing.T) {
	content := `|Date|Score|Speed|Difficulty|
|----|-----|-----|----------| `
	entries := ParseHighScores(content)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for header-only content, got %d", len(entries))
	}
}

func TestLoadHighScores_MissingFile(t *testing.T) {
	// Use a temp dir to ensure no HIGHSCORES.md exists.
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	entries := LoadHighScores()
	if entries != nil && len(entries) != 0 {
		t.Fatalf("expected nil/empty entries for missing file, got %d", len(entries))
	}
}

func TestSaveAndLoadHighScoreEntry(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Save a score.
	err := SaveHighScoreEntry(100, 15, "Hard")
	if err != nil {
		t.Fatalf("SaveHighScoreEntry failed: %v", err)
	}

	// Verify file was created.
	data, err := os.ReadFile(highScoresFile)
	if err != nil {
		t.Fatalf("cannot read created file: %v", err)
	}
	if !strings.Contains(string(data), "100") {
		t.Error("saved file does not contain score 100")
	}
	if !strings.Contains(string(data), "Hard") {
		t.Error("saved file does not contain difficulty Hard")
	}

	// Save another score.
	err = SaveHighScoreEntry(50, 12, "Normal")
	if err != nil {
		t.Fatalf("second SaveHighScoreEntry failed: %v", err)
	}

	// Load and verify.
	entries := LoadHighScores()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Should be sorted: 100 first, then 50.
	if entries[0].Score != 100 {
		t.Errorf("expected top score 100, got %d", entries[0].Score)
	}
	if entries[1].Score != 50 {
		t.Errorf("expected second score 50, got %d", entries[1].Score)
	}
}

func TestSaveHighScoreEntry_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Pre-create the file with existing content.
	existing := "|Date|Score|Speed|Difficulty|\n|----|-----|-----|----------| \n|01-01-2025 00:00:00|25|12|Normal|  \n"
	os.WriteFile(highScoresFile, []byte(existing), 0644)

	err := SaveHighScoreEntry(75, 20, "Hard")
	if err != nil {
		t.Fatalf("SaveHighScoreEntry failed: %v", err)
	}

	entries := LoadHighScores()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Score != 75 {
		t.Errorf("expected top score 75, got %d", entries[0].Score)
	}
}

func TestFormatHighScores_WithEntries(t *testing.T) {
	entries := []HighScoreEntry{
		{Date: "06-13-2026 15:30:00", Score: 42, Speed: 20, Difficulty: "Hard"},
		{Date: "01-01-2025 08:00:00", Score: 10, Speed: 8, Difficulty: "Easy"},
	}
	lines := FormatHighScores(entries, 10)
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines (header+separator+entries), got %d", len(lines))
	}
	// Check that the first entry appears in the output.
	found := false
	for _, line := range lines {
		if strings.Contains(line, "42") && strings.Contains(line, "Hard") {
			found = true
			break
		}
	}
	if !found {
		t.Error("formatted output does not contain the top score entry")
	}
}

func TestFormatHighScores_EmptyEntries(t *testing.T) {
	lines := FormatHighScores(nil, 10)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line for empty entries, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "No high scores") {
		t.Errorf("expected 'No high scores' message, got: %s", lines[0])
	}
}

func TestFormatHighScores_MaxEntriesLimit(t *testing.T) {
	entries := make([]HighScoreEntry, 20)
	for i := range entries {
		entries[i] = HighScoreEntry{
			Date:       "01-01-2025 00:00:00",
			Score:      100 - i,
			Speed:      12,
			Difficulty: "Normal",
		}
	}
	lines := FormatHighScores(entries, 5)
	// header + separator + 5 entries = 7 lines.
	if len(lines) != 7 {
		t.Fatalf("expected 7 lines (header+separator+5 entries), got %d", len(lines))
	}
}

func TestEnsureHighScoresFile_CreatesWithHeader(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	err := ensureHighScoresFile()
	if err != nil {
		t.Fatalf("ensureHighScoresFile failed: %v", err)
	}

	data, err := os.ReadFile(highScoresFile)
	if err != nil {
		t.Fatalf("cannot read file: %v", err)
	}
	if !strings.Contains(string(data), "|Date|Score|Speed|Difficulty|") {
		t.Error("created file missing header row")
	}
}

func TestEnsureHighScoresFile_DoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Pre-create with custom content.
	custom := "custom content"
	os.WriteFile(filepath.Join(dir, highScoresFile), []byte(custom), 0644)

	err := ensureHighScoresFile()
	if err != nil {
		t.Fatalf("ensureHighScoresFile failed: %v", err)
	}

	data, err := os.ReadFile(highScoresFile)
	if err != nil {
		t.Fatalf("cannot read file: %v", err)
	}
	if string(data) != custom {
		t.Errorf("ensureHighScoresFile overwrote existing file: got %q, want %q", string(data), custom)
	}
}
