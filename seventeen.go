package main

import (
	"strings"
)

/**
--- Day 17: Trick Shot ---

You finally decode the Elves' message. HI, the message says. You continue searching for the sleigh keys. Ahead of you is
what appears to be a large ocean trench. Could the keys have fallen into it? You'd better send a probe to investigate.
The probe launcher on your submarine can fire the probe with any integer velocity in the x (forward) and y (upward, or
downward if negative) directions. For example, an initial x,y velocity like 0,10 would fire the probe straight up, while
an initial velocity like 10,-1 would fire the probe forward at a slight downward angle. The probe's x,y position starts
at 0,0. Then, it will follow some trajectory by moving in steps. On each step, these changes occur in the following
order:

    The probe's x position increases by its x velocity.
    The probe's y position increases by its y velocity.
    Due to drag, the probe's x velocity changes by 1 toward the value 0; that is, it decreases by 1 if it is greater
		than 0, increases by 1 if it is less than 0, or does not change if it is already 0.
    Due to gravity, the probe's y velocity decreases by 1.

For the probe to successfully make it into the trench, the probe must be on some trajectory that causes it to be within
a target area after any step. The submarine computer has already calculated this target area (your puzzle input).
For example:

target area: x=20..30, y=-10..-5

This target area means that you need to find initial x,y velocity values such that after any step, the probe's x
position is at least 20 and at most 30, and the probe's y position is at least -10 and at most -5. Given this target
area, one initial velocity that causes the probe to be within the target area after any step is 7,2:

.............#....#............
.......#..............#........
...............................
S........................#.....
...............................
...............................
...........................#...
...............................
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................TTTTTTTT#TT
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................TTTTTTTTTTT

In this diagram, S is the probe's initial position, 0,0. The x coordinate increases to the right, and the y coordinate
increases upward. In the bottom right, positions that are within the target area are shown as T. After each step (until
the target area is reached), the position of the probe is marked with #. (The bottom-right # is both a position the
probe reaches and a position in the target area.) Another initial velocity that causes the probe to be within the target
area after any step is 6,3:

...............#..#............
...........#........#..........
...............................
......#..............#.........
...............................
...............................
S....................#.........
...............................
...............................
...............................
.....................#.........
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................T#TTTTTTTTT
....................TTTTTTTTTTT

Another one is 9,0:

S........#.....................
.................#.............
...............................
........................#......
...............................
....................TTTTTTTTTTT
....................TTTTTTTTTT#
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................TTTTTTTTTTT
....................TTTTTTTTTTT

One initial velocity that doesn't cause the probe to be within the target area after any step is 17,-4:

S..............................................................
...............................................................
...............................................................
...............................................................
.................#.............................................
....................TTTTTTTTTTT................................
....................TTTTTTTTTTT................................
....................TTTTTTTTTTT................................
....................TTTTTTTTTTT................................
....................TTTTTTTTTTT..#.............................
....................TTTTTTTTTTT................................
...............................................................
...............................................................
...............................................................
...............................................................
................................................#..............
...............................................................
...............................................................
...............................................................
...............................................................
...............................................................
...............................................................
..............................................................#

The probe appears to pass through the target area, but is never within it after any step. Instead, it continues down and
to the right - only the first few steps are shown. If you're going to fire a highly scientific probe out of a super cool
probe launcher, you might as well do it with style. How high can you make the probe go while still reaching the target
area? In the above example, using an initial velocity of 6,9 is the best you can do, causing the probe to reach a
maximum y position of 45. (Any higher initial y velocity causes the probe to overshoot the target area entirely.) Find
the initial velocity that causes the probe to reach the highest y position and still eventually be within the target
area after any step. What is the highest y position it reaches on this trajectory?

-- Soln --
Clearly, the greatest height will be that with the highest initial vertical velocity. With a positive v_y, it will arrive back at zero with a velocity -(v_y+1) at some x co-ordinate x'. Velocities up will
be: v_y, v_y-1, v_y-2, ..., 0. If -(v_y+1) < y_min, then it will move from (0, x') to (-v_y+1, x') and miss the box. But
if we set v_y = 1-y_min, we get (0, x') to (-v_y+1, x') == (y_min, x'). Therefore maximum v_y == 1-y_min (this is for a
-ve y_min)

Now we need to get an v_x satisfactory. v_x will be, given an initial v_x >0 , v_x(t) = max(0, initV_x - t) (Will not
go negative). One can use this to compute maximum distance x, for given v_x.
	sum_{t=0)^(initV_x) initV_x - t
  = initV_x**2 - initV_x*(initV_x+1)/2
  = (initV_x**2)/2 -initV_x/2

Now find where it will be >= left hand side of box, x0
	x0 = (initV_x**2)/2 -initV_x/2
	0 = (initV_x**2)/2 -initV_x/2 - x0
	initV_x = ... (solutions to root of quadratic)
If necessary and round up (must be integer solution).

We have shown we can find an initial velocity (v_x, v_y) that will hit the box, and is of maximum height. Maximum height
is summation of +ve velocities: sum_{t=0}^{v_y} v_y-t ==  sum_{t=0}^{v_y} t == v_y*(v_y + 1) / 2
													  == (1-y_min)*((1-y_min) + 1) / 2
													  == (1-y_min)*(2-y_min) / 2
*/

func Seventeen() interface{} {
	s := GetScanner(17)
	if !s.Scan() {
		return nil
	}
	b := CreateBox(s.Text())

	// (1-y_min)*(y_min) / 2
	return ((1 - b.y0) * (2 - b.y0)) / 2
}

func CreateBox(t string) *Box {
	// 'target area: x=81..129, y=-150..-108' ->  'x=81..129, y=-150..-108'
	p := strings.SplitN(t, ":", 2)[1]

	// 'x=81..129, y=-150..-108' -> [' x=81..129', ' y=-150..-108']
	pos := strings.SplitN(p, ",", 2)

	// ' x=81..129' -> ['81', '129']
	xx := strings.SplitN(pos[0][3:], "..", 2)
	yy := strings.SplitN(pos[1][3:], "..", 2)

	return &Box{
		x0: unsafeIntParse(xx[0]),
		x1: unsafeIntParse(xx[1]),
		y0: unsafeIntParse(yy[0]),
		y1: unsafeIntParse(yy[1]),
	}
}

func PathIntersectsBox(p [][2]int, b *Box) bool {
	for _, d := range p {
		if b.IncludesPoint(d[0], d[1]) {
			return true
		}
	}
	return false
}

// ComputePath calculates the path of a projectile fired from (0,0) given initial velocity (initVelX, initVelY) until it
//is below (in Y dimension), lowestY. Kinematic rules follow:
//  - Position (x, y) change by its velocity (vX, vY).
//  - |vX_{t}| = |vX_{t-1}| -1, unless |vX_{t}| == 0 then no change.
//  - vY_{t} = vY_{t-1} - 1
func ComputePath(initVelX int, initVelY int, lowestY int) [][2]int {
	x, y := 0, 0
	velX, velY := initVelX, initVelY
	points := [][2]int{{0, 0}}

	for y > lowestY {
		// Update position
		x += velX
		y += velY
		points = append(points, [2]int{x, y})

		// Update velocity
		if velX > 0 {
			velX--
		} else if velX < 0 {
			velX++
		}
		velY--
	}
	return points
}

type Box struct {
	x0, x1, y0, y1 int
}

func (b *Box) IncludesPoint(x int, y int) bool {
	return b.x0 <= x && x <= b.x1 && b.y0 <= y && y <= b.y1
}
