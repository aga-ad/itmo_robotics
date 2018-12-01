package main

import "testing"

func BenchmarkSample(b *testing.B) {
	gameMap, err := ParseMap("mess.json")
	if err != nil {
		panic(err)
	}

	Strategy(&gameMap)
}
