// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates the grid and row widgets.

import (
	"flag"
	"fmt"
	"image/color"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"gioui.org/widget"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var (
	th        = material.NewTheme(gofont.Collection())
	progress  float32
	count     float64
	startTime time.Time
)

func main() {
	flag.Parse()
	startProfile()
	progressIncrementer := make(chan float32)
	startTime = time.Now()
	go func() {
		for {
			// Max out at 500Hz
			time.Sleep(time.Millisecond * 2)
			progressIncrementer <- 0.002
		}
	}()

	go func() {
		w := app.NewWindow(app.Title("Gio grid demo"), app.Size(unit.Dp(600), unit.Dp(500)))
		var ops op.Ops
		for {
			select {
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					endProfile()
					os.Exit(0)
				case system.FrameEvent:
					count++
					gtx := layout.NewContext(&ops, e)
					// Lay out the form
					formLayout(gtx, th, true)
					// Apply the actual screen drawing
					e.Frame(gtx.Ops)
				}
			case pg := <-progressIncrementer:
				progress += pg
				if progress > 1 {
					progress = 0
				}
				w.Invalidate()
			}
		}
	}()
	app.Main()
}

func startProfile() {
	fmt.Printf("To extract data: go tool pprof --pdf cpu.prof > cpu-prof.pdf\n")
	fmt.Printf("To view realtime data, point browser to http://localhost:6060/debug/pprof/\n")
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic("could not create CPU profile, " + err.Error())
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic("could not start profiling, " + err.Error())
	}

	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			fmt.Printf("Error in http.ListAndServe %s\n", err.Error())
		}
	}()
}

func endProfile() {
	f, err := os.Create("mem.prof")
	if err != nil {
		panic("could not create memory profile: " + err.Error())
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		panic("could not write memory profile: " + err.Error())
	}
	_ = f.Close()
	pprof.StopCPUProfile()
}

type person struct {
	Selected bool
	Name     string
	Age      int
	Address  string
	Status   int
}

type personTable = []person

var (
	data = personTable{
		{Name: "First Person", Age: 21, Address: "Storgata 3", Status: 1},
		{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 33", Selected: true, Status: 1},
		{Name: "Nils Aure", Age: 23, Address: "Brogata 34"},
		{Name: "Kai Oppdal", Age: 28, Address: "Soleieveien 12"},
		{Name: "Gro Arneberg", Age: 29, Address: "Blomsterveien 22"},
		{Name: "Ole Kolås", Age: 21, Address: "Blåklokkevikua 33"},
		{Name: "Per Pedersen", Age: 22, Address: "Gamleveien 35"},
		{Name: "Nils Vukubråten", Age: 23, Address: "Nygata 64"},
		{Name: "Sindre Gratangen", Age: 28, Address: "Brosundet 34"},
		{Name: "Gro Nilsasveen", Age: 29, Address: "Blomsterveien 22"},
		{Name: "Petter Olsen", Age: 21, Address: "Katavågen 44"},
		{Name: "Per Pedersen", Age: 22, Address: "Nidelva 43"},
		{Name: "First Person", Age: 21, Address: "Storgata 3", Status: 1},
		{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 33", Selected: true, Status: 1},
		{Name: "Nils Aure", Age: 23, Address: "Brogata 34"},
		{Name: "Kai Oppdal", Age: 28, Address: "Soleieveien 12"},
		{Name: "Gro Arneberg", Age: 29, Address: "Blomsterveien 22"},
		{Name: "Ole Kolås", Age: 21, Address: "Blåklokkevikua 33"},
		{Name: "Per Pedersen", Age: 22, Address: "Gamleveien 35"},
		{Name: "Nils Vukubråten", Age: 23, Address: "Nygata 64"},
		{Name: "Sindre Gratangen", Age: 28, Address: "Brosundet 34"},
		{Name: "Gro Nilsasveen", Age: 29, Address: "Blomsterveien 22"},
		{Name: "Petter Olsen", Age: 21, Address: "Katavågen 44"},
		{Name: "Last Person", Age: 22, Address: "Nidelva 43"},
	}
	hdrPad = layout.Inset{Top: unit.Dp(8.0), Right: unit.Dp(0.0), Bottom: unit.Dp(2.0), Left: unit.Dp(6.0)}
)

func label(TextSize unit.Value, lblPad layout.Inset, s string) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		return lblPad.Layout(gtx, material.Label(th, TextSize, s).Layout)
	}
}

