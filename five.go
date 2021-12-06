package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

/**
--- Day 5: Hydrothermal Venture ---
You come across a field of hydrothermal vents on the ocean floor! These vents constantly produce large, opaque clouds,
so it would be best to avoid them if possible. They tend to form in lines; the submarine helpfully produces a list of
nearby lines of vents (your puzzle input) for you to review. For example:

0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2

Each line of vents is given as a line segment in the format x1,y1 -> x2,y2 where x1,y1 are the coordinates of one end
the line segment and x2,y2 are the coordinates of the other end. These line segments include the points at both ends.
In other words:

An entry like 1,1 -> 1,3 covers points 1,1, 1,2, and 1,3.
An entry like 9,7 -> 7,7 covers points 9,7, 8,7, and 7,7.

For now, only consider horizontal and vertical lines: lines where either x1 = x2 or y1 = y2. So, the horizontal and
vertical lines from the above list would produce the following diagram:

.......1..
..1....1..
..1....1..
.......1..
.112111211
..........
..........
..........
..........
222111....

In this diagram, the top left corner is 0,0 and the bottom right corner is 9,9. Each position is shown as the number of
lines which cover that point or . if no line covers that point. The top-left pair of 1s, for example, comes from
2,2 -> 2,1; the very bottom row is formed by the overlapping lines 0,9 -> 5,9 and 0,9 -> 2,9. To avoid the most
dangerous areas, you need to determine the number of points where at least two lines overlap. In the above example,
this is anywhere in the diagram with a 2 or larger - a total of 5 points. Consider only horizontal and vertical lines.
At how many points do at least two lines overlap?
*/

func Five() interface{} {
	lines := make(chan LineSegment)
	z2s := make(chan Z2)

	g := Grid{
		points:  sync.Map{},
		counter: 0,
	}

	// Goroutine: Send out line segments
	go createLines(lines)

	// Goroutine: Process line segment -> points
	go createPoints(lines, z2s)

	// Goroutine: Mark grid points as line covered
	wg := sync.WaitGroup{}
	wg.Add(1)
	go markGridPoints(z2s, &g, &wg)

	wg.Wait()
	return g.counter
}

func markGridPoints(in chan Z2, g *Grid, wg *sync.WaitGroup) {
	for z := range in {
		g.increment(z)
	}
	wg.Done()
}

func createPoints(in chan LineSegment, out chan Z2) {
	for l := range in {
		if l.a.x == l.b.x {
			a, b := getMinMax(l.a.y, l.b.y)
			for i := a; i <= b; i++ {
				out <- Z2{
					x: l.a.x,
					y: i,
				}
			}
		} else if l.a.y == l.b.y {
			a, b := getMinMax(l.a.x, l.b.x)
			for i := a; i <= b; i++ {
				out <- Z2{
					x: i,
					y: l.a.y,
				}
			}
		}
	}
	close(out)
}

func createLines(out chan LineSegment) {
	s := GetScanner(5)
	// Line format: 976,453 -> 561,38
	for s.Scan() {
		points := strings.SplitN(s.Text(), " -> ", 2)
		a := LineSegment{
			a: stringToZ2(points[0]),
			b: stringToZ2(points[1]),
		}
		out <- a
	}
	close(out)
}

func stringToZ2(s string) Z2 {
	pair := strings.SplitN(s, ",", 2)
	x, _ := strconv.ParseUint(pair[0], 10, 32)
	y, _ := strconv.ParseUint(pair[1], 10, 32)

	return Z2{
		x: uint32(x),
		y: uint32(y),
	}
}

type LineSegment struct {
	a Z2
	b Z2
}

type Z2 struct {
	x uint32
	y uint32
}

func (z Z2) toString() string {
	return fmt.Sprintf("%d,%d", z.x, z.y)
}

type Grid struct {
	points  sync.Map // Z2 -> uint32
	counter uint8
}

func (g *Grid) increment(z Z2) {
	a := uint32(0)
	v, _ := g.points.LoadOrStore(z.toString(), &a)

	r := atomic.AddUint32(v.(*uint32), 1)

	// Because adding to g.points is atomic, the value returned will be unique (for given z2).
	if r == 2 {
		g.counter++
	}
}
