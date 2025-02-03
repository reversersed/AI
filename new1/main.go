package main

import (
	"fmt"
	"math/rand"
)

const (
	N      = 20   // размерность клеток
	Teta   = 3    // порог
	DeltaW = 0.15 // вес обучения
)

type LC struct { // координаты
	L int // строка
	C int // столбец
}

var (
	R     [N][N]int     // слой рецепторов
	A     [N][N]int     // ассоциативный слой
	S     [N][N][N]LC   // связи
	W     [N][N]float64 // веса связей
	fTest = false       // флаг проверки работоспособности
)

func InitLayers() {
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			W[L][C] = 0              // очистка весов
			for K := 0; K < N; K++ { // генерация связей
				S[L][C][K].L = rand.Intn(N)
				S[L][C][K].C = rand.Intn(N)
			}
		}
	}
}

func Obraz(H int) int {
	// чистка рецепторов
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			R[L][C] = 0
		}
	}
	sd := H % 2            // выбор образа: 0-#, 1-=
	rr := rand.Intn(N / 2) // выбор радиуса образа
	if rr < 3 {
		rr = 3
	}

	// выбор места расположения образа
	Lc := rr + rand.Intn(N-2*rr) // строка центра
	Cc := rr + rand.Intn(N-2*rr) // столбец центра

	// рисование выбранного образа
	switch sd {
	case 0: // рисование #
		for i := -rr; i <= rr; i++ {
			L := Lc - rr/2
			C := Cc + i
			R[L][C] = 1
			L = Lc + rr/2
			R[L][C] = 1

			L = Lc + i
			C = Cc - rr/2 - i/3
			R[L][C] = 1
			C = Cc + rr/2 - i/3
			R[L][C] = 1
		}
	case 1: // рисование =
		for i := -rr; i <= rr; i++ {
			L := Lc - rr/3
			C := Cc + i
			R[L][C] = 1
			L = Lc + rr/3
			R[L][C] = 1
		}
	}

	// вывод образа на экран во время проверки работоспособности
	if fTest {
		if sd == 0 {
			fmt.Println("Рисую #")
		} else {
			fmt.Println("Рисую =")
		}

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

	return sd
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
					A[Lv][Cv]++
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

func Reak() int {
	E := 0.0 // эффекторный слой
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			E += float64(A[L][C]) * W[L][C]
		}
	}
	// вывод результата при проверке работоспособности
	if fTest {
		fmt.Println()
		if E > 0 {
			fmt.Println("Я думаю что это =")
		} else {
			fmt.Println("Я думаю что это #")
		}
	}
	if E > 0 {
		return 1
	}
	return 0
}

func Teach(sd int) {
	for L := 0; L < N; L++ {
		for C := 0; C < N; C++ {
			if A[L][C] == 1 { // обучение виноватых
				if sd == 0 {
					W[L][C] -= DeltaW
				} else {
					W[L][C] += DeltaW
				}
			}
		}
	}
}

func main() {
	fmt.Println("Обучение перцептрона распознаванию двух образов: # и =")

	lmax := 80000 // число шагов обучения
	nOk := 0      // число удачных ответов

	InitLayers() // инициализация слоёв
	for Step := 1; Step <= lmax; Step++ {
		dic := Obraz(Step) // новый образ: 0-#, 1-=
		Otobr()            // отображение
		dic1 := Reak()     // опознавание
		if dic == dic1 {
			nOk++
		} else {
			Teach(dic)
		}
		// вывод текущей информации на экран
		fmt.Printf("Шаг %d : Доля удачных ответов %.2f %%\n", Step, float64(nOk)/float64(Step)*100)
		// тестирование за 20 шагов до конца обучения
		if Step == (lmax - 20) {
			fTest = true
		}
	}
	fmt.Println("Конец")
}
