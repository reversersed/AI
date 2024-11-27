package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func Init() {
	for l := 0; l < 3; l++ { // очистка карты
		for y := 0; y < N; y++ {
			for x := 0; x < N; x++ {
				g.mapArray[l][y][x] = 0
			}
		}
	}
	for p := 0; p < Pmax; p++ { // посадка растений
		plants[p] = AddInEmptyCell(PLANT_PLANE)
	}
	for a := 0; a < Amax; a++ { // инициализация агентов
		agents[a].Type = CARNIVORE
		if a < (Amax / 2) {
			agents[a].Type = HERBIVORE
		}
		InitAgent(&agents[a])
	}
}

func AddInEmptyCell(Level int) TXY {
	var res TXY = TXY{
		X: getRand(N - 1),
		Y: getRand(N - 1),
	}

	for g.mapArray[Level][res.Y][res.X] != 0 {
		res.X = getRand(N - 1)
		res.Y = getRand(N - 1)
	}

	g.mapArray[Level][res.Y][res.X]++
	return res
}
func AgentToMap(agent *TAgent) {
	agent.location = AddInEmptyCell(agent.Type)
	agent.direction = getRand(3)
}

func InitAgent(agent *TAgent) {
	agent.energy = (EAmax / 2)
	agent.age = 0
	agent.generation = 1
	agentTypeCounts[agent.Type]++
	AgentToMap(agent)
	for i := 0; i < (Nin * Nout); i++ {
		agent.weights[i] = getWeight()
	}
	for i := 0; i < Nout; i++ {
		agent.biass[i] = getWeight()
	}
}

func Clip(z int) int {
	if z > N-1 {
		z = (z % N)
	} else if z < 0 {
		z = (N + z)
	}
	return z
}

func Percept(x, y int, inputs *[Nin]int, vectorStartPoint int, offsets []TXY, neg int) {
	for p := HERB_PLANE; p <= PLANT_PLANE; p++ {
		i := 0
		(*inputs)[vectorStartPoint+p] = 0
		for offsets[i].X != 9 {
			xoff := Clip(x + (offsets[i].X * neg))
			yoff := Clip(y + (offsets[i].Y * neg))
			if g.mapArray[p][yoff][xoff] != 0 {
				(*inputs)[vectorStartPoint+p]++
			}
			i++
		}
	}
}

func Turn(action int, agent *TAgent) {
	if agent.direction == NORTH {
		if action == ACTION_LEFT {
			agent.direction = WEST
		} else {
			agent.direction = EAST
		}
	}
	if agent.direction == SOUTH {
		if action == ACTION_LEFT {
			agent.direction = EAST
		} else {
			agent.direction = WEST
		}
	}

	if agent.direction == EAST {
		if action == ACTION_LEFT {
			agent.direction = NORTH
		} else {
			agent.direction = SOUTH
		}
	}

	if agent.direction == WEST {
		if action == ACTION_LEFT {
			agent.direction = SOUTH
		} else {
			agent.direction = NORTH
		}
	}
}

