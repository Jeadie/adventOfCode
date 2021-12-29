package main

import (
	"bufio"
	"strings"
)

/**
--- Day 13: Transparent Origami ---

You reach another volcanically active part of the cave. It would be nice if you could do some kind of thermal imaging so
you could tell ahead of time which caves are too hot to safely enter. Fortunately, the submarine seems to be equipped
with a thermal camera! When you activate it, you are greeted with:

Congratulations on your purchase! To activate this infrared thermal imaging
camera system, please enter the code found on page 1 of the manual.

Apparently, the Elves have never used this feature. To your surprise, you manage to find the manual; as you go to open
it, page 1 falls out. It's a large sheet of transparent paper! The transparent paper is marked with random dots and
includes instructions on how to fold it up (your puzzle input). For example:

6,10
0,14
9,10
0,3
10,4
4,11
6,0
6,12
4,1
0,13
10,12
3,4
3,0
8,4
1,10
2,14
8,10
9,0

fold along y=7
fold along x=5

The first section is a list of dots on the transparent paper. 0,0 represents the top-left coordinate. The first value,
x, increases to the right. The second value, y, increases downward. So, the coordinate 3,0 is to the right of 0,0, and
the coordinate 0,7 is below 0,0. The coordinates in this example form the following pattern, where # is a dot on the
paper and . is an empty, unmarked position:

...#..#..#.
....#......
...........
#..........
...#....#.#
...........
...........
...........
...........
...........
.#....#.##.
....#......
......#...#
#..........
#.#........

Then, there is a list of fold instructions. Each instruction indicates a line on the transparent paper and wants you to
fold the paper up (for horizontal y=... lines) or left (for vertical x=... lines). In this example, the first fold
instruction is fold along y=7, which designates the line formed by all of the positions where y is 7 (marked here
with -):

...#..#..#. 0
....#...... 1
........... 2
#.......... 3
...#....#.# 4
........... 5
........... 6
----------- y=7
........... 7
........... 8
.#....#.##. 9
....#...... 10
......#...# 11
#.......... 12
#.#........ 13

Because this is a horizontal line, fold the bottom half up. Some of the dots might end up overlapping after the fold is
complete, but dots will never appear exactly on a fold line. The result of doing this fold looks like this:

#.##..#..#.
#...#......
......#...#
#...#......
.#.#..#.###
...........
...........

Now, only 17 dots are visible. Notice, for example, the two dots in the bottom left corner before the transparent paper
is folded; after the fold is complete, those dots appear in the top left corner (at 0,0 and 0,1). Because the paper is
transparent, the dot just below them in the result (at 0,3) remains visible, as it can be seen through the transparent
paper. Also notice that some dots can end up overlapping; in this case, the dots merge together and become a single dot.
The second fold instruction is fold along x=5, which indicates this line:

#.##.|#..#.
#...#|.....
.....|#...#
#...#|.....
.#.#.|#.###
.....|.....
.....|.....

Because this is a vertical line, fold left:

#####
#...#
#...#
#...#
#####
.....
.....

The instructions made a square!

The transparent paper is pretty big, so for now, focus on just completing the first fold. After the first fold in the
example above, 17 dots are visible - dots that end up overlapping after the fold is completed count as a single dot.
How many dots are visible after completing just the first fold instruction on your transparent paper?
*/

func Thirteen() interface{} {
	// Create TransparentPaper and FoldInstructions
	// FoldInstructions is below TransparentPaper in input file.
	s, tp := createTransparentPaper()
	fi := createFoldInstructions(s)

	inst := fi[0]
	var tp1 TransparentPaper
	if inst.IsHorizontal() {
		tp1 = tp.ApplyHorizontal(inst)
	} else {
		tp1 = tp.ApplyVertical(inst)
	}

	// 0
	return tp1.getTotalCount()
}

func createFoldInstructions(s *bufio.Scanner) []FoldInstructions {
	fi := []FoldInstructions{}

	// Fold instructions next
	for s.Scan() {
		// get "x=655" from "fold along x=655"
		t := strings.SplitN(s.Text(), " ", 3)[2]

		// Get "x", "655" from "x=655"
		f := strings.SplitN(t, "=", 2)
		d, i := f[0], f[1]
		fi = append(fi, FoldInstructions{
			dimension: d == "y",
			index:     uint(unsafeUintParse(i)),
		})
	}
	return fi
}

//
func createTransparentPaper() (*bufio.Scanner, *TransparentPaper) {
	s := GetScanner(13)
	m := make(map[int][]int)

	// Start with points
	for s.Scan() {
		l := s.Text()
		// Empty line is end of the points.
		if len(l) == 0 {
			break
		}

		t := strings.SplitN(l, ",", 2)
		a, b := unsafeIntParse(t[0]), unsafeIntParse(t[1])
		if _, ok := m[a]; !ok {
			m[a] = []int{}
		}
		m[a] = append(m[a], b)
	}
	tp := TransparentPaper{marks: m}

	return s, &tp
}

type TransparentPaper struct {
	// Map row to index of dots in that row. Overlap becomes row intersection. Optimised for horizontal fold.
	//TODO: Use second map to optimise for vertical line fold.
	marks map[int][]int
}

func (tp *TransparentPaper) ApplyHorizontal(fi FoldInstructions) TransparentPaper {
	line := int(fi.index)
	result := TransparentPaper{marks: make(map[int][]int)}
	for i := 0; i < int(fi.index); i++ {
		result.marks[line-i-1] = tp.Intersection(line-i-1, line+i)
	}
	return result
}

func (tp *TransparentPaper) ApplyVertical(fi FoldInstructions) TransparentPaper {
	// TODO: Implement
	return TransparentPaper{marks: make(map[int][]int)}
}

func (tp *TransparentPaper) getTotalCount() int {
	count := 0
	for _, dots := range tp.marks {
		count += len(dots)
	}
	return count
}

func (tp *TransparentPaper) Intersection(i int, j int) []int {
	a := tp.marks[i]
	b := tp.marks[j]

	frequencyTable := make(map[int]int)
	result := []int{}

	for _, e := range a {
		frequencyTable[e] += 1
	}

	for _, e := range b {
		frequencyTable[e] += 1
	}

	for k, v := range result {
		if v >= 2 {
			result = append(result, k)
		}
	}
	return result
}

type FoldInstructions struct {
	dimension bool // false -> horizontal, true -> Columns
	index     uint //
}

func (fi FoldInstructions) IsHorizontal() bool {
	return !fi.dimension
}
