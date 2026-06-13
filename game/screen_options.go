package trisnake

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
)

const (
	colorMinIdx = 10
	colorMaxIdx = 22
)

func newOptionsScreen(g *Game) *OptionsScreen {
	os := &OptionsScreen{game: g}
	os.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	// Background panels.
	os.colorPanelBackground = tl.NewRectangle(43, 3, 33, 21, tl.ColorWhite)
	os.difficultyBackground = tl.NewRectangle(5, 3, 33, 10, tl.ColorWhite)
	os.objectBackground = tl.NewRectangle(5, 15, 33, 9, tl.ColorWhite)

	// Labels.
	os.startText = tl.NewText(2, 1, "Press Enter to start!", tl.ColorWhite, tl.ColorBlack)
	os.difficultyLabel = tl.NewText(6, 4,
		fmt.Sprintf("Current difficulty: %s", g.difficulty.String()),
		tl.ColorBlack, tl.ColorWhite)
	os.colorLabel = tl.NewText(44, 4,
		fmt.Sprintf("Current Object: %s", g.colors.Target.String()),
		tl.ColorBlack, tl.ColorWhite)
	os.colorSelectedIcon = tl.NewText(73, g.colors.ActiveIdx(), "■", tl.ColorBlack, tl.ColorWhite)

	// Color panel options.
	os.colorPanelOptions = []string{
		"Use ↑ ↓ to change colors",
		"White", "Red", "Green", "Blue", "Yellow", "Magenta", "Cyan",
	}

	// Difficulty options.
	os.difficultyOptions = []*tl.Text{
		tl.NewText(6, 7, "Press F1 for Easy (8 speed)", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(6, 9, "Press F2 for Normal (12 speed)", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(6, 11, "Press F3 for Hard (25 speed)", tl.ColorBlack, tl.ColorWhite),
	}

	// Color object options.
	os.colorObjectOptions = []*tl.Text{
		tl.NewText(6, 16, "Press F4 for Snake (Colors)", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(6, 18, "Press F6 for Arena (Colors)", tl.ColorBlack, tl.ColorWhite),
	}

	// Register entities.
	os.AddEntity(os.difficultyBackground)
	os.AddEntity(os.colorPanelBackground)
	os.AddEntity(os.objectBackground)
	os.AddEntity(os.difficultyLabel)
	os.AddEntity(os.colorLabel)
	os.AddEntity(os.colorSelectedIcon)
	os.AddEntity(os.startText)

	for _, opt := range os.difficultyOptions {
		os.AddEntity(opt)
	}

	y := 6
	for _, label := range os.colorPanelOptions {
		y += 2
		os.AddEntity(tl.NewText(44, y, label, tl.ColorBlack, tl.ColorWhite))
	}
	for _, opt := range os.colorObjectOptions {
		os.AddEntity(opt)
	}

	return os
}

// Tick handles key input on the options screen.
func (os *OptionsScreen) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyF1:
		os.game.difficulty = DiffEasy
		os.difficultyLabel.SetText(fmt.Sprintf("Current difficulty: %s", os.game.difficulty.String()))

	case tl.KeyF2:
		os.game.difficulty = DiffNormal
		os.difficultyLabel.SetText(fmt.Sprintf("Current difficulty: %s", os.game.difficulty.String()))

	case tl.KeyF3:
		os.game.difficulty = DiffHard
		os.difficultyLabel.SetText(fmt.Sprintf("Current difficulty: %s", os.game.difficulty.String()))

	case tl.KeyF4:
		os.game.colors.Target = TargetSnake
		os.colorLabel.SetText(fmt.Sprintf("Current object: %s", os.game.colors.Target.String()))
		os.colorSelectedIcon.SetPosition(73, os.game.colors.ActiveIdx())

	case tl.KeyF5:
		// Food color target is reserved for future use.

	case tl.KeyF6:
		os.game.colors.Target = TargetArena
		os.colorLabel.SetText(fmt.Sprintf("Current object: %s", os.game.colors.Target.String()))
		os.colorSelectedIcon.SetPosition(73, os.game.colors.ActiveIdx())

	case tl.KeyArrowUp:
		idx := os.game.colors.ActiveIdx()
		if idx > colorMinIdx {
			idx -= 2
			os.game.colors.SetActiveIdx(idx)
			os.colorSelectedIcon.SetPosition(73, idx)
		}

	case tl.KeyArrowDown:
		idx := os.game.colors.ActiveIdx()
		if idx < colorMaxIdx {
			idx += 2
			os.game.colors.SetActiveIdx(idx)
			os.colorSelectedIcon.SetPosition(73, idx)
		}

	case tl.KeyEnter:
		os.game.startNewGame()
	}
}
