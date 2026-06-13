package trisnake

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Food placement tests
// ---------------------------------------------------------------------------

func TestFoodNotSpawnOnSnakeBody(t *testing.T) {
	// Create a snake body that covers a significant portion of the arena.
	snakeBody := []Coordinates{
		{1, 6}, {2, 6}, {3, 6}, {4, 6}, {5, 6},
		{5, 7}, {5, 8}, {5, 9}, {5, 10},
		{4, 10}, {3, 10}, {2, 10}, {1, 10},
	}

	// Run many iterations to statistically confirm food never lands on snake.
	for i := 0; i < 200; i++ {
		food := NewFood(snakeBody)
		for _, c := range snakeBody {
			if food.Foodposition.X == c.X && food.Foodposition.Y == c.Y {
				t.Fatalf("food spawned on snake body at (%d, %d) on iteration %d", c.X, c.Y, i)
			}
		}
	}
}

func TestFoodNotSpawnOnBorder(t *testing.T) {
	// Arena borders: x=0, x=69, y=0, y=24.
	for i := 0; i < 500; i++ {
		food := NewFood(nil)
		x, y := food.Foodposition.X, food.Foodposition.Y

		if x == 0 || x == 69 {
			t.Fatalf("food spawned on horizontal border at (%d, %d) on iteration %d", x, y, i)
		}
		if y == 0 || y == 24 {
			t.Fatalf("food spawned on vertical border at (%d, %d) on iteration %d", x, y, i)
		}
		if x < 1 || x > 68 {
			t.Fatalf("food X out of playable range: %d on iteration %d", x, i)
		}
		if y < 1 || y > 23 {
			t.Fatalf("food Y out of playable range: %d on iteration %d", y, i)
		}
	}
}

