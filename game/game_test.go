package trisnake

import (
	"testing"

	tl "github.com/JoelOtter/termloop"
)

// --- Game constructor tests ---

func TestNewGame_Defaults(t *testing.T) {
	g := NewGame()
	if g.difficulty != DiffNormal {
		t.Errorf("expected DiffNormal, got %d", g.difficulty)
	}
	if g.colors.Target != TargetSnake {
		t.Errorf("expected TargetSnake, got %d", g.colors.Target)
	}
	if g.colors.SnakeIdx != 10 {
		t.Errorf("expected SnakeIdx=10, got %d", g.colors.SnakeIdx)
	}
	if g.colors.ArenaIdx != 10 {
		t.Errorf("expected ArenaIdx=10, got %d", g.colors.ArenaIdx)
	}
	if g.state != nil {
		t.Error("expected nil state before game starts")
	}
}

// --- Difficulty tests ---

func TestDifficulty_String(t *testing.T) {
	cases := map[Difficulty]string{
		DiffEasy:   "Easy",
		DiffNormal: "Normal",
		DiffHard:   "Hard",
	}
	for d, expected := range cases {
		if d.String() != expected {
			t.Errorf("Difficulty(%d).String() = %q, want %q", d, d.String(), expected)
		}
	}
}

func TestDifficulty_BaseFPS(t *testing.T) {
	cases := map[Difficulty]float64{
		DiffEasy:   8,
		DiffNormal: 12,
		DiffHard:   25,
	}
	for d, expected := range cases {
		if d.baseFPS() != expected {
			t.Errorf("Difficulty(%d).baseFPS() = %f, want %f", d, d.baseFPS(), expected)
		}
	}
}

// --- ColorConfig tests ---

func TestColorConfig_ActiveIdx(t *testing.T) {
	cc := ColorConfig{Target: TargetSnake, SnakeIdx: 14, ArenaIdx: 18}
	if cc.ActiveIdx() != 14 {
		t.Errorf("ActiveIdx for Snake = %d, want 14", cc.ActiveIdx())
	}
	cc.Target = TargetArena
	if cc.ActiveIdx() != 18 {
		t.Errorf("ActiveIdx for Arena = %d, want 18", cc.ActiveIdx())
	}
}

func TestColorConfig_SetActiveIdx(t *testing.T) {
	cc := ColorConfig{Target: TargetSnake, SnakeIdx: 10, ArenaIdx: 10}
	cc.SetActiveIdx(16)
	if cc.SnakeIdx != 16 {
		t.Errorf("SnakeIdx = %d, want 16", cc.SnakeIdx)
	}
	cc.Target = TargetArena
	cc.SetActiveIdx(20)
	if cc.ArenaIdx != 20 {
		t.Errorf("ArenaIdx = %d, want 20", cc.ArenaIdx)
	}
	// Snake index should not have changed.
	if cc.SnakeIdx != 16 {
		t.Errorf("SnakeIdx changed unexpectedly: got %d, want 16", cc.SnakeIdx)
	}
}

// --- Snake unit tests ---

func TestNewSnake_DefaultState(t *testing.T) {
	g := NewGame()
	snake := newSnake(g)

	if snake.Direction != DirRight {
		t.Errorf("expected DirRight, got %d", snake.Direction)
	}
	if len(snake.Body) != 3 {
		t.Fatalf("expected body length 3, got %d", len(snake.Body))
	}
	// Head should be the last element.
	head := snake.Head()
	if head.X != 3 || head.Y != 6 {
		t.Errorf("expected head at (3,6), got (%d,%d)", head.X, head.Y)
	}
}

func TestSnake_Head(t *testing.T) {
	g := NewGame()
	snake := newSnake(g)
	head := snake.Head()
	if *head != snake.Body[len(snake.Body)-1] {
		t.Error("Head() should return last body element")
	}
}

