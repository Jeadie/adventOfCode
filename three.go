package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sync"
)

/**
--- Day 3: Binary Diagnostic ---

The submarine has been making some odd creaking noises, so you ask it to produce a diagnostic report just in case. The
diagnostic report (your puzzle input) consists of a list of binary numbers which, when decoded properly, can tell you
many useful things about the conditions of the submarine. The first parameter to check is the power consumption. You
need to use the binary numbers in the diagnostic report to generate two new binary numbers (called the gamma rate and
the epsilon rate). The power consumption can then be found by multiplying the gamma rate by the epsilon rate. Each bit
in the gamma rate can be determined by finding the most common bit in the corresponding position of all numbers in the
diagnostic report. For example, given the following diagnostic report:

00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010

Considering only the first bit of each number, there are five 0 bits and seven 1 bits. Since the most common bit is 1,
the first bit of the gamma rate is 1. The most common second bit of the numbers in the diagnostic report is 0, so the
second bit of the gamma rate is 0. The most common value of the third, fourth, and fifth bits are 1, 1, and 0,
respectively, and so the final three bits of the gamma rate are 110. So, the gamma rate is the binary number 10110,
or 22 in decimal. The epsilon rate is calculated in a similar way; rather than use the most common bit, the least common
bit from each position is used. So, the epsilon rate is 01001, or 9 in decimal. Multiplying the gamma rate (22) by the
epsilon rate (9) produces the power consumption, 198. Use the binary numbers in your diagnostic report to calculate the
gamma rate and epsilon rate, then multiply them together. What is the power consumption of the submarine? (Be sure to
represent your answer in decimal, not binary.)
*/
func Three() interface{} {
	// Make go routine that processes specific column of binary input.
	ins := [12]chan string{}
	outputs := make(chan Result)
	wg := sync.WaitGroup{}
	wg.Add(13) // 12 ins & input goroutines. Output requires separate.

	for i := range ins {
		ins[i] = make(chan string)
		go func(in chan string, out chan Result, wg *sync.WaitGroup, i int) {
			for v := range in {
				out <- Result{
					i: i,
					v: parseAsciiBit(int32(v[i])),
				}
			}
			wg.Done()
		}(ins[i], outputs, &wg, i)
	}

	// For each line in file, send to channels, then close inputs.
	n := 0
	go func(filename string, wg *sync.WaitGroup, n *int) {
		f, _ := os.Open(filename)
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			s := s.Text()
			for _, c := range ins {
				c <- s
			}
			*n++
		}
		for _, in := range ins {
			close(in)
		}
		wg.Done()
	}("input/three.txt", &wg, &n)

	// Put each result into correct column
	rate := [12]int64{0,0,0,0,0,0,0,0,0,0,0,0}
	reduceWg := sync.WaitGroup{}
	reduceWg.Add(1)
	go func(out chan Result, wg *sync.WaitGroup, result *[12]int64) {
		for r := range out {
			rate[r.i] += r.v
		}
		wg.Done()
	}(outputs, &reduceWg, &rate)

	wg.Wait()
	// Close outputs and wait for final goroutine to finish processing results.
	close(outputs)
	reduceWg.Wait()


	// 63753
	return computeResult(computeGamma(rate, n))
}

func computeResult(gamma [12]int64) int64 {
	g := int64(0)
	e := int64(0)

	for i := range gamma {
		pow := math.Pow(float64(i), 2)
		g += int64(pow)*gamma[i]

		// gamma is most common, epsilon is least common.
		e += int64(pow)*(1-gamma[i])
	}
	return g*e
}

func computeGamma(rates [12]int64, n int) [12]int64 {
	g := [12]int64{}
	for i, r := range rates {
		if 2*r >= int64(n) {
			g[i] = 1
		} else {
			g[i] = 0
		}
	}
	return g
}

func parseAsciiBit(b int32) int64{
	if b == 49 {
		return 1
	} else if b == 48 {
		return 0
	}
	panic(fmt.Sprintf("Ascii bit %d was neither a '0' (48) or '1' (49)", b))
}

type Result struct {
	i int
	v int64
}

