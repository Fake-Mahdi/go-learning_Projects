package main

import "fmt"

type Person struct {
	Name string
	Age  int
	City string
}

func new_main() {
	list := []Person{}
	people := []map[string]interface{}{
		{"Name": "Alice", "Age": 25, "City": "Paris"},
		{"Name": "Bob", "Age": 30, "City": "London"},
		{"Name": "Charlie", "Age": 22, "City": "Berlin"},
		{"Name": "David", "Age": 28, "City": "Madrid"},
		{"Name": "Eve", "Age": 26, "City": "Rome"},
		{"Name": "Frank", "Age": 33, "City": "Lisbon"},
		{"Name": "Grace", "Age": 27, "City": "Vienna"},
		{"Name": "Hannah", "Age": 24, "City": "Prague"},
		{"Name": "Ian", "Age": 29, "City": "Amsterdam"},
		{"Name": "Judy", "Age": 31, "City": "Copenhagen"},
	}

	for _, item := range people {
		p := Person{Name: item["Name"].(string), Age: item["Age"].(int), City: item["City"].(string)}
		list = append(list, p)

	}
	fmt.Println(list)
}

func main() {
	new_main()
}
