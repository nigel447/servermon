package gui

import (
	//"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var NoSubmit func() = nil

func BootDialog(win fyne.Window, dims [2]int, bus chan<- bool, serverOnline bool) {
	var layOutContent *fyne.Container
	okBttn := &widget.Button{
		Text: "  OK  ",
		OnTapped: func() {
			bus <- serverOnline
		},
	}

	okBttn.Importance = widget.HighImportance

	if serverOnline {
		mssg := canvas.NewText("Server onlne proceed to metrics screen", color.White)
		mssg.TextSize = 24
		layOutContent = container.New(layout.NewVBoxLayout(),
			container.New(layout.NewCenterLayout(), mssg), layout.NewSpacer(),
			container.New(layout.NewCenterLayout(), okBttn), layout.NewSpacer())
	} else {
		mssg := canvas.NewText("Server offline, please check remote connection", color.White)
		mssg.TextSize = 24
		layOutContent = container.New(layout.NewVBoxLayout(), mssg, layout.NewSpacer(),
			container.New(layout.NewCenterLayout(), okBttn), layout.NewSpacer())
	}

	win.SetContent(layOutContent)
	win.Resize(fyne.NewSize(float32(dims[0]/2), float32(dims[1]/3)))
}
