package trisnake

import (
	"os"
	"path/filepath"
	"testing"

	tl "github.com/JoelOtter/termloop"
)

// setupTestGlobals initializes the minimal global state needed for game logic tests.
// This avoids needing a real terminal or game loop.
func setupTestGlobals() {
	sg = tl.NewGame()
	ts = &Titlescreen{GameDifficulty: normal}
	gs = &Gamescreen{
		Level:       tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack}),
		FPS:         12,
		Score:       0,
		ArenaEntity: NewArena(70, 25),
	}
	sp = &Sidepanel{
		ScoreText: tl.NewText(72, 1, "Score: 0", tl.ColorBlack, tl.ColorWhite),
		SpeedText: tl.NewText(72, 3, "Speed: 12", tl.ColorBlack, tl.ColorWhite),
	}
}

// TestRandomInsideArenaBounds verifies that RandomInsideArena always returns
// values within [iMin, iMax) range.
func TestRandomInsideArenaBounds(t *testing.T) {
	for i := 0; i < 10000; i++ {
		v := RandomInsideArena(69, 1)
		if v < 1 || v >= 69 {
			t.Fatalf("RandomInsideArena(69, 1) returned %d, expected [1, 69)", v)
		}
	}
	for i := 0; i < 10000; i++ {
		v := RandomInsideArena(24, 1)
		if v < 1 || v >= 24 {
			t.Fatalf("RandomInsideArena(24, 1) returned %d, expected [1, 24)", v)
		}
	}
}

// TestFoodNotOnSnakeBody verifies that food never spawns on any segment of the snake body.
func TestFoodNotOnSnakeBody(t *testing.T) {
	setupTestGlobals()

	// Create a snake with a known body
	gs.SnakeEntity = NewSnake()
	// Add the snake entity to the level so FoodEntity.MoveFood can access it
	gs.AddEntity(gs.SnakeEntity)

	// Run food spawning many times and verify no overlap
	for i := 0; i < 500; i++ {
		food := NewFood()
		for _, seg := range gs.SnakeEntity.Bodylength {
			if food.Foodposition.X == seg.X && food.Foodposition.Y == seg.Y {
				t.Fatalf("Food spawned on snake body at (%d, %d)", seg.X, seg.Y)
			}
		}
	}
}

// TestFoodNotOnArenaBorder verifies that food never spawns on the arena border.
func TestFoodNotOnArenaBorder(t *testing.T) {
	setupTestGlobals()
	gs.SnakeEntity = NewSnake()

	for i := 0; i < 500; i++ {
		food := NewFood()
		c := Coordinates{food.Foodposition.X, food.Foodposition.Y}
		if gs.ArenaEntity.Contains(c) {
			t.Fatalf("Food spawned on arena border at (%d, %d)", c.X, c.Y)
		}
	}
}

// TestFoodNotOnLargeSnake verifies food spawning works even when the snake is very long.
func TestFoodNotOnLargeSnake(t *testing.T) {
	setupTestGlobals()
	gs.SnakeEntity = NewSnake()

	// Grow the snake to fill a significant portion of the arena
	for y := 6; y <= 23; y++ {
		for x := 1; x <= 68; x++ {
			gs.SnakeEntity.Bodylength = append(gs.SnakeEntity.Bodylength, Coordinates{x, y})
		}
	}

	// Even with a large snake, food should find a valid spot (rows 1-5 are free)
	food := NewFood()
	for _, seg := range gs.SnakeEntity.Bodylength {
		if food.Foodposition.X == seg.X && food.Foodposition.Y == seg.Y {
			t.Fatalf("Food spawned on large snake body at (%d, %d)", seg.X, seg.Y)
		}
	}
	if gs.ArenaEntity.Contains(Coordinates{food.Foodposition.X, food.Foodposition.Y}) {
		t.Fatalf("Food spawned on arena border at (%d, %d)", food.Foodposition.X, food.Foodposition.Y)
	}
}

// TestSnakeContainsCoord verifies the ContainsCoord method on Snake.
func TestSnakeContainsCoord(t *testing.T) {
	setupTestGlobals()
	snake := NewSnake()

	// The default snake has body at {1,6}, {2,6}, {3,6}
	if !snake.ContainsCoord(Coordinates{1, 6}) {
		t.Error("Expected snake to contain {1,6}")
	}
	if !snake.ContainsCoord(Coordinates{3, 6}) {
		t.Error("Expected snake to contain {3,6} (head)")
	}
	if snake.ContainsCoord(Coordinates{0, 0}) {
		t.Error("Snake should not contain {0,0}")
	}
	if snake.ContainsCoord(Coordinates{5, 5}) {
		t.Error("Snake should not contain {5,5}")
	}
}

