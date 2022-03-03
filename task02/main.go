package main

import (
	"fmt"
	"math/rand"
	"sync"
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

type car struct {
	carFuel fuelType
	SPZ     string
}

type fuelType int64

const (
	fuelGas fuelType = iota
	fuelDiesel
	fuelLpg
	fuelElectric
)

type fuelPumpProps struct {
	pumpFuelType fuelType
	waitTime     [2]float32
}

func createPump(fuel fuelType, waitLow float32, waitHigh float32) fuelPumpProps {
	return fuelPumpProps{
		pumpFuelType: fuel,
		waitTime:     [2]float32{waitLow, waitHigh},
	}
}

func reFuelJob(jobProps fuelPumpProps, registerQ chan car) func(wg *sync.WaitGroup, queue chan car) {
	return func(wg *sync.WaitGroup, queue chan car) {
		defer wg.Done()
		for job := range queue {
			time.Sleep(time.Duration((float32(time.Second) * randRange(jobProps.waitTime[0], jobProps.waitTime[1]))))
			fmt.Println("Car refueled", job, " Lets go pay it!")
			registerQ <- job
		}
		fmt.Printf("reFuelJob type[%d] Done!\n", int(jobProps.pumpFuelType))
	}
}

type registerProps struct {
	waitTime [2]float32
}

func createRegister(waitLow float32, waitHigh float32) registerProps {
	return registerProps{
		waitTime: [2]float32{waitLow, waitHigh},
	}
}

func registerJob(jobProps registerProps) func(wg *sync.WaitGroup, queue chan car) {
	return func(wg *sync.WaitGroup, queue chan car) {
		defer wg.Done()
		for job := range queue {
			time.Sleep(time.Duration((float32(time.Second) * randRange(jobProps.waitTime[0], jobProps.waitTime[1]))))
			fmt.Println("Payed: ", job)
		}
		fmt.Println("registerJob Done!")
	}
}

func drivewayJob(fuelPump ...chan car) func(wg *sync.WaitGroup, queue chan car) {
	return func(wg *sync.WaitGroup, queue chan car) {
		defer wg.Done()
		for job := range queue {
			fmt.Println("New Car arrived", job)
			fuelPump[0] <- job
		}
		fmt.Printf("Station driveway empty\n")
	}
}

func simulation(numCustomers int, minArriveTime float32, maxArriveTime float32) {
	registerWG := new(sync.WaitGroup)
	pumpWG := new(sync.WaitGroup)
	drivewayWG := new(sync.WaitGroup)

	jobCreator := func(job func(wg *sync.WaitGroup, queue chan car), wg *sync.WaitGroup, count int) chan car {
		jobQueue := make(chan car, count)
		for i := 0; i < count; i++ {
			wg.Add(1)
			go job(wg, jobQueue)
		}
		return jobQueue
	}
	registerQ := jobCreator(registerJob(createRegister(0.5, 2)), registerWG, 2)

	gasQ := jobCreator(reFuelJob(createPump(fuelGas, 1, 5), registerQ), pumpWG, 4)
	dieselQ := jobCreator(reFuelJob(createPump(fuelDiesel, 1, 5), registerQ), pumpWG, 4)
	lpgQ := jobCreator(reFuelJob(createPump(fuelLpg, 1, 5), registerQ), pumpWG, 1)
	electricQ := jobCreator(reFuelJob(createPump(fuelElectric, 3, 10), registerQ), pumpWG, 8)

	drivewayQ := jobCreator(drivewayJob(gasQ, dieselQ, lpgQ, electricQ), drivewayWG, 1)

	for i := 1; i <= numCustomers; i++ {
		randfuel := fuelType(rand.Intn(int(fuelElectric)))
		newCustomer := car{
			carFuel: randfuel,
			SPZ:     fmt.Sprintf("5L%1d %04d", randfuel, i),
		}
		drivewayQ <- newCustomer
		time.Sleep(time.Duration((float32(time.Second) * randRange(minArriveTime, maxArriveTime))))
	}

	close(drivewayQ)
	drivewayWG.Wait()
	close(gasQ)
	close(lpgQ)
	close(dieselQ)
	close(electricQ)
	pumpWG.Wait()
	close(registerQ)
	registerWG.Wait()
}

func randRange(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func main() {
	fmt.Println("-- Start --")
	// simulation(5, 0.001, 0.1)
	simulation(100, 0.001, 0.1)
	// simulation(1000, 0.001, 0.1)
	fmt.Println("--  End  --")

}
