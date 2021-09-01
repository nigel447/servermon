package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func LoginForm(win fyne.Window, dims [2]int, bus chan<- [2]string) {
	username := widget.NewEntry()
	psswd := widget.NewPasswordEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{{Text: "User Name", Widget: username}, {Text: "Password  ", Widget: psswd}},
		OnSubmit: func() {
			go func() {
				var a [2]string
				a[0] = username.Text
				a[1] = psswd.Text
				bus <- a
			}()
		},
	}

	layOutContent := container.New(layout.NewPaddedLayout(), form)
	win.SetContent(layOutContent)
	win.Resize(fyne.NewSize(float32(dims[0]/2), float32(dims[1]/3)))

}
