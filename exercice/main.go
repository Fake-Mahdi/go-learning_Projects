package main

import (
	"fmt"
	"sync"
)

type Order struct {
	ID       int
	Costumer string
	Dish     string
}
type Worker struct {
	ID int
}

func createWorker(number int) []Worker {
	s_Worker := []Worker{}
	for i := 1; i <= number; i++ {
		worker := Worker{ID: i}
		s_Worker = append(s_Worker, worker)
	}
	return s_Worker
}

func (w *Worker) ProcessOrder(channel chan Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for ch := range channel {
		fmt.Printf("Worker %d is processing Order %d for Customer %s \n", w.ID, ch.ID, ch.Costumer)
	}
}
func main() {
	worker := createWorker(3)
	fmt.Println(worker)
	orders := []Order{
		{ID: 1, Costumer: "Alice", Dish: "Pasta"},
		{ID: 2, Costumer: "Bob", Dish: "Burger"},
		{ID: 3, Costumer: "Charlie", Dish: "Sushi"},
	}
	channel := make(chan Order, len(orders))
	var wg sync.WaitGroup
	for i, order := range orders {
		wg.Add(1)
		go worker[i].ProcessOrder(channel, &wg)

		channel <- order
	}
	close(channel)
	wg.Wait()
}
