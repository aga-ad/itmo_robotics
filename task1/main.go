package main

import (
	"fmt"
	"log"
	"os"
)

func run(mapName string) {
	gameMap, err := ParseMap(mapName + ".json")
	if err != nil {
		panic(err)
	}
	_, _, a := Strategy(&gameMap)

	out, err := os.Create(mapName + "_out.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer out.Close()
	for _, p := range a {
		fmt.Fprintln(out, p.X(), ",", p.Y())
	}
	fmt.Println(mapName+":", len(a))
}

func main() {
	/*c1 := CircleId{RelatedPoint{Point{0, 0}, 42}, 1, 1}

	//c2 := CircleId{RelatedPoint{6, 0, 43}, 3, 9}
	c2 := CircleId{RelatedPoint{Point{2, 0}, 43}, 1, 1}

	fmt.Println(TangentsPoints(&c1, &c2))*/
	maps := [...]string{"trivial", "gigant", "bumps", "harder", "cheese", "asteroids", "mess"}
	run(maps[1])
	// for _, mapName := range maps {
	// 	run(mapName)
	// }

	/*gameMap, err := ParseMap(mapName + ".json")
	if err != nil {
		panic(err)
	}



	x, v, a := Strategy(&gameMap)

	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer f.Close()

	out, err := os.Create(mapName + "_out.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer out.Close()

	fmt.Fprintln(f, "X")
	for _, p := range x {
		fmt.Fprintln(f, p)
	}
	fmt.Fprintln(f, "\n\n\nV")
	for _, p := range v {
		fmt.Fprintln(f, p)
	}
	fmt.Fprintln(f, "\n\n\nA")
	for _, p := range a {
		fmt.Fprintln(f, p)
		fmt.Fprintln(out, p.X(), ",", p.Y())
	}*/
}
