package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ttacon/chalk"
)

/*
Gas station simulation
auto se přidá do fronty stanice
	z fronty stanice se snaži zaparkovat do fronty svojí pumpy
		po natakování přejde do fronty pokladny
		po zaplacení uvolní pokladnu
	poté uvolní místo u dané pumpy

*/
var deltaDriveway chan int = make(chan int)
var deltaFuel chan fuelType = make(chan fuelType)
var deltaRegister chan int = make(chan int)

func animator() {
	qcDriveway, qcRegister := 0, 0
	qcGas, qcDiesel, qcLPG, qcElectric := 0, 0, 0, 0
	fmt.Print("\033[s")
	for {
		select {
		case valDeltaDriveway := <-deltaDriveway:
			qcDriveway += valDeltaDriveway
		case valDeltaFuel := <-deltaFuel:
			switch valDeltaFuel {
			case fuelGas:
				qcGas += 1
			case (fuelGas + 10):
				qcGas -= 1

			case fuelDiesel:
				qcDiesel += 1
			case (fuelDiesel + 10):
				qcDiesel -= 1

			case fuelLpg:
				qcLPG += 1
			case (fuelLpg + 10):
				qcLPG -= 1

			case fuelElectric:
				qcElectric += 1
			case (fuelElectric + 10):
				qcElectric -= 1
			}
		case valDeltaRegister := <-deltaRegister:
			qcRegister += valDeltaRegister
		default:
			fmt.Println("\033[u\033[0J")
			fmt.Printf("Cars in station \t %3d %s \nIn-use\n  Gas: \t\t\t %3d %s \n  Diesel: \t\t %3d %s \n  LPG: \t\t\t %3d %s \n  Electric: \t\t %3d %s \n\nRegister \t\t %3d %s \n\n",
				qcDriveway, strings.ReplaceAll(fmt.Sprintf("%s%-50s%s", chalk.Blue, strings.Repeat("█", qcDriveway), chalk.Reset), " ", "░"),
				qcGas, strings.ReplaceAll(fmt.Sprintf("%s%-10s%s", chalk.Blue, strings.Repeat("█", qcGas), chalk.Reset), " ", "░"),
				qcDiesel, strings.ReplaceAll(fmt.Sprintf("%s%-10s%s", chalk.Blue, strings.Repeat("█", qcDiesel), chalk.Reset), " ", "░"),
				qcLPG, strings.ReplaceAll(fmt.Sprintf("%s%-10s%s", chalk.Blue, strings.Repeat("█", qcLPG), chalk.Reset), " ", "░"),
				qcElectric, strings.ReplaceAll(fmt.Sprintf("%s%-10s%s", chalk.Blue, strings.Repeat("█", qcElectric), chalk.Reset), " ", "░"),
				qcRegister, strings.ReplaceAll(fmt.Sprintf("%s%-2s%s", chalk.Blue, strings.Repeat("█", qcRegister), chalk.Reset), " ", "░"))
			time.Sleep(time.Millisecond * 50)
		}
	}
}

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
			deltaDriveway <- -1 // Animator
			time.Sleep(time.Duration((float32(time.Second) * randRange(jobProps.waitTime[0], jobProps.waitTime[1]))))
			// fmt.Println("Car refueled", job, " Lets go pay it!")
			registerQ <- job
			deltaFuel <- job.carFuel + 10 // Animator
		}
		// fmt.Printf("reFuelJob type[%d] Done!\n", int(jobProps.pumpFuelType))

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
			deltaRegister <- 1 // Animator
			time.Sleep(time.Duration((float32(time.Second) * randRange(jobProps.waitTime[0], jobProps.waitTime[1]))))
			// fmt.Println("Payed: ", job)
			job.carFuel = job.carFuel + job.carFuel - job.carFuel // Animator
			deltaRegister <- -1                                   // Animator
		}
		// fmt.Println("registerJob Done!")
	}
}

func drivewayJob(fuelPump ...chan car) func(wg *sync.WaitGroup, queue chan car) {
	return func(wg *sync.WaitGroup, queue chan car) {
		defer wg.Done()
		for job := range queue {
			// fmt.Println("New Car arrived", job)
			deltaFuel <- job.carFuel // Animator
			fuelPump[0] <- job
		}
		// fmt.Printf("Station driveway empty\n")
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
		randfuel := fuelType(rand.Intn(int(fuelElectric) + 1))
		newCustomer := car{
			carFuel: randfuel,
			SPZ:     fmt.Sprintf("5L%1d %04d", randfuel, i),
		}
		deltaDriveway <- 1 // Animator
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
	fmt.Println(chalk.Red, "Writing in colors", chalk.Cyan, "is so much fun", chalk.Reset)
	go animator()
	// simulation(5, 0.001, 0.1)
	simulation(100, 0.001, 0.1)
	// simulation(1000, 0.001, 0.1)
	fmt.Println("--  End  --")

}
