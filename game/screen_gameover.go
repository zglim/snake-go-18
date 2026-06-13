package trisnake

import (
	"fmt"
	"os"

	tl "github.com/JoelOtter/termloop"
	tb "github.com/nsf/termbox-go"
)

func newGameOverScreen(g *Game) *GameOverScreen {
	gos := &GameOverScreen{game: g}
	gos.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	logoBytes, _ := os.ReadFile("util/gameover-logo.txt")
	gos.logo = tl.NewEntityFromCanvas(10, 3, tl.CanvasFromString(string(logoBytes)))

	gos.finalStats = []*tl.Text{
		tl.NewText(10, 13, fmt.Sprintf("Score: %d", g.currentScore()), tl.ColorWhite, tl.ColorBlack),
		tl.NewText(10, 15, fmt.Sprintf("Speed: %.0f", g.currentFPS()), tl.ColorWhite, tl.ColorBlack),
		tl.NewText(10, 17, fmt.Sprintf("Difficulty: %s", g.difficulty.String()), tl.ColorWhite, tl.ColorBlack),
	}
	gos.optionsBackground = tl.NewRectangle(45, 12, 45, 7, tl.ColorWhite)
	gos.optionsText = []*tl.Text{
		tl.NewText(47, 13, "Press \"Home\" to restart!", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(47, 15, "Press \"Delete\" to quit!", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(47, 17, "Press \"Spacebar\" to save your score!", tl.ColorBlack, tl.ColorWhite),
	}

	for _, stat := range gos.finalStats {
		gos.AddEntity(stat)
	}
	gos.AddEntity(gos.logo)
	gos.AddEntity(gos.optionsBackground)
	for _, opt := range gos.optionsText {
		gos.AddEntity(opt)
	}

	return gos
}

// Tick handles key input on the game over screen.
func (gos *GameOverScreen) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyHome:
		gos.game.startNewGame()
	case tl.KeyDelete:
		tb.Close()
	case tl.KeySpace:
		gos.game.saveHighScore()
	}
}
