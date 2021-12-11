package main

import (
	"fmt"
	"strings"
	"sync"
)

/**
--- Day 9: Smoke Basin ---

These caves seem to be lava tubes. Parts are even still volcanically active; small hydrothermal vents release smoke
into the caves that slowly settles like rain. If you can model how the smoke flows through the caves, you might be able
to avoid it and be that much safer. The submarine generates a heightmap of the floor of the nearby caves for you (your
puzzle input). Smoke flows to the lowest point of the area it's in. For example, consider the following heightmap:

2199943210
3987894921
9856789892
8767896789
9899965678

Each number corresponds to the height of a particular location, where 9 is the highest and 0 is the lowest a location
can be. Your first goal is to find the low points - the locations that are lower than any of its adjacent locations.
Most locations have four adjacent locations (up, down, left, and right); locations on the edge or corner of the map
have three or two adjacent locations, respectively. (Diagonal locations do not count as adjacent.) In the above example,
there are four low points, all highlighted: two are in the first row (a 1 and a 0), one is in the third row (a 5), and
one is in the bottom row (also a 5). All other locations on the heightmap have some lower adjacent location, and so are
not low points. The risk level of a low point is 1 plus its height. In the above example, the risk levels of the low
points are 2, 1, 6, and 6. The sum of the risk levels of all low points in the heightmap is therefore 15. Find all of
the low points on your heightmap. What is the sum of the risk levels of all low points on your heightmap?
*/

// maxASCIINumeric > all 0-9 as ASCII characters
const maxASCIINumeric = ":"
const lowPointGoroutines = 1

// 48 is '0' in ASCII.
const zeroASCIIValue = 48
const pointsBuffer = 20
const lowPointsBuffer = 20

func Nine() interface{} {

	points := make(chan PointAndBorders, pointsBuffer)
	lowPoints := make(chan uint8, lowPointsBuffer)

	// Construct map (with float.max border)
	b := ConstructBorderedMap()
	// Send out points with four borders (goroutine for each row)
	wg := sync.WaitGroup{}
	wg.Add(1)                // len(b.grid)-2)
	for i := 1; i < 2; i++ { //len(b.grid)-1; i++ {
		go sendPointsFromLine(i, b, points, &wg)
	}

	// Return low points
	for i := 0; i < lowPointGoroutines; i++ {
		go filterLowPoints(points, lowPoints)
	}

	// Wait for all grid points to be constructed, then wait for all points to be processed.
	wg.Wait()
	close(points)

	// Sum over low points
	result := uint8(0)
	for p := range lowPoints {
		result += (p - zeroASCIIValue) + 1
	}
	// 37
	return result
}

func filterLowPoints(in chan PointAndBorders, out chan uint8) {
	for p := range in {
		if p.x < p.left && p.x < p.right && p.x < p.up && p.x < p.down {
			out <- p.x
		}
	}
	close(out)
}

func sendPointsFromLine(i int, b *BorderedMap, output chan PointAndBorders, wg *sync.WaitGroup) {
	l := b.grid[i]
	for j := 1; j < len(l)-1; j++ {
		row := l[j-1 : j+2]
		output <- PointAndBorders{
			x:     row[1],
			left:  row[0],
			right: row[2],
			up:    b.grid[i-1][j],
			down:  b.grid[i+1][j],
		}
	}
	wg.Done()
}

func ConstructBorderedMap() *BorderedMap {
	s := GetScanner(9)
	// Create first row with maximums
	s.Scan()
	first := s.Text()
	g := []string{
		strings.Repeat(maxASCIINumeric, len(first)+2),
	}

	for s.Scan() {
		l := s.Text()
		g = append(
			g,
			fmt.Sprintf("%s%s%s", maxASCIINumeric, l, maxASCIINumeric),
		)
	}
	g = append(
		g,
		strings.Repeat(maxASCIINumeric, len(first)+2),
	)
	return &BorderedMap{grid: g}
}

type PointAndBorders struct {
	x     uint8
	left  uint8
	right uint8
	up    uint8
	down  uint8
}

type BorderedMap struct {
	grid []string
}
