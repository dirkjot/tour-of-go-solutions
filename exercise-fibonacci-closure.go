package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacciGen() func() int {
	alpha, beta := 1, 1
	inner := func() (result int) {
		result = alpha
		alpha, beta = beta, alpha+beta
		return 	}
	return inner
}

func main() {
	f := fibonacciGen()
	g := fibonacciGen()
	g()
	for i := 0; i < 10; i++ {
		fmt.Println(f(), g())
	}
}
