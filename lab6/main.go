package main

import (
	"fmt"
)

const (
	MaxSellers  = 8
	MaxShops    = 10
	MaxClusters = 4
)

var (
	Beta      = 1.0
	Rho       = 0.6
	ShopsName = [MaxShops]string{
		"«Магнит»", "«Лента»", "«Пятерочка»", "«Шоколад»",
		"«Окей»", "«Как раз»", "«Перекресток»", "«Дикси»",
		"«Яндекс.Маркет»", "«Сбермегамаркет»",
	}
	SellerName = [MaxSellers]string{
		"ООО «Ореховый Алтай»", "ООО «НУТРИЛЕНТ»", "ООО «Грин Поинт»",
		"ООО «7 ФРУКТОВ»", "ООО «ФудВэй»", "ООО «ВИШНЯ»", "«АГРОМИКС»", "ООО «ТРИАДА-ПОКОТОРГ»",
	}
	Data = [MaxShops][MaxSellers]int{
		{0, 1, 1, 0, 0, 1, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 1},
		{1, 0, 0, 1, 0, 0, 1, 0},
		{0, 1, 0, 1, 0, 0, 0, 1},
	}
	Clusters [MaxClusters][MaxSellers]int
	Sum      [MaxClusters][MaxSellers]int
	Members  [MaxClusters]int
	Group    [MaxShops]int
	N        int
)

func Initialize() {
	for i := 0; i < MaxClusters; i++ {
		for j := 0; j < MaxSellers; j++ {
			Clusters[i][j] = 0
			Sum[i][j] = 0
		}
		Members[i] = 0
	}
	for j := 0; j < MaxShops; j++ {
		Group[j] = -1
	}
}

func AndVectors(R, V, W *[MaxSellers]int) {
	for i := 0; i < MaxSellers; i++ {
		R[i] = V[i] & W[i]
	}
}

func UpdateVectors(K int) {
	if K < 0 || K >= MaxClusters {
		return
	}
	f := true
	for i := 0; i < MaxSellers; i++ {
		Clusters[K][i] = 0
		Sum[K][i] = 0
	}
	for j := 0; j < MaxShops; j++ {
		if Group[j] == K {
			if f {
				copy(Clusters[K][:], Data[j][:])
				copy(Sum[K][:], Data[j][:])
				f = false
			} else {
				AndVectors(&Clusters[K], &Clusters[K], &Data[j])
				for i := 0; i < MaxSellers; i++ {
					Sum[K][i] += Data[j][i]
				}
			}
		}
	}
}

func CreateVector(V *[MaxSellers]int) int {
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

func OnesVector(V *[MaxSellers]int) int {
	k := 0
	for j := 0; j < MaxSellers; j++ {
		if V[j] == 1 {
			k++
		}
	}
	return k
}

func ExecuteART1() {
	var R [MaxSellers]int
	var PE, P, E int
	var Test bool
	var i, j int
	count := 90
	exit := false
	var s int

	for {
		exit = true
		for i = 0; i < MaxShops; i++ {
			for j = 0; j < MaxClusters; j++ {
				if Members[j] > 0 {
					AndVectors(&R, &Data[i], &Clusters[j])
					PE = OnesVector(&R)
					P = OnesVector(&Clusters[j])
					E = OnesVector(&Data[i])
					Test = (float64(PE)/(Beta+float64(P)) > float64(E)/(Beta+float64(MaxSellers)))
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
		fmt.Printf("\nВектор-прототип %2d:\t", i+1)
		for j := 0; j < MaxSellers; j++ {
			fmt.Printf("%d ", Clusters[i][j])
		}
		fmt.Println()
		for k := 0; k < MaxShops; k++ {
			if Group[k] == i {
				fmt.Printf("Магазин %s:\t", ShopsName[k])
				for j := 0; j < MaxSellers; j++ {
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
	for i := 0; i < MaxSellers; i++ {
		if Data[p][i] == 0 && Sum[Group[p]][i] > max {
			best = i
			max = Sum[Group[p]][i]
		}
	}
	fmt.Printf("\nДля магазина %s ", ShopsName[p])
	if best >= 0 {
		fmt.Printf("лучший поставщик - %s\n", SellerName[best])
	} else {
		fmt.Println("нет лучших поставщиков")
	}
	fmt.Print("Уже выбраны: | ")
	for i := 0; i < MaxSellers; i++ {
		if Data[p][i] != 0 {
			fmt.Printf("%s | ", SellerName[i])
		}
	}
	fmt.Println()
}

func main() {
	Initialize()
	ExecuteART1()
	ShowClusters()
	for p := 0; p < MaxShops; p++ {
		MakeAdvise(p)
	}
}
