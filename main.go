package main

import (
	"fmt"
)

func main() {

	fmt.Println("Enter Your City Name: ")
	var city string
	fmt.Scanln(&city)
	var days int
	fmt.Println("Select the number of days you would like to see the weather forecast for: ")
	fmt.Scanln(&days)

	fmt.Println("select the propriate option")
	fmt.Println("1. Will it rain?")
	fmt.Println("2. When is the next rain day")
	fmt.Println("3. Avg temp")
	fmt.Println("4. min / max temp for each day")
	fmt.Println("5. Today Weather Summary")

	var option int
	fmt.Scanln(&option)

	switch option {
	case 1:
		fmt.Println(willItRainClient(city, days))
	case 2:
		fmt.Println(nextRainDay(city, days))
	case 3:
		fmt.Println(avgTempretureClient(city, days))
	case 4:
		fmt.Println(minMaxTempTill(city, days))
	case 5:
		fmt.Println(todaySum(city))
	default:
		fmt.Printf("please enter a propriate number")
	}

	// mainCal(city, days)

}
