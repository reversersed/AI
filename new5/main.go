package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	N                       = 27
	border                  = 10
	lmax                    = 10000000
	nOkInProgressiterations = 100000
	INPUT_NEURONS           = N
	HIDDEN_NEURONS_AMMOUNT  = 4
	OUTPUT_NEURONS          = 2
	FIGURE_AMMOUNT          = 4
	LEARN_RATE              = 0.2
)

var (
	inputs [INPUT_NEURONS][INPUT_NEURONS]float64
	hidden [HIDDEN_NEURONS_AMMOUNT]float64
	actual [OUTPUT_NEURONS]float64
	target [OUTPUT_NEURONS]float64

	wih [INPUT_NEURONS][INPUT_NEURONS][HIDDEN_NEURONS_AMMOUNT]float64
	who [HIDDEN_NEURONS_AMMOUNT][OUTPUT_NEURONS]float64

	moveh, moveo float64
	erro         [OUTPUT_NEURONS]float64
	errh         [HIDDEN_NEURONS_AMMOUNT]float64

	fTest = false
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func sigmoidFunction(val float64) float64 {
	return 1.0 / (1.0 + math.Exp(-val))
}

func reversedSigmoidFunction(val float64) float64 {
	return val * (1.0 - val)
}

func assignRandomWeights() {
	for inp1 := range INPUT_NEURONS {
		for inp2 := range INPUT_NEURONS {
			for hid1 := range HIDDEN_NEURONS_AMMOUNT {
				wih[inp1][inp2][hid1] = rng.Float64() - 0.5
			}
		}
	}
	moveh = rng.Float64() - 0.5

	for hid1 := range HIDDEN_NEURONS_AMMOUNT {
		for out := range OUTPUT_NEURONS {
			who[hid1][out] = rng.Float64() - 0.5
		}
	}
	moveo = rng.Float64() - 0.5
}

func feedForward() {
	for hid1 := range HIDDEN_NEURONS_AMMOUNT {
		sum := 0.0
		for inp1 := range INPUT_NEURONS {
			for inp2 := range INPUT_NEURONS {
				sum += inputs[inp1][inp2] * wih[inp1][inp2][hid1]
			}
		}
		sum += moveh
		hidden[hid1] = sigmoidFunction(sum)
	}

	for out := range OUTPUT_NEURONS {
		sum := 0.0
		for hid1 := range HIDDEN_NEURONS_AMMOUNT {
			sum += hidden[hid1] * who[hid1][out]
		}
		sum += moveo
		actual[out] = sigmoidFunction(sum)
	}
}

func backPropagate() {
	for out := range OUTPUT_NEURONS {
		erro[out] = (target[out] - actual[out]) * reversedSigmoidFunction(actual[out])
	}

	for hid1 := range HIDDEN_NEURONS_AMMOUNT {
		errh[hid1] = 0.0
		for out := range OUTPUT_NEURONS {
			errh[hid1] += erro[out] * who[hid1][out]
		}
		errh[hid1] *= reversedSigmoidFunction(hidden[hid1])
	}

	for out := range OUTPUT_NEURONS {
		for hid1 := range HIDDEN_NEURONS_AMMOUNT {
			who[hid1][out] += LEARN_RATE * erro[out] * hidden[hid1]
		}
		moveo += LEARN_RATE * erro[out]
	}

	for hid1 := range HIDDEN_NEURONS_AMMOUNT {
		for inp1 := range INPUT_NEURONS {
			for inp2 := range INPUT_NEURONS {
				wih[inp1][inp2][hid1] += LEARN_RATE * errh[hid1] * inputs[inp1][inp2]
			}
		}
		moveh += LEARN_RATE * errh[hid1]
	}
}

func nameDic(dic string) string {
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

func chooseOption() string {
	answer := ""
	for i := range OUTPUT_NEURONS {
		if actual[i] >= 0.5 {
			answer += "1"
		} else {
			answer += "0"
		}
	}
	if fTest {
		fmt.Printf("\nЯ думаю что это %s\n", nameDic(answer))
	}
	return answer
}

func obraz(H int) string {
	for L := range N {
		for C := range N {
			inputs[L][C] = 0
		}
	}
	sd := H % 4
	res := "00"
	LeftFigurePart := rng.Intn(N-2) + 1
	if LeftFigurePart < border {
		LeftFigurePart = border
	}
	TopFigurePart := LeftFigurePart
	X := rng.Intn(N-LeftFigurePart) + 1
	Y := rng.Intn(N-TopFigurePart) + 1

	switch sd {
	case 0:
		res = "00"
		for i := X; i < X+LeftFigurePart; i++ {
			inputs[i][Y] = 1
			inputs[i][Y+TopFigurePart-1] = 1
			inputs[i][Y+TopFigurePart-(i-X)-1] = 1
		}
	case 1:
		res = "01"
		count := -1
		for i := X; i < X+LeftFigurePart; i++ {
			inputs[i][Y] = 1
		}
		for i := 0; i <= LeftFigurePart/2; i++ {
			inputs[X+i][Y+TopFigurePart/2-count] = 1
			inputs[X+LeftFigurePart-i-1][Y+TopFigurePart/2-count] = 1
			count++
		}
	case 2:
		res = "10"
		count := 0
		for i := X; i < X+LeftFigurePart; i++ {
			inputs[i][Y+TopFigurePart-1] = 1
		}
		for i := X + (LeftFigurePart)/2; i < X+LeftFigurePart; i++ {
			inputs[i][Y] = 1
		}
		count = 0
		for i := 0; i < LeftFigurePart/2; i++ {
			inputs[X+i][Y+TopFigurePart/2-count] = 1
			count++
		}
		for j := 0; j < TopFigurePart/2; j++ {
			inputs[X][Y+TopFigurePart/2+j] = 1
		}
	case 3:
		res = "11"
		for i := Y; i < Y+TopFigurePart; i++ {
			inputs[i][X] = 1
			inputs[i][X+LeftFigurePart-1] = 1
		}
		for i := X; i < X+LeftFigurePart; i++ {
			inputs[Y+TopFigurePart/2-1][i] = 1
		}
	}

	if fTest {
		fmt.Printf("Рисую %s\n", nameDic(res))
		for L := range N {
			for C := range N {
				if inputs[L][C] == 0 {
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

func setTarget(drawnFigure int) {
	switch drawnFigure {
	case 0:
		target[0], target[1] = 0, 0
	case 1:
		target[0], target[1] = 0, 1
	case 2:
		target[0], target[1] = 1, 0
	case 3:
		target[0], target[1] = 1, 1
	}
}

func main() {
	nOk := 0
	nOkInProgress := 0

	assignRandomWeights()
	for Step := 1; Step <= lmax; Step++ {
		dic := obraz(Step)
		setTarget(Step % FIGURE_AMMOUNT)
		feedForward()
		dic1 := chooseOption()
		if dic == dic1 {
			nOk++
			nOkInProgress++
		} else {
			backPropagate()
		}
		if Step%nOkInProgressiterations == 0 {
			fmt.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\nШаг %d:\n\tДоля удачных ответов\t\t%.2f%%\n\tТочность за %d итераций\t%.2f%%\n\n",
				Step, float64(nOk)/float64(Step)*100, nOkInProgressiterations, float64(nOkInProgress)/float64(nOkInProgressiterations)*100)
			nOkInProgress = 0
		}
		if Step == (lmax - 20) {
			fTest = true
		}
	}
}
