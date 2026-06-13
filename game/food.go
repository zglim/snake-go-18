package trisnake

import (
	"math/rand"

	tl "github.com/JoelOtter/termloop"
)

// Food is a collectible item that spawns at random positions inside the arena.
type Food struct {
	*tl.Entity
	pos   Coordinates
	emoji rune
}

// newFood creates a food item at a random position.
func newFood() *Food {
	f := &Food{
		Entity: tl.NewEntity(1, 1, 1, 1),
	}
	f.respawn()
	return f
}

// respawn moves the food to a new random position inside the arena.
func (f *Food) respawn() {
	f.pos.X = randInRange(1, arenaWidth-1)
	f.pos.Y = randInRange(1, arenaHeight-1)
	f.emoji = randomFoodEmoji()
	f.SetPosition(f.pos.X, f.pos.Y)
}

// Draw renders the food on screen.
func (f *Food) Draw(screen *tl.Screen) {
	screen.RenderCell(f.pos.X, f.pos.Y, &tl.Cell{
		Ch: f.emoji,
	})
}

// contains returns true if the given coordinates match the food position.
func (f *Food) contains(c Coordinates) bool {
	return c.X == f.pos.X && c.Y == f.pos.Y
}

// foodEmojis defines the weighted food spawn table.
// 'R' = bonus dish, 'S' = skull penalty, '■' = normal (×10).
var foodEmojis = []rune{
	'R',
	'■', '■', '■', '■', '■', '■', '■', '■', '■', '■',
	'S',
}

// randomFoodEmoji picks a random food type from the weighted table.
func randomFoodEmoji() rune {
	return foodEmojis[rand.Intn(len(foodEmojis))]
}

// randInRange returns a random int in [min, max).
func randInRange(min, max int) int {
	return rand.Intn(max-min) + min
}
