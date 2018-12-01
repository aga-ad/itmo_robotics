package main

import (
	"math"
	"math/rand"
)

var eps float64 = 0.0000000001

func sqr(a float64) float64 {
	return a * a
}

type IPoint interface {
	X() float64
	Y() float64
	ScalarMult(float64)
	Add(IPoint)
}

type Line struct {
	a, b, c float64
}

type Segment struct {
	p1, p2 Point
}

func LineFromPoints(p1 IPoint, p2 IPoint) Line {
	var line Line
	line.a = p2.Y() - p1.Y()
	line.b = p1.X() - p2.X()
	line.c = -line.a*p1.X() - line.b*p1.Y()
	return line
}

type Point struct {
	x, y float64
}

type RelatedPoint struct {
	Point
	id int
}

func (p *Point) X() float64 {
	return p.x
}

func (p *Point) Y() float64 {
	return p.y
}

func Norm(p IPoint) float64 {
	return math.Sqrt(sqr(p.X()) + sqr(p.Y()))
}

func Direction(p1 IPoint, p2 IPoint) Point {
	res := Point{p2.X() - p1.X(), p2.Y() - p1.Y()}
	res.ScalarMult(1 / Norm(&res))
	return res
}

func (p *Point) ScalarMult(a float64) {
	p.x *= a
	p.y *= a
}

func (p1 *Point) Add(p2 IPoint) {
	p1.x += p2.X()
	p1.y += p2.Y()
}

func Dist2(p1 IPoint, p2 IPoint) float64 {
	return sqr(p1.X()-p2.X()) + sqr(p1.Y()-p2.Y())
}

func Dist(p1 IPoint, p2 IPoint) float64 {
	return math.Sqrt(Dist2(p1, p2))
}

type ICircle interface {
	IPoint
	R() float64
	R2() float64
}

type Circle struct {
	Point
	r, r2 float64
}

type CircleId struct {
	RelatedPoint
	r, r2 float64
}

func (c *CircleId) R() float64 {
	return c.r
}

func (c *CircleId) R2() float64 {
	return c.r2
}

func (c *Circle) R() float64 {
	return c.r
}

func (c *Circle) R2() float64 {
	return c.r2
}

func (p *CircleId) GetPoint(angle float64, add float64) RelatedPoint {
	return RelatedPoint{Point{p.X() + (p.R()+add)*math.Cos(angle), p.Y() + (p.R()+add)*math.Sin(angle)}, p.id}
}

func Intersect(c1 ICircle, c2 ICircle) bool {
	return Dist2(c1, c2) < sqr(c1.R()+c2.R())+eps
}

func (c1 *Circle) IntersectCircle(c2 ICircle) bool {
	return Dist2(c1, c2) < sqr(c1.R()+c2.R())+eps
}

func (c1 *CircleId) IntersectCircle(c2 ICircle) bool {
	return Dist2(c1, c2) < sqr(c1.R()+c2.R())+eps
}

func AnyIn(c1 ICircle, c2 ICircle) bool {
	return Dist2(c1, c2) < sqr(c1.R()-c2.R())+eps
}

func In(c ICircle, p IPoint) bool {
	return Dist2(c, p) < c.R()+eps
}

func (p *Point) IntersectCircle(c ICircle) bool {
	return Dist2(c, p) < c.R()+eps
}

func (p *RelatedPoint) IntersectCircle(c ICircle) bool {
	return Dist2(c, p) < c.R()+eps
}

