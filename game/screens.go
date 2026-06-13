package trisnake

import (
	"fmt"
	"os"

	tl "github.com/JoelOtter/termloop"
)

// ---------------------------------------------------------------------------
// TitleScreen
// ---------------------------------------------------------------------------

// TitleScreen is the first screen shown when the game starts.
type TitleScreen struct {
	tl.Level
	game    *Game
	logo    *tl.Entity
	options []*tl.Text
}

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

// Tick handles title screen input.
func (ts *TitleScreen) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyEnter:
		ts.game.startPlay()
	case tl.KeyInsert:
		ts.game.showOptions()
	}
}

// ---------------------------------------------------------------------------
// OptionsScreen
// ---------------------------------------------------------------------------

// OptionsScreen lets the player configure difficulty and colors.
type OptionsScreen struct {
	tl.Level
	game *Game

	startText      *tl.Text
	diffText       *tl.Text
	diffBg         *tl.Rectangle
	diffOptions    []*tl.Text
	colorText      *tl.Text
	colorBg        *tl.Rectangle
	objectOptions  []*tl.Text
	panelBg        *tl.Rectangle
	panelLabels    []string
	colorIcon      *tl.Text

	// Color cursor Y-positions for snake and arena.
	snakeColorY int
	arenaColorY int
}

func newOptionsScreen(g *Game) *OptionsScreen {
	opts := &OptionsScreen{
		game:        g,
		snakeColorY: 10, // default: White
		arenaColorY: 10,
	}
	opts.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	// Layout elements
	opts.diffBg = tl.NewRectangle(5, 3, 33, 10, tl.ColorWhite)
	opts.colorBg = tl.NewRectangle(5, 15, 33, 9, tl.ColorWhite)
	opts.panelBg = tl.NewRectangle(43, 3, 33, 21, tl.ColorWhite)

	opts.startText = tl.NewText(2, 1, "Press Enter to start!", tl.ColorWhite, tl.ColorBlack)
	opts.diffText = tl.NewText(6, 4, fmt.Sprintf("Current difficulty: %s", g.difficulty), tl.ColorBlack, tl.ColorWhite)
	opts.colorText = tl.NewText(44, 4, fmt.Sprintf("Current Object: %s", g.colorTarget), tl.ColorBlack, tl.ColorWhite)

	// Position the color cursor on the correct row based on current color.
	opts.syncColorCursorFromConfig()
	cursorY := opts.activeColorY()
	opts.colorIcon = tl.NewText(73, cursorY, "■", tl.ColorBlack, tl.ColorWhite)

	opts.panelLabels = []string{
		"Use ↑ ↓ to change colors",
		"White", "Red", "Green", "Blue", "Yellow", "Magenta", "Cyan",
	}

	opts.diffOptions = []*tl.Text{
		tl.NewText(6, 7, "Press F1 for Easy (8 speed)", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(6, 9, "Press F2 for Normal (12 speed)", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(6, 11, "Press F3 for Hard (25 speed)", tl.ColorBlack, tl.ColorWhite),
	}

	opts.objectOptions = []*tl.Text{
		tl.NewText(6, 16, "Press F4 for Snake (Colors)", tl.ColorBlack, tl.ColorWhite),
		tl.NewText(6, 18, "Press F6 for Arena (Colors)", tl.ColorBlack, tl.ColorWhite),
	}

	// Add entities
	opts.AddEntity(opts.diffBg)
	opts.AddEntity(opts.colorBg)
	opts.AddEntity(opts.panelBg)
	opts.AddEntity(opts.diffText)
	opts.AddEntity(opts.colorText)
	opts.AddEntity(opts.colorIcon)
	opts.AddEntity(opts.startText)

	for _, d := range opts.diffOptions {
		opts.AddEntity(d)
	}

	// Color panel labels
	y := 6
	for _, label := range opts.panelLabels {
		y += 2
		opts.AddEntity(tl.NewText(44, y, label, tl.ColorBlack, tl.ColorWhite))
	}

	for _, o := range opts.objectOptions {
		opts.AddEntity(o)
	}

	return opts
}

// syncColorCursorFromConfig sets the cursor Y positions from the current color config.
func (opts *OptionsScreen) syncColorCursorFromConfig() {
	opts.snakeColorY = attrToY(opts.game.snakeColor)
	opts.arenaColorY = attrToY(opts.game.arenaColor)
}

// activeColorY returns the cursor Y for the currently targeted object.
func (opts *OptionsScreen) activeColorY() int {
	if opts.game.colorTarget == "Arena" {
		return opts.arenaColorY
	}
	return opts.snakeColorY
}

// setActiveColorY sets the cursor Y for the currently targeted object.
func (opts *OptionsScreen) setActiveColorY(y int) {
	if opts.game.colorTarget == "Arena" {
		opts.arenaColorY = y
		opts.game.arenaColor = colorAtY(y)
	} else {
		opts.snakeColorY = y
		opts.game.snakeColor = colorAtY(y)
	}
}

// Tick handles options screen input.
func (opts *OptionsScreen) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyF1:
		opts.game.difficulty = DiffEasy
		opts.diffText.SetText(fmt.Sprintf("Current difficulty: %s", opts.game.difficulty))
	case tl.KeyF2:
		opts.game.difficulty = DiffNormal
		opts.diffText.SetText(fmt.Sprintf("Current difficulty: %s", opts.game.difficulty))
	case tl.KeyF3:
		opts.game.difficulty = DiffHard
		opts.diffText.SetText(fmt.Sprintf("Current difficulty: %s", opts.game.difficulty))

	case tl.KeyF4:
		opts.game.colorTarget = "Snake"
		opts.colorText.SetText(fmt.Sprintf("Current Object: %s", opts.game.colorTarget))
		opts.colorIcon.SetPosition(73, opts.activeColorY())
	case tl.KeyF5:
		opts.game.colorTarget = "Food"
		opts.colorText.SetText(fmt.Sprintf("Current Object: %s", opts.game.colorTarget))
	case tl.KeyF6:
		opts.game.colorTarget = "Arena"
		opts.colorText.SetText(fmt.Sprintf("Current Object: %s", opts.game.colorTarget))
		opts.colorIcon.SetPosition(73, opts.activeColorY())

	case tl.KeyArrowUp:
		y := opts.activeColorY()
		if y <= 10 {
			return
		}
		y -= 2
		opts.setActiveColorY(y)
		opts.colorIcon.SetPosition(73, y)
	case tl.KeyArrowDown:
		y := opts.activeColorY()
		if y >= 22 {
			return
		}
		y += 2
		opts.setActiveColorY(y)
		opts.colorIcon.SetPosition(73, y)

	case tl.KeyEnter:
		opts.game.startPlay()
	}
}

