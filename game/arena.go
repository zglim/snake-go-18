package trisnake

import tl "github.com/JoelOtter/termloop"

// newArena creates the play area with border walls at the given dimensions.
// Width and height are decreased by 1 to place corner cells correctly.
func newArena(g *Game, w, h int) *Arena {
	arena := &Arena{
		Entity: tl.NewEntity(1, 1, 1, 1),
		game:   g,
		Width:  w - 1,
		Height: h - 1,
		Border: make(map[Coordinates]int),
	}

	// Top and bottom borders.
	for x := 0; x < arena.Width; x++ {
		arena.Border[Coordinates{x, 0}] = 1
		arena.Border[Coordinates{x, arena.Height}] = 1
	}

	// Left and right borders.
	for y := 0; y < arena.Height+1; y++ {
		arena.Border[Coordinates{0, y}] = 1
		arena.Border[Coordinates{arena.Width, y}] = 1
	}

	return arena
}

// Contains returns true if the given coordinates are on the arena border.
func (arena *Arena) Contains(c Coordinates) bool {
	_, exists := arena.Border[c]
	return exists
}

// Draw renders the arena border using the configured arena color.
func (arena *Arena) Draw(screen *tl.Screen) {
	color := colorIndex(arena.game.colors.ArenaIdx)
	for c := range arena.Border {
		screen.RenderCell(c.X, c.Y, &tl.Cell{
			Bg: color,
		})
	}
}
