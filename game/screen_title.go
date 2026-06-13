package trisnake

import (
	"os"

	tl "github.com/JoelOtter/termloop"
)

func newTitleScreen(g *Game) *TitleScreen {
	ts := &TitleScreen{game: g}
	ts.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	logoBytes, _ := os.ReadFile("util/titlescreen-logo.txt")
	ts.logo = tl.NewEntityFromCanvas(10, 3, tl.CanvasFromString(string(logoBytes)))

	ts.options = []*tl.Text{
		tl.NewText(10, 15, "Press ENTER to start!", tl.ColorWhite, tl.ColorBlack),
		tl.NewText(10, 17, "Press INSERT for options!", tl.ColorWhite, tl.ColorBlack),
	}

	ts.AddEntity(ts.logo)
	for _, opt := range ts.options {
		ts.AddEntity(opt)
	}

	return ts
}

// Tick handles key input on the title screen.
func (ts *TitleScreen) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyEnter:
		ts.game.startNewGame()
	case tl.KeyInsert:
		ts.game.showOptions()
	}
}