// attrToY converts a termloop Attr back to the Y-position index used in the color picker.
func attrToY(attr tl.Attr) int {
	for i, c := range colorPalette {
		if c == attr {
			return 10 + i*2
		}
	}
	return 10 // default: White
}

// ---------------------------------------------------------------------------
// PlayScreen
// ---------------------------------------------------------------------------

// PlayScreen is the main gameplay screen.
type PlayScreen struct {
	tl.Level
	game    *Game
	snake   *Snake
	food    *Food
	arena   *Arena
	sidebar *Sidepanel
}

func newPlayScreen(g *Game) *PlayScreen {
	ps := &PlayScreen{game: g}
	ps.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	ps.arena = newArena(arenaWidth, arenaHeight, g.arenaColor)
	ps.snake = newSnake(g, g.snakeColor)
	ps.food = newFood()
	ps.sidebar = newSidepanel(g)

	// Add entities in the correct z-order.
	ps.AddEntity(ps.food)
	ps.AddEntity(ps.sidebar.bg)
	ps.AddEntity(ps.sidebar.scoreText)
	ps.AddEntity(ps.sidebar.speedText)
	ps.AddEntity(ps.sidebar.diffText)
	ps.AddEntity(ps.snake)
	ps.AddEntity(ps.arena)

	// Instruction labels
	y := 7
	for _, instr := range ps.sidebar.instructions {
		y += 2
		ps.AddEntity(tl.NewText(72, y, instr, tl.ColorBlack, tl.ColorWhite))
	}

	return ps
}

