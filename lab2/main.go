package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Product struct {
	Stamping   int
	Finishing  int
	Assembling int
}
type Solution struct {
	Pots    Product
	Coffees Product
	Samovar Product
}

const (
	PotPrice     = 100.0
	CoffeePrice  = 80.0
	SamovarPrice = 200.0
)

var (
	initialTable = Solution{
		Pots: Product{
			Stamping:   100,
			Finishing:  70,
			Assembling: 80,
		},
		Coffees: Product{
			Stamping:   120,
			Finishing:  120,
			Assembling: 100,
		},
		Samovar: Product{
			Stamping:   50,
			Finishing:  40,
			Assembling: 25,
		},
	}
	ticketRate  = 2500   // Размер штрафа энергии
	minTemp     = 0.2    // нижний предел температуры
	startTemp   = 50.0   // стартовая температура
	coolingRate = 0.9999 // степень остывания
	swapTimes   = 3      // количество перестановок за итерацию
)

func (solution *Solution) CalculateEnergy() float64 {
	energy := PotPrice*float64(solution.Pots.Assembling) + CoffeePrice*float64(solution.Coffees.Assembling) + SamovarPrice*float64(solution.Samovar.Assembling)

	if solution.Pots.Assembling > solution.Pots.Finishing {
		energy -= (float64(solution.Pots.Assembling) - float64(solution.Pots.Finishing)) * float64(ticketRate)
	}
	if solution.Coffees.Assembling > solution.Coffees.Finishing {
		energy -= (float64(solution.Coffees.Assembling) - float64(solution.Coffees.Finishing)) * float64(ticketRate)
	}
	if solution.Samovar.Assembling > solution.Samovar.Finishing {
		energy -= (float64(solution.Samovar.Assembling) - float64(solution.Samovar.Finishing)) * float64(ticketRate)
	}

	if solution.Pots.Finishing > solution.Pots.Stamping {
		energy -= (float64(solution.Pots.Finishing) - float64(solution.Pots.Stamping)) * float64(ticketRate)
	}
	if solution.Coffees.Finishing > solution.Coffees.Stamping {
		energy -= (float64(solution.Coffees.Finishing) - float64(solution.Coffees.Stamping)) * float64(ticketRate)
	}
	if solution.Samovar.Finishing > solution.Samovar.Stamping {
		energy -= (float64(solution.Samovar.Finishing) - float64(solution.Samovar.Stamping)) * float64(ticketRate)
	}

	PotsAssembly := math.Floor((1 - float64(solution.Coffees.Assembling)/float64(initialTable.Coffees.Assembling) - float64(solution.Samovar.Assembling)/float64(initialTable.Samovar.Assembling)) * float64(initialTable.Pots.Assembling))
	CoffeeAssembly := math.Floor((1 - float64(solution.Pots.Assembling)/float64(initialTable.Pots.Assembling) - float64(solution.Samovar.Assembling)/float64(initialTable.Samovar.Assembling)) * float64(initialTable.Coffees.Assembling))
	SamovarAssembly := math.Floor((1 - float64(solution.Pots.Assembling)/float64(initialTable.Pots.Assembling) - float64(solution.Coffees.Assembling)/float64(initialTable.Coffees.Assembling)) * float64(initialTable.Samovar.Assembling))

	if PotsAssembly < float64(solution.Pots.Assembling) {
		energy -= (float64(solution.Pots.Assembling) - PotsAssembly) * float64(ticketRate)
	}
	if CoffeeAssembly < float64(solution.Coffees.Assembling) {
		energy -= (float64(solution.Coffees.Assembling) - CoffeeAssembly) * float64(ticketRate)
	}
	if SamovarAssembly < float64(solution.Samovar.Assembling) {
		energy -= (float64(solution.Samovar.Assembling) - SamovarAssembly) * float64(ticketRate)
	}

	Pot := math.Floor((1 - float64(solution.Pots.Stamping)/float64(initialTable.Pots.Stamping)) * float64(initialTable.Pots.Finishing))
	Coffee := math.Floor((1 - float64(solution.Coffees.Stamping)/float64(initialTable.Coffees.Stamping)) * float64(initialTable.Coffees.Finishing))
	Samovar := math.Floor((1 - float64(solution.Samovar.Stamping)/float64(initialTable.Samovar.Stamping)) * float64(initialTable.Samovar.Finishing))

	if Pot < float64(solution.Pots.Finishing) {
		energy -= (float64(solution.Pots.Finishing) - Pot) * float64(ticketRate)
	}
	if Coffee < float64(solution.Coffees.Finishing) {
		energy -= (float64(solution.Coffees.Finishing) - Coffee) * float64(ticketRate)
	}
	if Samovar < float64(solution.Samovar.Finishing) {
		energy -= (float64(solution.Samovar.Finishing) - Samovar) * float64(ticketRate)
	}
	return energy
}
func RandomSolution() Solution {
	return Solution{
		Pots: Product{
			Stamping:   rand.IntN(initialTable.Pots.Stamping),
			Finishing:  rand.IntN(initialTable.Pots.Finishing),
			Assembling: rand.IntN(initialTable.Pots.Assembling),
		},
		Coffees: Product{
			Stamping:   rand.IntN(initialTable.Coffees.Stamping),
			Finishing:  rand.IntN(initialTable.Coffees.Finishing),
			Assembling: rand.IntN(initialTable.Coffees.Assembling),
		},
		Samovar: Product{
			Stamping:   rand.IntN(initialTable.Samovar.Stamping),
			Finishing:  rand.IntN(initialTable.Samovar.Finishing),
			Assembling: rand.IntN(initialTable.Samovar.Assembling),
		},
	}
}
func (D Solution) Swap(times int) Solution {
	input := D
	addition := 0
	for i := 0; i < times; i++ {
		if rand.IntN(2) == 0 {
			addition = 1
		} else {
			addition = -1
		}

		ran := rand.IntN(3)
		switch ran {
		case 0:
			switch rand.IntN(3) {
			case 0:
				input.Pots.Assembling += addition
			case 1:
				input.Coffees.Assembling += addition
			case 2:
				input.Samovar.Assembling += addition
			}
		case 1:
			switch rand.IntN(3) {
			case 0:
				input.Pots.Finishing += addition
			case 1:
				input.Coffees.Finishing += addition
			case 2:
				input.Samovar.Finishing += addition
			}
		case 2:
			switch rand.IntN(3) {
			case 0:
				input.Pots.Stamping += addition
			case 1:
				input.Coffees.Stamping += addition
			case 2:
				input.Samovar.Stamping += addition
			}
		}
		input.Pots.Stamping = int(math.Abs(float64(input.Pots.Stamping)))
		input.Pots.Finishing = int(math.Abs(float64(input.Pots.Finishing)))
		input.Pots.Assembling = int(math.Abs(float64(input.Pots.Assembling)))
		input.Coffees.Stamping = int(math.Abs(float64(input.Coffees.Stamping)))
		input.Coffees.Finishing = int(math.Abs(float64(input.Coffees.Finishing)))
		input.Coffees.Assembling = int(math.Abs(float64(input.Coffees.Assembling)))
		input.Samovar.Stamping = int(math.Abs(float64(input.Samovar.Stamping)))
		input.Samovar.Finishing = int(math.Abs(float64(input.Samovar.Finishing)))
		input.Samovar.Assembling = int(math.Abs(float64(input.Samovar.Assembling)))
	}
	return input
}
func (D Solution) Valid() bool {
	if D.Pots.Assembling < 0 || D.Pots.Finishing < 0 || D.Pots.Stamping < 0 ||
		D.Coffees.Assembling < 0 || D.Coffees.Finishing < 0 || D.Coffees.Stamping < 0 ||
		D.Samovar.Assembling < 0 || D.Samovar.Finishing < 0 || D.Samovar.Stamping < 0 {
		return false
	}
	if D.Pots.Assembling > D.Pots.Finishing {
		return false
	}
	if D.Pots.Finishing > D.Pots.Stamping {
		return false
	}
	if D.Coffees.Assembling > D.Coffees.Finishing {
		return false
	}
	if D.Coffees.Finishing > D.Coffees.Stamping {
		return false
	}
	if D.Samovar.Assembling > D.Samovar.Finishing {
		return false
	}
	if D.Samovar.Finishing > D.Samovar.Stamping {
		return false
	}
	return true
}
func (initial Solution) Simulate(minTemp, startingTemp float64, coolingRate float64) Solution {
	current := initial
	best := current
	temp := startingTemp
	var P float64

	fmt.Printf("temp\tprob\t\tbest\t\tcurrent\t\tnew\n")
	for temp > minTemp {
		next := current.Swap(swapTimes)

		nextEnergy := next.CalculateEnergy()
		currentEnergy := current.CalculateEnergy()
		bestEnergy := best.CalculateEnergy()

		if nextEnergy > currentEnergy {
			current = next
		} else {
			P = math.Exp(-math.Abs(float64(nextEnergy-currentEnergy)) / temp)
			if P > rand.Float64() {
				current = next
			}
		}

		if currentEnergy > bestEnergy {
			best = current
		}

		temp *= coolingRate
		fmt.Printf("%.2f\t%.2f\t\t%.1f\t\t%.1f\t\t%.1f\n", temp, P, bestEnergy, currentEnergy, nextEnergy)
	}
	fmt.Printf("temp\tprob\t\tbest\t\tcurrent\t\tnew\n")
	return best
}
func main() {
	best := RandomSolution().Simulate(minTemp, startTemp, coolingRate)

	fmt.Printf("\nОптимальный план производства:\n")
	fmt.Printf("\t\tШтамповка\tОтделка\t\tСборка\n")
	fmt.Printf("Кастрюли\t%d\t\t%d\t\t%d\n", best.Pots.Stamping, best.Pots.Finishing, best.Pots.Assembling)
	fmt.Printf("Кофеварки\t%d\t\t%d\t\t%d\n", best.Coffees.Stamping, best.Coffees.Finishing, best.Coffees.Assembling)
	fmt.Printf("Самовары\t%d\t\t%d\t\t%d\n", best.Samovar.Stamping, best.Samovar.Finishing, best.Samovar.Assembling)
	fmt.Printf("\nРасчетная прибыль: %.2f\n", best.CalculateEnergy())
}
