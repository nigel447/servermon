package gui

import (
	"encoding/json"
	"fmt"
	"fyne-client/meter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var Slot *fyne.Container

var Start, Stop *widget.Button

var winFU fyne.Window
var screenSize, panelSize fyne.Size

var cpuText = widget.NewMultiLineEntry()
var memText = widget.NewMultiLineEntry()
var diskText = widget.NewMultiLineEntry()

func initSlot() {
	Slot = container.New(layout.NewPaddedLayout(), layout.NewSpacer())
}

func MainScreen(win fyne.Window, screenDims [2]int) {
	winFU = win
	screenSize = fyne.NewSize(float32(screenDims[0]), float32(screenDims[1]))
	panelSize = fyne.NewSize(float32(screenDims[0]/2), float32(screenDims[1]/2))

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
	// create meter widgets here
	// cpu total = usp +sys+idle then => [0, 60] snd to meter Show
	// mem total then => [0, 60] snd to meter Show
	// dsk ect ...
	// each meter will havetext blck beow summary
	fmt.Println("begin switch on type ", dataType)
	cpuText.Wrapping = fyne.TextWrapWord
	memText.Wrapping = fyne.TextWrapWord
	diskText.Wrapping = fyne.TextWrapWord
	switch dataType {
	case "cpu":
		total := aggregateCpuValue(data)
		fmt.Println("total cpu ", total)
		cpuText.SetText(fmt.Sprintf("CPU: userspace %s, system %s, idle %s", data["user"], data["sys"], data["idle"]))
	case "mem":
		memText.SetText(fmt.Sprintf("Mem: total %s, used %s, cached %s, free %s",
			data["totalMem"], data["usedMem"], data["cachedMem"], data["freeMem"]))
	case "disk":
		sData := data["data"].(map[string]interface{})
		diskText.SetText(fmt.Sprintf("Disk: total %s, used %s, free %s, fs %s,pct %s, mt %s",
			sData["total"], sData["used"], sData["free"], sData["fs"], sData["pct"], sData["mt"]))
	}

	meterSize := fyne.NewSize(float32(screenSize.Width/4), float32(screenSize.Width/4))

	col1 := container.New(layout.NewVBoxLayout(), meter.Show(meterSize), cpuText)
	col2 := container.New(layout.NewVBoxLayout(), meter.Show(meterSize), memText)
	col3 := container.New(layout.NewVBoxLayout(), meter.Show(meterSize), diskText)

	return container.New(layout.NewGridLayoutWithColumns(3), col1, col2, col3)

}

func aggregateCpuValue(data map[string]interface{}) float64 {
	user, err := strconv.ParseFloat(data["user"].(string), 64)
	handleError(err)
	sys, err := strconv.ParseFloat(data["sys"].(string), 64)
	handleError(err)
	idle, err := strconv.ParseFloat(data["idle"].(string), 64)
	handleError(err)
	return user + sys + idle

}

func SetSlot(s *fyne.Container) {
	header := container.New(layout.NewPaddedLayout(), createHeaderButtons())
	content := container.New(layout.NewVBoxLayout(), header, s)
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
