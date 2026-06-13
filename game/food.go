package trisnake

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// foodEmojis defines the weighted pool of possible food types.
var foodEmojis = []rune{
	'R',   // Reward: 5 points + speed boost
	'■', '■', '■', '■', '■', '■', '■', '■', '■', '■', // Normal: 1 point + grow
	'S', // Speed penalty
}

// insideBorderW and insideBorderH define the playable area bounds (excluding border).
const (
	insideBorderW = arenaWidth - 1
	insideBorderH = arenaHeight - 1
)

// newFood creates a new food item at a random position inside the arena.
func newFood() *Food {
	food := &Food{
		Entity: tl.NewEntity(1, 1, 1, 1),
	}
	food.Randomize()
	return food
}

// Randomize repositions the food to a new random location with a random emoji.
func (food *Food) Randomize() {
	food.Position.X = randIntRange(1, insideBorderW)
	food.Position.Y = randIntRange(1, insideBorderH)
	food.Emoji = randomFoodEmoji()
	food.SetPosition(food.Position.X, food.Position.Y)
}

// Draw renders the food on the screen.
func (food *Food) Draw(screen *tl.Screen) {
	screen.RenderCell(food.Position.X, food.Position.Y, &tl.Cell{
		Ch: food.Emoji,
	})
}

// Contains returns true if the given coordinates match the food position.
func (food *Food) Contains(c Coordinates) bool {
	return c == food.Position
}

// randomFoodEmoji picks a random food type from the weighted pool.
func randomFoodEmoji() rune {
	rand.Seed(time.Now().UnixNano())
	return foodEmojis[rand.Intn(len(foodEmojis))]
}

// randIntRange returns a random int in [min, max).
func randIntRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
