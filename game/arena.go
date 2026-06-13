package trisnake

import tl "github.com/JoelOtter/termloop"

// Arena is the bordered play area.
type Arena struct {
	*tl.Entity
	width  int
	height int
	border map[Coordinates]bool
	color  tl.Attr
}

// newArena creates an arena with the given dimensions and border color.
// Width and height are decreased by 1 so that the border corners align.
func newArena(w, h int, color tl.Attr) *Arena {
	a := &Arena{
		Entity: tl.NewEntity(1, 1, 1, 1),
		width:  w - 1,
		height: h - 1,
		border: make(map[Coordinates]bool),
		color:  color,
	}

	// Top and bottom borders
	for x := 0; x < a.width; x++ {
		a.border[Coordinates{x, 0}] = true
		a.border[Coordinates{x, a.height}] = true
	}
	// Left and right borders
	for y := 0; y <= a.height; y++ {
		a.border[Coordinates{0, y}] = true
		a.border[Coordinates{a.width, y}] = true
	}

	return a
}

// contains returns true if the given coordinates are on the arena border.
func (a *Arena) contains(c Coordinates) bool {
	return a.border[c]
}

// Draw renders the arena border on screen.
func (a *Arena) Draw(screen *tl.Screen) {
	for c := range a.border {
		screen.RenderCell(c.X, c.Y, &tl.Cell{
			Bg: a.color,
		})
	}
}
