package main

import (
	"fmt"
	"testing"
)

func TestLine(t *testing.T) {
	line := LineFromPoints(&RelatedPoint{Point{3, 8}, 1}, &RelatedPoint{Point{5, 7}, 2})
	if line != (Line{-1, -2, 19}) {
		t.Error("Expected -1, -2, 19, got ", line)
	}
}

func TestMinPath(t *testing.T) {
	//v := []RelatedPoint{RelatedPoint{Point{0, 0}, 0}, RelatedPoint{Point{1, 1}, 1}}
	dist := [][]float64{{0, 1}, {1, 0}}
	start := 0
	end := 1
	//v, dist, start, end := MakeGraph(gameMap.Circles)
	path := MinPath(dist, start, end)
	fmt.Println("path:", path)
}