func TestSnake_Advance(t *testing.T) {
	g := NewGame()
	snake := newSnake(g)

	// Start heading right from (3,6).
	snake.Direction = DirRight
	next := snake.advance()
	if next.X != 4 || next.Y != 6 {
		t.Errorf("DirRight from (3,6) -> (%d,%d), want (4,6)", next.X, next.Y)
	}

	snake.Direction = DirLeft
	next = snake.advance()
	if next.X != 2 || next.Y != 6 {
		t.Errorf("DirLeft from (3,6) -> (%d,%d), want (2,6)", next.X, next.Y)
	}

	snake.Direction = DirUp
	next = snake.advance()
	if next.X != 3 || next.Y != 5 {
		t.Errorf("DirUp from (3,6) -> (%d,%d), want (3,5)", next.X, next.Y)
	}

	snake.Direction = DirDown
	next = snake.advance()
	if next.X != 3 || next.Y != 7 {
		t.Errorf("DirDown from (3,6) -> (%d,%d), want (3,7)", next.X, next.Y)
	}
}

func TestSnake_BorderCollision(t *testing.T) {
	g := NewGame()
	arena := newArena(g, arenaWidth, arenaHeight)
	snake := newSnake(g)

	// Manually set up a minimal game state so collision checks work.
	g.state = &GameState{arena: arena, snake: snake}

	// Snake starts at (3,6) which is inside the arena.
	if snake.headCollidesWithBorder() {
		t.Error("snake at (3,6) should not collide with border")
	}

	// Move head to border position (0, 6).
	snake.Body[len(snake.Body)-1] = Coordinates{0, 6}
	if !snake.headCollidesWithBorder() {
		t.Error("snake at (0,6) should collide with border")
	}

	// Move head to top border (5, 0).
	snake.Body[len(snake.Body)-1] = Coordinates{5, 0}
	if !snake.headCollidesWithBorder() {
		t.Error("snake at (5,0) should collide with border")
	}
}

func TestSnake_BodyCollision(t *testing.T) {
	g := NewGame()
	snake := newSnake(g)

	// Normal body: no self-collision.
	if snake.headCollidesWithBody() {
		t.Error("fresh snake should not self-collide")
	}

	// Force head to overlap a body segment.
	snake.Body[len(snake.Body)-1] = snake.Body[0]
	if !snake.headCollidesWithBody() {
		t.Error("snake with head on tail should self-collide")
	}
}

func TestSnake_FoodCollision(t *testing.T) {
	g := NewGame()
	snake := newSnake(g)
	food := newFood()
	g.state = &GameState{snake: snake, food: food}

	// Place food at head position.
	food.Position = *snake.Head()
	if !snake.headCollidesWithFood() {
		t.Error("food at head position should trigger collision")
	}

	// Move food away.
	food.Position = Coordinates{50, 50}
	if snake.headCollidesWithFood() {
		t.Error("food far from head should not trigger collision")
	}
}

// --- Food unit tests ---

func TestNewFood_HasPosition(t *testing.T) {
	food := newFood()
	if food.Position.X < 1 || food.Position.X >= insideBorderW {
		t.Errorf("food X=%d out of range [1, %d)", food.Position.X, insideBorderW)
	}
	if food.Position.Y < 1 || food.Position.Y >= insideBorderH {
		t.Errorf("food Y=%d out of range [1, %d)", food.Position.Y, insideBorderH)
	}
	if food.Emoji == 0 {
		t.Error("food emoji should not be zero")
	}
}

func TestFood_Randomize(t *testing.T) {
	food := newFood()
	oldPos := food.Position
	// Re-randomize several times; at least one should differ.
	different := false
	for i := 0; i < 20; i++ {
		food.Randomize()
		if food.Position != oldPos {
			different = true
			break
		}
	}
	if !different {
		t.Error("food should eventually randomize to a different position")
	}
}

func TestFood_Contains(t *testing.T) {
	food := &Food{Position: Coordinates{10, 10}}
	if !food.Contains(Coordinates{10, 10}) {
		t.Error("food should contain its own position")
	}
	if food.Contains(Coordinates{10, 11}) {
		t.Error("food should not contain a different position")
	}
}

// --- Arena unit tests ---

