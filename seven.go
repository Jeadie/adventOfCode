package main

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

/**
--- Day 7: The Treachery of Whales ---

A giant whale has decided your submarine is its next meal, and it's much faster than you are. There's nowhere to run!
Suddenly, a swarm of crabs (each in its own tiny submarine - it's too deep for them otherwise) zooms in to rescue you!
They seem to be preparing to blast a hole in the ocean floor; sensors indicate a massive underground cave system just
beyond where they're aiming! The crab submarines all need to be aligned before they'll have enough power to blast a
large enough hole for your submarine to get through. However, it doesn't look like they'll be aligned before the whale
catches you! Maybe you can help? There's one major catch - crab submarines can only move horizontally. You quickly make
a list of the horizontal position of each crab (your puzzle input). Crab submarines have limited fuel, so you need to
find a way to make all of their horizontal positions match while requiring them to spend as little fuel as possible.
For example, consider the following horizontal positions:

16,1,2,0,4,2,7,1,2,14

This means there's a crab with horizontal position 16, a crab with horizontal position 1, and so on. Each change of 1
step in horizontal position of a single crab costs 1 fuel. You could choose any horizontal position to align them all
on, but the one that costs the least fuel is horizontal position 2:

    Move from 16 to 2: 14 fuel
    Move from 1 to 2: 1 fuel
    Move from 2 to 2: 0 fuel
    Move from 0 to 2: 2 fuel
    Move from 4 to 2: 2 fuel
    Move from 2 to 2: 0 fuel
    Move from 7 to 2: 5 fuel
    Move from 1 to 2: 1 fuel
    Move from 2 to 2: 0 fuel
    Move from 14 to 2: 12 fuel

This costs a total of 37 fuel. This is the cheapest possible outcome; more expensive outcomes include aligning at
position 1 (41 fuel), position 3 (39 fuel), or position 10 (71 fuel). Determine the horizontal position that the crabs
can align to using the least fuel possible. How much fuel must they spend to align to that position?

minimise sum_{i=0}^n |x_n - k| --> dE/dk = sum_{i=0}^n k > x_i ? 1 : -1 --> dE/dk = sum_{i=0}^n 2*I(k > x_i) - 1
Start with d = mean(x_n),
		d_{i+1} = d_i - l* dE/dk


							*
	x_0			x_1 	x_2 k		x_3 	x_4
*/

type DerivativeInput struct {
	k float64
	x float64
}

type DerivativeResponse struct {
	d float64
}

func Seven() interface{} {
	const computeGoroutines = 5

	s := GetScanner(7)
	if !s.Scan() {
		fmt.Println("File had no content")
		return nil
	}

	x := strings.Split(s.Text(), ",")

	crabs := []float64{unsafeFloatParse(x[0])}
	for _, c := range x[1:] {
		crabs = append(crabs, unsafeFloatParse(c))
	}

	inputs := make(chan DerivativeInput)
	outputs := make(chan DerivativeResponse)

	// Goroutine: computeDerivative
	for i := 0; i < computeGoroutines; i++ {
		go computeDerivative(inputs, outputs)
	}

	threshold := 2.0
	k0 := computeAverage(crabs)
	mover0 := computeEarthMoverToPoint(crabs, k0)
	moverPrev := math.MaxFloat64

	for math.Abs(mover0-moverPrev) > threshold {
		moverPrev = mover0
		k0 = gradientStep(crabs, k0, inputs, outputs)
		mover0 = computeEarthMoverToPoint(crabs, k0)
	}
	return k0
}

func computeEarthMoverToPoint(x []float64, p float64) float64 {
	n := 0.0
	for _, f := range x {
		n += math.Abs(p - f)
	}
	return n
}

func gradientStep(crabs []float64, k0 float64, inputs chan DerivativeInput, outputs chan DerivativeResponse) float64 {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(crabs []float64, k0 float64, inputs chan DerivativeInput) {
		// Send out all inputs
		for c := range crabs {
			inputs <- DerivativeInput{
				k: k0,
				x: float64(c),
			}
		}
		wg.Done()
	}(crabs, k0, inputs)

	derivativeSum := 0.0

	go func(crabs []float64, outputs chan DerivativeResponse, result *float64) {
		n := len(crabs)
		// Get equal amount of responses to compute descent
		for n > 0 {
			r := <-outputs
			derivativeSum += r.d
			n--
		}
		wg.Done()
	}(crabs, outputs, &derivativeSum)

	wg.Wait()
	if derivativeSum > 0 {
		return k0 - 1
	}
	return k0 + 1
}

func computeAverage(crabs []float64) float64 {
	a, n := 0.0, 0.0
	for _, c := range crabs {
		a += c
		n++
	}
	return math.RoundToEven(a / n)
}

func computeDerivative(i chan DerivativeInput, o chan DerivativeResponse) {
	for in := range i {
		ind := float64(0)
		if in.k > in.x {
			ind = -1
		} else if in.k < in.x {
			ind = +1
		} else {
			ind = 0
		}
		o <- DerivativeResponse{d: ind}
	}
}
