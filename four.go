package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

/**
You're already almost 1.5km (almost a mile) below the surface of the ocean, already so deep that you can't see any
sunlight. What you can see, however, is a giant squid that has attached itself to the outside of your submarine. Maybe
it wants to play bingo? Bingo is played on a set of boards each consisting of a 5x5 grid of numbers. Numbers are chosen
at random, and the chosen number is marked on all boards on which it appears. (Numbers may not appear on all boards.)
If all numbers in any row or any column of a board are marked, that board wins. (Diagonals don't count.) The submarine
has a bingo subsystem to help passengers (currently, you and the giant squid) pass the time. It automatically generates
a random order in which to draw numbers and a random set of boards (your puzzle input). For example:

7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7

After the first five numbers are drawn (7, 4, 9, 5, and 11), there are no winners, but the boards are marked as follows
(shown here adjacent to each other to save space):

22 13 17 11  0         3 15  0  2 22        14 21 17 24  4
 8  2 23  4 24         9 18 13 17  5        10 16 15  9 19
21  9 14 16  7        19  8  7 25 23        18  8 23 26 20
 6 10  3 18  5        20 11 10 24  4        22 11 13  6  5
 1 12 20 15 19        14 21 16 12  6         2  0 12  3  7

After the next six numbers are drawn (17, 23, 2, 0, 14, and 21), there are still no winners:

22 13 17 11  0         3 15  0  2 22        14 21 17 24  4
 8  2 23  4 24         9 18 13 17  5        10 16 15  9 19
21  9 14 16  7        19  8  7 25 23        18  8 23 26 20
 6 10  3 18  5        20 11 10 24  4        22 11 13  6  5
 1 12 20 15 19        14 21 16 12  6         2  0 12  3  7

Finally, 24 is drawn:

22 13 17 11  0         3 15  0  2 22        14 21 17 24  4
 8  2 23  4 24         9 18 13 17  5        10 16 15  9 19
21  9 14 16  7        19  8  7 25 23        18  8 23 26 20
 6 10  3 18  5        20 11 10 24  4        22 11 13  6  5
 1 12 20 15 19        14 21 16 12  6         2  0 12  3  7

At this point, the third board wins because it has at least one complete row or column of marked numbers (in this case,
the entire top row is marked: 14 21 17 24 4). The score of the winning board can now be calculated. Start by finding
the sum of all unmarked numbers on that board; in this case, the sum is 188. Then, multiply that sum by the number that
was just called when the board won, 24, to get the final score, 188 * 24 = 4512. To guarantee victory against the giant
squid, figure out which board will win first. What will your final score be if you choose that board?
*/
func Four() interface{} {
	f, err := os.Open("input/four.txt")
	if err != nil {
		fmt.Printf("input file does not exist. Error %v\n", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	const boardProcessingCount = 1

	wg := sync.WaitGroup{}
	wg.Add(1 + boardProcessingCount) // input & output

	// Construct Draw
	s.Scan()
	d := constructDraw(s.Text())

	// Send out boards
	boards := make(chan *Board) // closed in inputBoards
	go inputBoards(s, boards, &wg)

	// Each Board goroutine
	candidates := make(chan *FourCandidate) // closed after wg.Wait()
	for i := 0; i < boardProcessingCount; i++ {
		go boardProcessing(boards, candidates, d, &wg)
	}

	// Find the quickest candidate board
	best := make(chan *FourCandidate) // Closed in reduceCandidates
	go reduceCandidates(candidates, best)

	// Wait for inputs to be processed.
	wg.Wait()
	close(candidates)

	// Wait for best, compute numeric result.
	// ?
	return computeOutput(<-best, d)
}

func computeOutput(c *FourCandidate, d *Draw) interface{} {
	output := make(chan uint8, 2)
	go d.unMarkedAtI(int(c.i), output)

	// Sum list of marked numbers
	x := uint64(0)
	for u := range output {
		y, _ := strconv.ParseUint(string(u), 10, 64)
		x += y
	}

	// Multiply by end number
	end, _ := strconv.ParseUint(d.line[c.i], 10, 64)
	return x * end
}

func reduceCandidates(in chan *FourCandidate, out chan *FourCandidate) {
	defer close(out)
	candidate := &FourCandidate{i: math.MaxUint8}
	for c := range in {
		if c.i < candidate.i {
			candidate = c
		}
	}
	out <- candidate
}

func boardProcessing(in chan *Board, out chan *FourCandidate, d *Draw, w *sync.WaitGroup) {
	defer w.Done()
	for b := range in {
		// Construct row and column of index i simultaneously.
		for i := uint8(0); i < 5; i++ {
			row := [5]string{}
			col := [5]string{}

			for j := uint8(0); j < 5; j++ {
				row[j] = b.Get(i, j)
				col[j] = b.Get(j, i)
			}
			out <- &FourCandidate{
				b: b,
				i: getIndexOfLine(row, d),
			}
			out <- &FourCandidate{
				b: b,
				i: getIndexOfLine(col, d),
			}
		}
	}
}

// getIndexOfLine finds the index into the draw which by this line of the Bingo board would be complete.
func getIndexOfLine(l [5]string, d *Draw) uint8 {
	max := uint8(0)
	for _, s := range l {
		i := d.get(s)
		if i > max {
			max = i
		}
	}
	return max
}

// inputBoards constructs Board structs from the input file and sends them into the channel.
func inputBoards(s *bufio.Scanner, out chan *Board, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	const size = 5

	for s.Scan() {
		rows := [size][size]string{}
		for i := 0; i < size; i++ {
			if s.Scan() {

				copy(rows[i][:], strings.SplitN(strings.Replace(s.Text(), "  ", " ", -1), " ", size))
			}
		}
		out <- createBoard(rows)
	}
}

// FourCandidate does not need to return line in question, only index.
type FourCandidate struct {
	b *Board
	i uint8
}

func constructDraw(raw string) *Draw {
	d := Draw{
		order: make(map[string]uint8),
		line:  strings.Split(raw, ","),
	}
	for i, s := range d.line {
		d.order[s] = uint8(i)
	}
	return &d
}

type Draw struct {
	// maps a token, 0 <= t <= 99, to the index, i of when it was drawn 0 <= i <= 99
	order map[string]uint8

	line []string
}

func (d Draw) get(x string) uint8 {
	return d.order[x]
}

func (d Draw) unMarkedAtI(i int, out chan uint8) {
	for _, s := range d.line[i+1:] {
		u, _ := strconv.ParseUint(s, 10, 8)
		out <- uint8(u)
	}
	close(out)
}

type Board struct {
	b [5][5]string
}

func createBoard(rows [5][5]string) *Board {
	return &Board{
		b: rows,
	}
}

func (b Board) Get(x uint8, y uint8) string {
	return b.b[x][y]
}
