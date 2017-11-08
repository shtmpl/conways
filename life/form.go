package life

var (
	// Still

	Block = []Point{
		{0, 0}, {0, 1},
		{1, 0}, {1, 1},
	}

	// Oscillators

	Blinker = []Point{
		{0, 1}, {1, 1}, {2, 1},
	}

	Pulsar = []Point{
		{1, 3}, {1, 4}, {1, 5}, {1, 9}, {1, 10}, {1, 11},
		{3, 1}, {3, 6}, {3, 8}, {3, 13},
		{4, 1}, {4, 6}, {4, 8}, {4, 13},
		{5, 1}, {5, 6}, {5, 8}, {5, 13},
		{6, 3}, {6, 4}, {6, 5}, {6, 9}, {6, 10}, {6, 11},
		{8, 3}, {8, 4}, {8, 5}, {8, 9}, {8, 10}, {8, 11},
		{9, 1}, {9, 6}, {9, 8}, {9, 13},
		{10, 1}, {10, 6}, {10, 8}, {10, 13},
		{11, 1}, {11, 6}, {11, 8}, {11, 13},
		{13, 3}, {13, 4}, {13, 5}, {13, 9}, {13, 10}, {13, 11},
	}

	// Spaceships

	Glider = []Point{
		{0, 0}, {0, 2},
		{1, 0}, {1, 1},
		{2, 1},
	}

	LightweightSpaceship = []Point{
		{0, 1}, {0, 3},
		{1, 0},
		{2, 0},
		{3, 0}, {3, 3},
		{4, 0}, {4, 1}, {4, 2},
	}

	// Methuselahs

	Rpentomino = []Point{
		{0, 1},
		{1, 0}, {1, 1}, {1, 2},
		{2, 2},
	}

	Diehard = []Point{
		{0, 1},
		{1, 0}, {1, 1},
		{5, 0},
		{6, 0}, {6, 2},
		{7, 0},
	}

	Acorn = []Point{
		{0, 0},
		{1, 0}, {1, 2},
		{3, 1},
		{4, 0},
		{5, 0},
		{6, 0},
	}

	// Guns

	GosperGliderGun = []Point{
		{0, 3}, {0, 4},
		{1, 3}, {1, 4},
		{10, 2}, {10, 3}, {10, 4},
		{11, 1}, {11, 5},
		{12, 0}, {12, 6},
		{13, 0}, {13, 6},
		{14, 3},
		{15, 1}, {15, 5},
		{16, 2}, {16, 3}, {16, 4},
		{17, 3},
		{20, 4}, {20, 5}, {20, 6},
		{21, 4}, {21, 5}, {21, 6},
		{22, 3}, {22, 7},
		{24, 2}, {24, 3}, {24, 7}, {24, 8},
		{34, 5}, {34, 6},
		{35, 5}, {35, 6},
	}
)

type Form []*Point

func NewForm(points []Point) *Form {
	result := make(Form, len(points))
	for i, point := range points {
		result[i] = &Point{point.X, point.Y}
	}

	return &result
}

//func FromString(points string) *Form {
//	strings.Split(points, "\n")
//}

//func min(xs ...int) (m int) {
//	for _, x := range xs {
//		if x < m {
//			m = x
//		}
//	}
//
//	return
//}
//
//func max(xs ...int) (m int) {
//	for _, x := range xs {
//		if m < x {
//			m = x
//		}
//	}
//
//	return
//}
//
//func (form *Form) bottomLeft() (x, y int) {
//	for _, p := range *form {
//		x, y = min(x, p.X), min(y, p.Y)
//	}
//
//	return
//}
//
//func (form *Form) bottomRight() (x, y int) {
//	for _, p := range *form {
//		x, y = max(x, p.X), min(y, p.Y)
//	}
//
//	return
//}
//
//func (form *Form) topLeft() (x, y int) {
//	for _, p := range *form {
//		x, y = min(x, p.X), max(y, p.Y)
//	}
//
//	return
//}
//
//func (form *Form) topRight() (x, y int) {
//	for _, p := range *form {
//		x, y = max(x, p.X), max(y, p.Y)
//	}
//
//	return
//}
//
//func (form *Form) center() (x, y int) {
//	loX, loY := form.bottomLeft()
//	hiX, hiY := form.topRight()
//
//	return loX + (hiX-loX)/2, loY + (hiY-loY)/2
//}

func (form *Form) Translate(x, y int) *Form {
	for _, p := range *form {
		p.X, p.Y = p.X+x, p.Y+y
	}

	return form
}

// FIXME: Remove duplicates. Care for pointers
func (form *Form) Comp(rest ...*Form) *Form {
	for _, other := range rest {
		*form = append(*form, *other...)
	}

	return form
}
