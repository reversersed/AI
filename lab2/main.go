package main

import (
	"fmt"
	"math/rand/v2"
)

type Solution struct {
	Pots    int
	Coffees int
	Samovar int
	Energy  float64
}

var (
	products = []struct {
		Name            string
		Profit          float64
		ProductionLimit int
	}{
		{"Кастрюли", 3.8, 120},
		{"Кофеварки", 5.0, 70},
		{"Самовары", 4.5, 50},
	}
	generalCapacity  = 200    // лимит продукции в день на штамповке+отделке
	assemblyCapacity = 180    // лимит продукции в день на сборочном оборудовании
	ticketRate       = 20     // Размер штрафа энергии
	minTemp          = 0.05   // нижний предел температуры
	startTemp        = 100.0  // стартовая температура
	coolingRate      = 0.9999 // степень остывания
	outputAll        = true   // делать ли полный вывод. true = вывод всех current, false = только новые значения
)

func (solution *Solution) CalculateEnergy() {
	solution.Energy = float64(solution.Pots)*products[0].Profit +
		float64(solution.Coffees)*products[1].Profit +
		float64(solution.Samovar)*products[2].Profit

	if products[0].ProductionLimit-solution.Pots < 0 {
		solution.Energy += float64(products[0].ProductionLimit-solution.Pots) * float64(ticketRate)
	}
	if products[1].ProductionLimit-solution.Coffees < 0 {
		solution.Energy += float64(products[1].ProductionLimit-solution.Coffees) * float64(ticketRate)
	}
	if products[2].ProductionLimit-solution.Samovar < 0 {
		solution.Energy += float64(products[2].ProductionLimit-solution.Samovar) * float64(ticketRate)
	}

	total := solution.Pots + solution.Coffees + solution.Samovar
	if total > generalCapacity {
		solution.Energy += (float64(generalCapacity) - float64(total)) * float64(ticketRate)
	}
	if total > assemblyCapacity {
		solution.Energy += (float64(assemblyCapacity) - float64(total)) * float64(ticketRate)
	}
}
func RandomSolution() Solution {
	sol := Solution{
		Pots:    rand.IntN(products[0].ProductionLimit + 1),
		Coffees: rand.IntN(products[1].ProductionLimit + 1),
		Samovar: rand.IntN(products[2].ProductionLimit + 1),
	}
	sol.CalculateEnergy()
	return sol
}
func (D Solution) Swap(S Solution) Solution {
	input := D
	ran := rand.IntN(3)

	switch ran {
	case 0:
		input.Pots = S.Pots
	case 1:
		input.Coffees = S.Coffees
	case 2:
		input.Samovar = S.Samovar
	}
	input.CalculateEnergy()
	return input
}
func (initial Solution) Simulate(minTemp, startingTemp float64, coolingRate float64) Solution {
	current := initial
	best := current
	temp := startingTemp
	var P float64
	last := current.Energy

	fmt.Printf("temp\tprob\t\tbest\t\tcurrent\t\tnew\n")
	for temp > minTemp {
		next := RandomSolution().Swap(current)

		if next.Energy > current.Energy {
			current = next
		} else {
			P = rand.Float64() * startingTemp
			if P < temp {
				current = next
			}
		}

		if current.Energy > best.Energy {
			best = current
		}

		temp *= coolingRate
		if last != current.Energy || outputAll {
			last = current.Energy
			fmt.Printf("%.2f\t%.2f\t\t%.1f\t\t%.1f\t\t%.1f\n", temp, P, best.Energy, current.Energy, next.Energy)
		}
	}
	fmt.Printf("temp\tprob\t\tbest\t\tcurrent\t\tnew\n")
	return best
}
func main() {
	initial := RandomSolution()

	best := initial.Simulate(minTemp, startTemp, coolingRate)

	fmt.Printf("\nОптимальный план производства:\n")
	fmt.Printf("Кастрюлей: %d, кофеварок: %d, самоваров: %d\n", best.Pots, best.Coffees, best.Samovar)
	fmt.Printf("Расчетная прибыль: %.2f\n", best.Energy)
}
