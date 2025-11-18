package main

import (
	"fmt"
	"sync"
)

type Task struct {
	ID       string
	Name     string
	Priority int
}
type Worker struct {
	ID int
}

var buckets = make([][]Task, 4)

func (w *Worker) handleProcess(channel chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range channel {
		handleTask(task.Name, task.Priority, w.ID)
	}
}

func handleTask(taskName string, priority int, workerID int) {
	switch taskName {
	case "Task 1":
		fmt.Printf("Worker %d: Executing Task 1 logic (Priority %d)\n", workerID, priority)
	case "Task 2":
		fmt.Printf("Worker %d: Executing Task 2 logic (Priority %d)\n", workerID, priority)
	case "Task 3":
		fmt.Printf("Worker %d: Executing Task 3 logic (Priority %d)\n", workerID, priority)
	case "Task 4":
		fmt.Printf("Worker %d: Executing Task 4 logic (Priority %d)\n", workerID, priority)
	case "Task 5":
		fmt.Printf("Worker %d: Executing Task 5 logic (Priority %d)\n", workerID, priority)
	case "Task 6":
		fmt.Printf("Worker %d: Executing Task 6 logic (Priority %d)\n", workerID, priority)
	case "Task 7":
		fmt.Printf("Worker %d: Executing Task 7 logic (Priority %d)\n", workerID, priority)
	case "Task 8":
		fmt.Printf("Worker %d: Executing Task 8 logic (Priority %d)\n", workerID, priority)
	case "Task 9":
		fmt.Printf("Worker %d: Executing Task 9 logic (Priority %d)\n", workerID, priority)
	case "Task 10":
		fmt.Printf("Worker %d: Executing Task 10 logic (Priority %d)\n", workerID, priority)
	default:
		fmt.Printf("Worker %d: Processing normal task %s (Priority %d)\n", workerID, taskName, priority)
	}
}

func createWorker(number int) []Worker {
	slWorker := []Worker{}
	for i := 1; i <= number; i++ {
		worker := Worker{ID: i}
		slWorker = append(slWorker, worker)
	}
	return slWorker
}

func SperatePriorityBuckets(tasks []Task) {
	for priority := 1; priority <= 3; priority++ {
		for _, task := range tasks {
			if task.Priority == priority {
				buckets[priority-1] = append(buckets[priority-1], task)
			}
		}
	}
}
func main() {
	var wg sync.WaitGroup
	tasks := []Task{
		{ID: "1", Name: "Task 1", Priority: 1},
		{ID: "2", Name: "Task 2", Priority: 2},
		{ID: "3", Name: "Task 3", Priority: 3},
		{ID: "4", Name: "Task 4", Priority: 1},
		{ID: "5", Name: "Task 5", Priority: 2},
		{ID: "6", Name: "Task 6", Priority: 3},
		{ID: "7", Name: "Task 7", Priority: 1},
		{ID: "8", Name: "Task 8", Priority: 2},
		{ID: "9", Name: "Task 9", Priority: 3},
		{ID: "10", Name: "Task 10", Priority: 1},
	}
	channel := make([]chan Task, 3)
	slWorker := createWorker(5)
	fmt.Println(slWorker)
	SperatePriorityBuckets(tasks)
	fmt.Println(buckets)

	for i := 1; i <= 3; i++ {
		channel[i-1] = make(chan Task, len(tasks))
	}

	for priority := 1; priority <= 3; priority++ {
		for i := 0; i < len(buckets[priority-1]); i++ {
			wg.Add(1)
			go slWorker[i].handleProcess(channel[priority-1], &wg)
		}
		for _, bucket := range buckets[priority-1] {
			channel[priority-1] <- bucket
		}
		close(channel[priority-1])
		wg.Wait()
	}

}
