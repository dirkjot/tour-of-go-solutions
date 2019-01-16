package main

import (
	"fmt"
)

type FloatError struct {
	value float64
	msg string
}

func (f FloatError) Error() string {
	return fmt.Sprintf("Cannot handle this float input (%v) because %v", f.value, f.msg)
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return -1, FloatError{x, "less than zero"}
	}
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
	return guess, nil
}

func Sqrtxx(x float64) (float64, error) {
	return 0, nil 
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
