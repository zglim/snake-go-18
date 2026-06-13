package trisnake

import tl "github.com/JoelOtter/termloop"

// newSnake creates a snake at the default starting position heading right.
func newSnake(g *Game) *Snake {
	snake := &Snake{
		Entity:    tl.NewEntity(5, 5, 1, 1),
		Direction: DirRight,
		game:      g,
		Body: []Coordinates{
			{1, 6}, // Tail
			{2, 6}, // Body
			{3, 6}, // Head
		},
	}
	return snake
}

// Head returns a pointer to the snake's head coordinates (the last element).
func (snake *Snake) Head() *Coordinates {
	return &snake.Body[len(snake.Body)-1]
}

// headCollidesWithBorder returns true if the snake's head is on an arena border cell.
func (snake *Snake) headCollidesWithBorder() bool {
	arena := snake.game.state.arena
	_, exists := arena.Border[*snake.Head()]
	return exists
}

// headCollidesWithFood returns true if the snake's head is on the food position.
func (snake *Snake) headCollidesWithFood() bool {
	food := snake.game.state.food
	return food.Position == *snake.Head()
}

// headCollidesWithBody returns true if the snake's head overlaps any body segment.
func (snake *Snake) headCollidesWithBody() bool {
	head := *snake.Head()
	for i := 0; i < len(snake.Body)-1; i++ {
		if head == snake.Body[i] {
			return true
		}
	}
	return false
}

// advance computes the next head position based on the current direction.
func (snake *Snake) advance() Coordinates {
	next := *snake.Head()
	switch snake.Direction {
	case DirUp:
		next.Y--
	case DirDown:
		next.Y++
	case DirLeft:
		next.X--
	case DirRight:
		next.X++
	}
	return next
}

// Draw is called by termloop each tick. It moves the snake, handles collisions,
// and renders the snake body on the screen.
func (snake *Snake) Draw(screen *tl.Screen) {
	g := snake.game
	nextHead := snake.advance()

	// Handle food collision before updating body position.
	if snake.headCollidesWithFood() {
		snake.handleFoodEffect(g, nextHead)
		g.moveFood()
	} else {
		// No food: slide forward (add new head, drop tail).
		snake.Body = append(snake.Body[1:], nextHead)
	}

	snake.SetPosition(nextHead.X, nextHead.Y)

	// Check for fatal collisions (border or self).
	if snake.headCollidesWithBorder() || snake.headCollidesWithBody() {
		g.showGameOver()
		return
	}

	// Render all body segments.
	color := colorIndex(g.colors.SnakeIdx)
	for _, c := range snake.Body {
		screen.RenderCell(c.X, c.Y, &tl.Cell{
			Fg: color,
			Ch: '░',
		})
	}
}

// handleFoodEffect processes the food type effect and grows the snake.
// nextHead is the pre-computed next head position.
func (snake *Snake) handleFoodEffect(g *Game, nextHead Coordinates) {
	food := g.state.food

	switch food.Emoji {
	case 'R':
		// Reward food: 5 points + speed boost (varies by difficulty).
		g.addScore(5)
		switch g.difficulty {
		case DiffEasy:
			if g.state.fps-3 > 8 {
				g.adjustFPS(-3)
			}
		case DiffNormal:
			if g.state.fps-2 > 12 {
				g.adjustFPS(-2)
			}
		case DiffHard:
			if g.state.fps-1 > 20 {
				g.adjustFPS(-1)
			}
		}
		snake.Body = append(snake.Body, nextHead)

	case 'S':
		// Speed food: increases speed (penalty), snake does not grow or move.
		switch g.difficulty {
		case DiffEasy:
			g.adjustFPS(1)
		case DiffNormal:
			g.adjustFPS(3)
		case DiffHard:
			g.adjustFPS(5)
		}

	default:
		// Normal food: 1 point + grow.
		g.addScore(1)
		snake.Body = append(snake.Body, nextHead)
	}
}

// Tick handles directional key input for the snake.
func (snake *Snake) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	switch event.Key {
	case tl.KeyArrowRight:
		if snake.Direction != DirLeft {
			snake.Direction = DirRight
		}
	case tl.KeyArrowLeft:
		if snake.Direction != DirRight {
			snake.Direction = DirLeft
		}
	case tl.KeyArrowUp:
		if snake.Direction != DirDown {
			snake.Direction = DirUp
		}
	case tl.KeyArrowDown:
		if snake.Direction != DirUp {
			snake.Direction = DirDown
		}
	}
}