// TangentsPoints returns slice of zero, four of eight points.
func TangentsPoints(c1 *CircleId, c2 *CircleId, add float64) []RelatedPoint {
	if AnyIn(c1, c2) {
		return []RelatedPoint{}
	}
	alpha := math.Atan2(math.Abs(c1.R()-c2.R()), Dist(c1, c2))
	beta := math.Atan2(c2.Y()-c1.Y(), c2.X()-c1.X())
	res := make([]RelatedPoint, 0, 4)
	res = append(res, c1.GetPoint(math.Pi/2+alpha+beta, add))
	res = append(res, c2.GetPoint(math.Pi/2+alpha+beta, add))
	res = append(res, c1.GetPoint(-math.Pi/2-alpha+beta, add))
	res = append(res, c2.GetPoint(-math.Pi/2-alpha+beta, add))
	if !Intersect(c1, c2) {
		gamma := math.Acos((c1.R() + c2.R()) / Dist(c1, c2))
		//fmt.Println("gamma: ", gamma/math.Pi*180)
		res = append(res, c1.GetPoint(math.Pi/2-gamma+beta, add))
		res = append(res, c2.GetPoint(math.Pi/2-gamma+beta, add))
		res = append(res, c1.GetPoint(-math.Pi/2+gamma+beta, add))
		res = append(res, c2.GetPoint(-math.Pi/2+gamma+beta, add))
	}
	return res
}

func TangentsPointsFromPoint(c *CircleId, p IPoint, add float64) []RelatedPoint {
	if In(c, p) {
		return []RelatedPoint{}
	}
	alpha := math.Atan2(c.R(), Dist(c, p))
	beta := math.Atan2(p.Y()-c.Y(), p.X()-c.X())
	res := make([]RelatedPoint, 0, 2)
	res = append(res, c.GetPoint(math.Pi/2+alpha+beta, add))
	res = append(res, c.GetPoint(-math.Pi/2-alpha+beta, add))
	return res
}

func (s *Segment) IntersectCircle(c ICircle) bool {
	x01 := s.p1.x - c.X()
	y01 := s.p1.y - c.Y()
	x02 := s.p2.x - c.X()
	y02 := s.p2.y - c.Y()
	dx := x02 - x01
	dy := y02 - y01
	a := dx*dx + dy*dy
	b := 2.0 * (x01*dx + y01*dy)
	cc := x01*x01 + y01*y01 - c.R2()
	if -b < eps {
		return cc < eps
	}
	if -b < (2.0 * a) {
		return 4.0*a*cc-b*b < eps
	}
	return a+b+cc < eps
}

type CircleStorage struct {
	circles  [][]CircleId
	clusters []Circle
	All      []CircleId
}

func newCircleStorage(circles []CircleId) *CircleStorage {
	cs := CircleStorage{}
	cs.clusters = make([]Circle, int(math.Sqrt(float64(len(circles)))))
	cs.circles = make([][]CircleId, len(cs.clusters))
	cs.All = circles
	for i := 0; i < len(cs.clusters); i++ {
		cs.clusters[i].x = rand.Float64()
		cs.clusters[i].y = rand.Float64()
		cs.circles[i] = make([]CircleId, 0)
	}
	for _, c := range circles {
		nearest := 0
		for i := 1; i < len(cs.clusters); i++ {
			if Dist2(&cs.clusters[nearest], &c.Point) > Dist2(&cs.clusters[i], &c.Point) {
				nearest = i
			}
		}
		cs.circles[nearest] = append(cs.circles[nearest], c)
	}
	for i := 0; i < len(cs.clusters); i++ {
		if len(cs.circles[i]) == 0 {
			cs.circles[i] = cs.circles[len(cs.circles)-1]
			cs.circles = cs.circles[:len(cs.circles)-1]
			cs.clusters[i] = cs.clusters[len(cs.clusters)-1]
			cs.clusters = cs.clusters[:len(cs.clusters)-1]
			i--
			continue
		}
		for j := 0; j < len(cs.circles[i]); j++ {
			cs.clusters[i].r = math.Max(cs.clusters[i].r, Dist(&cs.clusters[i].Point, &cs.circles[i][j].Point)+cs.circles[i][j].R())
		}
		cs.clusters[i].r2 = sqr(cs.clusters[i].r)
	}
	return &cs
}

func (cs *CircleStorage) Intersect(figure interface{ IntersectCircle(c ICircle) bool }) bool {
	for i := 0; i < len(cs.clusters); i++ {
		if figure.IntersectCircle(&cs.clusters[i]) {
			for j := 0; j < len(cs.circles[i]); j++ {
				if figure.IntersectCircle(&cs.circles[i][j]) {
					return true
				}
			}
		}
	}
	return false
}
