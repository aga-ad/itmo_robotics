package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	/*c1 := CircleId{RelatedPoint{Point{0, 0}, 42}, 1, 1}

	//c2 := CircleId{RelatedPoint{6, 0, 43}, 3, 9}
	c2 := CircleId{RelatedPoint{Point{2, 0}, 43}, 1, 1}

	fmt.Println(TangentsPoints(&c1, &c2))*/
	gameMap, err := ParseMap("mess.json")
	if err != nil {
		panic(err)
	}

	//cs := newCircleStorage(gameMap.Circles)
	/*for _, a := range cs.circles {
		if len(a) < 4 {
			fmt.Println(a)
		}
	}*/

	x, v, a := Strategy(&gameMap)

	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer f.Close()

	out, err := os.Create("res.txt")
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
	}
}
