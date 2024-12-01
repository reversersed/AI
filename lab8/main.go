package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func Init() {
	for l := 0; l < 2; l++ { // очистка карты
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
	agent.location = AddInEmptyCell(HERBIVORE)
	agent.direction = getRand(3)
}
func createRange(agent *TAgent) {
	agent.vectorRange = [9]int{
		getRand(float64(tableLimit[0])),
		getRand(float64(tableLimit[1])),
		getRand(float64(tableLimit[2])),
		getRand(float64(tableLimit[3])),
		getRand(float64(tableLimit[4])),
		getRand(float64(tableLimit[5])),
		getRand(float64(tableLimit[6])),
		getRand(float64(tableLimit[7])),
		getRand(float64(tableLimit[8])),
	}
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			agent.vectorRange[i] = getRand(float64(tableLimit[i/3]))
		} else {
			agent.vectorRange[i] = getRand(float64(agent.vectorRange[i-1]))
		}
	}
	calculatePrice(agent)
}
func mutateRange(agent *TAgent) {
	agent.weights[getRand(Nw)] = getWeight()
	first := getRand(9)
	second := getRand(9)
	for first == second {
		second = getRand(9)
	}
	agent.vectorRange[first]--
	agent.vectorRange[second]++
	calculatePrice(agent)
}
func calculatePrice(agent *TAgent) {
	agent.price = agent.vectorRange[2]*prices[0] + agent.vectorRange[5]*prices[1] + agent.vectorRange[8]*prices[2]
	for i := 0; i < 9; i += 3 {
		max := int((float64(agent.vectorRange[i]) / float64(tableLimit[i])) * float64(agent.vectorRange[i]))
		if agent.vectorRange[i+1] > max {
			agent.price -= (agent.vectorRange[i+1] - max) * ticketValue
		}
	}
}
func InitAgent(agent *TAgent) {
	agent.energy = (EAmax / 2)
	agent.age = 0
	agent.generation = 1
	agentTypeCounts++
	AgentToMap(agent)
	for i := 0; i < (Nin * Nout); i++ {
		agent.weights[i] = getWeight()
	}
	for i := 0; i < Nout; i++ {
		agent.biass[i] = getWeight()
	}
	createRange(agent)
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
	if agent.Type == DEAD {
		return
	}
	offsets := [4]TXY{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	g.mapArray[HERBIVORE][agent.location.Y][agent.location.X]--
	agent.location.X = Clip(agent.location.X + offsets[agent.direction].X)
	agent.location.Y = Clip(agent.location.Y + offsets[agent.direction].Y)
	g.mapArray[HERBIVORE][agent.location.Y][agent.location.X]++
}
func KillAgent(agent *TAgent) {
	if agent.Type == DEAD {
		return
	}
	agentDeaths++
	g.mapArray[HERBIVORE][agent.location.Y][agent.location.X]--
	agentTypeCounts--
	if bestAgent == nil || agent.price > bestAgent.price {
		bestAgent = agent
	}
	if agentTypeCounts < (Amax / 2) {
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

	if agentTypeCounts < (Amax / 2) {
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
				mutateRange(child)
			}
			child.generation = child.generation + 1
			child.age = 0
			if agentMaxGen < child.generation {
				agentMaxGen = child.generation
			}
			child.energy = (EAmax / 2) // энергия
			agent.energy = (EAmax / 2)
			agentTypeCounts++
			agentTypeReproductions++
		}
	}
}

