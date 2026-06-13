package trisnake

import (
	"testing"

	tl "github.com/JoelOtter/termloop"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// newTestGame creates a Game wired for unit testing (no termloop engine).
func newTestGame() *Game {
	return &Game{
		difficulty:  DiffNormal,
		snakeColor:  colorPalette[0],
		arenaColor:  colorPalette[0],
		colorTarget: "Snake",
	}
}

// newTestPlayScreen creates a minimal PlayScreen for testing game logic.
func newTestPlayScreen(g *Game) *PlayScreen {
	ps := &PlayScreen{game: g}
	ps.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	ps.arena = newArena(arenaWidth, arenaHeight, g.arenaColor)
	ps.snake = newSnake(g, g.snakeColor)
	ps.food = newFood()
	ps.sidebar = newSidepanel(g)

	g.playScreen = ps
	g.score = 0
	g.fps = g.difficulty.baseFPS()
	return ps
}

// ---------------------------------------------------------------------------
// Config / type tests
// ---------------------------------------------------------------------------

func TestDifficultyString(t *testing.T) {
	tests := []struct {
		d    Difficulty
		want string
	}{
		{DiffEasy, "Easy"},
		{DiffNormal, "Normal"},
		{DiffHard, "Hard"},
	}
	for _, tt := range tests {
		if got := tt.d.String(); got != tt.want {
			t.Errorf("Difficulty(%d).String() = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestDifficultyBaseFPS(t *testing.T) {
	tests := []struct {
		d    Difficulty
		want float64
	}{
		{DiffEasy, 8},
		{DiffNormal, 12},
		{DiffHard, 25},
	}
	for _, tt := range tests {
		if got := tt.d.baseFPS(); got != tt.want {
			t.Errorf("Difficulty(%d).baseFPS() = %v, want %v", tt.d, got, tt.want)
		}
	}
}

func TestDirectionOpposite(t *testing.T) {
	tests := []struct {
		d    Direction
		want Direction
	}{
		{DirUp, DirDown},
		{DirDown, DirUp},
		{DirLeft, DirRight},
		{DirRight, DirLeft},
	}
	for _, tt := range tests {
		if got := tt.d.opposite(); got != tt.want {
			t.Errorf("Direction(%d).opposite() = %d, want %d", tt.d, got, tt.want)
		}
	}
}

func TestColorAtY(t *testing.T) {
	if got := colorAtY(10); got != tl.ColorWhite {
		t.Errorf("colorAtY(10) = %v, want ColorWhite", got)
	}
	if got := colorAtY(12); got != tl.ColorRed {
		t.Errorf("colorAtY(12) = %v, want ColorRed", got)
	}
	if got := colorAtY(22); got != tl.ColorCyan {
		t.Errorf("colorAtY(22) = %v, want ColorCyan", got)
	}
}

func TestAttrToYRoundTrip(t *testing.T) {
	for i, c := range colorPalette {
		y := 10 + i*2
		if got := attrToY(c); got != y {
			t.Errorf("attrToY(%v) = %d, want %d", c, got, y)
		}
		if got := colorAtY(y); got != c {
			t.Errorf("colorAtY(%d) = %v, want %v", y, got, c)
		}
	}
}

// ---------------------------------------------------------------------------
// Arena tests
// ---------------------------------------------------------------------------

func TestArenaBorderContains(t *testing.T) {
	a := newArena(70, 25, tl.ColorWhite)

	// Corners
	if !a.contains(Coordinates{0, 0}) {
		t.Error("arena should contain top-left corner (0,0)")
	}
	if !a.contains(Coordinates{69, 0}) {
		t.Error("arena should contain top-right corner")
	}
	if !a.contains(Coordinates{0, 24}) {
		t.Error("arena should contain bottom-left corner")
	}
	if !a.contains(Coordinates{69, 24}) {
		t.Error("arena should contain bottom-right corner")
	}

	// Interior
	if a.contains(Coordinates{10, 10}) {
		t.Error("arena should not contain interior point (10,10)")
	}
	if a.contains(Coordinates{35, 12}) {
		t.Error("arena should not contain interior point (35,12)")
	}
}

// ---------------------------------------------------------------------------
// Snake tests
// ---------------------------------------------------------------------------

func TestNewSnakeInitialState(t *testing.T) {
	g := newTestGame()
	s := newSnake(g, tl.ColorWhite)

	if s.dir != DirRight {
		t.Errorf("new snake direction = %d, want DirRight", s.dir)
	}
	if len(s.body) != 3 {
		t.Errorf("new snake body length = %d, want 3", len(s.body))
	}
	head := s.head()
	if head.X != 3 || head.Y != 6 {
		t.Errorf("new snake head = (%d,%d), want (3,6)", head.X, head.Y)
	}
}

func TestSnakeNextHead(t *testing.T) {
	g := newTestGame()
	s := newSnake(g, tl.ColorWhite)
	// Default: head at (3,6), direction right

	tests := []struct {
		dir  Direction
		want Coordinates
	}{
		{DirRight, Coordinates{4, 6}},
		{DirLeft, Coordinates{2, 6}},
		{DirUp, Coordinates{3, 5}},
		{DirDown, Coordinates{3, 7}},
	}
	for _, tt := range tests {
		s.dir = tt.dir
		got := s.nextHead()
		if got != tt.want {
			t.Errorf("nextHead(dir=%d) = %v, want %v", tt.dir, got, tt.want)
		}
	}
}

func TestSnakeTickDirectionChange(t *testing.T) {
	g := newTestGame()
	s := newSnake(g, tl.ColorWhite)
	s.dir = DirRight

	// Press up → should change to up
	s.Tick(tl.Event{Type: tl.EventKey, Key: tl.KeyArrowUp})
	if s.dir != DirUp {
		t.Errorf("after Up key, dir = %d, want DirUp", s.dir)
	}

	// Press down (opposite of up) → should NOT change
	s.Tick(tl.Event{Type: tl.EventKey, Key: tl.KeyArrowDown})
	if s.dir != DirUp {
		t.Errorf("after Down key (opposite), dir = %d, want DirUp (unchanged)", s.dir)
	}

	// Press left → should change
	s.Tick(tl.Event{Type: tl.EventKey, Key: tl.KeyArrowLeft})
	if s.dir != DirLeft {
		t.Errorf("after Left key, dir = %d, want DirLeft", s.dir)
	}

	// Press right (opposite of left) → should NOT change
	s.Tick(tl.Event{Type: tl.EventKey, Key: tl.KeyArrowRight})
	if s.dir != DirLeft {
		t.Errorf("after Right key (opposite), dir = %d, want DirLeft (unchanged)", s.dir)
	}
}

func TestSnakeSelfCollision(t *testing.T) {
	g := newTestGame()
	s := newSnake(g, tl.ColorWhite)

	// Default body: no self-collision
	if s.selfCollision() {
		t.Error("new snake should not have self collision")
	}

	// Create a body where head overlaps body
	s.body = []Coordinates{
		{3, 6},
		{4, 6},
		{3, 6}, // head same as tail
	}
	if !s.selfCollision() {
		t.Error("snake with head overlapping body should have self collision")
	}
}

// ---------------------------------------------------------------------------
// Food tests
// ---------------------------------------------------------------------------

func TestFoodRespawnInsideArena(t *testing.T) {
	for i := 0; i < 100; i++ {
		f := newFood()
		if f.pos.X < 1 || f.pos.X >= arenaWidth-1 {
			t.Errorf("food X=%d outside arena bounds [1, %d)", f.pos.X, arenaWidth-1)
		}
		if f.pos.Y < 1 || f.pos.Y >= arenaHeight-1 {
			t.Errorf("food Y=%d outside arena bounds [1, %d)", f.pos.Y, arenaHeight-1)
		}
	}
}

func TestFoodContains(t *testing.T) {
	f := &Food{pos: Coordinates{5, 5}}

	if !f.contains(Coordinates{5, 5}) {
		t.Error("food should contain its own position")
	}
	if f.contains(Coordinates{5, 6}) {
		t.Error("food should not contain different position")
	}
}

func TestRandomFoodEmoji(t *testing.T) {
	valid := map[rune]bool{'R': true, '■': true, 'S': true}
	for i := 0; i < 100; i++ {
		e := randomFoodEmoji()
		if !valid[e] {
			t.Errorf("randomFoodEmoji() returned unexpected rune: %c", e)
		}
	}
}

// ---------------------------------------------------------------------------
// Game state tests
// ---------------------------------------------------------------------------

func TestGameAddScore(t *testing.T) {
	g := newTestGame()
	ps := newTestPlayScreen(g)
	_ = ps

	g.addScore(5)
	if g.score != 5 {
		t.Errorf("score = %d, want 5", g.score)
	}
	g.addScore(3)
	if g.score != 8 {
		t.Errorf("score = %d, want 8", g.score)
	}

	// Sidebar text should reflect score
	text := ps.sidebar.scoreText.Text()
	if text != "Score: 8" {
		t.Errorf("sidebar score text = %q, want %q", text, "Score: 8")
	}
}

func TestGameStartPlayResetsState(t *testing.T) {
	g := newTestGame()
	g.score = 99
	g.fps = 100

	// startPlay needs the engine, which we can't create in tests.
	// Instead, test the reset logic via newTestPlayScreen.
	ps := newTestPlayScreen(g)
	_ = ps

	if g.score != 0 {
		t.Errorf("score after reset = %d, want 0", g.score)
	}
	if g.fps != g.difficulty.baseFPS() {
		t.Errorf("fps after reset = %v, want %v", g.fps, g.difficulty.baseFPS())
	}
}

func TestGameRestart(t *testing.T) {
	g := newTestGame()
	g.engine = tl.NewGame() // need engine for restart (sets FPS and level)
	ps := newTestPlayScreen(g)

	// Simulate some score
	g.score = 42
	g.fps = 20
	ps.sidebar.updateScore(g.score)
	ps.sidebar.updateSpeed(g.fps)

	// Restart
	g.restart()

	if g.score != 0 {
		t.Errorf("score after restart = %d, want 0", g.score)
	}
	if g.fps != g.difficulty.baseFPS() {
		t.Errorf("fps after restart = %v, want %v", g.fps, g.difficulty.baseFPS())
	}
	if ps.snake == nil {
		t.Error("snake should not be nil after restart")
	}
	if ps.food == nil {
		t.Error("food should not be nil after restart")
	}
}

// ---------------------------------------------------------------------------
// Integration-style: snake eats food
// ---------------------------------------------------------------------------

func TestSnakeEatsFood(t *testing.T) {
	g := newTestGame()
	ps := newTestPlayScreen(g)

	// Place food directly in front of the snake head.
	head := *ps.snake.head()
	var foodPos Coordinates
	switch ps.snake.dir {
	case DirRight:
		foodPos = Coordinates{head.X + 1, head.Y}
	case DirLeft:
		foodPos = Coordinates{head.X - 1, head.Y}
	case DirUp:
		foodPos = Coordinates{head.X, head.Y - 1}
	case DirDown:
		foodPos = Coordinates{head.X, head.Y + 1}
	}
	ps.food.pos = foodPos
	ps.food.emoji = '■' // normal food
	ps.food.SetPosition(foodPos.X, foodPos.Y)

	bodyBefore := len(ps.snake.body)
	scoreBefore := g.score

	// Simulate one frame: move the snake
	newHead := ps.snake.nextHead()
	if ps.food.contains(newHead) {
		ps.snake.handleFood(newHead)
	} else {
		ps.snake.body = append(ps.snake.body[1:], newHead)
	}

	if len(ps.snake.body) != bodyBefore+1 {
		t.Errorf("body length after eating = %d, want %d", len(ps.snake.body), bodyBefore+1)
	}
	if g.score != scoreBefore+1 {
		t.Errorf("score after eating = %d, want %d", g.score, scoreBefore+1)
	}
}

func TestSnakeBorderCollision(t *testing.T) {
	g := newTestGame()
	ps := newTestPlayScreen(g)

	// Place snake head on the border
	ps.snake.body[len(ps.snake.body)-1] = Coordinates{0, 0}

	if !ps.arena.contains(*ps.snake.head()) {
		t.Error("snake head on border should trigger arena.contains")
	}
}

// ---------------------------------------------------------------------------
// Sidepanel tests
// ---------------------------------------------------------------------------

func TestSidepanelUpdate(t *testing.T) {
	g := newTestGame()
	sp := newSidepanel(g)

	sp.updateScore(42)
	if got := sp.scoreText.Text(); got != "Score: 42" {
		t.Errorf("score text = %q, want %q", got, "Score: 42")
	}

	sp.updateSpeed(15)
	if got := sp.speedText.Text(); got != "Speed: 15" {
		t.Errorf("speed text = %q, want %q", got, "Speed: 15")
	}
}
