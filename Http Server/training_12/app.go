package main

import (
	"fmt"
	"sync"
)

func task1(number1, number2 float64, result chan float64) {
	result <- number1 + number2
}

func task2(waitingResult chan float64, number2 float64, waitingResult2 chan float64) {
	val := <-waitingResult
	val += number2
	waitingResult2 <- val
}

func task3(waitingResult2 chan float64, number2 float64) {
	val := <-waitingResult2
	val += number2
	fmt.Println(val)
}

func main() {
	var wg sync.WaitGroup

	result := make(chan float64, 3)
	result2 := make(chan float64, 1)

	number1 := 15.0
	number2 := 60.0

	sliceOfFunc := []func(){
		func() { task1(number1, number2, result) },
		func() { task2(result, number2, result2) },
		func() { task3(result2, number1) },
	}

	for _, function := range sliceOfFunc {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			function()
		}(function)
	}

	wg.Wait()
}
