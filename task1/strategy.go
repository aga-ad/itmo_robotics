package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

type node struct {
	d   float64
	par int
}

type node2 struct {
	d        float64
	ind, par int
}

type MPHeap []node2

func (h MPHeap) Len() int            { return len(h) }
func (h MPHeap) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h MPHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MPHeap) Push(x interface{}) { *h = append(*h, x.(node2)) }

func (h *MPHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func MinPath(dist [][]float64, start, end int) []int {
	min := make([]node, len(dist))
	free := &MPHeap{}
	for i := 0; i < len(dist); i++ {
		min[i] = node{dist[start][i], i}
		heap.Push(free, node2{dist[start][i], i, start})
	}
	//fmt.Println(min)
	//fmt.Println(free)
	//count := 0
	for free.Len() > 0 {
		//count++
		//fmt.Println(count)
		best := heap.Pop(free).(node2)
		min[best.ind] = node{best.d, best.par}
		if math.IsInf(best.d, 1) {
			break
		}
		for i := 0; i < free.Len(); i++ {
			if (*free)[i].d > best.d+dist[best.ind][(*free)[i].ind] {
				(*free)[i] = node2{best.d + dist[best.ind][(*free)[i].ind], (*free)[i].ind, best.ind}
				heap.Fix(free, i)
			}
		}
	}
	fmt.Println("dijkstra done")
	res := make([]int, 0)
	cur := end
	for cur != start {
		res = append(res, cur)
		if cur == min[cur].par {
			return make([]int, 0)
		}
		cur = min[cur].par
	}
	res = append(res, start)
	for i := len(res)/2 - 1; i >= 0; i-- {
		opp := len(res) - 1 - i
		res[i], res[opp] = res[opp], res[i]
	}
	return res
}

func MakeGraph(gm *Map, add float64) ([]RelatedPoint, [][]float64, int, int) {
	fmt.Println("making graph. circles:", len(gm.Circles))
	v := make([]RelatedPoint, 0)
	v = append(v, RelatedPoint{Point{0, 0}, -1})
	v = append(v, RelatedPoint{Point{1, 1}, -2})
	for _, c := range gm.Circles {
		for k := 0; k < 2; k++ {
			for _, p := range TangentsPointsFromPoint(&c, &v[k], add) {
				s := Segment{v[k].Point, p.Point}
				if p.X() >= 0 && p.Y() >= 0 && p.X() <= 1 && p.Y() <= 1 && !gm.Storage.Intersect(&s) {
					v = append(v, p)
				}
			}
		}
	}
	for i := 0; i < len(gm.Circles); i++ {
		for j := i + 1; j < len(gm.Circles); j++ {
			tangents := TangentsPoints(&gm.Circles[i], &gm.Circles[j], add)
			for i := 0; i < len(tangents); i += 2 {
				s := Segment{tangents[i].Point, tangents[i+1].Point}
				f := !gm.Storage.Intersect(&s)
				for k := 0; f && k < 2; k++ {
					if f && tangents[i+k].X() >= 0 && tangents[i+k].Y() >= 0 && tangents[i+k].X() <= 1 && tangents[i+k].Y() <= 1 {
						v = append(v, tangents[i+k])
					}
				}
			}
		}
	}
	fmt.Println("making graph. vertices:", len(v))
	dist := make([][]float64, len(v))
	for i := 0; i < len(v); i++ {
		dist[i] = make([]float64, len(v))
		for j := 0; j < len(v); j++ {
			if i == j {
				dist[i][j] = 0
				continue
			}
			if i > j {
				dist[i][j] = dist[j][i]
				continue
			}
			s := Segment{v[i].Point, v[j].Point}
			if !gm.Storage.Intersect(&s) {
				dist[i][j] = Dist(&v[i], &v[j])
			} else {
				dist[i][j] = math.Inf(1)
			}
		}
	}
	return v, dist, 0, 1
}

type Map struct {
	Name    string     `json:"name"`
	Dt      float64    `json:"dt"`
	Fmax    float64    `json:"Fmax"`
	Circles []CircleId `json:"circles"`
	Storage CircleStorage
}

var circleCounter = 0

func (c *CircleId) UnmarshalJSON(b []byte) error {
	s := struct {
		X float64 `json:"X"`
		Y float64 `json:"Y"`
		R float64 `json:"R"`
	}{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*c = CircleId{RelatedPoint{Point{s.X, s.Y}, circleCounter}, s.R, s.R * s.R}
	circleCounter++
	return nil
}

/*func (cs *CircleStorage) UnmarshalJSON(b []byte) error {
	c := make([]CircleId, 0)
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}
	*cs = *newCircleStorage(c)
	return nil
}*/

func ParseMap(filename string) (Map, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Can't read")
		//return Map{}, err
	}
	res := Map{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		panic("Can't parse")
		//return Map{}, err
	}
	res.Storage = *newCircleStorage(res.Circles)
	return res, err
}

func LinearRun(start IPoint, end IPoint, gameMap *Map) []Point {
	length := Dist(start, end)
	vector := Direction(start, end)
	t := math.Sqrt(length/gameMap.Fmax) / gameMap.Dt
	n := int(math.Ceil(t))
	a := vector
	a.ScalarMult(length / sqr(gameMap.Dt*float64(n)))
	as := make([]Point, 2*n)
	for i := 0; i < n; i++ {
		as[i] = a
	}
	a.ScalarMult(-1)
	for i := n; i < 2*n; i++ {
		as[i] = a
	}
	//fmt.Println("LinearRun", t, n)
	return as
}

func Run(x0, v0 *Point, a []Point, gameMap *Map) ([]Point, []Point) {
	x := make([]Point, len(a))
	v := make([]Point, len(a))
	x[0] = *x0
	v[0] = a[0]
	v[0].ScalarMult(gameMap.Dt)
	v[0].Add(v0)
	for i := 1; i < len(a); i++ {
		x[i] = v[i-1]
		x[i].ScalarMult(gameMap.Dt)
		x[i].Add(&x[i-1])
		v[i] = a[i-1]
		v[i].ScalarMult(gameMap.Dt)
		v[i].Add(&v[i-1])
	}
	return x, v
}

func Strategy(gameMap *Map) ([]Point, []Point, []Point) {
	bestPathLength := math.MaxInt64
	bestX := make([]Point, 0)
	bestV := make([]Point, 0)
	bestA := make([]Point, 0)
	for _, add := range [...]float64{0.001, 0.005, 0.01, 0.03} {
		vertices, dist, start, end := MakeGraph(gameMap, add)
		fmt.Println("Made graph", len(vertices))
		path := MinPath(dist, start, end)
		fmt.Println("Made path", len(path))

		//fmt.Println(path)
		//fmt.Println("WAY:")

		/*for _, ind := range path {
			fmt.Println(vertices[ind])
		}*/
		resX := make([]Point, 0)
		resV := make([]Point, 0)
		resA := make([]Point, 0)
		zeroV := Point{0, 0}
		for i := 1; i < len(path); i++ {
			a := LinearRun(&vertices[path[i-1]], &vertices[path[i]], gameMap)
			x, v := Run(&vertices[path[i-1]].Point, &zeroV, a, gameMap)
			resA = append(resA, a...)
			resX = append(resX, x...)
			resV = append(resV, v...)
		}
		if len(resA) > 0 && bestPathLength > len(resA) {
			bestX = resX
			bestV = resV
			bestA = resA
			bestPathLength = len(resA)
		}
		fmt.Println("add = ", add, "  length = ", len(resA))
	}
	return bestX, bestV, bestA
}
