package trisnake

import tl "github.com/JoelOtter/termloop"

// Direction represents the snake's movement direction.
type Direction int

const (
	DirUp    Direction = iota
	DirDown
	DirLeft
	DirRight
)

// opposite returns the direction opposite to d.
func (d Direction) opposite() Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	default:
		return d
	}
}

// Difficulty represents the game difficulty setting.
type Difficulty int

const (
	DiffEasy   Difficulty = iota
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

// baseFPS returns the starting FPS for this difficulty.
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

// minFPS returns the minimum FPS floor for the R-food speed bonus.
func (d Difficulty) minFPS() float64 {
	switch d {
	case DiffEasy:
		return 8
	case DiffNormal:
		return 12
	case DiffHard:
		return 20
	default:
		return 12
	}
}

// rFoodSpeedReduction returns how much R-food reduces FPS for this difficulty.
func (d Difficulty) rFoodSpeedReduction() float64 {
	switch d {
	case DiffEasy:
		return 3
	case DiffNormal:
		return 2
	case DiffHard:
		return 1
	default:
		return 2
	}
}

// sFoodSpeedIncrease returns how much S-food increases FPS for this difficulty.
func (d Difficulty) sFoodSpeedIncrease() float64 {
	switch d {
	case DiffEasy:
		return 1
	case DiffNormal:
		return 3
	case DiffHard:
		return 5
	default:
		return 3
	}
}

// Coordinates represents an (X, Y) position on the game grid.
type Coordinates struct {
	X int
	Y int
}

// colorPalette maps color indices to termloop attributes.
// The order matches the options screen layout: White, Red, Green, Blue, Yellow, Magenta, Cyan.
var colorPalette = []tl.Attr{
	tl.ColorWhite,
	tl.ColorRed,
	tl.ColorGreen,
	tl.ColorBlue,
	tl.ColorYellow,
	tl.ColorMagenta,
	tl.ColorCyan,
}

// colorAtY returns the termloop color attribute for a given Y-position index in the options screen.
// The options screen uses Y positions 10, 12, 14, 16, 18, 20, 22 for the 7 colors.
func colorAtY(y int) tl.Attr {
	idx := (y - 10) / 2
	if idx < 0 || idx >= len(colorPalette) {
		return tl.ColorWhite
	}
	return colorPalette[idx]
}

// Arena dimensions (constants used across the game).
const (
	arenaWidth  = 70
	arenaHeight = 25
)
