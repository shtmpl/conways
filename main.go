package main

import (
	"fmt"
	"time"
	"flag"
	"bufio"
	"os"
	"github.com/shtmpl/conways/life"
)

func stdinBasedTicks() chan time.Time {
	out := make(chan time.Time)
	go func() {
		defer close(out)

		reader := bufio.NewReader(os.Stdin)
		for {
			_, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			out <- time.Now()
		}
	}()

	return out
}

func main() {
	interactive := flag.Bool("interactive", false, "Allows stepping through the iterations")
	interval := flag.Int("interval", 100, "Time interval between consecutive steps (in ms)")
	flag.Parse()

	form := life.NewForm(life.GosperGliderGun).Translate(80, 25)
	world := life.NewWorld(life.Torus, 210, 53, *form)
	fmt.Println(world)
	fmt.Println("Created")

	steps := make(chan int)
	result := world.Run(steps)

	var ticks <-chan time.Time
	if *interactive {
		ticks = stdinBasedTicks()
	} else {
		ticks = time.Tick(time.Duration(*interval) * time.Millisecond)
	}

	i := 0
	for range ticks {
		go func() { steps <- i }()
		fmt.Println(<-result)
		fmt.Println("Iteration:", i)
		i++
	}

	fmt.Println("Main")
}
