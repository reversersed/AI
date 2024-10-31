package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	MatrixSize = 9    // размерность матрицы 3*3
	MaxAnts    = 50   // число муравьев
	Alfa       = 1.0  // вес фермента
	Beta       = 5.0  // эвристика
	Rho        = 0.2  // испарение
	Q          = 1000 // константа
	InitOdor   = 1.0 / MatrixSize
	MaxTime    = MaxAnts * MatrixSize
)

type Ant struct {
	TekProduct int
	Tabu       [MatrixSize]int
	Path       [MatrixSize]int
	NumProduct int
	Price      float64
}

var (
	Products      [MatrixSize]int
	Ants          [MaxAnts]Ant
	OdorMap       [MatrixSize]float64
	Best          Ant
	ProductPrice  = []float64{5.0, 2.5, 7.0}
	MaxProduction = []int{
		100, 80, 60,
		200, 150, 160,
		30, 35, 20,
	}
)

// Генерация случайного числа от min до max
func Random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func MakeProducts() {
	for i := 0; i < MatrixSize; i++ {
		Products[i] = int(math.Floor(Random(0, float64(MaxProduction[i])/3)))

		OdorMap[i] = InitOdor
	}
}

func MakeAnts(recreate bool) {
	k := 0
	for i := 0; i < MaxAnts; i++ {
		if recreate && Ants[i].Price > Best.Price {
			Best = Ants[i]
		}
		Ants[i] = Ant{
			TekProduct: k,
			Tabu:       [MatrixSize]int{},
			Path:       [MatrixSize]int{k},
			NumProduct: 1,
			Price:      0,
		}
		Ants[i].Tabu[k] = 1
		k++
		if k >= MatrixSize {
			k = 0
		}
	}
}

func Chance(i int) float64 {
	return math.Pow(OdorMap[i], Alfa) * math.Pow(1.0/float64(Products[i]+1), Beta)
}

func NextProduct(k int) int {
	i := Ants[k].TekProduct
	d := 0.0
	for j := 0; j < MatrixSize; j++ {
		if Ants[k].Tabu[j] == 0 {
			d += Chance(j)
		}
	}
	if d > 0 {
		for {
			j := rand.Intn(MatrixSize)
			if i != j && Ants[k].Tabu[j] == 0 && Random(0, 1) <= Chance(j)/d {
				return j
			}
		}
	}
	return MatrixSize - 1
}

func AntsMoving() bool {
	moving := false
	for k := 0; k < MaxAnts; k++ {
		if Ants[k].NumProduct < MatrixSize {
			nextProduct := NextProduct(k)
			if nextProduct%3 != 2 {
				if Products[nextProduct] > Products[nextProduct+1] {
					if Products[nextProduct] -= 1; Products[nextProduct] < 0 {
						Products[nextProduct] = 0
					}
				} else {
					if Products[nextProduct] += 1; Products[nextProduct] > MaxProduction[nextProduct] {
						Products[nextProduct] = MaxProduction[nextProduct]
					}
				}
			} else {
				if Random(0, 1) > 0.5 {
					if Products[nextProduct] += 1; Products[nextProduct] > MaxProduction[nextProduct] {
						Products[nextProduct] = MaxProduction[nextProduct]
					} else {
						if Products[nextProduct] -= 1; Products[nextProduct] < 0 {
							Products[nextProduct] = 0
						}
					}
				}
			}
			Ants[k].Path[Ants[k].NumProduct] = nextProduct
			Ants[k].NumProduct++
			Ants[k].Tabu[nextProduct] = 1
			Ants[k].Price += float64(Products[nextProduct/3] * int(ProductPrice[nextProduct/3]))
			for t := 0; t < MatrixSize; t++ {
				current := (t / 3) * 3
				if Products[current] < Products[current+1] {
					Ants[k].Price -= float64(Products[current+1]-Products[current]) * ProductPrice[current/3]
				}
				if Products[current+1] < Products[current+2] {
					Ants[k].Price -= float64(Products[current+2]-Products[current+1]) * ProductPrice[current/3]
				}
			}
			if Ants[k].Price < 0 {
				Ants[k].Price = 0
			}
			Ants[k].TekProduct = nextProduct
			moving = true
		}
	}
	return moving
}

func UpdateOdors() {
	for i := 0; i < MatrixSize; i++ {

		OdorMap[i] *= (1 - Rho)
		if OdorMap[i] < InitOdor {
			OdorMap[i] = InitOdor
		}

	}

	for ant := 0; ant < MaxAnts; ant++ {
		for k := 0; k < MatrixSize; k++ {
			OdorMap[Ants[ant].Path[k]] += (Ants[ant].Price / Q)
		}
	}
}

func main() {
	Best.Price = math.SmallestNonzeroFloat64
	MakeProducts()
	MakeAnts(false)

	for curTime := 0; curTime < MaxTime; curTime++ {
		if !AntsMoving() {
			UpdateOdors()
			MakeAnts(true)
			fmt.Printf("Время = %d  Путь = %.2f\n", curTime, Best.Price)
		}
	}

	fmt.Printf("Оптимальный путь = %.2f\n", Best.Price)
	for i := 0; i < MatrixSize; i++ {
		fmt.Printf("%d\t", Products[i])
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}
}