func TestNewArena_BorderIntegrity(t *testing.T) {
	g := NewGame()
	arena := newArena(g, arenaWidth, arenaHeight)

	// Width is arenaWidth-1, Height is arenaHeight-1.
	if arena.Width != arenaWidth-1 {
		t.Errorf("arena.Width = %d, want %d", arena.Width, arenaWidth-1)
	}
	if arena.Height != arenaHeight-1 {
		t.Errorf("arena.Height = %d, want %d", arena.Height, arenaHeight-1)
	}

	// Top border should be present.
	if !arena.Contains(Coordinates{0, 0}) {
		t.Error("top-left corner should be border")
	}
	if !arena.Contains(Coordinates{35, 0}) {
		t.Error("top edge should be border")
	}

	// Bottom border.
	if !arena.Contains(Coordinates{35, arena.Height}) {
		t.Error("bottom edge should be border")
	}

	// Left border.
	if !arena.Contains(Coordinates{0, 12}) {
		t.Error("left edge should be border")
	}

	// Right border.
	if !arena.Contains(Coordinates{arena.Width, 12}) {
		t.Error("right edge should be border")
	}

	// Interior should not be border.
	if arena.Contains(Coordinates{35, 12}) {
		t.Error("interior should not be border")
	}
}

// --- Game state mutation tests (without termloop) ---

// newTestGameState creates a Game with a manually initialized GameState
// for testing score/FPS mutations without needing termloop.
func newTestGameState() *Game {
	g := NewGame()
	g.state = &GameState{
		score: 0,
		fps:   g.difficulty.baseFPS(),
		sidepanel: &Sidepanel{
			ScoreText: tl.NewText(0, 0, "Score: 0", tl.ColorWhite, tl.ColorBlack),
			SpeedText: tl.NewText(0, 0, "Speed: 12", tl.ColorWhite, tl.ColorBlack),
		},
	}
	return g
}

func TestGame_AddScore(t *testing.T) {
	g := newTestGameState()

	g.addScore(5)
	if g.state.score != 5 {
		t.Errorf("score = %d, want 5", g.state.score)
	}

	g.addScore(3)
	if g.state.score != 8 {
		t.Errorf("score = %d, want 8", g.state.score)
	}
}

func TestGame_AdjustFPS(t *testing.T) {
	g := newTestGameState()
	initialFPS := g.state.fps

	g.adjustFPS(-3)
	if g.state.fps != initialFPS-3 {
		t.Errorf("fps = %f, want %f", g.state.fps, initialFPS-3)
	}

	g.adjustFPS(5)
	if g.state.fps != initialFPS+2 {
		t.Errorf("fps = %f, want %f", g.state.fps, initialFPS+2)
	}
}

func TestGame_SetFPS(t *testing.T) {
	g := newTestGameState()

	g.setFPS(20)
	if g.state.fps != 20 {
		t.Errorf("fps = %f, want 20", g.state.fps)
	}
}

func TestGame_CurrentScoreAndFPS_NilState(t *testing.T) {
	g := NewGame()
	if g.currentScore() != 0 {
		t.Errorf("expected 0 score for nil state, got %d", g.currentScore())
	}
	if g.currentFPS() != 0 {
		t.Errorf("expected 0 fps for nil state, got %f", g.currentFPS())
	}
}

func TestGame_CurrentScoreAndFPS_WithState(t *testing.T) {
	g := newTestGameState()
	g.state.score = 42
	g.state.fps = 15

	if g.currentScore() != 42 {
		t.Errorf("currentScore() = %d, want 42", g.currentScore())
	}
	if g.currentFPS() != 15 {
		t.Errorf("currentFPS() = %f, want 15", g.currentFPS())
	}
}

// --- Build state integration test ---

