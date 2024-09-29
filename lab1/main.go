package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

const (
	N    = 20
	Tn   = 30
	Tk   = 0.5
	Alfa = 0.98
	ST   = 100
)

type TMember struct {
	Plan   []int
	Energy int
}

var (
	Current  TMember
	Working  TMember
	Best     TMember
	T        float64
	Delta    float64
	P        float64
	fNew     bool
	fBest    bool
	Time     int64
	Step     int64
	Accepted int64
)

func Swap(M *TMember) {
	var x, y, v int
	x = (rand.IntN(N))
	for {
		y = (rand.IntN(N))
		if y != x {
			break
		}
	}
	v = M.Plan[x]
	M.Plan[x] = M.Plan[y]
	M.Plan[y] = v
}
func New(M *TMember) {
	M.Plan = make([]int, N)
	for i := 0; i < N; i++ {
		M.Plan[i] = i
	}
	for i := 0; i < N; i++ {
		Swap(M)
	}
}
func CalcEnergy(M *TMember) {
	dx := []int{-1, 1, -1, 1}
	dy := []int{-1, 1, 1, -1}
	error := 0
	var tx, ty int
	for x := 0; x < N; x++ {
		for j := 0; j < 4; j++ {
			tx = x + dx[j]
			ty = M.Plan[x] + dy[j]
			for tx > 0 && tx < N && ty > 0 && ty < N {
				if M.Plan[tx] == ty {
					error++
				}
				tx = tx + dx[j]
				ty = ty + dy[j]
			}
		}
	}
	M.Energy = error
}
func Copy(MD, MS *TMember) {
	for i := 0; i < N; i++ {
		MD.Plan[i] = MS.Plan[i]
	}
	MD.Energy = MS.Energy
}
func Show(M *TMember) {
	fmt.Println("Решение:")
	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			if M.Plan[x] == y {
				fmt.Print("Q")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
func main() {
	T = Tn
	fBest = false
	Time = 0
	Best.Energy = math.MaxInt
	Current.Plan = make([]int, N)
	Working.Plan = make([]int, N)
	Best.Plan = make([]int, N)
	New(&Current)
	CalcEnergy(&Current)
	Copy(&Working, &Current)
	for T > Tk {
		Accepted = 0
		for Step = 0; Step < ST; Step++ {
			fNew = false
			Swap(&Working)
			CalcEnergy(&Working)
			if Working.Energy <= Current.Energy {
				fNew = true
			} else {
				Delta = float64(Working.Energy - Current.Energy)
				P = math.Exp(-Delta / P)
				if P > rand.Float64() {
					Accepted++
					fNew = true
				}
			}
			if fNew {
				fNew = false
				Copy(&Current, &Working)
				if Current.Energy < Best.Energy {
					Copy(&Best, &Current)
					fBest = true
				} else {
					Copy(&Working, &Current)
				}
			}
		}
		fmt.Printf("Temp = %.1f P = %f Energy = %d\n", T, P, Best.Energy)
		T = T * Alfa
		Time++
	}
	if fBest {
		Show(&Best)
	}
}
