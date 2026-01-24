package main

import "fmt"

type Employee struct {
	name    string
	age     int
	balance int
}

func (e *Employee) deposit(amount int) {
	e.balance += amount
}

func main() {
	employee1 := Employee{"Jorge", 45, 1000}
	employee1.deposit(2000)
	fmt.Println(employee1.balance) // 3000
}
