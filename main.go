package main

import (
	"fmt"
	"github.com/samuel/go-pcx/pcx"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"image"
	"log"
	"math"
	"os"
	"time"
)

var (
	tileset       map[string]image.Image
	offset        int
	scale         float64
	board         Board
	board_offset  = 2
	rows, columns = 20, 10
)

func init() {
	tilefiles := []string{"forest", "grass", "marsh", "village"}
	tileset = make(map[string]image.Image)

	for _, tf := range tilefiles {
		file, err := os.Open("u2_" + tf + ".pcx")
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		img, err := pcx.Decode(file)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
		tileset[tf] = img
	}

	img := tileset["village"]

	offset = img.Bounds().Dx() / 2
	scale = float64(img.Bounds().Dx())
}

func main() {
	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "Tile Map")
	if err != nil {
		log.Println(err)
		return
	}
	defer wnd.Destroy()

	var mx, my, action float64

	board := Board{
		OffsetX:   board_offset,
		OffsetY:   board_offset,
		Rows:      rows,
		Columns:   columns,
		Positions: make([]Position, rows*columns),
	}

	wnd.MouseMove = func(x, y int) {
		mx, my = float64(x), float64(y)
	}

	wnd.MouseDown = func(button, x, y int) {
		if button == 1 { /// mouse left == 1, mouse right == 3
			action = 1
			board.AddTile(x, y)
		}
		if button == 3 {
			action = 1
		}
	}

	wnd.KeyDown = func(scancode int, rn rune, name string) {
		switch name {
		case "Escape":
			wnd.Close()
		case "Space":
			action = 1

		case "Enter":
			action = 1

		}
	}
	wnd.SizeChange = func(w, h int) {
		cv.SetBounds(0, 0, w, h)
	}

	lastTime := time.Now()

	wnd.MainLoop(func() {
		now := time.Now()
		diff := now.Sub(lastTime)
		lastTime = now
		action -= diff.Seconds() * 3
		action = math.Max(0, action)

		w, h := float64(cv.Width()), float64(cv.Height())

		// Clear the screen
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)

		new_grid(cv)

		// Draw a circle around the cursor
		cv.SetStrokeStyle("#778899")
		cv.SetLineWidth(2)
		cv.BeginPath()

		tx, ty := fit_gridf(mx, my)
		open_tl, open_br := action*12, action*24
		cv.Rect(tx-open_tl, ty-open_tl, scale+open_br, scale+open_br)
		cv.Stroke()
		cv.FillText(fmt.Sprintf("x:", tx, "y:", ty), tx, ty-2)

		// Draw tiles where the user has clicked
		for _, p := range board.Positions {
			t := p.PTile
			if t != nil {
				cv.PutImageData(tileset[t.Type].(*image.RGBA), t.X, t.Y)
			}
		}
	})
}

func fit_gridf(mx, my float64) (tx, ty float64) {
	nxt := offset * 2
	nx, ny := int(mx), int(my)
	tx, ty = float64((nx/nxt)*nxt), float64((ny/nxt)*nxt)
	return
}

func new_grid(cv *canvas.Canvas) {
	penwidth := 1.0
	ix, iy := scale*2, scale*2
	vstep, hstep := scale, scale
	step := 1.0 * scale

	for x := ix; x <= hstep*step; x += step {
		cv.SetStrokeStyle("#1e90ff")
		cv.SetLineWidth(penwidth)
		cv.BeginPath()
		cv.MoveTo(x, 0)
		cv.LineTo(x, vstep*step)
		cv.Stroke()
	}

	for y := iy; y <= vstep*step; y += step {
		cv.SetStrokeStyle("#1e90ff")
		cv.SetLineWidth(penwidth)
		cv.BeginPath()
		cv.MoveTo(0, y)
		cv.LineTo(hstep*step, y)
		cv.Stroke()
	}
}
