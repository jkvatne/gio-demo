// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://todo.sr.ht/~eliasnaur/gio/285

import (
	"gioui.org/unit"
	"image"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
	"gioui.org/x/outlay"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

var in = layout.UniformInset(unit.Dp(1))

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var list widget.List
	list.Axis = layout.Vertical
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			// hold elements in a single-item list so it's scrollable if it wraps
			// below the window's visible area
			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(material.Body1(th, "Please select a device:").Layout),
				layout.Rigid(func(gtx C) D {
					return material.List(th, &list).Layout(gtx, 1, func(gtx C, _ int) D {
						// Device list
						hWrap := outlay.GridWrap{
							Axis:      layout.Horizontal,
							Alignment: layout.End,
						}
						return hWrap.Layout(gtx, 100, func(gtx C, i int) D {
							material.Body1(th, " "+strconv.Itoa(i)+" ").Layout(gtx)
							return D{
								Size: image.Pt(200, 200),
							}
						})
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
