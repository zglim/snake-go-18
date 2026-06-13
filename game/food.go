package trisnake

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// rng is a package-level random source, seeded once at init time.
// This avoids the deprecated rand.Seed() being called on every random operation,
// which could produce identical seeds when called in rapid succession.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// insideborderW and insideborderH define the safe spawn area inside the arena border.
var insideborderW = 70 - 1
var insideborderH = 25 - 1

// NewFood will create a new piece of food, this will only happen once when the game has started.
func NewFood() *Food {
	food := new(Food)
	// Create a new entity food with a standard position and 1x1 size
	food.Entity = tl.NewEntity(1, 1, 1, 1)
	// Call function MoveFood to move the food to a safe random position.
	food.MoveFood()

	return food
}

// MoveFood moves the food into a new random position that does not overlap
// with the snake body or the arena border.
func (food *Food) MoveFood() {
	const maxAttempts = 1000
	for attempt := 0; attempt < maxAttempts; attempt++ {
		NewX := RandomInsideArena(insideborderW, 1)
		NewY := RandomInsideArena(insideborderH, 1)
		candidate := Coordinates{NewX, NewY}

		// Skip if position is on the arena border
		if gs != nil && gs.ArenaEntity != nil && gs.ArenaEntity.Contains(candidate) {
			continue
		}
		// Skip if position is on the snake body
		if gs != nil && gs.SnakeEntity != nil && gs.SnakeEntity.ContainsCoord(candidate) {
			continue
		}

		food.Foodposition.X = NewX
		food.Foodposition.Y = NewY
		food.Emoji = RandomFood()
		food.SetPosition(food.Foodposition.X, food.Foodposition.Y)
		return
	}

	// Fallback: if we couldn't find a safe position after max attempts,
	// just place it somewhere (extremely unlikely with normal arena sizes)
	food.Foodposition.X = RandomInsideArena(insideborderW, 1)
	food.Foodposition.Y = RandomInsideArena(insideborderH, 1)
	food.Emoji = RandomFood()
	food.SetPosition(food.Foodposition.X, food.Foodposition.Y)
}

// RandomFood will use the ASCII-charset to pick a random rune from the slice and print it out as food.
func RandomFood() rune {
	// This slice contains all of the possible food icons.
	emoji := []rune{
		'R', // Favourite dish, extra points!!!
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'■', // 1 point
		'S', // You do not want to eat the skull
	}

	// Return a random rune picked from the slice
	return emoji[rng.Intn(len(emoji))]
}

// Draw will print out the food on the screen.
func (food *Food) Draw(screen *tl.Screen) {
	screen.RenderCell(food.Foodposition.X, food.Foodposition.Y, &tl.Cell{
		Ch: food.Emoji,
	})
}

// Contains checks if food contains the coordinates, if so this will return a bool.
func (food *Food) Contains(c Coordinates) bool {
	return c.X == food.Foodposition.X && c.Y == food.Foodposition.Y
}

// RandomInsideArena will return a random value between iMin (inclusive) and iMax (exclusive).
func RandomInsideArena(iMax int, iMin int) int {
	return rng.Intn(iMax-iMin) + iMin
}
