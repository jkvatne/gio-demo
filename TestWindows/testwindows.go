package main

import (
	"flag"
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

var testNo = 0

func main() {
	flag.IntVar(&testNo, "test", 0, "Test number 0..5")
	flag.Parse()

	go func() {
		// Test different startup configurations
		var w *app.Window
		switch testNo {
		case 1:
			w = app.NewWindow(app.Title("Maximized"), app.Maximized.Option())
		case 2:
			fmt.Println("NB: The window will now be minimized, and not visible on screen")
			w = app.NewWindow(app.Title("Minimized"), app.Minimized.Option())
		case 3:
			w = app.NewWindow(app.Title("Centered"), app.Centered())
		case 4:
			w = app.NewWindow(app.Title("Positioned"), app.Pos(unit.Dp(70), unit.Dp(70)))
		case 5:
			w = app.NewWindow(app.Title("Sized"), app.Size(unit.Dp(400), unit.Dp(400)))
		case 6:
			w = app.NewWindow(app.Title("Fullscreen"), app.Fullscreen.Option())
		default:
			fmt.Println("Specify test number on command line, -test=n, where n=1..6")
			fmt.Println("Example: go run testwindows.go -test=1")
			fmt.Println("1 = Maximized window")
			fmt.Println("2 = Minimized window")
			fmt.Println("3 = Centered window")
			fmt.Println("4 = Positioned window")
			fmt.Println("5 = Sized window")
			fmt.Println("6 = Fullscreen window")
			os.Exit(1)
		}
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
								return material.Button(th, button2, "Minimize").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button3, "Fullscreen").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button4, "700x400").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button5, "TopLeft 700x300").Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return in.Layout(gtx, func(gtx C) D {
								return material.Button(th, button6, "Center").Layout(gtx)
							})
						}),
					)
				},
				),
			)
			for button1.Clicked() {
				w.Maximize()
				//w.Option(app.Maximized.Option())
			}
			for button2.Clicked() {
				w.Minimize()
				//w.Option(app.Minimized.Option())
			}
			for button3.Clicked() {
				w.Fullscreen()
				//w.Option(app.Fullscreen.Option())
			}
			for button4.Clicked() {
				//w.Option(app.Size(unit.Dp(700), unit.Dp(200)), app.Pos(unit.Dp(200), unit.Dp(200)))
				w.Size(unit.Dp(700), unit.Dp(400))
				w.Pos(unit.Dp(200), unit.Dp(200))
			}
			for button5.Clicked() {
				//w.Option(app.Size(unit.Dp(700), unit.Dp(300)), app.Pos(unit.Dp(0), unit.Dp(0)))
				w.Size(unit.Dp(700), unit.Dp(300))
				w.Pos(unit.Dp(0), unit.Dp(0))
			}
			for button6.Clicked() {
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
