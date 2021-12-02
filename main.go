package main

import (
	"os"
	"strconv"
)

type adventCode struct {
	fn func()
}

func main() {
	days := [25]adventCode{{one}}

	adventDay, _ := strconv.ParseInt(os.Args[1], 10, 8)
	days[adventDay-1].fn()
}