package main

import (
	"image"
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

	img := image.NewRGBA(image.Rect(0, 0, int(dx), int(dy)))
	gradientOffset := img.Bounds().Dx() / 4
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			adjust := x - gradientOffset
			if adjust < 0 {
				adjust = 0
			}
			r := uint8(adjust * 255 / img.Bounds().Dx())
			img.Set(x, y, color.RGBA{r, 0, 0, 255})
		}
	}

	Hh = ebiten.NewImageFromImage(img)
}

func (g *Game) Update() error {

	if time.Now().UnixMilli()-lastSleepTime < latency {
		return nil
	}
	lastSleepTime = time.Now().UnixMilli()

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
				col = nil

				centerX := float64(Hh.Bounds().Dx()) / 2.0
				centerY := float64(Hh.Bounds().Dy()) / 2.0

				op := &ebiten.DrawImageOptions{}
				// Сначала переносим изображение в точку поворота
				op.GeoM.Translate(-centerX, -centerY) // Сдвигаем влево и вверх на половину ширины и высоты

				for _, a := range agents {
					if a.location.X == x && a.location.Y == y {
						switch a.direction {
						case EAST:
							op.GeoM.Rotate(0)
						case WEST:
							op.GeoM.Rotate((3.14 / 180.0) * 180.0)
						case NORTH:
							op.GeoM.Rotate((3.14 / 180.0) * 270.0)
						case SOUTH:
							op.GeoM.Rotate((3.14 / 180.0) * 90.0)
						}
					}
				}

				op.GeoM.Translate(float64(bx+float32(x)*dx)+centerX, float64(by+float32(y)*dy)+centerY)
				screen.DrawImage(Hh, op)
			} else if g.mapArray[PLANT_PLANE][y][x] != 0 {
				col = Hp
			} else {
				col = Hs
			}
			if col != nil {
				vector.DrawFilledRect(screen, bx+float32(x)*dx, by+float32(y)*dy, dx, dy, col, true)
			}
		}
	}
	vector.DrawFilledRect(screen, 0, 0, dx, dy, color.Black, true)
}
