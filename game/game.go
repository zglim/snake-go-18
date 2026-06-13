package trisnake

import (
	"fmt"
	"os"
	"time"

	tl "github.com/JoelOtter/termloop"
	tb "github.com/nsf/termbox-go"
)

// Game is the central state container. It owns all mutable game state
// and exposes it through methods so that screens and entities cannot
// modify it directly.
type Game struct {
	engine      *tl.Game // termloop engine (private)
	difficulty  Difficulty
	snakeColor  tl.Attr
	arenaColor  tl.Attr
	colorTarget string // "Snake" or "Arena" – which object the color picker edits

	score      int
	fps        float64
	playScreen *PlayScreen // current active play screen (nil outside gameplay)
}

// StartGame is the package-level entry point called from main.
func StartGame() {
	g := &Game{
		difficulty:  DiffNormal,
		snakeColor:  colorPalette[0],
		arenaColor:  colorPalette[0],
		colorTarget: "Snake",
	}
	g.engine = tl.NewGame()
	g.engine.Screen().SetFps(10)
	g.showTitle()
	g.engine.Start()
}

// showTitle builds and displays the title screen.
func (g *Game) showTitle() {
	ts := newTitleScreen(g)
	g.engine.Screen().SetLevel(ts)
}

// showOptions builds and displays the options screen.
func (g *Game) showOptions() {
	opts := newOptionsScreen(g)
	g.engine.Screen().SetLevel(opts)
}

// startPlay initialises a new play session and switches to the play screen.
func (g *Game) startPlay() {
	g.score = 0
	g.fps = g.difficulty.baseFPS()

	ps := newPlayScreen(g)
	g.playScreen = ps
	g.engine.Screen().SetFps(g.fps)
	g.engine.Screen().SetLevel(ps)
}

// gameOver transitions to the game-over screen.
func (g *Game) gameOver() {
	gos := newGameOverScreen(g)
	g.engine.Screen().SetLevel(gos)
}

// restart resets the play state and returns to the play screen.
func (g *Game) restart() {
	ps := g.playScreen

	ps.RemoveEntity(ps.snake)
	ps.RemoveEntity(ps.food)

	ps.snake = newSnake(g, g.snakeColor)
	ps.food = newFood()

	g.score = 0
	g.fps = g.difficulty.baseFPS()

	ps.sidebar.updateScore(g.score)
	ps.sidebar.updateSpeed(g.fps)

	ps.AddEntity(ps.snake)
	ps.AddEntity(ps.food)

	g.engine.Screen().SetFps(g.fps)
	g.engine.Screen().SetLevel(ps)
}

// addScore increases the score by n and updates the sidebar display.
func (g *Game) addScore(n int) {
	g.score += n
	if g.playScreen != nil && g.playScreen.sidebar != nil {
		g.playScreen.sidebar.updateScore(g.score)
	}
}

// setSpeed changes the game speed and updates both the engine and sidebar.
func (g *Game) setSpeed(fps float64) {
	g.fps = fps
	g.engine.Screen().SetFps(g.fps)
	if g.playScreen != nil && g.playScreen.sidebar != nil {
		g.playScreen.sidebar.updateSpeed(g.fps)
	}
}

// applyRFoodBonus handles the speed bonus when eating R-food.
func (g *Game) applyRFoodBonus() {
	reduction := g.difficulty.rFoodSpeedReduction()
	floor := g.difficulty.minFPS()
	if g.fps-reduction > floor {
		g.setSpeed(g.fps - reduction)
	}
}

// applySFoodPenalty handles the speed penalty when eating S-food.
func (g *Game) applySFoodPenalty() {
	increase := g.difficulty.sFoodSpeedIncrease()
	g.setSpeed(g.fps + increase)
}

// quit terminates the game.
func (g *Game) quit() {
	tb.Close()
}

// saveHighScore appends the current score to HIGHSCORES.md.
func (g *Game) saveHighScore() {
	ts := time.Now().Format("01-02-2006 15:04:05")
	row := fmt.Sprintf("\n|%s|%d|%.0f|%s|  ", ts, g.score, g.fps, g.difficulty)

	f, err := os.OpenFile("HIGHSCORES.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(row)
}
