package main

import (
	"math"
	"math/rand"
)

const (
	N            = 30
	screenWidth  = 900
	screenHeight = 600

	Pmax = 50
	Amax = 36

	Nin  = 12
	Nout = 4
	Nw   = (Nin * Nout)

	EAmax = 60
	EFmax = 20
	Erep  = 0.7

	NORTH = 0
	SOUTH = 1
	EAST  = 2
	WEST  = 3

	HERB_PLANE  = 0
	CARN_PLANE  = 1
	PLANT_PLANE = 2

	HERBIVORE = 0
	CARNIVORE = 1
	DEAD      = -1

	ACTION_LEFT  = 0
	ACTION_RIGHT = 1
	ACTION_MOVE  = 2
	ACTION_EAT   = 3

	HERB_FRONT      = 0  // травоядное впереди
	CARN_FRONT      = 1  // хищник впереди
	PLANT_FRONT     = 2  // еда впереди
	HERB_LEFT       = 3  // травоядное слева
	CARN_LEFT       = 4  // хищник слева
	PLANT_LEFT      = 5  // еда слева
	HERB_RIGHT      = 6  // травоядное справа
	CARN_RIGHT      = 7  // хищник справа
	PLANT_RIGHT     = 8  // еда справа
	HERB_PROXIMITY  = 9  // травоядное вблизи
	CARN_PROXIMITY  = 10 // хищник вблизи
	PLANT_PROXIMITY = 11 // еда вблизи

	maxSimulationCycles = 600
)

func getSRand() float64 {
	return rand.Float64()
}
func getRand(x float64) int {
	return int(math.Round(x * getSRand()))
}
func getWeight() int {
	return getRand(9)
}
