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

func createPump(fuel fuelType, waitLow int, waitHigh int, count int) fuel_pump {
	return fuel_pump{
		pumpFuelType: fuel,
		waitTime: [2]int{waitLow, waitHigh},
		count: count,
	}
};

func main() {
	fmt.Println("-- Start --");

	stations := [4]fuel_pump{
		createPump(fuelGas, 1, 5, 4),
		createPump(fuelDiesel, 1, 5, 4),
		createPump(fuelLpg, 1, 5, 1),
		createPump(fuelElectric, 3, 10, 8),
	}
	fmt.Println(stations)
}
