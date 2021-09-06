package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var NoSubmit func() = nil

func BootDialog(win fyne.Window, dims [2]int, bus chan<- bool, serverOnline bool) {

	form := &widget.Form{
		OnSubmit: func() {
			go func() {
				bus <- serverOnline
			}()
		},
	}

	fmt.Println("BootDialog serverOnline ", serverOnline)

	if serverOnline {

		form.Append("", widget.NewLabel("Server onlne proceed to metrics screen"))

	} else {
		form.Append("", widget.NewLabel("Server offline, please check remote connection"))

	}

	layOutContent := container.New(layout.NewPaddedLayout(), form)
	win.SetContent(layOutContent)
	win.Resize(fyne.NewSize(float32(dims[0]/2), float32(dims[1]/3)))

}
