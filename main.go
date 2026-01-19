package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	str := "Hello, 世界!"
	length := utf8.RuneCountInString(str)
	fmt.Printf("The number of characters in the string is: %d\n", length)
	fmt.Printf("The number of bytes in the string is: %d\n", len(str))
}

// The number of characters in the string is: 10
// The number of bytes in the string is: 14
