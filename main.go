package main

import (
	"fmt"
)

func main() {

	fmt.Println("Enter Your City Name: ")
	var city string
	fmt.Scanln(&city)
	fmt.Println("Select the number of days you would like to see the weather forecast for: ")
	var days int
	fmt.Scanln(&days)
	mainCal(city, days)

}
