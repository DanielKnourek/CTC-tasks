package main

import (
	"fmt"
	// "sync"
	"math/rand"
	"time"
)

type fuelType int64

const (
	fuelGas fuelType = iota
	fuelDiesel
	fuelLpg
	fuelElectric
)

type fuel_pump struct {
	pumpFuelType fuelType
	waitTime     [2]float32
}

type register struct {
	waitTime [2]float32
}

func main() {
	fmt.Println("-- Start --")
	// simulation(1000, 0.001, 0.1)
	fmt.Println(fuelType(randRange(0, float32(fuelElectric)+1)))
}

func simulation(numCustomers int, minArriveTime float32, maxArriveTime float32) {
	stations := createStation()
	registers := createRegisters()
	for i := numCustomers; i > 0; i-- {
		time.Sleep(time.Duration((float32(time.Second) * randRange(minArriveTime, maxArriveTime))))
		// go simDriver(fuelType())
		// totalTime += time.Duration((float32(time.Second) * randRange(minArriveTime, maxArriveTime)))
	}

	fmt.Println(stations)
	fmt.Println(registers)
}

func simDriver(carType fuelType) {
	// go to queue pump
	// sim fueling
	// go to queue register
	// sim paying

}

func simPump(simDuration time.Duration) {
	time.Sleep(simDuration)
}
func simRegister(simDuration time.Duration) {
	time.Sleep(simDuration)
}

func randRange(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func createStation() [17]fuel_pump {
	createPump := func(fuel fuelType, waitLow float32, waitHigh float32) fuel_pump {
		return fuel_pump{
			pumpFuelType: fuel,
			waitTime:     [2]float32{waitLow, waitHigh},
		}
	}
	return [17]fuel_pump{
		createPump(fuelGas, 1, 5),
		createPump(fuelGas, 1, 5),
		createPump(fuelGas, 1, 5),
		createPump(fuelGas, 1, 5),
		createPump(fuelDiesel, 1, 5),
		createPump(fuelDiesel, 1, 5),
		createPump(fuelDiesel, 1, 5),
		createPump(fuelDiesel, 1, 5),
		createPump(fuelLpg, 1, 5),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
		createPump(fuelElectric, 3, 10),
	}
}

func createRegisters() [2]register {
	createRegister := func(waitLow float32, waitHigh float32) register {
		return register{
			waitTime: [2]float32{waitLow, waitHigh},
		}
	}
	return [2]register{
		createRegister(0.5, 2),
		createRegister(0.5, 2),
	}
}
