package meter

import (
	"fmt"
	"math"
	"time"

	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var minSize fyne.Size

func (meter *MeterLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	meter.boxPad = float32(40)
	diameter := fyne.Min(size.Width-meter.boxPad, size.Height-meter.boxPad)
	boxDim := fyne.Min(size.Width, size.Height)
	radius := (diameter / 2)
	dotRadius := radius / 12
	smallDotRadius := dotRadius / 2

	smallStroke := diameter / 100
	largeStroke := diameter / 50

	circleSize := fyne.NewSize(diameter, diameter)
	boxSize := fyne.NewSize(boxDim, boxDim)
	middle := fyne.NewPos(float32(boxDim/2), float32(boxDim/2))
	topleft := fyne.NewPos(middle.X-radius, middle.Y-radius)
	boxTopLeft := fyne.NewPos(middle.X-radius-(meter.boxPad/2), middle.Y-radius-(meter.boxPad/2))

	meter.centerDot.StrokeWidth = smallStroke
	meter.box.StrokeWidth = largeStroke
	meter.indicator.StrokeWidth = smallStroke

	// label positions
	headerTop := fyne.NewPos((meter.boxPad / 3), meter.boxPad/3)
	valueCenterTop := fyne.NewPos(radius-(meter.boxPad/6), meter.boxPad/1.5)
	middleRight := fyne.NewPos(float32(boxDim/2)+radius-meter.boxPad, float32(boxDim/2)-meter.boxPad/2)
	middleBottom := fyne.NewPos(float32(boxDim/2)-meter.boxPad/2, float32(boxDim)-meter.boxPad*1.2)
	middleLeft := fyne.NewPos(meter.boxPad/2.2, float32(boxDim/2)-meter.boxPad/2)
	meter.zeroL.Move(headerTop)
	meter.valueLabel.Move(valueCenterTop)
	meter.perc25.Move(middleRight)
	meter.perc50.Move(middleBottom)
	meter.perc75.Move(middleLeft)

	meter.face.Resize(circleSize)
	meter.face.Move(topleft)
	meter.box.Resize(boxSize)
	meter.box.Move(boxTopLeft)

	meter.rotate(meter.indicator, middle, meter.dataValf64, 0, radius-3)

	meter.centerDot.Resize(fyne.NewSize(smallDotRadius*2, smallDotRadius*2))
	meter.centerDot.Move(fyne.NewPos(middle.X-smallDotRadius, middle.Y-smallDotRadius))
	meter.face.StrokeWidth = smallStroke
}

func (meter *MeterLayout) render() *fyne.Container {
	meter.centerDot = &canvas.Circle{StrokeColor: color.Black, StrokeWidth: 3}
	meter.face = &canvas.Circle{
		StrokeColor: theme.ForegroundColor(),
		StrokeWidth: 1,
		FillColor:   &color.RGBA{G: 0x66, A: 0xff},
	}
	meter.box = &canvas.Rectangle{
		StrokeColor: color.NRGBA{B: 0xfa, A: 0xff},
		FillColor:   &color.RGBA{R: 50, G: 50, B: 50, A: 0x0f},
	}
	// Resize sets a new bottom-right position for the line
	meter.indicator = &canvas.Line{StrokeColor: &color.RGBA{R: 0xfa, A: 0xff}, StrokeWidth: 1}
	meter.zeroL = widget.NewLabel(meter.header)
	meter.perc25 = widget.NewLabel("25%")
	meter.perc50 = widget.NewLabel("50%")
	meter.perc75 = widget.NewLabel("75%")
	meter.value = binding.NewString()
	meter.valueLabel = widget.NewLabelWithData(meter.value)

	container := container.NewWithoutLayout(meter.centerDot, meter.face, meter.indicator,
		meter.box, meter.zeroL, meter.valueLabel, meter.perc25, meter.perc50, meter.perc75)
	container.Layout = meter

	meter.canvas = container
	return container
}

func (meter *MeterLayout) rotate(hand fyne.CanvasObject, middle fyne.Position,
	facePosition float64, offset, length float32) {
	// facePosition== increment value
	rotation := (math.Pi * 2 / 60) * facePosition
	x2 := length * float32(math.Sin(rotation))
	y2 := -length * float32(math.Cos(rotation))

	offX := float32(0)
	offY := float32(0)
	if offset > 0 {
		offX += offset * float32(math.Sin(rotation))
		offY += -offset * float32(math.Cos(rotation))
	}

	hand.Move(fyne.NewPos(middle.X+offX, middle.Y+offY))
	hand.Resize(fyne.NewSize(x2, y2))
}

func (meter *MeterLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return minSize
}

// Show loads a clock example window for the specified app context
func Show(meterSize fyne.Size) fyne.CanvasObject {
	minSize = meterSize
	meter := &MeterLayout{dataValf64: 20, header: "Cpu"}
	meter.dataInput = make(chan float64, 10)
	content := meter.render()
	meter.value.Set("0")

	go func(m *MeterLayout) {
		for {
			rand.Seed(time.Now().UnixNano())
			dVal := 0 + rand.Intn(59)
			meter.dataInput <- float64(dVal)
			time.Sleep(5 * time.Second)
		}

	}(meter)

	go func(chan float64, *fyne.Container) {
		for {
			data := <-meter.dataInput
			meter.value.Set(fmt.Sprintf("%0.1f%%", (data/60)*100))
			meter.dataValf64 = data
			meter.Layout(nil, content.Size())
			canvas.Refresh(meter.canvas)
		}
	}(meter.dataInput, content)

	return content
}
