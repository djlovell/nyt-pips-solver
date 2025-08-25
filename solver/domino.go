package solver

type domino struct {
	Left  int
	Right int
}

func NewDomino(left, right int) *domino {
	if left < 0 || left > 6 || right < 0 || right > 6 {
		panic("Domino values must be between 0 and 6")
	}
	return &domino{
		Left:  left,
		Right: right,
	}
}
