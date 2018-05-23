package main

import "fmt"

func generateThreeValue() (int, int, int){
	return 10, 22, 33
}

func returnValue() (x int, y int) {
	x, y, _ = generateThreeValue()
	return
}

func main() {
	fmt.Print(returnValue())
}
