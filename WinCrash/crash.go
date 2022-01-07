package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"os"
)

type (
	// D is a shortcut
	D = layout.Dimensions
	// C is a shortcut
	C = layout.Context
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("Sized"), app.Size(unit.Dp(400), unit.Dp(400)))
		err := run(w)
		if err != nil {
			fmt.Println(err)
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
			in := layout.Inset{Top: unit.Dp(13), Right: unit.Dp(5), Bottom: unit.Dp(3), Left: unit.Dp(5)}
			layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D { return l.Layout(gtx) }),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button1, "Maximize").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button2, "Fullscreen").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button3, "700x400").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button4, "600x500").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button5, "Center").Layout(gtx)
							})
						}),
					)
				},
				),
			)
			for button1.Clicked() {
				w.Maximize()
			}
			for button2.Clicked() {
				w.Option(app.Fullscreen.Option())
			}
			for button3.Clicked() {
				w.Option(app.Size(unit.Dp(700), unit.Dp(400)))
			}
			for button4.Clicked() {
				w.Option(app.Size(unit.Dp(600), unit.Dp(500)))
			}
			for button5.Clicked() {
				w.Center()
			}

			e.Frame(gtx.Ops)
		}
	}
}

var (
	button1 = new(widget.Clickable)
	button2 = new(widget.Clickable)
	button3 = new(widget.Clickable)
	button4 = new(widget.Clickable)
	button5 = new(widget.Clickable)
	button6 = new(widget.Clickable)
)
