package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"golangbook.io/ch7/tempconv"
)

var (
	period = flag.Duration("period", 1*time.Second, "sleep period")
	temp   = tempconv.CelsiusFlag("temp", 20.0, "the temperature")
)

func main() {
	// fmt.Println("ch7")
	// var c bytecounter.ByteCounter
	// c.Write([]byte("hello"))
	//
	// fmt.Println(c)
	//
	// c = 0
	//
	// name := "Dolly"
	//
	// fmt.Fprintf(&c, "hello, %s\n", name)
	// fmt.Println(c)
	//
	// flag.Parse()
	// fmt.Println(*temp)
	// fmt.Printf("Sleeping for %v...", *period)
	// time.Sleep(*period)
	// fmt.Println()
	//
	// sorting.Exec()

	// http1.Exec()
	//
	var w io.Writer
	w = os.Stdout

	f, ok := w.(*os.File)
	fmt.Printf("%v, %v", f, ok)
	c, ok := w.(*bytes.Buffer)
	fmt.Printf("%v, %v", c, ok)
	println(c, ok)

	_, err := os.Open("/no/such/file")
	fmt.Println(os.IsNotExist(err)) // "true"
}

func sqlQuote(x interface{}) string {
	switch x := x.(type) {
	case nil:
		return "NULL"
	case int, uint:
		return fmt.Sprintf("%d", x) // x has type interface{} here.
	case bool:
		if x {
			return "TRUE"
		}
		return "FALSE"
	case string:
		return x // (not shown)
	default:
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
}
