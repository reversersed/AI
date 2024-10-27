package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	MaxTowns = 30  // число городов
	MaxAnts  = 30  // число муравьев
	Alfa     = 1.0 // вес фермента
	Beta     = 5.0 // эвристика
	Rho      = 0.5 // испарение
	Q        = 100 // константа
	InitOdor = 1.0 / MaxTowns
	MaxWay   = 100
	MaxTour  = MaxTowns * MaxWay
	MaxTime  = 20 * MaxTowns
)

type Town struct {
	x, y float64
}

type Ant struct {
	TekTown int
	Tabu    [MaxTowns]int
	Path    [MaxTowns]int
	NumTown int
	Len     float64
}

var (
	Towns   [MaxTowns]Town
	Ants    [MaxAnts]Ant
	DistMap [MaxTowns][MaxTowns]float64
	OdorMap [MaxTowns][MaxTowns]float64
	Best    Ant
)

// Генерация случайного числа от min до max
func Random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func MakeTowns() {
	for i := 0; i < MaxTowns; i++ {
		Towns[i] = Town{Random(0, MaxWay-1), Random(0, MaxWay-1)}
		for j := 0; j < MaxTowns; j++ {
			OdorMap[i][j] = InitOdor
		}
		DistMap[i][i] = 0.0
	}

	for i := 0; i < MaxTowns-1; i++ {
		for j := i + 1; j < MaxTowns; j++ {
			xd := Towns[i].x - Towns[j].x
			yd := Towns[i].y - Towns[j].y
			DistMap[i][j] = math.Sqrt(xd*xd + yd*yd)
			DistMap[j][i] = DistMap[i][j]
		}
	}
}

func MakeAnts(recreate bool) {
	k := 0
	for i := 0; i < MaxAnts; i++ {
		if recreate && Ants[i].Len < Best.Len {
			Best = Ants[i]
		}
		Ants[i] = Ant{
			TekTown: k,
			Tabu:    [MaxTowns]int{},
			Path:    [MaxTowns]int{k},
			NumTown: 1,
			Len:     0,
		}
		Ants[i].Tabu[k] = 1
		k++
		if k >= MaxTowns {
			k = 0
		}
	}
}

func Chance(i, j int) float64 {
	return math.Pow(OdorMap[i][j], Alfa) * math.Pow(1.0/DistMap[i][j], Beta)
}

func NextTown(k int) int {
	i := Ants[k].TekTown
	d := 0.0
	for j := 0; j < MaxTowns; j++ {
		if Ants[k].Tabu[j] == 0 {
			d += Chance(i, j)
		}
	}

	if d > 0 {
		for {
			j := rand.Intn(MaxTowns)
			if i != j && Ants[k].Tabu[j] == 0 && Random(0, 1) <= Chance(i, j)/d {
				return j
			}
		}
	}
	return MaxTowns - 1
}

func AntsMoving() bool {
	moving := false
	for k := 0; k < MaxAnts; k++ {
		if Ants[k].NumTown < MaxTowns {
			nextTown := NextTown(k)
			Ants[k].Path[Ants[k].NumTown] = nextTown
			Ants[k].NumTown++
			Ants[k].Tabu[nextTown] = 1
			Ants[k].Len += DistMap[Ants[k].TekTown][nextTown]
			if Ants[k].NumTown == MaxTowns {
				Ants[k].Len += DistMap[Ants[k].Path[MaxTowns-1]][Ants[k].Path[0]]
			}
			Ants[k].TekTown = nextTown
			moving = true
		}
	}
	return moving
}

func UpdateOdors() {
	for i := 0; i < MaxTowns; i++ {
		for j := 0; j < MaxTowns; j++ {
			if i != j {
				OdorMap[i][j] *= (1 - Rho)
				if OdorMap[i][j] < InitOdor {
					OdorMap[i][j] = InitOdor
				}
			}
		}
	}

	for ant := 0; ant < MaxAnts; ant++ {
		for k := 0; k < MaxTowns; k++ {
			i := Ants[ant].Path[k]
			j := Ants[ant].Path[(k+1)%MaxTowns]
			OdorMap[i][j] += Q / Ants[ant].Len
			OdorMap[j][i] = OdorMap[i][j]
		}
	}
}

func main() {
	Best.Len = MaxTour
	MakeTowns()
	MakeAnts(false)

	for curTime := 0; curTime < MaxTime; curTime++ {
		if !AntsMoving() {
			UpdateOdors()
			MakeAnts(true)
			fmt.Printf("Время = %d  Путь = %.2f\n", curTime, Best.Len)
		}
	}

	fmt.Printf("Оптимальный путь = %.2f\n", Best.Len)
}
