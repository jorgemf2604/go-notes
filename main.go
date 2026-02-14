package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5}
	b := a[2:]
	fmt.Printf("Value of A: %v, length of A: %v, capacity of A: %v\n", a, len(a), len(a))
	fmt.Printf("Value of B: %v, length of B: %v, capacity of B: %v\n", b, len(b), len(b))
	b[0] = 99
	fmt.Printf("Value of A: %v, length of A: %v, capacity of A: %v\n", a, len(a), len(a))
	fmt.Printf("Value of B: %v, length of B: %v, capacity of B: %v\n", b, len(b), len(b))

	// Value of A: [1 2 3 4 5], length of A: 5, capacity of A: 5
	// Value of B: [3 4 5], length of B: 3, capacity of B: 3
	// Value of A: [1 2 99 4 5], length of A: 5, capacity of A: 5
	// Value of B: [99 4 5], length of B: 3, capacity of B: 3
}
