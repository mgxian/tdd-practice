package marsrover

import (
	"reflect"
	"testing"
)

func TestTurn(t *testing.T) {
	t.Run("turn left 90 degrees", func(t *testing.T) {
		marsRover := newMarsRover()
		marsRover.setDirection(North)
		marsRover.turn90DegreeLeft()
		assertDirection(t, marsRover.direction, West)
		marsRover.turn90DegreeLeft()
		marsRover.turn90DegreeLeft()
		marsRover.turn90DegreeLeft()
		assertDirection(t, marsRover.direction, North)
	})

	t.Run("turn right 90 degrees", func(t *testing.T) {
		marsRover := newMarsRover()
		marsRover.setDirection(North)
		marsRover.turn90DegreeRight()
		assertDirection(t, marsRover.direction, East)
		marsRover.turn90DegreeRight()
		marsRover.turn90DegreeRight()
		marsRover.turn90DegreeRight()
		assertDirection(t, marsRover.direction, North)
	})
}

func TestMove(t *testing.T) {
	moveTests := []struct {
		testName       string
		maxX, maxY     int
		startx, starty int
		startDirection Direction
		distance       int
		endx, endy     int
	}{
		{"forward 3 from (1, 3) head north direction", 10, 10, 1, 3, North, 3, 1, 6},
		{"forward 9 from (1, 3) head north direction", 10, 10, 1, 3, North, 9, 1, 10},
		{"forward 3 from (1, 3) head east direction", 10, 10, 1, 3, East, 3, 4, 3},
		{"forward 9 from (1, 3) head east direction", 10, 10, 1, 3, East, 9, 10, 3},
	}
	for _, tt := range moveTests {
		t.Run(tt.testName, func(t *testing.T) {
			marsRover := newMarsRover()
			marsRover.limitArea(tt.maxX, tt.maxY)
			marsRover.setStartPostion(tt.startx, tt.starty)
			marsRover.setDirection(tt.startDirection)
			marsRover.forward(tt.distance)
			assertPostion(t, marsRover.postion(), Postion{tt.endx, tt.endy})
		})
	}
}

func assertPostion(t *testing.T, got, want Postion) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got (%d, %d), want (%d, %d)", got.x, got.y, want.x, want.y)
	}
}

func assertDirection(t *testing.T, got, want Direction) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
