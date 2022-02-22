package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	// nk, nk+1
	nk, nk1 := 0, 1
	return func() int {
		result := nk
		nk, nk1 = nk1, nk+nk1
		return result
	}
}

func main() {
	f := fibonacci()
	f2 := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println("f\t", f())
		if i > 2 {
			fmt.Println("f2\t\t", f2())
		}
	}
}