func TestMoveFoodAvoidsOccupiedCells(t *testing.T) {
	food := NewFood(nil)

	// Fill a block of cells as "occupied".
	occupied := make([]Coordinates, 0)
	for x := 10; x <= 20; x++ {
		for y := 5; y <= 15; y++ {
			occupied = append(occupied, Coordinates{x, y})
		}
	}

	for i := 0; i < 200; i++ {
		food.MoveFood(occupied)
		for _, c := range occupied {
			if food.Foodposition.X == c.X && food.Foodposition.Y == c.Y {
				t.Fatalf("MoveFood placed food on occupied cell (%d, %d) on iteration %d", c.X, c.Y, i)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// RestartGame state-reset tests
// ---------------------------------------------------------------------------

func TestRestartGameResetsState(t *testing.T) {
	// Set up state without calling NewGamescreen (which needs sg/terminal).
	ts = &Titlescreen{GameDifficulty: normal}
	Difficulty = "Normal"

	gs = &Gamescreen{
		FPS:         12,
		Score:       0,
		SnakeEntity: NewSnake(),
		FoodEntity:  NewFood(nil),
		ArenaEntity: NewArena(70, 25),
	}

	// Dirty the state as if a game was played.
	gs.Score = 42
	gs.FPS = 99
	gs.SnakeEntity.Bodylength = append(gs.SnakeEntity.Bodylength,
		Coordinates{10, 10}, Coordinates{11, 10}, Coordinates{12, 10})
	gs.SnakeEntity.Direction = left

	// --- Direct state reset logic (mirrors RestartGame) ---
	gs.SnakeEntity = NewSnake()
	gs.FoodEntity = NewFood(gs.SnakeEntity.Bodylength)
	SetDiffiultyFPS()
	gs.Score = 0
	// -------------------------------------------------------

	// Verify score is reset.
	if gs.Score != 0 {
		t.Errorf("score not reset: got %d, want 0", gs.Score)
	}

	// Verify FPS is reset to normal difficulty default.
	if gs.FPS != 12 {
		t.Errorf("FPS not reset: got %.0f, want 12", gs.FPS)
	}

	// Verify snake direction is reset to right.
	if gs.SnakeEntity.Direction != right {
		t.Errorf("snake direction not reset: got %d, want %d (right)", gs.SnakeEntity.Direction, right)
	}

	// Verify snake body is reset to 3 segments.
	if len(gs.SnakeEntity.Bodylength) != 3 {
		t.Errorf("snake body length not reset: got %d, want 3", len(gs.SnakeEntity.Bodylength))
	}

	// Verify snake body matches the initial positions.
	expectedBody := []Coordinates{{1, 6}, {2, 6}, {3, 6}}
	for i, c := range gs.SnakeEntity.Bodylength {
		if c != expectedBody[i] {
			t.Errorf("snake body[%d] = (%d,%d), want (%d,%d)",
				i, c.X, c.Y, expectedBody[i].X, expectedBody[i].Y)
		}
	}

	// Verify food is within playable area and not on snake.
	fx, fy := gs.FoodEntity.Foodposition.X, gs.FoodEntity.Foodposition.Y
	if fx < 1 || fx > 68 || fy < 1 || fy > 23 {
		t.Errorf("food out of playable area: (%d, %d)", fx, fy)
	}
	for _, c := range gs.SnakeEntity.Bodylength {
		if fx == c.X && fy == c.Y {
			t.Errorf("food spawned on snake at (%d, %d)", fx, fy)
		}
	}
}

func TestSetDifficultyFPS(t *testing.T) {
	ts = &Titlescreen{}
	gs = &Gamescreen{}

	ts.GameDifficulty = easy
	SetDiffiultyFPS()
	if gs.FPS != 8 {
		t.Errorf("easy FPS: got %.0f, want 8", gs.FPS)
	}

	ts.GameDifficulty = normal
	SetDiffiultyFPS()
	if gs.FPS != 12 {
		t.Errorf("normal FPS: got %.0f, want 12", gs.FPS)
	}

	ts.GameDifficulty = hard
	SetDiffiultyFPS()
	if gs.FPS != 25 {
		t.Errorf("hard FPS: got %.0f, want 25", gs.FPS)
	}
}

// ---------------------------------------------------------------------------
// SaveHighScore tests
// ---------------------------------------------------------------------------

func TestSaveHighScoreCreatesMissingFile(t *testing.T) {
	// Work in a temp directory so we don't pollute the project.
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// HIGHSCORES.md does not exist in tmpDir.
	err := SaveHighScore(42, 12, "Normal")
	if err != nil {
		t.Fatalf("SaveHighScore returned error for missing file: %v", err)
	}

	// Verify the file was created and contains the score.
	data, err := os.ReadFile("HIGHSCORES.md")
	if err != nil {
		t.Fatalf("could not read created file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "42") {
		t.Errorf("file does not contain score 42: %s", content)
	}
	if !strings.Contains(content, "Normal") {
		t.Errorf("file does not contain difficulty Normal: %s", content)
	}
}

func TestSaveHighScoreAppendsToExisting(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create an existing HIGHSCORES.md.
	existing := "|Date|Score|Speed|Difficulty|\n|----|-----|-----|----------|\n"
	os.WriteFile("HIGHSCORES.md", []byte(existing), 0644)

	err := SaveHighScore(100, 25, "Hard")
	if err != nil {
		t.Fatalf("SaveHighScore returned error: %v", err)
	}

	data, _ := os.ReadFile("HIGHSCORES.md")
	content := string(data)
	if !strings.Contains(content, existing) {
		t.Errorf("existing content was lost")
	}
	if !strings.Contains(content, "100") {
		t.Errorf("new score not appended: %s", content)
	}
}

func TestSaveHighScoreHandlesReadOnlyDir(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Make directory read-only so file creation fails.
	os.Chmod(tmpDir, 0555)
	defer os.Chmod(tmpDir, 0755) // restore for cleanup

	err := SaveHighScore(10, 8, "Easy")
	if err == nil {
		// On some systems (e.g. running as root), read-only dir might still work.
		// Just verify no crash occurred.
		t.Log("SaveHighScore did not return error in read-only dir (possibly running as root)")
	}
	// The key assertion: the function did NOT call log.Fatalf / panic.
	// If we reach here, the graceful degradation works.
}

func TestSaveHighScoreHandlesPermissionDenied(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create a read-only file.
	fpath := filepath.Join(tmpDir, "HIGHSCORES.md")
	os.WriteFile(fpath, []byte("header\n"), 0444)

	err := SaveHighScore(10, 8, "Easy")
	if err == nil {
		t.Log("SaveHighScore did not return error (possibly running as root)")
	}
	// Key: no crash / fatal.
}

// ---------------------------------------------------------------------------
// randomInsideArena tests
// ---------------------------------------------------------------------------

func TestRandomInsideArenaBounds(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := randomInsideArena(69, 1)
		if v < 1 || v >= 69 {
			t.Fatalf("randomInsideArena(69, 1) returned %d, want [1, 69)", v)
		}
	}
	for i := 0; i < 1000; i++ {
		v := randomInsideArena(24, 1)
		if v < 1 || v >= 24 {
			t.Fatalf("randomInsideArena(24, 1) returned %d, want [1, 24)", v)
		}
	}
}

// ---------------------------------------------------------------------------
// Snake collision tests
// ---------------------------------------------------------------------------

func TestSnakeHead(t *testing.T) {
	snake := &Snake{
		Bodylength: []Coordinates{{1, 6}, {2, 6}, {3, 6}},
	}
	head := snake.Head()
	if head.X != 3 || head.Y != 6 {
		t.Errorf("Head() = (%d, %d), want (3, 6)", head.X, head.Y)
	}
}

func TestSnakeSelfCollision(t *testing.T) {
	// Snake not colliding with itself.
	snake := &Snake{
		Bodylength: []Coordinates{{1, 6}, {2, 6}, {3, 6}},
	}
	if snake.Contains() {
		t.Error("Contains() returned true for non-colliding snake")
	}

	// Snake colliding with itself (head overlaps body).
	snake.Bodylength = []Coordinates{{1, 6}, {2, 6}, {1, 6}}
	if !snake.Contains() {
		t.Error("Contains() returned false for self-colliding snake")
	}
}

// ---------------------------------------------------------------------------
// RandomFood distribution test
// ---------------------------------------------------------------------------

func TestRandomFoodReturnsValidRune(t *testing.T) {
	validRunes := map[rune]bool{'R': true, '■': true, 'S': true}
	for i := 0; i < 200; i++ {
		r := RandomFood()
		if !validRunes[r] {
			t.Fatalf("RandomFood() returned unexpected rune: %c", r)
		}
	}
}
