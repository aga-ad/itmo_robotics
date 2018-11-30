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

func TestLineIntersectsCircle(t *testing.T) {
	c := CircleId{RelatedPoint{Point{4, 5}, 1}, 4, 16}
	var s Segment

	s = Segment{Point{1, 2}, Point{2, 1}}
	if SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
	}
	s = Segment{Point{0, 3}, Point{3, 0}}
	if SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
	}
	s = Segment{Point{0, 0}, Point{1, 1}}
	if SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
	}

	s = Segment{Point{3, 1}, Point{5, 1}}
	if !SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
	}
	s = Segment{Point{4, 5}, Point{5, 5}}
	if !SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
	}
	s = Segment{Point{4, 1}, Point{4, 2}}
	if !SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
	}
	s = Segment{Point{0, 0}, Point{10, 10}}
	if !SegmentIntersectsCircle(&s, &c) {
		t.Error("FUU ", s)
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
