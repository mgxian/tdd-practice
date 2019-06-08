package marsrover

import "testing"

func TestTurn(t *testing.T) {
	t.Run("turn left 90 degrees", func(t *testing.T) {
		marsRover := newMarsRover(North)
		marsRover.turn90DegreeLeft()
		assertDirection(t, marsRover.direction, West)
		marsRover.turn90DegreeLeft()
		marsRover.turn90DegreeLeft()
		marsRover.turn90DegreeLeft()
		assertDirection(t, marsRover.direction, North)
	})

	t.Run("turn right 90 degrees", func(t *testing.T) {
		marsRover := newMarsRover(North)
		marsRover.turn90DegreeRight()
		assertDirection(t, marsRover.direction, East)
		marsRover.turn90DegreeRight()
		marsRover.turn90DegreeRight()
		marsRover.turn90DegreeRight()
		assertDirection(t, marsRover.direction, North)
	})
}

func assertDirection(t *testing.T, got, want Direction) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