// PlayScreen does not override Tick; BaseLevel.Tick forwards events to all
// entities (Snake handles direction input via its own Tick).

// ---------------------------------------------------------------------------
// GameOverScreen
// ---------------------------------------------------------------------------

// GameOverScreen is shown when the player dies.
type GameOverScreen struct {
	tl.Level
	game *Game
}

func newGameOverScreen(g *Game) *GameOverScreen {
	gos := &GameOverScreen{game: g}
	gos.Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorBlack})

	logoBytes, _ := os.ReadFile("util/gameover-logo.txt")
	logo := tl.NewEntityFromCanvas(10, 3, tl.CanvasFromString(string(logoBytes)))

	stats := []*tl.Text{
		tl.NewText(10, 13, fmt.Sprintf("Score: %d", g.score), tl.ColorWhite, tl.ColorBlack),
		tl.NewText(10, 15, fmt.Sprintf("Speed: %.0f", g.fps), tl.ColorWhite, tl.ColorBlack),
		tl.NewText(10, 17, fmt.Sprintf("Difficulty: %s", g.difficulty), tl.ColorWhite, tl.ColorBlack),
	}

	optionBg := tl.NewRectangle(45, 12, 45, 7, tl.ColorWhite)
	options := []*tl.Text{
		tl.NewText(47, 13, `Press "Home" to restart!`, tl.ColorBlack, tl.ColorWhite),
		tl.NewText(47, 15, `Press "Delete" to quit!`, tl.ColorBlack, tl.ColorWhite),
		tl.NewText(47, 17, `Press "Spacebar" to save your score!`, tl.ColorBlack, tl.ColorWhite),
	}

	for _, s := range stats {
		gos.AddEntity(s)
	}
	gos.AddEntity(logo)
	gos.AddEntity(optionBg)
	for _, o := range options {
		gos.AddEntity(o)
	}

	return gos
}

// Tick handles game-over screen input.
func (gos *GameOverScreen) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyHome:
		gos.game.restart()
	case tl.KeyDelete:
		gos.game.quit()
	case tl.KeySpace:
		gos.game.saveHighScore()
	}
}

// ---------------------------------------------------------------------------
// Sidepanel
// ---------------------------------------------------------------------------

// Sidepanel displays score, speed, difficulty, and instructions during gameplay.
type Sidepanel struct {
	bg           *tl.Rectangle
	scoreText    *tl.Text
	speedText    *tl.Text
	diffText     *tl.Text
	instructions []string
}

func newSidepanel(g *Game) *Sidepanel {
	sp := &Sidepanel{
		instructions: []string{
			"Instructions:",
			"Use ← → ↑ ↓ to move the snake around",
			"Pick up the food to grow bigger",
			"■: 1 point/growth",
			"R: 5 points (removes some speed!)",
			"S: 1 point (increased speed!!)",
		},
	}

	sp.bg = tl.NewRectangle(71, 0, 45, 25, tl.ColorWhite)
	sp.scoreText = tl.NewText(72, 1, fmt.Sprintf("Score: %d", g.score), tl.ColorBlack, tl.ColorWhite)
	sp.speedText = tl.NewText(72, 3, fmt.Sprintf("Speed: %.0f", g.fps), tl.ColorBlack, tl.ColorWhite)
	sp.diffText = tl.NewText(72, 5, fmt.Sprintf("Difficulty: %s", g.difficulty), tl.ColorBlack, tl.ColorWhite)

	return sp
}

// updateScore refreshes the score display.
func (sp *Sidepanel) updateScore(score int) {
	sp.scoreText.SetText(fmt.Sprintf("Score: %d", score))
}

// updateSpeed refreshes the speed display.
func (sp *Sidepanel) updateSpeed(fps float64) {
	sp.speedText.SetText(fmt.Sprintf("Speed: %.0f", fps))
}
