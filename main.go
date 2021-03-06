package main

import (
	"fmt"
	"os"
	"strconv"
)

type adventCode struct {
	fn func() interface{}
}

// To consider: https://felixge.de/2021/12/01/advent-of-go-profiling-2021-day-1-1/
func main() {
	days := [25]adventCode{
		{One}, {Two}, {Three}, {Four}, {Five}, {Six}, {Seven}, {Eight}, {Nine},
		{Ten}, {Eleven}, {Twelve}, {Thirteen}, {Fourteen}, {Fifteen}, {Sixteen},
		{Seventeen}, {Eighteen}, {Nineteen}, {Twenty},
	}

	adventDay, _ := strconv.ParseInt(os.Args[1], 10, 8)
	x := days[adventDay-1].fn()
	fmt.Println(x)
}