func TestBuildState_CreatesValidRound(t *testing.T) {
	g := NewGame()
	gs := newGameScreen(g)
	state := gs.buildState()

	if state.score != 0 {
		t.Errorf("new round score = %d, want 0", state.score)
	}
	if state.fps != DiffNormal.baseFPS() {
		t.Errorf("new round fps = %f, want %f", state.fps, DiffNormal.baseFPS())
	}
	if state.snake == nil {
		t.Error("new round should have a snake")
	}
	if state.food == nil {
		t.Error("new round should have food")
	}
	if state.arena == nil {
		t.Error("new round should have arena")
	}
	if state.sidepanel == nil {
		t.Error("new round should have sidepanel")
	}

	// Snake should have 3 body segments.
	if len(state.snake.Body) != 3 {
		t.Errorf("snake body length = %d, want 3", len(state.snake.Body))
	}

	// Food should be inside arena.
	if state.food.Position.X < 1 || state.food.Position.X >= insideBorderW {
		t.Errorf("food X=%d out of range", state.food.Position.X)
	}
}

func TestBuildState_DifferentDifficulties(t *testing.T) {
	for _, diff := range []Difficulty{DiffEasy, DiffNormal, DiffHard} {
		g := NewGame()
		g.difficulty = diff
		gs := newGameScreen(g)
		state := gs.buildState()
		if state.fps != diff.baseFPS() {
			t.Errorf("difficulty %v: fps = %f, want %f", diff, state.fps, diff.baseFPS())
		}
	}
}

// --- Snake food effect tests ---

func TestSnake_FoodEffect_Normal(t *testing.T) {
	g := newTestGameState()
	snake := newSnake(g)
	food := newFood()
	food.Emoji = '■'
	food.Position = Coordinates{4, 6} // Right next to head at (3,6)
	g.state.snake = snake
	g.state.food = food

	initialLen := len(snake.Body)
	snake.Direction = DirRight
	nextHead := snake.advance()
	snake.handleFoodEffect(g, nextHead)

	if g.state.score != 1 {
		t.Errorf("score = %d, want 1", g.state.score)
	}
	if len(snake.Body) != initialLen+1 {
		t.Errorf("body length = %d, want %d (should grow)", len(snake.Body), initialLen+1)
	}
}

func TestSnake_FoodEffect_Reward(t *testing.T) {
	g := newTestGameState()
	snake := newSnake(g)
	food := newFood()
	food.Emoji = 'R'
	g.state.snake = snake
	g.state.food = food

	initialLen := len(snake.Body)
	nextHead := snake.advance()
	snake.handleFoodEffect(g, nextHead)

	if g.state.score != 5 {
		t.Errorf("score = %d, want 5", g.state.score)
	}
	if len(snake.Body) != initialLen+1 {
		t.Errorf("body length = %d, want %d (should grow)", len(snake.Body), initialLen+1)
	}
}

func TestSnake_FoodEffect_Speed(t *testing.T) {
	g := newTestGameState()
	snake := newSnake(g)
	food := newFood()
	food.Emoji = 'S'
	g.state.snake = snake
	g.state.food = food

	initialLen := len(snake.Body)
	initialFPS := g.state.fps
	nextHead := snake.advance()
	snake.handleFoodEffect(g, nextHead)

	// Score should not change for 'S' food.
	if g.state.score != 0 {
		t.Errorf("score = %d, want 0 (S food gives no points)", g.state.score)
	}
	// Body should NOT change for 'S' food.
	if len(snake.Body) != initialLen {
		t.Errorf("body length = %d, want %d (S food should not grow)", len(snake.Body), initialLen)
	}
	// FPS should increase.
	if g.state.fps <= initialFPS {
		t.Errorf("fps = %f, should be > %f after S food", g.state.fps, initialFPS)
	}
}

// --- Color index tests ---

func TestColorIndex(t *testing.T) {
	cases := map[int]tl.Attr{
		10: tl.ColorWhite,
		12: tl.ColorRed,
		14: tl.ColorGreen,
		16: tl.ColorBlue,
		18: tl.ColorYellow,
		20: tl.ColorMagenta,
		22: tl.ColorCyan,
	}
	for idx, expected := range cases {
		got := colorIndex(idx)
		if got != expected {
			t.Errorf("colorIndex(%d) = %d, want %d", idx, got, expected)
		}
	}
}
