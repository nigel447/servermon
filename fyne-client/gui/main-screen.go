package gui

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var Slot *fyne.Container

var Start, Stop *widget.Button

var winFU fyne.Window

var cpuText = widget.NewMultiLineEntry()
var memText = widget.NewMultiLineEntry()
var diskText = widget.NewMultiLineEntry()

func initSlot() {
	Slot = container.New(layout.NewPaddedLayout(), layout.NewSpacer())
}

func MainScreen(win fyne.Window, screenDims [2]int) {
	winFU = win
	initSlot()
	header := container.New(layout.NewPaddedLayout(), createHeaderButtons())
	mainScreen := container.New(layout.NewPaddedLayout(), Slot)
	content := container.New(layout.NewVBoxLayout(), header, mainScreen)
	winFU.SetContent(content)
}

func toggleStartStopImportance(bType string) {
	debug := &widget.Label{Alignment: fyne.TextAlignTrailing}
	debug.TextStyle.Monospace = true

	switch bType {
	case "start":
		Stop.Importance = widget.LowImportance
		Start.Importance = widget.HighImportance
		StartStop <- "start"
		go func() {
			//var pdata ProfileData
			var pdata map[string]interface{}
			for {
				// should block until data
				inputData := <-DataPipe
				err := json.Unmarshal(inputData, &pdata)
				handleError(err)
				SetSlot(MetricsDisplay(pdata))
			}

		}()

	case "stop":
		debug.Text = "update slot for stop"
		SetSlot(container.NewHBox(debug))
		Start.Importance = widget.LowImportance
		Stop.Importance = widget.HighImportance
		StartStop <- "stop"
	}

	Start.Refresh()
	Stop.Refresh()

}

func MetricsDisplay(data map[string]interface{}) (s *fyne.Container) {
	dataType := data["type"]

	switch dataType {
	case "cpu":
		cpuText.SetText(fmt.Sprintf("CPU: userspace %s, system %s, idle %s", data["user"], data["sys"], data["idle"]))
	case "mem":
		memText.SetText(fmt.Sprintf("Mem: total %s, used %s, cached %s, free %s",
			data["totalMem"], data["usedMem"], data["cachedMem"], data["freeMem"]))
	case "disk":
		sData := data["data"].(map[string]interface{})
		diskText.SetText(fmt.Sprintf("Disk: total %s, used %s, free %s, fs %s,pct %s, mt %s",
			sData["total"], sData["used"], sData["free"], sData["fs"], sData["pct"], sData["mt"]))
	}
	return container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewPaddedLayout(), cpuText),
		container.New(layout.NewPaddedLayout(), memText),
		container.New(layout.NewPaddedLayout(), diskText))

}

func SetSlot(s *fyne.Container) {
	Slot = container.New(layout.NewPaddedLayout(), s)
	header := container.New(layout.NewPaddedLayout(), createHeaderButtons())
	mainScreen := container.New(layout.NewPaddedLayout(), Slot)
	content := container.New(layout.NewVBoxLayout(), header, mainScreen)
	winFU.SetContent(content)
}

func createHeaderButtons() *fyne.Container {
	Start = widget.NewButton("   Start   ", func() {
		fmt.Println("on start")
		toggleStartStopImportance("start")

	})
	Stop = widget.NewButton("   Stop   ", func() {
		fmt.Println("on stop")
		toggleStartStopImportance("stop")
	})

	return container.NewHBox(Start, layout.NewSpacer(), Stop, layout.NewSpacer())

}
