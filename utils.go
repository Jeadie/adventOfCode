package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func scannerToChannel(s *bufio.Scanner, output chan string) {
	defer close(output)
	for s.Scan() {
		output <- s.Text()
	}
}
