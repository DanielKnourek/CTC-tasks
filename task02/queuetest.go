package pain

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func register(wg *sync.WaitGroup, queue <-chan int) {
	defer wg.Done()
	for job := range queue {
		time.Sleep(time.Duration((float32(time.Second) * randRange(0.5, 2))))
		fmt.Println(job)
	}
	fmt.Println("Register worker Done!")

}
func main() {
	fmt.Println("-- Start --")
	// fmt.Println(time.Duration((float32(time.Second) * randRange(0.5, 2))))

	myqueue := make(chan int, 2)
	registerWG := new(sync.WaitGroup)
	registerWG.Add(2)

	go register(registerWG, myqueue)
	go register(registerWG, myqueue)

	myqueue <- 1
	fmt.Println("sent1")
	myqueue <- 2
	fmt.Println("sent2")
	myqueue <- 3
	fmt.Println("sent3")
	myqueue <- 4
	fmt.Println("sent4")
	myqueue <- 5
	fmt.Println("sent5")

	close(myqueue)
	registerWG.Wait()
	fmt.Println("--  End  --")
}

func randRange(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
