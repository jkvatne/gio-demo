package fps

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"strconv"
	"time"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type FrameTiming struct {
	Start, End      time.Time
	FrameCount      int
	FramesPerSecond float64
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var headingText = []string{"Start", "End", "Frames", "FPS"}
var timings []FrameTiming
var frameCounter int
var TimingStart = time.Time{}
var grid widget.Grid
var TimingWindow = time.Second * 5

func LayoutFpsTable(th *material.Theme, gtx C) D {
	if TimingStart == (time.Time{}) {
		TimingStart = gtx.Now
	}
	if interval := gtx.Now.Sub(TimingStart); interval >= TimingWindow {
		timings = append(timings, FrameTiming{
			Start:           TimingStart,
			End:             gtx.Now,
			FrameCount:      frameCounter,
			FramesPerSecond: float64(frameCounter) / interval.Seconds(),
		})
		frameCounter = 0
		TimingStart = gtx.Now
	}
	frameCounter++

	// Configure width based on available space and a minimum size.
	widthUnit := float32(max(gtx.Constraints.Max.X/3, gtx.Px(unit.Dp(200))))
	widths := []unit.Value{
		unit.Px(widthUnit),
		unit.Px(widthUnit),
		unit.Px(widthUnit * .5),
		unit.Px(widthUnit * .5),
	}
	border := widget.Border{
		Color: color.NRGBA{A: 255},
		Width: unit.Px(1),
	}

	inset := layout.UniformInset(unit.Dp(2))

	// Configure a label styled to be a heading.
	headingLabel := material.Body1(th, "")
	headingLabel.Font.Weight = text.Bold
	headingLabel.Alignment = text.Middle
	headingLabel.MaxLines = 1

	// Configure a label styled to be a data element.
	dataLabel := material.Body1(th, "")
	dataLabel.Font.Variant = "Mono"
	dataLabel.MaxLines = 1
	dataLabel.Alignment = text.End

	// Measure the height of a heading row.
	orig := gtx.Constraints
	gtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, headingLabel.Layout)
	_ = macro.Stop()
	cellHeight := unit.Px(float32(dims.Size.Y))
	gtx.Constraints = orig

	return material.Table(th, &grid).Layout(gtx, len(timings), cellHeight, widths, func(gtx C, row, col int) D {
		return inset.Layout(gtx, func(gtx C) D {
			if row >= len(timings) {
				dataLabel.Text = "Fatal error"
				return dataLabel.Layout(gtx)
			}
			timing := timings[row]
			switch col {
			case 0:
				dataLabel.Text = timing.Start.Format("15:04:05.000000")
			case 1:
				dataLabel.Text = timing.End.Format("15:04:05.000000")
			case 2:
				dataLabel.Text = strconv.Itoa(timing.FrameCount)
			case 3:
				dataLabel.Text = strconv.FormatFloat(timing.FramesPerSecond, 'f', 2, 64)
			}
			return dataLabel.Layout(gtx)
		})
	}, func(gtx C, col int) D {
		return border.Layout(gtx, func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				headingLabel.Text = headingText[col]
				return headingLabel.Layout(gtx)
			})
		})
	})
}
