package life

import (
	"fmt"
	"sync"
)

type Point struct {
	X, Y int
}

func (point *Point) String() string {
	return fmt.Sprintf("[%d %d]", point.X, point.Y)
}

func (point *Point) AdjacentPoints() [8]Point {
	return [...]Point{
		{point.X, point.Y + 1},
		{point.X + 1, point.Y + 1},
		{point.X + 1, point.Y},
		{point.X + 1, point.Y - 1},
		{point.X, point.Y - 1},
		{point.X - 1, point.Y - 1},
		{point.X - 1, point.Y},
		{point.X - 1, point.Y + 1},
	}
}

type State struct {
	IsAlive bool
}

type Thing struct {
	Point
	State
}

type WorldLayout int

const (
	Unbounded WorldLayout = iota
	Bounded
	Torus
)

type World struct {
	Layout        WorldLayout
	Width, Height int
	things        [][]*Thing
}

func wrap(v, limit int) int {
	switch {
	case 0 <= v && v < limit:
		return v
	case v < 0:
		if v%limit == 0 {
			return 0
		} else {
			return limit + (v % limit)
		}
	case limit <= v:
		return v % limit
	default:
		panic(fmt.Sprintf("There's some weird math happening. Unable to compare two ints: %d, %d", v, limit))
	}
}

func offLimit(v, limit int) bool {
	return v < 0 || limit <= v
}

func (world *World) FindThing(point Point) (*Thing, bool) {
	switch world.Layout {
	case Unbounded:
		panic("Not implemented") // TODO: Implement
	case Bounded:
		if offLimit(point.X, world.Width) || offLimit(point.Y, world.Height) {
			return nil, false
		} else {
			return world.things[point.X][point.Y], true
		}
	case Torus:
		return world.things[wrap(point.X, world.Width)][wrap(point.Y, world.Height)], true
	default:
		return nil, false
	}
}

func (world *World) Neighbours(thing *Thing) []*Thing {
	result := make([]*Thing, 0)
	for _, point := range thing.Point.AdjacentPoints() {
		if thing, found := world.FindThing(point); found {
			result = append(result, thing)
		}
	}

	return result
}

func countAlive(things []*Thing) (count int) {
	for _, thing := range things {
		if thing.IsAlive {
			count++
		}
	}

	return
}

func (thing *Thing) String() string {
	if thing.IsAlive {
		return `.`
	}

	return ` `
}

func (world *World) String() string {
	result := ""
	for j := world.Height - 1; j >= 0; j-- {
		for i := 0; i < world.Width; i++ {
			result += world.things[i][j].String()
		}

		result += "\n"
	}

	return result
}

func NewThing(point Point, state State) *Thing {
	return &Thing{point, state}
}

func NewWorld(layout WorldLayout, width, height int, points []*Point) *World {
	things := make([][]*Thing, width)
	for i := 0; i < width; i++ {
		things[i] = make([]*Thing, height)
		for j := 0; j < height; j++ {
			things[i][j] = NewThing(Point{i, j}, State{false})
		}
	}

	for _, point := range points {
		x, y := wrap(point.X, width), wrap(point.Y, height)

		things[x][y].IsAlive = true
	}

	return &World{layout, width, height, things}
}

func (thing *Thing) Run(in chan *World, waitGroup *sync.WaitGroup) chan State {
	out := make(chan State)
	go func() {
		defer close(out)
		for world := range in {
			neighbours := world.Neighbours(thing)
			aliveCount := countAlive(neighbours)
			waitGroup.Done()
			switch {
			case aliveCount < 2:
				out <- State{IsAlive: false}
			case aliveCount == 2:
				out <- thing.State
			case aliveCount == 3:
				out <- State{IsAlive: true}
			case aliveCount > 3:
				out <- State{IsAlive: false}
			default:
				out <- thing.State
			}
		}
	}()

	return out
}

func (world *World) Run(in chan int) chan *World {
	out := make(chan *World)
	go func() {
		defer close(out)

		var waitGroup sync.WaitGroup

		ins, outs := make([][]chan *World, world.Width), make([][]chan State, world.Width)
		for i := 0; i < world.Width; i++ {
			ins[i], outs[i] = make([]chan *World, world.Height), make([]chan State, world.Height)
			for j := 0; j < world.Height; j++ {
				ins[i][j] = make(chan *World)
				outs[i][j] = world.things[i][j].Run(ins[i][j], &waitGroup)
			}
		}

		for range in {
			waitGroup.Add(world.Width * world.Height)
			for i := 0; i < world.Width; i++ {
				for j := 0; j < world.Height; j++ {
					ins[i][j] <- world
				}
			}

			waitGroup.Wait()
			for i := 0; i < world.Width; i++ {
				for j := 0; j < world.Height; j++ {
					world.things[i][j].State = <-outs[i][j]
				}
			}

			out <- world
		}
	}()

	return out
}