// TestSaveHighScoreCreatesFile verifies SaveHighScore creates the file if missing.
func TestSaveHighScoreCreatesFile(t *testing.T) {
	// Use a temp directory to avoid polluting the project
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// File should not exist yet
	scoreFile := filepath.Join(tmpDir, "HIGHSCORES.md")
	if _, err := os.Stat(scoreFile); err == nil {
		t.Fatal("HIGHSCORES.md should not exist before test")
	}

	// This should NOT panic or crash even though the file doesn't exist
	SaveHighScore(100, 12, "Normal")

	// Verify file was created and contains the score
	data, err := os.ReadFile(scoreFile)
	if err != nil {
		t.Fatalf("Expected HIGHSCORES.md to be created: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Fatal("HIGHSCORES.md should not be empty after saving score")
	}
}

// TestSaveHighScoreAppends verifies that multiple saves append correctly.
func TestSaveHighScoreAppends(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	SaveHighScore(100, 12, "Normal")
	SaveHighScore(200, 25, "Hard")

	data, err := os.ReadFile(filepath.Join(tmpDir, "HIGHSCORES.md"))
	if err != nil {
		t.Fatalf("Failed to read HIGHSCORES.md: %v", err)
	}
	content := string(data)
	// Should contain both scores
	if !contains(content, "100") || !contains(content, "200") {
		t.Fatalf("Expected both scores in file, got: %s", content)
	}
}

// TestSaveHighScoreReadOnlyDir verifies SaveHighScore handles write failure gracefully.
func TestSaveHighScoreReadOnlyDir(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create the file first, then make it read-only
	os.WriteFile("HIGHSCORES.md", []byte("existing"), 0644)
	os.Chmod("HIGHSCORES.md", 0444)

	// This should NOT panic — it should silently handle the write error
	SaveHighScore(999, 8, "Easy")

	// Restore permissions for cleanup
	os.Chmod("HIGHSCORES.md", 0644)
}

// TestRestartGameResetsState verifies that RestartGame fully resets game state.
func TestRestartGameResetsState(t *testing.T) {
	setupTestGlobals()
	gs.SnakeEntity = NewSnake()
	gs.FoodEntity = NewFood()

	// Simulate a played game: grow the snake, change score and speed
	for i := 0; i < 10; i++ {
		gs.SnakeEntity.Bodylength = append(gs.SnakeEntity.Bodylength, Coordinates{10 + i, 10})
	}
	gs.Score = 50
	gs.FPS = 20
	gs.SnakeEntity.Direction = up

	gs.AddEntity(gs.SnakeEntity)
	gs.AddEntity(gs.FoodEntity)

	// Restart
	RestartGame()

	// Verify score is reset
	if gs.Score != 0 {
		t.Errorf("Expected score 0 after restart, got %d", gs.Score)
	}

	// Verify FPS is reset to normal difficulty default
	if gs.FPS != 12 {
		t.Errorf("Expected FPS 12 after restart, got %.0f", gs.FPS)
	}

	// Verify snake body is reset to initial 3 segments
	if len(gs.SnakeEntity.Bodylength) != 3 {
		t.Errorf("Expected snake body length 3 after restart, got %d", len(gs.SnakeEntity.Bodylength))
	}

	// Verify snake direction is reset to right
	if gs.SnakeEntity.Direction != right {
		t.Errorf("Expected snake direction right after restart, got %d", gs.SnakeEntity.Direction)
	}

	// Verify food is not on snake body
	for _, seg := range gs.SnakeEntity.Bodylength {
		if gs.FoodEntity.Foodposition.X == seg.X && gs.FoodEntity.Foodposition.Y == seg.Y {
			t.Errorf("Food at (%d,%d) overlaps snake body after restart", seg.X, seg.Y)
		}
	}
}

// TestNewSnakeDefaultState verifies NewSnake creates a clean default snake.
func TestNewSnakeDefaultState(t *testing.T) {
	snake := NewSnake()

	if snake.Direction != right {
		t.Errorf("Expected default direction right, got %d", snake.Direction)
	}
	if len(snake.Bodylength) != 3 {
		t.Errorf("Expected 3 body segments, got %d", len(snake.Bodylength))
	}
	expected := []Coordinates{{1, 6}, {2, 6}, {3, 6}}
	for i, exp := range expected {
		if snake.Bodylength[i] != exp {
			t.Errorf("Body[%d] = %v, expected %v", i, snake.Bodylength[i], exp)
		}
	}
}

// TestSetDifficultyFPS verifies FPS is set correctly for each difficulty.
func TestSetDifficultyFPS(t *testing.T) {
	setupTestGlobals()

	ts.GameDifficulty = easy
	SetDiffiultyFPS()
	if gs.FPS != 8 {
		t.Errorf("Easy FPS: expected 8, got %.0f", gs.FPS)
	}

	ts.GameDifficulty = normal
	SetDiffiultyFPS()
	if gs.FPS != 12 {
		t.Errorf("Normal FPS: expected 12, got %.0f", gs.FPS)
	}

	ts.GameDifficulty = hard
	SetDiffiultyFPS()
	if gs.FPS != 25 {
		t.Errorf("Hard FPS: expected 25, got %.0f", gs.FPS)
	}
}

// TestArenaBorderContains verifies the arena border detection.
func TestArenaBorderContains(t *testing.T) {
	arena := NewArena(70, 25)

	// Corners should be border
	if !arena.Contains(Coordinates{0, 0}) {
		t.Error("Expected (0,0) to be border")
	}
	if !arena.Contains(Coordinates{69, 0}) {
		t.Error("Expected (69,0) to be border")
	}
	if !arena.Contains(Coordinates{0, 24}) {
		t.Error("Expected (0,24) to be border")
	}

	// Interior should NOT be border
	if arena.Contains(Coordinates{5, 5}) {
		t.Error("Expected (5,5) to NOT be border")
	}
	if arena.Contains(Coordinates{35, 12}) {
		t.Error("Expected (35,12) to NOT be border")
	}
}

// helper
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
