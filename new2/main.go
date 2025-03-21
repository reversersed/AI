package main

import (
	"fmt"
	"math/rand"
)

const (
	N      = 20
	Teta   = 3
	DeltaW = 0.05
)

type LC struct {
	L int
	C int
}

var (
	R     [N][N]int
	A     [N][N]int
	S     [N][N][N]LC
	W     [N][N][4]float64
	fTest = false
)

func InitLayers() {
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			for K := 0; K < N; K++ {
				S[L][C][K].L = rand.Intn(N)
				S[L][C][K].C = rand.Intn(N)
			}
		}
	}
	for K := 0; K < 4; K++ {
		for L := 0; L < N; L++ {
			for C := 0; C < N; C++ {
				W[L][C][K] = 0
			}
		}
	}
}

func NameDic(dic string) string {
	switch dic {
	case "1000":
		return "И"
	case "0100":
		return "К"
	case "0010":
		return "Л"
	case "0001":
		return "Н"
	default:
		return "не знаю"
	}
}

func Obraz(H int) string {
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			R[L][C] = 0
		}
	}
	sd := H % 4
	rr := rand.Intn(N / 3)
	if rr < 5 {
		rr = 5
	}
	res := "0000"
	Lc := rr + rand.Intn(N-2*rr)
	Cc := rr + rand.Intn(N-2*rr)

	switch sd {
	case 0: // И
		res = "1000"
		for i := Lc - rr; i < Lc+rr; i++ {
			R[i][Cc-rr] = 1
			R[i][Cc+rr] = 1
			R[i+1][Cc+rr] = 1
		}
		L := Lc + rr
		for i := Cc - rr; i < Cc+rr; i++ {
			R[L][i] = 1
			L--
		}
	case 1: // К
		res = "0100"
		for i := Lc - rr; i < Lc+rr; i++ {
			R[i][Cc-rr] = 1
		}
		for i := 0; i <= rr; i++ {
			R[Lc-i][Cc-rr+i] = 1
			if i != rr {
				R[Lc+i][Cc-rr+i] = 1
			}
		}
	case 2: // Л
		res = "0010"
		for i := Lc - rr; i < Lc+rr; i++ {
			R[i][Cc+rr] = 1
		}
		L := Lc + rr - 1
		for i := Cc - rr; i < Cc+rr; i++ {
			R[L][i] = 1
			L--
		}
	case 3: // Н
		res = "0001"
		for i := Lc - rr; i < Lc+rr; i++ {
			R[i][Cc-rr] = 1
			R[i][Cc+rr] = 1
		}
		for i := Cc - rr; i < Cc+rr; i++ {
			R[Lc][i] = 1
		}
	}

	if fTest {
		fmt.Printf("Рисую %s\n", NameDic(res))
		for L := 0; L < N; L++ {
			for C := 0; C < N; C++ {
				if R[L][C] == 0 {
					fmt.Print("_")
				} else {
					fmt.Print("*")
				}
			}
			fmt.Println()
		}
	}
	return res
}

func Otobr() {
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			A[L][C] = 0
		}
	}
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			if R[L][C] == 1 {
				for K := 0; K < N; K++ {
					Lv := S[L][C][K].L
					Cv := S[L][C][K].C
					A[Lv][Cv]++
				}
			}
		}
	}
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			if A[L][C] > Teta {
				A[L][C] = 1
			} else {
				A[L][C] = 0
			}
		}
	}
}

func Reak() string {
	E := make([]float64, 4)
	res := ""
	for K := 0; K < 4; K++ {
		E[K] = 0
		for L := 0; L < N; L++ {
			for C := 0; C < N; C++ {
				E[K] += float64(A[L][C]) * W[L][C][K]
			}
		}
		if E[K] > 0 {
			res += "1"
		} else {
			res += "0"
		}
	}
	if fTest {
		fmt.Printf("\nЯ думаю что это %s\n", NameDic(res))
	}
	return res
}

func Teach(sd, sd1 string) {
	for K := 0; K < 4; K++ {
		if sd[K] != sd1[K] {
			for L := 0; L < N; L++ {
				for C := 0; C < N; C++ {
					if A[L][C] == 1 {
						if sd[K] == '0' {
							W[L][C][K] -= DeltaW
						} else {
							W[L][C][K] += DeltaW
						}
					}
				}
			}
		}
	}
}

func main() {
	fmt.Println("Обучение перцептрона распознаванию четырех образов")

	limit := 96.68
	nOk := 0
	percent := 0.0
	Step := 1

	InitLayers()
	for percent < limit {
		dic := Obraz(Step)
		Otobr()
		dic1 := Reak()
		if dic == dic1 {
			nOk++
		} else {
			Teach(dic, dic1)
		}
		percent = float64(nOk) / float64(Step) * 100
		fmt.Printf("Шаг %d : Доля удачных ответов %.2f %%\n", Step, percent)
		Step++
	}
	fTest = true

	for i := 0; i < 20; i++ {
		dic := Obraz(Step)
		Otobr()
		dic1 := Reak()
		if dic == dic1 {
			nOk++
		} else {
			Teach(dic, dic1)
		}

		percent = float64(nOk) / float64(Step) * 100
		fmt.Printf("Шаг %d : Доля удачных ответов %.2f %%\n", Step, percent)
		Step++
	}
	fmt.Println("Конец")
}
