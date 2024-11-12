package component

type Direction struct {
	X int
	Y int
}

var Directions = map[int]Direction{
	0: {1, 0},  // right
	1: {-1, 0}, // left
	2: {0, -1}, // up
	3: {0, 1},  // down
}

type Location struct {
	X float32
	Y float32
}

type Movement struct {
	CurrentDirection Direction
	Velocity         float32
	CurrentLocation  Location
}

func (Movement) Name() string {
	return "Movement"
}
