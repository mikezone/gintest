package main

//import "fmt"
//

//fmt.Println(f(1, 2))

import "strconv"
import "fmt"

func main() {
	f := func(x, y int) int {
		return x + y
	}
	fmt.Println(f(1, 2))
	fmt.Println(strconv.Atoi("10"))
}
