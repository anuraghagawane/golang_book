package main

import (
	"fmt"
	"golangbook/ch6.io/geometry"
	"golangbook/ch6.io/intset"
)

func main() {
	perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
	// fmt.Println(geometry.PathDistance(perim))
	fmt.Println(perim.Distance())

	p := geometry.Point{1, 1}
	q := geometry.Point{2, 2}
	p.ScaleBy(5)
	fmt.Println(p)
	fmt.Println(p.Distance(q))

	var x, y intset.IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"
	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"
	x.UnionWith(&y)
	fmt.Println(x.String())           // "{1 9 42 144}"
	fmt.Println(x.Has(9), x.Has(123)) // "true false"
}
