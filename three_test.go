package main

import "testing"

func TestThree(t *testing.T) {
	for i := 0; i < 1000; i++ {
		if Three().(int64) != 63753 {
			t.Error("Incorrect!")
		}
	}
}
