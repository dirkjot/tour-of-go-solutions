package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	guess := 1.0
	tolerance := 0.001
	distance := 1.0
	count := 0
	for distance > tolerance && count < 40 {
		guess -= (guess*guess - x) / (2 * guess)
		if distance = x - (guess * guess); distance < 0 {
			distance = -distance
		}
		count++
		fmt.Printf("%v Guess: %v, distance %v\n", count, guess, distance)

	}
	return guess
}

func main() {
	fmt.Println(Sqrt(2435433))
}
