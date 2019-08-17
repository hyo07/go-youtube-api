package main

import (
	"api_test/callAPI"
	"fmt"
)

func main() {
	a := callAPI.GroupList()
	fmt.Println(a)

}
