package main

import (
	"fmt"
	"gioui.org/unit"
	"gioui.org/widget"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type (
	// D is a shortcut
	D = layout.Dimensions
	// C is a shortcut
	C = layout.Context
)

func main() {
	go func() {
		var w *app.Window
		w = app.NewWindow(app.Title("Sized"), app.Size(unit.Dp(400), unit.Dp(400)))
		// w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			fmt.Printf("ConfigEvent, mode=%s, Size=%d,%d\n",
				e.Config.Mode.String(), e.Config.Size.X, e.Config.Size.Y)
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			l := material.H1(th, "Hello, Gio")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			l.Color = maroon
			l.Alignment = text.Middle
			layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D { return l.Layout(gtx) }),
				layout.Rigid(func(gtx C) D {
					return material.Button(th, button1, "Hover here").Layout(gtx)
				},
				),
			)
			e.Frame(gtx.Ops)
		}
	}
}

var (
	button1 = new(widget.Clickable)
)
