package marsrover

// Direction mars rover's direction
type Direction int

// four direction
const (
	North Direction = iota
	East
	South
	West
)

// MarsRover is mars rover.
type MarsRover struct {
	direction Direction
}

func (mr *MarsRover) turn90DegreeLeft() {
	direction := mr.direction - 1
	if direction < 0 {
		direction += 4
	}
	mr.direction = direction
}

func (mr *MarsRover) turn90DegreeRight() {
	mr.direction = (mr.direction + 1) % 4
}

func newMarsRover(d Direction) *MarsRover {
	mr := new(MarsRover)
	mr.direction = d
	return mr
}
