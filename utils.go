package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func GetScanner(x uint8) *bufio.Scanner {
	f, err := os.Open(intToName(x))
	if err != nil {
		fmt.Printf("input file does not exist. Error %v\n", err)
	}
	return bufio.NewScanner(f)
}

func intToName(x uint8) string {
	FilenameArray := [25]string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten",
		"eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty",
		"twentyOne", "twentyTwo", "twentyThree", "twentyFour",
	}
	return fmt.Sprintf("input/%s.txt", FilenameArray[x-1])
}

func getMinMax(x uint32, y uint32) (uint32, uint32) {
	var a, b uint32
	if x < y {
		a = x
		b = y
	} else {
		a = y
		b = x
	}
	return a, b
}

func unsafeUintParse(x string) uint64 {
	y, _ := strconv.ParseUint(x, 10, 64)
	return y
}

func unsafeIntParse(x string) int {
	y, _ := strconv.ParseInt(x, 10, 64)
	return int(y)
}

func unsafeFloatParse(x string) float64 {
	y, _ := strconv.ParseFloat(x, 10)
	return y
}

func maxInt(x int, y int) int {
	if x < y {
		return y
	}
	return x
}
func minInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func maxUint(x uint, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func makeIntChan(a int, b int) chan int {
	c := make(chan int, b-a)
	defer close(c)
	for i := a; i < b; i++ {
		c <- i
	}
	return c
}

func scannerToChannel(s *bufio.Scanner, output chan string) {
	defer close(output)
	for s.Scan() {
		output <- s.Text()
	}
}

type Appender struct {
	l []string
	m sync.Mutex
}

func (a *Appender) Append(s string) {
	a.m.Lock()
	a.l = append(a.l, s)
	a.m.Unlock()
}

func (a *Appender) GetSet() []string {
	k := map[string]bool{}
	r := []string{}
	a.m.Lock()
	for _, s := range a.l {
		exists, ok := k[s]
		if !(ok && exists) {
			r = append(r, s)
			k[s] = true
		}
	}
	a.m.Unlock()
	return r
}

func CreateAppender() *Appender {
	return &Appender{
		l: []string{},
		m: sync.Mutex{},
	}
}
