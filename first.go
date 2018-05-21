package main

import (
	"fmt"
	"os"
	"math"
	"github.com/mikezone/goext"
)

func main() {
	f := func(x, y int) int {
		return x + y
	}
	fmt.Println(f(1, 2))

	const Pi float64 = 3.1415926535897932384
	fmt.Println(Pi)
	LOCALMODE := os.Getenv("LOCAL_MODE")
	fmt.Println(LOCALMODE + " ")

	const name = iota
	const name2 = iota

	const (
		name3 = iota
		name4 = iota
		name5 = iota
	)
	fmt.Println(name, name2, name3, name4, name5)
	fmt.Println(IsEqual(1.0, 1.01, 0.01))
	fmt.Println(math.Dim(2.0, 1.01))
	fmt.Println(math.Dim(1.0, 1.01))

	str := "hello world"
	str = "中国"
	c := str[0]
	fmt.Println(str, len(str), rune(c))

	str = "Hello, 世界"
	for i, ch := range str {
		fmt.Println(i, ch)
	}

	array := [5]int{11, 22, 33, 44, 55}
	fmt.Println(array)
	for i, val := range array {
		fmt.Println(i, val)
	}

	// ways to create slice
	// 1. from array
	//var myArray [10]int = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//var mySlice []int = myArray[:5]

	// 2. use `make`
	//mySlice1 := make([]int, 5)

	// 3. direct create
	//mySlice3 := []int{1, 2, 3, 4, 5}

	//mySlice1 = append(mySlice1, mySlice3...)

	// 4. from slice
	//oldSlice := []int{1, 2, 3, 4, 5}
	//newSlice := oldSlice[:3] // 基于oldSlice的前3个元素构建新数组切片

	// 5. copy
	//slice1 := []int{1, 2, 3, 4, 5}
	//slice2 := []int{5, 4, 3}
	//copy(slice2, slice1) // 只会复制slice1的前3个元素到slice2中 copy(slice1, slice2) // 只会复制slice2的3个元素到slice1的前3个位置

	ClosureTest()
	//os.PathError{}
	goext.Hello()
	//flag.String()
	//flag.String()

	// 类型查询
	var val interface{} = 10
	switch val.(type) {
	case int:
		fmt.Println("整型")
	default:
		fmt.Println("未知")
	}
	// 查询接口

	var xx interface{} = tt{10}
	if _, ok := xx.(Reader); ok {
		fmt.Println(ok)
	}
	var yy interface{} = &tt{10}
	if _, ok := yy.(*tt); ok {
		fmt.Println(ok)
	}
	var char interface{} = '中'
	if _, ok := char.(rune); ok {
		fmt.Println("rune")
	}

}

type Reader interface {
	Read(buf []byte) (n int, err error)
}

type tt struct {
	a int
}


func (this tt) Read(buf []byte) (n int, err error) {
	return 1, nil
}

func ClosureTest() {
	var j int = 5
	a := func() (func()) {
		var i int = 10
		return func() {
			fmt.Println("i, j： %d, %d\n", i, j)
		}
	}()

	a()
	j *= 2
	a()
	//error
	//logger := &log.Logger{}
	//logger.SetOutput()
}

func IsEqual(f1, f2, p float64) bool {
	return math.Dim(f1, f2) < p
}