func headingCell(th *material.Theme) layout.ListElement {
	return func(gtx layout.Context, col int) layout.Dimensions {
		paint.ColorOp{Color: color.NRGBA{R: 88, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		switch col {
		case 0:
			return material.Label(th, th.TextSize, "No").Layout(gtx)
		case 1:
			return material.Label(th, th.TextSize, "Name").Layout(gtx)
		case 2:
			return material.Label(th, th.TextSize, "Address").Layout(gtx)
		case 3:
			return material.Label(th, th.TextSize, "Age").Layout(gtx)
		}
		return layout.Dimensions{}
	}
}

func gridCell(th *material.Theme, tbl personTable) layout.Cell {
	return func(gtx layout.Context, col, row int) layout.Dimensions {
		if col < len(tbl) {
			if row&1 == 0 {
				paint.ColorOp{Color: color.NRGBA{R: 244, G: 244, B: 244, A: 255}}.Add(gtx.Ops)
			} else {
				paint.ColorOp{Color: color.NRGBA{R: 225, G: 225, B: 225, A: 225}}.Add(gtx.Ops)
			}
			paint.PaintOp{}.Add(gtx.Ops)
			switch col {
			case 0:
				return material.Label(th, th.TextSize, fmt.Sprintf("%d", row)).Layout(gtx)
			case 1:
				return material.Label(th, th.TextSize, tbl[row].Name).Layout(gtx)
			case 2:
				return material.Label(th, th.TextSize, tbl[row].Address).Layout(gtx)
			case 3:
				return material.Label(th, th.TextSize, fmt.Sprintf("%d", tbl[row].Age)).Layout(gtx)
			}
		}
		return layout.Dimensions{}
	}
}

func ConfigureScrollbar(s *material.ScrollbarStyle) {
	s.Track.Color = color.NRGBA{0, 0, 0, 55}
	s.Indicator.MinorWidth = unit.Dp(10)
}

var (
	grid      = &widget.Grid{Grid: layout.Grid{}}
	oldAlloc  uint64
	alloc     uint64
	colWidths = []float32{50, 350, 350, 100}
	rowHeight = unit.Dp(23)
)

func formLayout(gtx layout.Context, th *material.Theme, showGrid bool) layout.Dimensions {
	// Read memory statistics to determine allocated memory size. First garbage collect.
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	oldAlloc = m.Alloc

	// Fixed layout with heading and statistics
	var children []layout.FlexChild
	children = append(children,
		layout.Rigid(label(th.TextSize.Scale(1.5), hdrPad, "Grid demo with widget.grid")),
		layout.Rigid(label(th.TextSize, hdrPad, "Running at "+
			fmt.Sprintf(" %0.1f frames/second, %v allocations", count/time.Since(startTime).Seconds(), alloc))))
	//layout.Rigid(material.Separator(th.Fg, unit.Dp(4), unit.Dp(4), unit.Dp(1))))

	// Add the grid itself only when showGrid is true. This is done to separate timing of grids from the rest
	if showGrid {
		myGrid := material.Table(th, grid)
		//myGrid.AnchorStrategy = material.Overlay
		ConfigureScrollbar(&myGrid.VScrollbarStyle)
		ConfigureScrollbar(&myGrid.HScrollbarStyle)
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				//return myGrid.Layout(gtx, len(data), rowHeight, colWidths, gridCell(th, data))
				return myGrid.Layout(gtx, len(data), rowHeight, colWidths, gridCell(th, data), headingCell(th))
			}))
	}

	// Then do actual layout
	d := layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx, children...)

	// Read memory statistics and calculate allocated memory size
	runtime.ReadMemStats(&m)
	alloc = m.Alloc - oldAlloc
	return d
}