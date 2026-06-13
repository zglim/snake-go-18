package trisnake

import (
	tl "github.com/JoelOtter/termloop"
)

// Snake is the player-controlled entity.
type Snake struct {
	*tl.Entity
	game  *Game
	dir   Direction
	body  []Coordinates
	color tl.Attr
}

// newSnake creates a snake at the default starting position.
func newSnake(g *Game, color tl.Attr) *Snake {
	s := &Snake{
		Entity: tl.NewEntity(5, 5, 1, 1),
		game:   g,
		dir:    DirRight,
		color:  color,
		body: []Coordinates{
			{1, 6}, // tail
			{2, 6}, // body
			{3, 6}, // head
		},
	}
	return s
}

// Tick handles direction input from key events.
func (s *Snake) Tick(event tl.Event) {
	if event.Type != tl.EventKey {
		return
	}
	var requested Direction
	switch event.Key {
	case tl.KeyArrowRight:
		requested = DirRight
	case tl.KeyArrowLeft:
		requested = DirLeft
	case tl.KeyArrowUp:
		requested = DirUp
	case tl.KeyArrowDown:
		requested = DirDown
	default:
		return
	}
	// Prevent reversing into yourself.
	if requested != s.dir.opposite() {
		s.dir = requested
	}
}

// Draw is called every frame. It moves the snake, resolves collisions,
// and renders the body.
func (s *Snake) Draw(screen *tl.Screen) {
	newHead := s.nextHead()

	if s.game.playScreen.food.contains(newHead) {
		s.handleFood(newHead)
	} else {
		// No food: move forward by dropping the tail.
		s.body = append(s.body[1:], newHead)
	}

	s.SetPosition(newHead.X, newHead.Y)

	if s.game.playScreen.arena.contains(newHead) || s.selfCollision() {
		s.game.gameOver()
		return
	}

	s.render(screen)
}

// head returns a pointer to the last element of the body (the head).
func (s *Snake) head() *Coordinates {
	return &s.body[len(s.body)-1]
}

// nextHead computes where the head will be on the next frame.
func (s *Snake) nextHead() Coordinates {
	h := *s.head()
	switch s.dir {
	case DirUp:
		h.Y--
	case DirDown:
		h.Y++
	case DirLeft:
		h.X--
	case DirRight:
		h.X++
	}
	return h
}

// handleFood processes a food collision.
func (s *Snake) handleFood(head Coordinates) {
	food := s.game.playScreen.food
	switch food.emoji {
	case 'R':
		s.game.addScore(5)
		s.game.applyRFoodBonus()
		s.body = append(s.body, head)
	case 'S':
		s.game.applySFoodPenalty()
		// S-food: score but no growth (original behaviour)
	default:
		s.game.addScore(1)
		s.body = append(s.body, head)
	}
	food.respawn()
}

// selfCollision returns true if the head overlaps any body segment.
func (s *Snake) selfCollision() bool {
	h := *s.head()
	for i := 0; i < len(s.body)-1; i++ {
		if s.body[i] == h {
			return true
		}
	}
	return false
}

// render draws every body segment on the screen.
func (s *Snake) render(screen *tl.Screen) {
	for _, c := range s.body {
		screen.RenderCell(c.X, c.Y, &tl.Cell{
			Fg: s.color,
			Ch: '░',
		})
	}
}
