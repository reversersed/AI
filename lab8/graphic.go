package main

import (
	"errors"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	mapArray [2][N][N]int
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func GraphInit() {
	dx = float32(screenWidth) / float32(N)
	dy = float32(screenHeight) / float32(N)
	bx = (float32(screenWidth) - dx*float32(N)) / 2
	by = (float32(screenHeight) - dy*float32(N)) / 2
}

func (g *Game) Update() error {

	if time.Now().UnixMilli()-lastSleepTime < latency {
		return nil
	}
	lastSleepTime = time.Now().UnixMilli()

	simulationCycle++
	if simulationCycle == maxSimulationCycles {
		return errors.New("")
	}

	for i := 0; i < Amax; i++ {
		Simulate(&agents[i])
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			var col color.Color
			if g.mapArray[HERB_PLANE][y][x] != 0 {
				col = Hh
			} else if g.mapArray[PLANT_PLANE][y][x] != 0 {
				col = Hp
			} else {
				col = Hs
			}
			vector.DrawFilledRect(screen, bx+float32(x)*dx, by+float32(y)*dy, dx, dy, col, true)
		}
	}
	vector.DrawFilledRect(screen, 0, 0, dx, dy, color.Black, true)
}
