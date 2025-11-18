package main

import (
	"fmt"
	"sync"
)

type Account struct {
	Owner   string
	Balance float64
	mu      sync.Mutex
}

func (a *Account) Deposit(amount float64) {
	a.mu.Lock()
	a.Balance += amount
	fmt.Println(a.Balance)
	a.mu.Unlock()

}
func (a *Account) Withdraw(amount float64) {
	a.mu.Lock()
	if amount > a.Balance {
		fmt.Println("U don't have enough balance")
		a.mu.Unlock()
		return
	}
	a.Balance -= amount
	fmt.Println(a.Balance)
	a.mu.Unlock()

}
func main() {
	var counter int
	var data float64
	sliceOfData := []float64{}
	fmt.Println("Enter how much element u want to test")
	n, err := fmt.Scanln(&counter)
	if err != nil || n != 1 {
		fmt.Println("Invalid input! Please enter a number.")
		return
	}

	if counter == 0 {
		fmt.Println("You entered zero, nothing to test.")
		return
	}

	for i := 0; i < counter; i++ {
		fmt.Println("Enter The values that u want Mahdi : ")
		n, err := fmt.Scanln(&data)
		if err != nil || n != 1 {
			fmt.Println("Invalid input! Please enter a number.")
			return
		}
		sliceOfData = append(sliceOfData, data)
	}
	fmt.Println("-------------------------------------------------------------------------------------------------")
	user := Account{Owner: "Mahdi", Balance: 27500}
	var wg sync.WaitGroup
	tasks := []func(float64){user.Deposit, user.Withdraw}
	for i, amount := range sliceOfData {
		wg.Add(1)
		go func(i int, money float64) {
			defer wg.Done()
			task := tasks[i%len(tasks)]
			task(amount)

		}(i, amount)
	}
	wg.Wait()
}
