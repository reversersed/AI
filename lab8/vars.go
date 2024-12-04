package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type TXY struct {
	X, Y int
}

type TAgent struct {
	Type        int
	energy      int
	age         int
	generation  int
	location    TXY
	direction   int
	inputs      [Nin]int
	weights     [Nin * Nout]int
	biass       [Nout]int
	actions     [Nout]int
	price       int
	vectorRange [9]int
}

var (
	dx, dy, bx, by float32
	Hh             *ebiten.Image
	Hp             color.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	Hs             color.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	NGrassEated = color.RGBA{R: 0, G: 255, B: 0, A: 255}     // Зеленый
	NAgentEat   = color.RGBA{R: 255, G: 255, B: 0, A: 255}   // Желтый
	NAgentEated = color.RGBA{R: 255, G: 0, B: 0, A: 255}     // Красный
	NBirth      = color.RGBA{R: 255, G: 192, B: 203, A: 255} // Розовый
	NDisabled   = color.RGBA{R: 192, G: 192, B: 192, A: 255} // Серый

	plants [Pmax]TXY
	agents [Amax]TAgent

	northFront = []TXY{{-2, -2}, {-2, -1}, {-2, 0}, {-2, 1}, {-2, 2}}
	northLeft  = []TXY{{0, -2}, {-1, -2}}
	northRight = []TXY{{0, 2}, {-1, 2}}
	northProx  = []TXY{{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}}
	westFront  = []TXY{{2, -2}, {1, -2}, {0, -2}, {-1, -2}, {-2, -2}}
	westLeft   = []TXY{{2, 0}, {2, -1}}
	westRight  = []TXY{{-2, 0}, {-2, -1}}
	westProx   = []TXY{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}}

	agentTypeCounts        = 0     // коиличество агентов
	agentMaxAge            = 0     // возраст агентов по типам
	agentBirths            = 0     // количество рождений по типам
	agentDeaths            = 0     // количество гибелей по типам
	agentMaxPtr            *TAgent // старейшие агенты по типам
	agentTypeReproductions = 0     // количество репродукций по типам
	bestAgent              *TAgent // старейшие погибшие агенты
	agentMaxGen            = 0     // наибольшие поколения по типам
	eated                  = [2]int{0, 0}

	g               *Game
	lastSleepTime   int64
	simulationCycle       = 0
	firstStart      int64 = 0

	tableLimit = [9]int{
		80, 50, 70,
		40, 45, 30,
		20, 25, 18,
	}
	prices = []int{
		20,
		30,
		50,
	}
	ticketValue = max(prices) * 2
)
