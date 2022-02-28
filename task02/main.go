package main

import "fmt"

type fuelType int64;

const (
	fuelGas fuelType = iota
	fuelDiesel
	fuelLpg
	fuelElectric
)

type fuel_pump struct {
	pumpFuelType fuelType;
	waitTime[2] int;
	count int;
}

func createPump(fuel fuelType, waitLow int, waitHigh int) fuel_pump {
	return fuel_pump{
		pumpFuelType: fuel,
		waitTime: [2]int{waitLow, waitHigh},
	}
};

func main() {
	fmt.Println("-- Start --");

	stations := [17]fuel_pump{
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
	fmt.Println(stations)
}
