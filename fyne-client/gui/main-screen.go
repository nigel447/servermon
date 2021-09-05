package gui

import (
	"encoding/json"
	"fmt"
	"fyne-client/meter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var Slot *fyne.Container

var Start, Stop *widget.Button

var winFU fyne.Window
var screenSize fyne.Size

var valueCPU, valueMem, valueDisk binding.String

var cpuText, memText, diskText *widget.Label

func initSlot() {
	Slot = container.New(layout.NewPaddedLayout(), layout.NewSpacer())
}

func updateSlot() {
	valueCPU = binding.NewString()
	valueMem = binding.NewString()
	valueDisk = binding.NewString()
	cpuText = widget.NewLabelWithData(valueCPU)
	memText = widget.NewLabelWithData(valueMem)
	diskText = widget.NewLabelWithData(valueDisk)
	meterSize := fyne.NewSize(float32(screenSize.Width/4), float32(screenSize.Width/4))
	col1 := container.New(layout.NewVBoxLayout(), meter.Show(meterSize, "Cpu", CpuPipe), cpuText)
	col2 := container.New(layout.NewVBoxLayout(), meter.Show(meterSize, "Memory", MemPipe), memText)
	col3 := container.New(layout.NewVBoxLayout(), meter.Show(meterSize, "Disk", DiskPipe), diskText)
	Slot = container.New(layout.NewGridLayoutWithColumns(3), col1, col2, col3)
}

func MainScreen(win fyne.Window, screenDims [2]int) {
	winFU = win
	screenSize = fyne.NewSize(float32(screenDims[0]), float32(screenDims[1]))
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
		updateSlot()
		SetSlot()
		go func() {
			//var pdata ProfileData
			var pdata map[string]interface{}
			for {
				// should block until data
				inputData := <-DataPipe
				err := json.Unmarshal(inputData, &pdata)
				handleError(err)
				MetricsDisplay(pdata)
			}

		}()

	case "stop":
		debug.Text = "update slot for stop"
		initSlot()
		SetSlot()
		Start.Importance = widget.LowImportance
		Stop.Importance = widget.HighImportance
		StartStop <- "stop"
	}

	Start.Refresh()
	Stop.Refresh()

}

func MetricsDisplay(data map[string]interface{}) {
	dataType := data["type"]

	fmt.Println("begin switch on type ", dataType)

	switch dataType {
	case "cpu":
		CpuPipe <- agregateCpuValue(data)
		valueCPU.Set(fmt.Sprintf("user %s, sys %s, \nidle %s", data["user"], data["sys"], data["idle"]))
	case "mem":
		MemPipe <- agregateUsedMemValue(data)
		valueMem.Set(fmt.Sprintf("total %s, used %s, \ncached %s, free %s",
			data["totalMem"], data["usedMem"], data["cachedMem"], data["freeMem"]))
	case "disk":
		DiskPipe <- agregateDiskValue(data["data"].(map[string]interface{}))
		sData := data["data"].(map[string]interface{}) // data here is raw json
		valueDisk.Set(fmt.Sprintf("total %s, used %s, free %s, \nfs %s, pct %s, mt %s",
			sData["total"], sData["used"], sData["free"], sData["fs"], sData["pct"], sData["mt"]))
	}

}

func agregateUsedMemValue(data map[string]interface{}) float64 {
	tMem := data["totalMem"].(string)
	totalMemStr := tMem[:len(tMem)-2]
	totalMem, err := strconv.ParseFloat(totalMemStr, 64)
	handleError(err)
	uMem := data["usedMem"].(string)
	usedMemStr := uMem[:len(uMem)-2]
	usedMem, err := strconv.ParseFloat(usedMemStr, 64)
	handleError(err)
	cMem := data["cachedMem"].(string)
	cachedMemStr := cMem[:len(cMem)-2]
	cachedMem, err := strconv.ParseFloat(cachedMemStr, 64)
	handleError(err)
	// freeMem, err := strconv.ParseFloat(data["freeMem"].(string), 64)
	// handleError(err)

	return ((cachedMem + usedMem) / totalMem) * 100

}

// this will be a used cpu percentage value
func agregateCpuValue(data map[string]interface{}) float64 {
	user, err := strconv.ParseFloat(data["user"].(string), 64)
	handleError(err)
	sys, err := strconv.ParseFloat(data["sys"].(string), 64)
	handleError(err)
	// idle, err := strconv.ParseFloat(data["idle"].(string), 64)
	// handleError(err)
	return user + sys

}

func agregateDiskValue(sData map[string]interface{}) float64 {
	val := sData["pct"].(string)
	ret, err := strconv.ParseFloat(val[:1], 64)
	handleError(err)
	return ret

}

func SetSlot() {
	header := container.New(layout.NewPaddedLayout(), createHeaderButtons())
	content := container.New(layout.NewVBoxLayout(), header, Slot)
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