func ChooseObject(agent *TAgent, ax, ay int, offsets []TXY, neg int, ox, oy *int) int {
	xoff := 0
	yoff := 0
	i := 0

	for offsets[i].X != 9 {
		xoff = Clip(ax + (offsets[i].X * neg))
		yoff = Clip(ay + (offsets[i].Y * neg))
		if g.mapArray[HERB_PLANE][yoff][xoff] != 0 {
			for _, a := range agents {
				if a.location.X == xoff && a.location.Y == yoff && a.Type != DEAD && a.price <= agent.price {
					*ox = xoff
					*oy = yoff
					return HERB_PLANE
				}
			}
		}
		if g.mapArray[PLANT_PLANE][yoff][xoff] != 0 {
			*ox = xoff
			*oy = yoff
			return PLANT_PLANE
		}
		i++
	}
	return 0
}
func Eat(agent *TAgent) {
	var ox, oy int
	var ret int

	ax := agent.location.X
	ay := agent.location.Y

	// Определение направления и выбор объекта
	switch agent.direction {
	case NORTH:
		ret = ChooseObject(agent, ax, ay, northProx, 1, &ox, &oy)
	case SOUTH:
		ret = ChooseObject(agent, ax, ay, northProx, -1, &ox, &oy)
	case WEST:
		ret = ChooseObject(agent, ax, ay, westProx, 1, &ox, &oy)
	case EAST:
		ret = ChooseObject(agent, ax, ay, westProx, -1, &ox, &oy)
	}

	if ret != 0 {
		if ret == PLANT_PLANE {
			for i := 0; i < Pmax; i++ {
				if plants[i].X == ox && plants[i].Y == oy {
					agent.energy += EFmax
					if agent.energy > EAmax {
						agent.energy = EAmax
					}
					// Логика уменьшения количества растений на карте
					plants[i] = AddInEmptyCell(PLANT_PLANE)
					fmt.Println("Была съедена трава")
					eated[0]++
					break
				}
			}
		} else if ret == HERB_PLANE {
			for i := 0; i < Amax; i++ {
				if agents[i].location.X == ox && agents[i].location.Y == oy {
					agent.energy += EFmax * 2
					if agent.energy > EAmax {
						agent.energy = EAmax
					}
					KillAgent(&agents[i])
					eated[1]++
					fmt.Println("Был съеден агент")
					break
				}
			}
		}

		if agent.energy > (Erep * EAmax) {
			ReproduceAgent(agent)
			agentBirths++
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
	agent.energy -= 1

	// Проверка на гибель
	if agent.energy <= 0 || (bestAgent != nil && agent.price < bestAgent.price/10) {
		KillAgent(agent) // Гибель агента
	} else {
		agent.age++
		if agent.age > agentMaxAge {
			// Фиксируем старейшего агента
			agentMaxAge = agent.age
			agentMaxPtr = agent
		}
	}
}

func ShowStat() {
	fmt.Println("Результаты:")
	fmt.Printf("Агентов всего                   - %d\n", agentTypeCounts)
	fmt.Printf("Возраст агентов                 - %d\n", agentMaxAge)
	fmt.Printf("Рождений агентов                - %d\n", agentBirths)
	fmt.Printf("Гибелей агентов                 - %d\n", agentDeaths)
	fmt.Printf("Репродукций агентов             - %d\n", agentTypeReproductions)
	fmt.Printf("Наибольшие поколения агентов    - %d\n", agentMaxGen)
	fmt.Printf("Съедено травы                   - %d\n", eated[0])
	fmt.Printf("Съедено агентов                 - %d\n", eated[1])
	fmt.Println()

	if agentMaxPtr != nil {
		fmt.Println("Веса старейшего агента:")
		for i := 0; i < Nout; i++ {
			fmt.Printf("%4d ", agentMaxPtr.biass[i])
		}
		fmt.Println()
		for o := 0; o < Nout; o++ {
			for i := 0; i < Nin; i++ {
				fmt.Printf("%4d ", agentMaxPtr.weights[o*Nin+i])
			}
			fmt.Println()
		}
		fmt.Println("Цена старейшего агента: ", agentMaxPtr.price)
		fmt.Println(agentMaxPtr.vectorRange[0], "\t", agentMaxPtr.vectorRange[1], "\t", agentMaxPtr.vectorRange[2])
		fmt.Println(agentMaxPtr.vectorRange[3], "\t", agentMaxPtr.vectorRange[4], "\t", agentMaxPtr.vectorRange[5])
		fmt.Println(agentMaxPtr.vectorRange[6], "\t", agentMaxPtr.vectorRange[7], "\t", agentMaxPtr.vectorRange[8])
	}
	if bestAgent != nil {
		fmt.Println("\n\nВеса лучшего агента:")
		for i := 0; i < Nout; i++ {
			fmt.Printf("%4d ", bestAgent.biass[i])
		}
		fmt.Println()
		for o := 0; o < Nout; o++ {
			for i := 0; i < Nin; i++ {
				fmt.Printf("%4d ", bestAgent.weights[o*Nin+i])
			}
			fmt.Println()
		}
		fmt.Println("Цена лучшего агента: ", bestAgent.price)
		fmt.Println(bestAgent.vectorRange[0], "\t", bestAgent.vectorRange[1], "\t", bestAgent.vectorRange[2])
		fmt.Println(bestAgent.vectorRange[3], "\t", bestAgent.vectorRange[4], "\t", bestAgent.vectorRange[5])
		fmt.Println(bestAgent.vectorRange[6], "\t", bestAgent.vectorRange[7], "\t", bestAgent.vectorRange[8])
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
