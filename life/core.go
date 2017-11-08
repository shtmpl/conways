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
	Torus WorldLayout = iota
	Unbounded
)

type World struct {
	Layout WorldLayout
	Width, Height int
	things        [][]*Thing
}

func wrap(v, limit int) (int, bool) {
	switch {
	case 0 <= v && v < limit:
		return v, true
	case limit <= v:
		return v % limit, true
	case v < 0:
		if v%limit == 0 {
			return 0, true
		} else {
			return limit + (v % limit), true
		}
	default:
		return -1, false
	}
}

func (world *World) FindThing(point Point) (*Thing, bool) {
	switch world.Layout {
	case Torus:
		x, appropriate := wrap(point.X, world.Width)
		if !appropriate {
			return nil, false
		}

		y, appropriate := wrap(point.Y, world.Height)
		if !appropriate {
			return nil, false
		}

		return world.things[x][y], true
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
		x, normalised := wrap(point.X, width)
		if !normalised {
			continue
		}

		y, normalised := wrap(point.Y, height)
		if !normalised {
			continue
		}

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
