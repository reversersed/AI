package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	N              = 50
	TAU            = 0.05
	TAU_STEP       = 0.006
	D_H            = 0.01
	ALPHA          = 0.002
	START_INTERVAL = -2.0
	END_INTERVAL   = 8.0
	R              = 15
	A0             = 9
	A1             = 37.2
	A2             = 30.4
	A3             = 5.2
	A4             = -234.6
	A5             = 76
	screenWidth    = 800
	screenHeight   = 600
)

type Vars struct {
	c_x  float64
	u_t  [N]float64
	u    [N]float64
	u_n  [N]float64
	g    [R]float64
	n_t  int
	step int
	err  int
}

type Graph struct {
	w  float64
	h  float64
	dx float64
	bx float64
	rx float64
}

var (
	vars  *Vars
	graph Graph
)

func sign(x float64) float64 {
	if x > 0 {
		return 1.0
	} else if x < 0 {
		return -1.0
	}
	return 0.0
}

func f(x float64) float64 {
	return A0*math.Pow(x, 5) + A1*math.Pow(x, 4) + A2*math.Pow(x, 3) + A3*math.Pow(x, 2) + A4*x + A5
}

func teach(u [N]float64, g *[R]float64, n_t int, err *int, delta float64) {
	var delta_g [R]float64

	for r := 0; r < R; r++ {
		x_r := u[n_t-1-r]
		delta_g[r] = ALPHA * sign(x_r) * sign(delta)
	}
	for r := 0; r < R; r++ {
		g[r] += delta_g[r]
	}
	(*err)++
}

func newStep() {
	// Новое положение F
	for i := 0; i < N; i++ {
		vars.u_n[i] = f(vars.c_x + float64(i)*TAU)
	}
	vars.c_x += TAU_STEP
	if vars.c_x > END_INTERVAL {
		vars.c_x = START_INTERVAL
	}

	// Прогноз положения F
	for i := R; i < N; i++ {
		for r := 0; r < R; r++ {
			vars.u_t[i] = vars.g[r] * vars.u[i-1-r]
		}
	}
	for i := 0; i < R; i++ {
		for r := 0; r < R; r++ {
			vars.u_t[i] = vars.u[i]
		}
	}

	// Вычисление отклонения в обучаемой точке
	d_u := vars.u_n[vars.n_t] - vars.u_t[vars.n_t]
	if d_u < -D_H || d_u > D_H {
		teach(vars.u, &vars.g, vars.n_t, &vars.err, d_u)
	}
	vars.n_t++
	if vars.n_t >= N-1 {
		vars.n_t = R
	}
	vars.step++
}

func updateVars() {
	newStep()
	for i := 0; i < N; i++ {
		vars.u[i] = vars.u_n[i]
	}
}

type Game struct{}

func (g *Game) Update() error {
	updateVars()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	minVal, maxVal := vars.u[0], vars.u[0]
	for i := 1; i < N-1; i++ {
		if vars.u[i] < minVal {
			minVal = vars.u[i]
		}
		if vars.u[i] > maxVal {
			maxVal = vars.u[i]
		}
	}
	mid := (maxVal + minVal) / 2.0
	diff := maxVal - minVal
	minVal = mid - diff*1.5
	maxVal = mid + diff*1.5

	mtt := graph.h / (maxVal - minVal)
	if math.Abs(maxVal-minVal) < 1e-10 {
		mtt = 1.0
	}
	y_c := mtt * maxVal

	axis := y_c - mtt*vars.u[0]
	// Draw axes
	ebitenutil.DrawLine(screen, float64(graph.bx), 0, float64(graph.bx), graph.h, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, float64(graph.rx), 0, float64(graph.rx), graph.h, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, float64(graph.bx), 0, float64(graph.rx), 0, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, float64(graph.bx), graph.h, float64(graph.rx), graph.h, color.RGBA{0, 0, 255, 255})
	ebitenutil.DrawLine(screen, float64(graph.bx), axis, float64(graph.rx), axis, color.RGBA{0, 0, 255, 255})

	// Draw ticks
	for i := 1; i < N-1; i++ {
		x := graph.bx + float64(i)*graph.dx
		ebitenutil.DrawLine(screen, x, axis-3, x, axis+3, color.RGBA{0, 0, 255, 255})
	}

	// Draw vertical line at learning point
	x_nt := graph.bx + float64(vars.n_t-1)*graph.dx
	ebitenutil.DrawLine(screen, x_nt, 0, x_nt, graph.h, color.RGBA{0, 0, 255, 255})

	// Draw F curve (u array) in green
	x0 := graph.bx
	y0 := y_c - mtt*vars.u[0]
	for i := 1; i < N; i++ {
		x := graph.bx + float64(i)*graph.dx
		y := y_c - mtt*vars.u[i]
		ebitenutil.DrawLine(screen, x0, y0, x, y, color.RGBA{0, 255, 0, 255})
		x0, y0 = x, y
	}

	// Draw forecast (u_t array) in red
	x0 = graph.bx
	y0 = y_c - mtt*vars.u_t[0]
	for i := 1; i < N; i++ {
		x := graph.bx + float64(i)*graph.dx
		y := y_c - mtt*vars.u_t[i]
		ebitenutil.DrawLine(screen, x0, y0, x, y, color.RGBA{255, 0, 0, 255})
		x0, y0 = x, y
	}

	// Draw text with parameters in red
	font := basicfont.Face7x13
	textY := 20
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "Step   %d\n", vars.step)

	fmt.Fprintf(buf, "dH     %.2f\n", D_H)

	fmt.Fprintf(buf, "dG     %.3f\n", ALPHA)

	for i := 0; i < R; i++ {
		fmt.Fprintf(buf, "G[%d]   %.3f\n", i, vars.g[i])
	}
	fmt.Fprintf(buf, "x      %.3f\n", vars.c_x)
	fmt.Fprintf(buf, "y      %.3f\n", f(vars.c_x))

	funcDiff := 0.0
	for i := 0; i < N; i++ {
		funcDiff += math.Abs(vars.u[i] - vars.u_t[i])
	}
	fmt.Fprintf(buf, "diff   %.3f\n", funcDiff)

	text.Draw(screen, buf.String(), font, int(graph.rx)-220, textY, color.RGBA{255, 0, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	graph.w = screenWidth
	graph.h = screenHeight
	graph.dx = (graph.w + 1.0) / N
	graph.bx = (graph.w + 1.0 - graph.dx*float64(N-1)) / 2.0
	graph.rx = graph.bx + float64(N-1)*graph.dx

	vars = new(Vars)
	vars.c_x = START_INTERVAL
	for i := 0; i < N; i++ {
		vars.u[i] = 0.0
		vars.u_t[i] = 0.0
		vars.u_n[i] = 0.0
	}
	for i := 0; i < R; i++ {
		vars.g[i] = 0.0
	}

	vars.n_t = N/2 + 1
	vars.step = 1
	vars.err = 0

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Lab 4")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