func Move(agent *TAgent) {
	offsets := [4]TXY{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	g.mapArray[agent.Type][agent.location.Y][agent.location.X]--
	agent.location.X = Clip(agent.location.X + offsets[agent.direction].X)
	agent.location.Y = Clip(agent.location.Y + offsets[agent.direction].Y)
	g.mapArray[agent.Type][agent.location.Y][agent.location.X]++
}
func KillAgent(agent *TAgent) {
	agentDeaths[agent.Type]++
	g.mapArray[agent.Type][agent.location.Y][agent.location.X]--
	agentTypeCounts[agent.Type]--
	if bestAgent[agent.Type] == nil || agent.age > bestAgent[agent.Type].age {
		bestAgent[agent.Type] = agent
	}
	if agentTypeCounts[agent.Type] < (Amax / 4) {
		InitAgent(agent) // инициализация агента
	} else { // конец агента
		agent.location.X = -1
		agent.location.Y = -1
		agent.Type = DEAD
	}
}
func ReproduceAgent(agent *TAgent) {
	var child *TAgent
	i := 0

	if agentTypeCounts[agent.Type] < (Amax / 2) {
		for i < Amax {
			if agents[i].Type == DEAD {
				break
			}
			i++
		}
		if i < Amax {
			child = new(TAgent)
			*child = *agent
			AgentToMap(child)
			if getSRand() <= 0.4 {
				child.weights[getRand(Nw)] = getWeight()
			}
			child.generation = child.generation + 1
			child.age = 0
			if agentMaxGen[child.Type] < child.generation {
				agentMaxGen[child.Type] = child.generation
			}
			child.energy = (EAmax / 2) // энергия
			agent.energy = (EAmax / 2)
			agentTypeCounts[child.Type]++
			agentTypeReproductions[child.Type]++
		}
	}
}

func ChooseObject(plane, ax, ay int, offsets []TXY, neg int, ox, oy *int) int {
	xoff := 0
	yoff := 0
	i := 0

	for offsets[i].X != 9 {
		xoff = Clip(ax + (offsets[i].X * neg))
		yoff = Clip(ay + (offsets[i].Y * neg))
		if g.mapArray[plane][yoff][xoff] != 0 {
			*ox = xoff
			*oy = yoff
			return 1
		}
		i++
	}
	return 0
}
func Eat(agent *TAgent) {
	var plane int
	var ox, oy int
	var ret int

	if agent.Type == CARNIVORE {
		plane = HERB_PLANE
	} else if agent.Type == HERBIVORE {
		plane = PLANT_PLANE
	}

	ax := agent.location.X
	ay := agent.location.Y

	// Определение направления и выбор объекта
	switch agent.direction {
	case NORTH:
		ret = ChooseObject(plane, ax, ay, northProx, 1, &ox, &oy)
	case SOUTH:
		ret = ChooseObject(plane, ax, ay, northProx, -1, &ox, &oy)
	case WEST:
		ret = ChooseObject(plane, ax, ay, westProx, 1, &ox, &oy)
	case EAST:
		ret = ChooseObject(plane, ax, ay, westProx, -1, &ox, &oy)
	}

	if ret != 0 {
		if plane == PLANT_PLANE {
			for i := 0; i < Pmax; i++ {
				if plants[i].X == ox && plants[i].Y == oy {
					agent.energy += EFmax
					if agent.energy > EAmax {
						agent.energy = EAmax
					}
					// Логика уменьшения количества растений на карте
					plants[i] = AddInEmptyCell(PLANT_PLANE)
					eated[0]++
					fmt.Println("Была съедена трава")
					break
				}
			}
		} else if plane == HERB_PLANE {
			for i := 0; i < Amax; i++ {
				if agents[i].location.X == ox && agents[i].location.Y == oy {
					agent.energy += EFmax * 2
					if agent.energy > EAmax {
						agent.energy = EAmax
					}
					KillAgent(&agents[i])
					eated[1]++
					fmt.Println("Было съедено травоядное")
					break
				}
			}
		}

		if agent.energy > (Erep * EAmax) {
			ReproduceAgent(agent)
			agentBirths[agent.Type]++
		}
	}
}

func Simulate(agent *TAgent) {
	x := agent.location.X
	y := agent.location.Y

	// Восприятие по направлению
	switch agent.direction {
	case NORTH:
		Percept(x, y, &agent.inputs, HERB_FRONT, northFront, 1)
		Percept(x, y, &agent.inputs, HERB_LEFT, northLeft, 1)
		Percept(x, y, &agent.inputs, HERB_RIGHT, northRight, 1)
		Percept(x, y, &agent.inputs, HERB_PROXIMITY, northProx, 1)
	case SOUTH:
		Percept(x, y, &agent.inputs, HERB_FRONT, northFront, -1)
		Percept(x, y, &agent.inputs, HERB_LEFT, northLeft, -1)
		Percept(x, y, &agent.inputs, HERB_RIGHT, northRight, -1)
		Percept(x, y, &agent.inputs, HERB_PROXIMITY, northProx, -1)
	case WEST:
		Percept(x, y, &agent.inputs, HERB_FRONT, westFront, 1)
		Percept(x, y, &agent.inputs, HERB_LEFT, westLeft, 1)
		Percept(x, y, &agent.inputs, HERB_RIGHT, westRight, 1)
		Percept(x, y, &agent.inputs, HERB_PROXIMITY, westProx, 1)
	case EAST:
		Percept(x, y, &agent.inputs, HERB_FRONT, westFront, -1)
		Percept(x, y, &agent.inputs, HERB_LEFT, westLeft, -1)
		Percept(x, y, &agent.inputs, HERB_RIGHT, westRight, -1)
		Percept(x, y, &agent.inputs, HERB_PROXIMITY, westProx, -1)
	}
	// Расчет решений
	for out := 0; out < Nout; out++ {
		agent.actions[out] = agent.biass[out] // Инициализация выхода смещением
		for in := 0; in < Nin; in++ {
			agent.actions[out] += (agent.inputs[in] * agent.weights[(out*Nin)+in]) // Взвешенные входы
		}
	}

	// Принятие решения
	largest := -9
	winner := -1
	for out := 0; out < Nout; out++ {
		if agent.actions[out] >= largest {
			largest = agent.actions[out]
			winner = out
		}
	}
	// Выполнение решения
	switch winner {
	case ACTION_LEFT:
		Turn(ACTION_LEFT, agent)
	case ACTION_RIGHT:
		Turn(ACTION_RIGHT, agent)
	case ACTION_MOVE:
		Move(agent)
	case ACTION_EAT:
		Eat(agent)
	}

	// Затраты энергии
	if agent.Type == HERBIVORE {
		agent.energy -= 2
	} else {
		agent.energy -= 1
	}

	// Проверка на гибель
	if agent.energy <= 0 {
		KillAgent(agent) // Гибель агента
	} else {
		agent.age++
		if agent.age > agentMaxAge[agent.Type] {
			// Фиксируем старейшего агента
			agentMaxAge[agent.Type] = agent.age
			agentMaxPtr[agent.Type] = agent
		}
	}
}

func ShowStat() {
	fmt.Println("Результаты:")
	fmt.Printf("Травоядных всего                - %d\n", agentTypeCounts[HERBIVORE])
	fmt.Printf("Хищников всего                  - %d\n", agentTypeCounts[CARNIVORE])
	fmt.Printf("Возраст травоядных              - %d\n", agentMaxAge[HERBIVORE])
	fmt.Printf("Возраст хищников                - %d\n", agentMaxAge[CARNIVORE])
	fmt.Printf("Рождений травоядных             - %d\n", agentBirths[HERBIVORE])
	fmt.Printf("Рождений хищников               - %d\n", agentBirths[CARNIVORE])
	fmt.Printf("Гибелей травоядных              - %d\n", agentDeaths[HERBIVORE])
	fmt.Printf("Гибелей хищников                - %d\n", agentDeaths[CARNIVORE])
	fmt.Printf("Репродукций травоядных          - %d\n", agentTypeReproductions[HERBIVORE])
	fmt.Printf("Репродукций хищников            - %d\n", agentTypeReproductions[CARNIVORE])
	fmt.Printf("Наибольшие поколения травоядных - %d\n", agentMaxGen[HERBIVORE])
	fmt.Printf("Наибольшие поколения хищников   - %d\n", agentMaxGen[CARNIVORE])
	fmt.Printf("Съедено травы                   - %d\n", eated[0])
	fmt.Printf("Съедено травоядных              - %d\n", eated[1])
	fmt.Println()

	if bestAgent[HERBIVORE] != nil {
		fmt.Println("Веса лучшего травоядного:")
		for i := 0; i < Nout; i++ {
			fmt.Printf("%4d ", bestAgent[HERBIVORE].biass[i])
		}
		fmt.Println()
		for o := 0; o < Nout; o++ {
			for i := 0; i < Nin; i++ {
				fmt.Printf("%4d ", bestAgent[HERBIVORE].weights[o*Nin+i])
			}
			fmt.Println()
		}
	}

	if bestAgent[CARNIVORE] != nil {
		fmt.Println("Веса лучшего хищника:")
		for i := 0; i < Nout; i++ {
			fmt.Printf("%4d ", bestAgent[CARNIVORE].biass[i])
		}
		fmt.Println()
		for o := 0; o < Nout; o++ {
			for i := 0; i < Nin; i++ {
				fmt.Printf("%4d ", bestAgent[CARNIVORE].weights[o*Nin+i])
			}
			fmt.Println()
		}
	}
}

func main() {
	g = &Game{}
	GraphInit()
	Init()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game Example")
	ebiten.RunGame(g)

	ShowStat()
}
