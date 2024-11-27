package main

import (
	"image/color"
)

type TXY struct {
	X, Y int
}

type TAgent struct {
	Type       int
	energy     int
	age        int
	generation int
	location   TXY
	direction  int
	inputs     [Nin]int
	weights    [Nin * Nout]int
	biass      [Nout]int
	actions    [Nout]int
}

var (
	dx, dy, bx, by float32
	Hc             color.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Hh             color.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	Hp             color.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	Hs             color.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	plants [Pmax]TXY
	agents [Amax]TAgent

	northFront = []TXY{{-2, -2}, {-2, -1}, {-2, 0}, {-2, 1}, {-2, 2}, {9, 9}}
	northLeft  = []TXY{{0, -2}, {-1, -2}, {9, 9}}
	northRight = []TXY{{0, 2}, {-1, 2}, {9, 9}}
	northProx  = []TXY{{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {9, 9}}
	westFront  = []TXY{{2, -2}, {1, -2}, {0, -2}, {-1, -2}, {-2, -2}, {9, 9}}
	westLeft   = []TXY{{2, 0}, {2, -1}, {9, 9}}
	westRight  = []TXY{{-2, 0}, {-2, -1}, {9, 9}}
	westProx   = []TXY{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {9, 9}}

	agentTypeCounts        = [2]int{0, 0} // количество агентов по типам
	agentMaxAge            = [2]int{0, 0} // возраст агентов по типам
	agentBirths            = [2]int{0, 0} // количество рождений по типам
	agentDeaths            = [2]int{0, 0} // количество гибелей по типам
	agentMaxPtr            [2]*TAgent     // старейшие агенты по типам
	agentTypeReproductions = [2]int{0, 0} // количество репродукций по типам
	bestAgent              [2]*TAgent     // старейшие погибшие агенты
	agentMaxGen            = [2]int{0, 0} // наибольшие поколения по типам
	eated                  = [2]int{0, 0}

	g               *Game
	simulationCycle = 0
	lastSleepTime   int64
	latency         int64 = 50
)
