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

// Postion is postion
type Postion struct {
	x int
	y int
}

// MarsRover is mars rover.
type MarsRover struct {
	direction Direction
	Postion
	maxX int
	maxY int
}

func (mr *MarsRover) setStartPostion(x, y int) {
	mr.x = x
	mr.y = y
}

func (mr *MarsRover) setDirection(d Direction) {
	mr.direction = d
}

func (mr *MarsRover) limitArea(x, y int) {
	mr.maxX = x
	mr.maxY = y
}

func (mr *MarsRover) postion() Postion {
	return mr.Postion
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

func (mr *MarsRover) forward(d int) {
	switch mr.direction {
	case North:
		mr.y += d
	case East:
		mr.x += d
	case South:
		mr.y -= d
	case West:
		mr.x -= d
	}

	if mr.x > mr.maxX {
		mr.x = mr.maxX
	}

	if mr.y > mr.maxY {
		mr.y = mr.maxY
	}

	if mr.x < 0 {
		mr.x = 0
	}

	if mr.y < 0 {
		mr.y = 0
	}
}

func newMarsRover() *MarsRover {
	mr := new(MarsRover)
	mr.direction = North
	mr.x = 0
	mr.y = 0
	return mr
}
