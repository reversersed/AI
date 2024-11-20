package main

import (
	"fmt"
)

const (
	MaxItems    = 11
	MaxClients  = 10
	MaxClusters = 5
)

var (
	Beta     = 1.0
	Rho      = 0.9
	ItemName = [MaxItems]string{
		"Молоток", "Бумага", "Шоколадка",
		"Отвертка", "Ручка", "Кофе", "Гвоздодер", "Карандаш",
		"Конфеты", "Дрель", "Дырокол",
	}
	Data = [MaxClients][MaxItems]int{
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0},
		{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	}
	Clusters [MaxClusters][MaxItems]int
	Sum      [MaxClusters][MaxItems]int
	Members  [MaxClusters]int
	Group    [MaxClients]int
	N        int
)

func Initialize() {
	for i := 0; i < MaxClusters; i++ {
		for j := 0; j < MaxItems; j++ {
			Clusters[i][j] = 0
			Sum[i][j] = 0
		}
		Members[i] = 0
	}
	for j := 0; j < MaxClients; j++ {
		Group[j] = -1
	}
}

func AndVectors(R, V, W *[MaxItems]int) {
	for i := 0; i < MaxItems; i++ {
		R[i] = V[i] & W[i]
	}
}

func UpdateVectors(K int) {
	if K < 0 || K >= MaxClusters {
		return
	}
	f := true
	for i := 0; i < MaxItems; i++ {
		Clusters[K][i] = 0
		Sum[K][i] = 0
	}
	for j := 0; j < MaxClients; j++ {
		if Group[j] == K {
			if f {
				copy(Clusters[K][:], Data[j][:])
				copy(Sum[K][:], Data[j][:])
				f = false
			} else {
				AndVectors(&Clusters[K], &Clusters[K], &Data[j])
				for i := 0; i < MaxItems; i++ {
					Sum[K][i] += Data[j][i]
				}
			}
		}
	}
}

func CreateVector(V *[MaxItems]int) int {
	i := -1
	for {
		i++
		if i >= MaxClusters {
			return -1
		}
		if Members[i] == 0 {
			break
		}
	}
	N++
	Members[i] = 1
	copy(Clusters[i][:], V[:])
	return i
}

func OnesVector(V *[MaxItems]int) int {
	k := 0
	for j := 0; j < MaxItems; j++ {
		if V[j] == 1 {
			k++
		}
	}
	return k
}

func ExecuteART1() {
	var R [MaxItems]int
	var PE, P, E int
	var Test bool
	var i, j int
	count := 50
	exit := false
	var s int

	for {
		exit = true
		for i = 0; i < MaxClients; i++ {
			for j = 0; j < MaxClusters; j++ {
				if Members[j] > 0 {
					AndVectors(&R, &Data[i], &Clusters[j])
					PE = OnesVector(&R)
					P = OnesVector(&Clusters[j])
					E = OnesVector(&Data[i])
					Test = (float64(PE)/(Beta+float64(P)) > float64(E)/(Beta+float64(MaxItems)))
					if Test {
						Test = (float64(PE)/float64(E) < Rho)
					}
					if Test {
						Test = (Group[i] != j)
					}
					if Test {
						s = Group[i]
						Group[i] = j
						if s >= 0 {
							Members[s]--
							if Members[s] == 0 {
								N--
							}
						}
						Members[j]++
						UpdateVectors(s)
						UpdateVectors(j)
						exit = false
						break
					}
				}
			}
			if Group[i] == -1 {
				Group[i] = CreateVector(&Data[i])
				exit = false
			}
		}
		count--
		if exit || count <= 0 {
			break
		}
	}
}

func ShowClusters() {
	for i := 0; i < N; i++ {
		fmt.Printf("Вектор-прототип %2d : ", i+1)
		for j := 0; j < MaxItems; j++ {
			fmt.Printf("%d ", Clusters[i][j])
		}
		fmt.Println()
		for k := 0; k < MaxClients; k++ {
			if Group[k] == i {
				fmt.Printf("Покупатель  %2d : ", k+1)
				for j := 0; j < MaxItems; j++ {
					fmt.Printf("%d ", Data[k][j])
				}
				fmt.Println()
			}
		}
	}
	fmt.Println()
}

func MakeAdvise(p int) {
	best := -1
	max := 0
	for i := 0; i < MaxItems; i++ {
		if Data[p][i] == 0 && Sum[Group[p]][i] > max {
			best = i
			max = Sum[Group[p]][i]
		}
	}
	fmt.Printf("Для покупателя %d ", p+1)
	if best >= 0 {
		fmt.Printf("рекомендация - %s\n", ItemName[best])
	} else {
		fmt.Println("нет рекомендаций")
	}
	fmt.Print("Уже выбраны: ")
	for i := 0; i < MaxItems; i++ {
		if Data[p][i] != 0 {
			fmt.Printf("%s ", ItemName[i])
		}
	}
	fmt.Println()
}

func main() {
	Initialize()
	ExecuteART1()
	ShowClusters()
	for p := 0; p < MaxClients; p++ {
		MakeAdvise(p)
	}
}
