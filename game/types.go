package trisnake

import tl "github.com/JoelOtter/termloop"

// Direction represents the movement direction of the snake.
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// Difficulty represents the game difficulty level.
type Difficulty int

const (
	DiffEasy Difficulty = iota
	DiffNormal
	DiffHard
)

// String returns the display name of the difficulty.
func (d Difficulty) String() string {
	switch d {
	case DiffEasy:
		return "Easy"
	case DiffNormal:
		return "Normal"
	case DiffHard:
		return "Hard"
	default:
		return "Normal"
	}
}

// baseFPS returns the starting FPS for the given difficulty.
func (d Difficulty) baseFPS() float64 {
	switch d {
	case DiffEasy:
		return 8
	case DiffNormal:
		return 12
	case DiffHard:
		return 25
	default:
		return 12
	}
}

// ColorTarget selects which game element's color is being configured.
type ColorTarget int

const (
	TargetSnake ColorTarget = iota
	TargetArena
)

// String returns the display name of the color target.
func (ct ColorTarget) String() string {
	switch ct {
	case TargetSnake:
		return "Snake"
	case TargetArena:
		return "Arena"
	default:
		return "Snake"
	}
}

// colorIndex maps a numeric selector position to a termloop color attribute.
func colorIndex(idx int) tl.Attr {
	switch idx {
	case 10:
		return tl.ColorWhite
	case 12:
		return tl.ColorRed
	case 14:
		return tl.ColorGreen
	case 16:
		return tl.ColorBlue
	case 18:
		return tl.ColorYellow
	case 20:
		return tl.ColorMagenta
	case 22:
		return tl.ColorCyan
	default:
		return tl.ColorDefault
	}
}

// Coordinates represents an (X, Y) position on the game grid.
type Coordinates struct {
	X int
	Y int
}

// GameState holds all mutable runtime state for an active game round.
// It is reset each time a new game round starts.
type GameState struct {
	score    int
	fps      float64
	snake    *Snake
	food     *Food
	arena    *Arena
	sidepanel *Sidepanel
}

// ColorConfig holds the color selection state shared across screens.
type ColorConfig struct {
	Target    ColorTarget
	SnakeIdx  int
	ArenaIdx  int
}

// ActiveColorIdx returns the selector index for the currently targeted element.
func (cc *ColorConfig) ActiveIdx() int {
	if cc.Target == TargetArena {
		return cc.ArenaIdx
	}
	return cc.SnakeIdx
}

// SetActiveIdx sets the selector index for the currently targeted element.
func (cc *ColorConfig) SetActiveIdx(idx int) {
	if cc.Target == TargetArena {
		cc.ArenaIdx = idx
	} else {
		cc.SnakeIdx = idx
	}
}

// Snake is the player-controlled entity that moves on the game grid.
type Snake struct {
	*tl.Entity
	game      *Game
	Direction Direction
	Body      []Coordinates
}

// Food is a collectible item that spawns at random positions inside the arena.
type Food struct {
	*tl.Entity
	Position Coordinates
	Emoji    rune
}

// Arena defines the play area and its border walls.
type Arena struct {
	*tl.Entity
	game   *Game
	Width  int
	Height int
	Border map[Coordinates]int
}

// Sidepanel displays score, speed and instructions next to the arena.
type Sidepanel struct {
	Background     *tl.Rectangle
	Instructions   []string
	ScoreText      *tl.Text
	SpeedText      *tl.Text
	DifficultyText *tl.Text
}

// TitleScreen is the initial screen shown when the game starts.
type TitleScreen struct {
	tl.Level
	game *Game
	logo *tl.Entity
	options []*tl.Text
}

// OptionsScreen allows the player to configure difficulty and colors.
type OptionsScreen struct {
	tl.Level
	game *Game

	startText             *tl.Text
	difficultyLabel       *tl.Text
	difficultyBackground  *tl.Rectangle
	difficultyOptions     []*tl.Text

	colorLabel            *tl.Text
	objectBackground      *tl.Rectangle
	colorObjectOptions    []*tl.Text

	colorPanelBackground  *tl.Rectangle
	colorPanelOptions     []string
	colorSelectedIcon     *tl.Text
}

// GameScreen is the main play area level where the snake game runs.
type GameScreen struct {
	tl.Level
	game *Game
}

// GameOverScreen shows the final stats and restart/quit options.
type GameOverScreen struct {
	tl.Level
	game              *Game
	logo              *tl.Entity
	finalStats        []*tl.Text
	optionsBackground *tl.Rectangle
	optionsText       []*tl.Text
}
