package main

import (
	"fmt"
)

func main_menu() int {
	var choice int
	fmt.Println(`
	Welcome to GoBank!
	1. Create Account
	2. Deposit Money
	3. Withdraw Money
	4. Check Balance
	5. Transfer Money
	6. Exit
	Choose an option: 
	`)
	fmt.Scanln(&choice)
	return choice
}
func create_account() []map[string]interface{} {
	var rotation int
	var name string
	var lastname string
	var balance float32
	list := []map[string]interface{}{}
	fmt.Println("how many user are you going to add : ")
	fmt.Scanln(&rotation)

	for i := 0; i < rotation; i++ {
		fmt.Println("Write user name please : ")
		fmt.Scanln(&name)

		fmt.Println("Write lastname please : ")
		fmt.Scanln(&lastname)

		fmt.Println("how many balance u have please  : ")
		fmt.Scanln(&balance)

		amap := map[string]interface{}{
			"name":     name,
			"lastname": lastname,
			"balance":  balance,
		}

		list = append(list, amap)
	}
	return list
}
func deposit_money(account_list []map[string]interface{}) {
	var name string
	var depose_amount float32
	check_existence := false
	for _, account := range account_list {
		fmt.Printf("Name: %s, Lastname: %s, Balance: %.2f\n",
			account["name"].(string),
			account["lastname"].(string),
			account["balance"].(float32),
		)
	}
	fmt.Println("Write the name of the account that u want to depose ur money in : ")
	fmt.Scanln(&name)
	fmt.Println("How much u are Going to dipose Here")
	fmt.Scanln(&depose_amount)
	for _, account := range account_list {
		if account["name"].(string) == name {
			check_existence = true
			hold_balance := account["balance"].(float32)
			hold_balance += depose_amount
			account["balance"] = hold_balance
			break
		}
	}
	if !check_existence {
		fmt.Println("User account is not found")
	}
	for _, account := range account_list {
		fmt.Printf("Name: %s, Lastname: %s, Balance: %.2f\n",
			account["name"].(string),
			account["lastname"].(string),
			account["balance"].(float32),
		)
	}
}
func withdraw_money(account_list []map[string]interface{}) {
	var withdraw_balance float32
	var name string
	for i := 0; i < len(account_list); i++ {
		fmt.Printf("Name: %s, Lastname: %s, Balance: %.2f\n",
			account_list[i]["name"].(string),
			account_list[i]["lastname"].(string),
			account_list[i]["balance"].(float32),
		)
	}

	fmt.Println("How much money do you want to withdraw")
	fmt.Scanln(&withdraw_balance)
	fmt.Println("Write the account name from where you want to withdraw")
	fmt.Scanln(&name)

	for i := 0; i < len(account_list); i++ {
		if account_list[i]["name"].(string) == name {
			if account_list[i]["balance"].(float32) < withdraw_balance {
				fmt.Println("Sorry the withdraw amount is much bigger than your balance")
			} else {
				balance := account_list[i]["balance"].(float32)
				balance = balance - withdraw_balance
				account_list[i]["balance"] = balance
				fmt.Println("Operation has been done successfully")
			}
		}
	}
}

func check_balance(account_list []map[string]interface{}) {
	var account_name string
	for i := 0; i < len(account_list); i++ {
		fmt.Printf("Name: %s, Lastname: %s, Balance: %.2f\n",
			account_list[i]["name"].(string),
			account_list[i]["lastname"].(string),
			account_list[i]["balance"].(float32),
		)
	}
	fmt.Println("Enter the name of the Account that u are Looking for : ")
	fmt.Scanln(&account_name)

	found := false
	for i := 0; i < len(account_list); i++ {
		if account_list[i]["name"].(string) == account_name {
			fmt.Printf("Name: %s, Lastname: %s, Balance: %.2f\n",
				account_list[i]["name"].(string),
				account_list[i]["lastname"].(string),
				account_list[i]["balance"].(float32),
			)
			found = true
			break
		}
	}
	if !found {
		fmt.Println("No user Found")
	}
}

func transfer_money(account_list []map[string]interface{}) {
	var sender string
	var receiver string
	var transfer_balance float32

	fmt.Println("wich account user are you : ")
	fmt.Scanln(&sender)
	fmt.Println("wich account user will u send the money to : ")
	fmt.Scanln(&receiver)

	check_sender := false
	check_receiver := false
	for _, account_sender := range account_list {
		if account_sender["name"].(string) == sender {
			fmt.Printf("Your balance is : %.2f", account_sender["balance"].(float32))
			check_sender = true
			for _, account_receiver := range account_list {
				if account_receiver["name"].(string) == receiver {
					check_receiver = true
					fmt.Println("how much mony would u like to transfer : ")
					fmt.Scanln(&transfer_balance)
					if account_sender["balance"].(float32) < transfer_balance {
						fmt.Println("The transfer amount is larger than ur current balance")
					} else {
						balance := account_sender["balance"].(float32)
						balance -= transfer_balance
						account_sender["balance"] = balance
						receiver_balance := account_receiver["balance"].(float32)
						receiver_balance += transfer_balance
						account_receiver["balance"] = receiver_balance
					}
					break
				}
			}
			fmt.Printf("Your balance is : %.2f\n", account_sender["balance"].(float32))
			break
		}
	}
	if !check_sender {
		fmt.Println("No sender Found")
	} else if !check_receiver {
		fmt.Println("No receiver Found")
	}
}
func main() {
	var accounts []map[string]interface{}

	for {
		menu_choice := main_menu() // only call once!

		switch menu_choice {
		case 1:
			accounts = create_account()
		case 2:
			deposit_money(accounts)
		case 3:
			withdraw_money(accounts)
		case 4:
			check_balance(accounts)
		case 5:
			transfer_money(accounts)
		case 6:
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice! Please choose a number from 1 to 6.")
		}

		fmt.Println("\n------------------------------------")
	}
}
