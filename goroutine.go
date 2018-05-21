package main

import (
	"fmt"
)

//func Add(x, y int) {
//	z := x + y
//	fmt.Println(z)
//}
//
//func main() {
//	for i := 0; i < 10; i++ {
//		go Add(i, i)
//	}
//}

//func Count(ch chan int) {
//	fmt.Println("Counting Task")
//	ch <- 1
//}
//
//func main() {
//	chs := make([]chan int, 10)
//	for i := 0; i < 10; i++ {
//		chs[i] = make(chan int)
//		go Count(chs[i])
//	}
//
//	for _, ch := range(chs) {
//		// 在主routine阻塞，等待 ch被写入int，然后读出来
//		i := <- ch
//		fmt.Println(i)
//	}
//	fmt.Println("end")
//}

func Count(ch chan interface{}, i int) {
	fmt.Println("Counting Task")
	if i == 0 {
		ch <- "stop"
	} else {
		ch <- 1
	}
}

func main() {
	chs := make([]chan interface{}, 10)
	for i := 0; i < 10; i++ {
		chs[i] = make(chan interface{})
		go Count(chs[i], i)
	}

	for _, ch := range(chs) {
		// 在主routine阻塞，等待 ch被写入int，然后读出来
		i := <- ch
		if i == "stop" {
			fmt.Println("hehe")
			break
		}
		fmt.Println(i)
	}
	fmt.Println("end")
}