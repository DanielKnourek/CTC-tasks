package main

import (
	"fmt"
	"time"
)

var numlol chan int = make(chan int)

func main() {
	fmt.Println("Test 0")
	fmt.Print("\033[s") //one up, remove line (should work after the newline of the Println)

	fmt.Println("Test 1")
	fmt.Println("Test 2")
	// fmt.Print("\033[1A\033[0J") //one up, remove line (should work after the newline of the Println)
	// fmt.Print("\033[1A\033[K")  //one up, remove line (should work after the newline of the Println)
	fmt.Println("Test 3")
	fmt.Println("Test 4")
	time.Sleep(time.Second)
	fmt.Println("\033[u\033[0J")

}
