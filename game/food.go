package trisnake

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// Maximum number of attempts to find a valid food position before giving up.
const maxFoodPlacementAttempts = 1000

func init() {
	rand.Seed(time.Now().UnixNano())
}

// insideBorderW and insideBorderH define the playable area inside the arena border.
// The arena is created with NewArena(70, 25), giving arena.Width=69, arena.Height=24.
// The border occupies x=0, x=69, y=0, y=24, so playable cells are x=[1,68], y=[1,23].
var insideBorderW = 70 - 1
var insideBorderH = 25 - 1

// NewFood creates a new piece of food, avoiding the given occupied coordinates (e.g. snake body).
func NewFood(avoid []Coordinates) *Food {
	food := new(Food)
	food.Entity = tl.NewEntity(1, 1, 1, 1)
	food.MoveFood(avoid)
	return food
}

// MoveFood moves the food to a random position that is not on the border and not
// on any of the coordinates in the avoid list (typically the snake body).
func (food *Food) MoveFood(avoid []Coordinates) {
	avoidSet := make(map[Coordinates]struct{}, len(avoid))
	for _, c := range avoid {
		avoidSet[c] = struct{}{}
	}

	for attempt := 0; attempt < maxFoodPlacementAttempts; attempt++ {
		newX := randomInsideArena(insideBorderW, 1)
		newY := randomInsideArena(insideBorderH, 1)
		candidate := Coordinates{newX, newY}

		if _, occupied := avoidSet[candidate]; !occupied {
			food.Foodposition.X = newX
			food.Foodposition.Y = newY
			food.Emoji = RandomFood()
			food.SetPosition(food.Foodposition.X, food.Foodposition.Y)
			return
		}
	}

	// Fallback: if we couldn't find a free cell after many attempts, just place it
	// at the first random position (extremely unlikely — the arena has ~1500 cells).
	food.Foodposition.X = randomInsideArena(insideBorderW, 1)
	food.Foodposition.Y = randomInsideArena(insideBorderH, 1)
	food.Emoji = RandomFood()
	food.SetPosition(food.Foodposition.X, food.Foodposition.Y)
}

// RandomFood picks a random food rune from the available set.
func RandomFood() rune {
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
	return emoji[rand.Intn(len(emoji))]
}

// Draw prints the food on the screen.
func (food *Food) Draw(screen *tl.Screen) {
	screen.RenderCell(food.Foodposition.X, food.Foodposition.Y, &tl.Cell{
		Ch: food.Emoji,
	})
}

// Contains checks if the food is at the given coordinates.
func (food *Food) Contains(c Coordinates) bool {
	return c.X == food.Foodposition.X && c.Y == food.Foodposition.Y
}

// randomInsideArena returns a random int in [iMin, iMax).
func randomInsideArena(iMax int, iMin int) int {
	return rand.Intn(iMax-iMin) + iMin
}
