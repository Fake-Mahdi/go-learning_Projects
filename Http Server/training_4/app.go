package main

import (
	"fmt"
)

func taskA() {
	fmt.Println("Hello from Task A")
}

func taskB() {
	fmt.Println("Hello from Task B")
}

func taskC() {
	fmt.Println("Hello from Task C")
}
func main() {

	tasks := []func(){taskA, taskB, taskC}

	for _, t := range tasks {
		func(f func()) {
			f()
		}(t)
	}
}
