package main

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
)

/*
Gas station simulation
auto se přidá do fronty stanice
	z fronty stanice se snaži zaparkovat do fronty svojí pumpy
		po natakování přejde do fronty pokladny
		po zaplacení uvolní pokladnu
	poté uvolní místo u dané pumpy

*/


// type fuelType int64

// const (
// 	fuelGas fuelType = iota
// 	fuelDiesel
// 	fuelLpg
// 	fuelElectric
// )

// type fuel_pump struct {
// 	pumpFuelType fuelType
// 	waitTime     [2]float32
// }

// type register struct {
// 	waitTime [2]float32
// }

func registerJob(wg *sync.WaitGroup, queue <-chan int) {
	defer wg.Done()
	for job := range queue {
		time.Sleep(time.Duration((float32(time.Second) * randRange(0.5, 2))))
		fmt.Println(job)
	}
	fmt.Println("Register worker Done!")

}

func simulation(numCustomers int, minArriveTime float32, maxArriveTime float32) {
	jobWG := new(sync.WaitGroup)

	jobCreator := func(job func(wg *sync.WaitGroup) , count int) {
		for i := 0; i < count; i++ {
			jobWG.Add(1)
			go job(jobWG, )
		}
	}
	jobCreator(registerJob, 2)
	stations := createStation()
	registers := createRegisters()
	for i := numCustomers; i > 0; i-- {
		time.Sleep(time.Duration((float32(time.Second) * randRange(minArriveTime, maxArriveTime))))
		// go simDriver(fuelType())
		// totalTime += time.Duration((float32(time.Second) * randRange(minArriveTime, maxArriveTime)))
	}
}

func randRange(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func main() {
	fmt.Println("-- Start --")
	simulation(1000, 0.001, 0.1)
	fmt.Println("--  End  --")
	
}