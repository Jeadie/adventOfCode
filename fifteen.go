package main

import (
	"bufio"
	"fmt"
	"math"
)

/**
--- Day 15: Chiton ---

You've almost reached the exit of the cave, but the walls are getting closer together. Your submarine can barely still
fit, though; the main problem is that the walls of the cave are covered in chitons, and it would be best not to bump any
of them. The cavern is large, but has a very low ceiling, restricting your motion to two dimensions. The shape of the
cavern resembles a square; a quick scan of chiton density produces a map of risk level throughout the cave (your puzzle
input). For example:

1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581

You start in the top left position, your destination is the bottom right position, and you cannot move diagonally. The
number at each position is its risk level; to determine the total risk of an entire path, add up the risk levels of each
position you enter (that is, don't count the risk level of your starting position unless you enter it; leaving it adds
no risk to your total). Your goal is to find a path with the lowest total risk. In this example, a path with the lowest
total risk is highlighted here:

1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581

The total risk of this path is 40 (the starting position is never entered, so its risk is not counted). What is the
lowest total risk of any path from the top left to the bottom right?
*/

func Fifteen() interface{} {
	c := InitChiton(GetScanner(15))
	paths := make(chan ChitonPath, 100000) // TODO: This does not scale at all
	paths <- ChitonPath{
		x:       0,
		y:       1,
		netRisk: c.getRisk(0, 1),
	}
	paths <- ChitonPath{
		x:       1,
		y:       0,
		netRisk: c.getRisk(1, 0),
	}

	i := 2
	// TODO: BFS implementation has problems due to number of elements in path space. Convert to DFS
	for p := range paths {
		i--
		fmt.Println(p, i)
		isNewMin := c.updateChitonRisk(p)
		if !isNewMin {
			if i == 0 {
				break
			}
			continue
		}
		for _, pp := range c.GetAdjacent(p) {
			paths <- pp
			i++
		}
	}
	return c.getFullPathRisk()
}

func InitChiton(s *bufio.Scanner) *Chiton {
	var risks [][]uint
	var netRisk [][]uint

	for s.Scan() {
		l := s.Text()
		var risk []uint
		var net []uint

		for _, i := range l {
			risk = append(risk, uint(i-48))
			net = append(net, math.MaxUint)
		}
		risks = append(risks, risk)
		netRisk = append(netRisk, net)
	}
	return &Chiton{
		risk:          risks,
		minChitonRisk: netRisk,
	}
}

type Chiton struct {
	risk          [][]uint // the risk of each position (x, y)
	minChitonRisk [][]uint // the minimum risk to get to position (x, y)
}

func (c *Chiton) getRisk(x int, y int) uint {
	return c.risk[y][x]
}

func (c *Chiton) getMinRisk(x int, y int) uint {
	return c.minChitonRisk[y][x]
}

func (c *Chiton) updateChitonRisk(p ChitonPath) bool {
	if p.netRisk < c.getMinRisk(p.x, p.y) {
		c.minChitonRisk[p.y][p.x] = p.netRisk
		return true
	}
	return false
}

func (c *Chiton) getWidth() int {
	return len(c.minChitonRisk[0])
}

func (c *Chiton) getHeight() int {
	return len(c.minChitonRisk)
}

func (c *Chiton) getFullPathRisk() uint {
	return c.minChitonRisk[c.getHeight()-1][c.getWidth()-1]
}

func (c *Chiton) GetAdjacent(p ChitonPath) []ChitonPath {
	x := []int{p.x - 1, p.x + 1}
	y := []int{p.y - 1, p.y + 1}
	var adjs []ChitonPath

	for _, xx := range x {
		for _, yy := range y {
			if 0 <= xx && xx < c.getWidth() && 0 <= yy && yy < c.getHeight() {
				adjs = append(adjs, ChitonPath{
					x:       xx,
					y:       yy,
					netRisk: p.netRisk + c.getRisk(xx, yy), // Continue on path from ChitonPath
				})
			}
		}
	}
	return adjs
}

type ChitonPath struct {
	x       int
	y       int
	netRisk uint
}
