package trisnake

import (
	"fmt"
	"os"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// Game is the central state owner. It holds the termloop game instance,
// persistent configuration (difficulty, colors), and the active round state.
// All mutable state changes go through Game methods.
type Game struct {
	loop       *tl.Game
	difficulty Difficulty
	colors     ColorConfig
	state      *GameState
}

// NewGame creates a new Game with default configuration.
func NewGame() *Game {
	return &Game{
		difficulty: DiffNormal,
		colors: ColorConfig{
			Target:   TargetSnake,
			SnakeIdx: 10,
			ArenaIdx: 10,
		},
	}
}

// Run starts the game loop from the title screen.
func (g *Game) Run() {
	g.loop = tl.NewGame()
	g.loop.Screen().SetFps(10)
	g.loop.Screen().SetLevel(newTitleScreen(g))
	g.loop.Start()
}

// startNewGame creates a fresh game round and switches to the game screen.
// Used both for the initial start and for restarts.
func (g *Game) startNewGame() {
	screen := newGameScreen(g)
	g.state = screen.buildState()

	// Register all drawable entities with termloop.
	screen.AddEntity(g.state.food)
	screen.AddEntity(g.state.sidepanel.Background)
	screen.AddEntity(g.state.sidepanel.ScoreText)
	screen.AddEntity(g.state.sidepanel.SpeedText)
	screen.AddEntity(g.state.sidepanel.DifficultyText)
	screen.AddEntity(g.state.snake)
	screen.AddEntity(g.state.arena)

	y := 7
	for _, text := range g.state.sidepanel.Instructions {
		y += 2
		label := tl.NewText(72, y, text, tl.ColorBlack, tl.ColorWhite)
		screen.AddEntity(label)
	}

	g.loop.Screen().SetFps(g.state.fps)
	g.loop.Screen().SetLevel(screen)
}

// showGameOver switches to the game over screen.
func (g *Game) showGameOver() {
	g.loop.Screen().SetLevel(newGameOverScreen(g))
}

// showOptions switches to the options screen.
func (g *Game) showOptions() {
	g.loop.Screen().SetLevel(newOptionsScreen(g))
}

// currentFPS returns the current game FPS, or 0 if no active round.
func (g *Game) currentFPS() float64 {
	if g.state == nil {
		return 0
	}
	return g.state.fps
}

// currentScore returns the current game score, or 0 if no active round.
func (g *Game) currentScore() int {
	if g.state == nil {
		return 0
	}
	return g.state.score
}

// addScore increases the score and updates the sidepanel display.
func (g *Game) addScore(amount int) {
	if g.state == nil {
		return
	}
	g.state.score += amount
	g.state.sidepanel.ScoreText.SetText(fmt.Sprintf("Score: %d", g.state.score))
}

// adjustFPS changes the game speed by the given delta and updates the display.
func (g *Game) adjustFPS(delta float64) {
	if g.state == nil {
		return
	}
	g.state.fps += delta
	if g.loop != nil {
		g.loop.Screen().SetFps(g.state.fps)
	}
	g.state.sidepanel.SpeedText.SetText(fmt.Sprintf("Speed: %.0f", g.state.fps))
}

// setFPS sets the game speed to an exact value and updates the display.
func (g *Game) setFPS(fps float64) {
	if g.state == nil {
		return
	}
	g.state.fps = fps
	if g.loop != nil {
		g.loop.Screen().SetFps(g.state.fps)
	}
	g.state.sidepanel.SpeedText.SetText(fmt.Sprintf("Speed: %.0f", g.state.fps))
}

// moveFood repositions the food to a new random location inside the arena.
func (g *Game) moveFood() {
	if g.state == nil {
		return
	}
	g.state.food.Randomize()
}

// saveHighScore appends the current round's stats to HIGHSCORES.md.
func (g *Game) saveHighScore() {
	if g.state == nil {
		return
	}
	row := fmt.Sprintf("\n|%s|%d|%.0f|%s|  ",
		time.Now().Format("01-02-2006 15:04:05"),
		g.state.score,
		g.state.fps,
		g.difficulty.String(),
	)
	f, err := os.OpenFile("HIGHSCORES.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(row)
}
