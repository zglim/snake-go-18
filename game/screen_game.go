package trisnake

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
)

// arenaWidth and arenaHeight define the play area dimensions.
const (
	arenaWidth  = 70
	arenaHeight = 25
)

func newGameScreen(g *Game) *GameScreen {
	gs := &GameScreen{game: g}
	gs.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})
	return gs
}

// buildState creates a fresh GameState for a new game round.
// It initializes snake, arena, food, sidepanel and resets score/FPS.
func (gs *GameScreen) buildState() *GameState {
	state := &GameState{
		score: 0,
		fps:   gs.game.difficulty.baseFPS(),
	}
	state.snake = newSnake(gs.game)
	state.arena = newArena(gs.game, arenaWidth, arenaHeight)
	state.food = newFood()
	state.sidepanel = newSidepanel(state, gs.game.difficulty)
	return state
}

// newSidepanel creates the sidepanel UI next to the arena.
func newSidepanel(state *GameState, diff Difficulty) *Sidepanel {
	sp := &Sidepanel{}
	sp.Instructions = []string{
		"Instructions:",
		"Use ← → ↑ ↓ to move the snake around",
		"Pick up the food to grow bigger",
		"■: 1 point/growth",
		"R: 5 points (removes some speed!)",
		"S: 1 point (increased speed!!)",
	}

	sp.Background = tl.NewRectangle(71, 0, 45, arenaHeight, tl.ColorWhite)
	sp.ScoreText = tl.NewText(72, 1, fmt.Sprintf("Score: %d", state.score), tl.ColorBlack, tl.ColorWhite)
	sp.SpeedText = tl.NewText(72, 3, fmt.Sprintf("Speed: %.0f", state.fps), tl.ColorBlack, tl.ColorWhite)
	sp.DifficultyText = tl.NewText(72, 5, fmt.Sprintf("Difficulty: %s", diff.String()), tl.ColorBlack, tl.ColorWhite)

	return sp
}
