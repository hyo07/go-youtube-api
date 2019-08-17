package main

import (
	"api_test/callAPI"
	"fmt"
)

func main() {
	status, message := callAPI.GetVideo("https://www.youtube.com/watch?v=ad-6W3bazMs&list=PLpt61bADOMwW2aA9I1hWHyR8OixSuCXzc&index=2")
	fmt.Println(status)
	fmt.Println(message)
}
