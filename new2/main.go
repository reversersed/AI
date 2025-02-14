package main

import (
	"fmt"
	"math/rand"
)

const (
	N      = 20   // размерность клеток
	Teta   = 3    // порог
	DeltaW = 0.05 // вес обучения
)

type LC struct { // координаты
	L int // строка
	C int // столбец
}

var (
	R     [N][N]int        // слой рецепторов
	A     [N][N]int        // ассоциативный слой
	S     [N][N][N]LC      // связи
	W     [N][N][2]float64 // веса связей
	fTest = false          // флаг проверки работоспособности
)

func InitLayers() {
	// генерация связей
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			for K := 0; K < N; K++ {
				S[L][C][K].L = rand.Intn(N)
				S[L][C][K].C = rand.Intn(N)
			}
		}
	}
	// очистка весов
	for K := 0; K < 2; K++ {
		for L := 0; L < N; L++ {
			for C := 0; C < N; C++ {
				W[L][C][K] = 0
			}
		}
	}
}

func NameDic(dic string) string {
	switch dic {
	case "00":
		return "И"
	case "01":
		return "К"
	case "10":
		return "Л"
	case "11":
		return "Н"
	default:
		return "не знаю"
	}
}

func Obraz(H int) string {
	// чистка рецепторов
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			R[L][C] = 0
		}
	}

	sd := H % 4            // выбор образа: 0-И,1-К,2-Л,3-Н
	rr := rand.Intn(N / 3) // выбор радиуса образа
	if rr < 5 {
		rr = 5
	}
	res := "00" // код образа

	// выбор места расположения образа
	Lc := rr + rand.Intn(N-2*rr) // строка центра
	Cc := rr + rand.Intn(N-2*rr) // столбец центра

	// рисование выбранного образа
	switch sd {
	case 0: // И
		res = "00"
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
		res = "01"
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
		res = "10"
		for i := Lc - rr; i < Lc+rr; i++ {
			R[i][Cc+rr] = 1
		}
		L := Lc + rr - 1
		for i := Cc - rr; i < Cc+rr; i++ {
			R[L][i] = 1
			L--
		}
	case 3: // Н
		res = "11"
		for i := Lc - rr; i < Lc+rr; i++ {
			R[i][Cc-rr] = 1
			R[i][Cc+rr] = 1
		}
		for i := Cc - rr; i < Cc+rr; i++ {
			R[Lc][i] = 1
		}
	}

	// вывод образа на экран во время проверки работоспособности
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
	// очистка ассоциативного слоя A
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			A[L][C] = 0
		}
	}
	// отображение в слое A (входы)
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			if R[L][C] == 1 {
				for K := 0; K < N; K++ {
					Lv := S[L][C][K].L
					Cv := S[L][C][K].C
					if Lv >= 0 && Lv < N && Cv >= 0 && Cv < N {
						A[Lv][Cv]++
					}
				}
			}
		}
	}
	// отображение в слое A (выходы)
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
	var E [2]float64 // эффекторный слой
	res := ""        // код распознавания
	for K := 0; K < 2; K++ {
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
	// вывод результата при проверке работоспособности
	if fTest {
		fmt.Printf("\nЯ думаю что это %s\n", NameDic(res))
	}
	return res
}

func Teach(sd, sd1 string) {
	for K := 0; K < 2; K++ {
		if sd[K] != sd1[K] {
			for L := 0; L < N; L++ {
				for C := 0; C < N; C++ {
					if A[L][C] == 1 { // обучение виноватых
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
	fmt.Println("Обучение перцептрона распознаванию четырех образов: крестика, нолика, плюса, ромба")

	limit := 97.68 // предел обучения
	nOk := 0       // число удачных ответов
	percent := 0.0
	Step := 1

	InitLayers() // инициализация слоёв
	for percent < limit {
		dic := Obraz(Step) // новый образ
		Otobr()            // отображение
		dic1 := Reak()     // опознавание
		if dic == dic1 {
			nOk++
		} else {
			Teach(dic, dic1)
		}
		// вывод текущей информации на экран
		percent = float64(nOk) / float64(Step) * 100
		fmt.Printf("Шаг %d : Доля удачных ответов %.2f %%\n", Step, percent)
		Step++
	}
	fTest = true

	for i := 0; i < 20; i++ {
		dic := Obraz(Step) // новый образ
		Otobr()            // отображение
		dic1 := Reak()     // опознавание
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
