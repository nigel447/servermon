package meter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var (
	minSize fyne.Size
)

type (
	MeterLayout struct {
		indicator              *canvas.Line
		centerDot, face        *canvas.Circle
		box                    *canvas.Rectangle
		canvas                 fyne.CanvasObject
		dataInput              chan float64
		dataValf64             float64 // in degrees 1 to 60
		boxPad                 float32
		perc25, perc50, perc75 *widget.Label
		valueLabel             *widget.Label
		value                  binding.String
		header                 string
		headerText             *canvas.Text
	}
)
