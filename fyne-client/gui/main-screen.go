package gui

import (
	"fmt"

	"strconv"

	"fyne-client/icon"
	"fyne-client/meter"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//#region  vars
var (
	Slot *fyne.Container
	// Start, Stop                   *widget.Button
	winFU                         fyne.Window
	screenSize                    fyne.Size
	valueCPU, valueMem, valueDisk binding.String
	cpuText, memText, diskText    *widget.Label
)

//#endregion

//#region updateSlot
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

//#endregion

//#region MainScreen
func MainScreen(win fyne.Window, screenDims [2]int) {
	winFU = win
	screenSize = fyne.NewSize(float32(screenDims[0]), float32(screenDims[1]))
	initBootSlot()
	header := container.New(layout.NewPaddedLayout(), createHeaderButtons())
	mainScreen := container.New(layout.NewPaddedLayout(), Slot)
	content := container.New(layout.NewVBoxLayout(), header, mainScreen)
	winFU.SetContent(content)
}

//#endregion

//#region toggleStartStopImportance
func toggleStartStopImportance(bType string) {
	debug := &widget.Label{Alignment: fyne.TextAlignTrailing}
	debug.TextStyle.Monospace = true
	switch bType {
	case "start":
		// Stop.Importance = widget.LowImportance
		// Start.Importance = widget.HighImportance
		go func() {
			if len(DataPipe) > 1 {
				for maps := range DataPipe {
					fmt.Println("drain DataPipe channel dry", maps)
				}
			}
		}()

		StartStop <- "start"

		updateSlot()
		SetSlot()
		go func() {

			for {
				// should block until data
				inputData := <-DataPipe
				MetricsDisplay(inputData)
			}

		}()
	case "stop":
		debug.Text = "update slot for stop"
		SetSlot()
		// Start.Importance = widget.LowImportance
		// Stop.Importance = widget.HighImportance
		StartStop <- "stop"
	}
	// Start.Refresh()
	// Stop.Refresh()
}

//#endregion

///#region MetricsDisplay
func MetricsDisplay(data map[string]interface{}) {
	dataType := data["type"]
	// fmt.Println("begin switch on type ", dataType)
	switch dataType {
	case "cpu":
		CpuPipe <- agregateCpuValue(data)
		valueCPU.Set(fmt.Sprintf("idle %s", data["idle"]))
	case "mem":
		MemPipe <- agregateUsedMemValue(data)
		valueMem.Set(fmt.Sprintf("total %s, free %s", data["totalMem"], data["freeMem"]))
	case "disk":
		DiskPipe <- agregateDiskValue(data["data"].(map[string]interface{}))
		sData := data["data"].(map[string]interface{}) // data here is raw json
		valueDisk.Set(fmt.Sprintf("total %s, used %s, free %s, \nfs %s, mt %s",
			sData["total"], sData["used"], sData["free"], sData["fs"], sData["mt"]))
	}
}

//#endregion

//#region agregateUsedXXXValue
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

	return ((cachedMem + usedMem) / totalMem) * 100

}

// this will be a used cpu percentage value
func agregateCpuValue(data map[string]interface{}) float64 {
	user, err := strconv.ParseFloat(data["user"].(string), 64)
	handleError(err)
	sys, err := strconv.ParseFloat(data["sys"].(string), 64)
	handleError(err)
	return user + sys

}

func agregateDiskValue(sData map[string]interface{}) float64 {
	val := sData["pct"].(string)
	ret, err := strconv.ParseFloat(val[:1], 64)
	handleError(err)
	return ret

}

//#endregion

//#region SetSlot
func SetSlot() {
	header := createHeaderButtons()
	content := container.New(layout.NewVBoxLayout(),
		header, layout.NewSpacer(), Slot, layout.NewSpacer())
	winFU.SetContent(content)
}

func initBootSlot() {

	StartStop <- "onBoot"

	go func() {
		sysData := <-BootPipe
		headerText := canvas.NewText("  Remote Server", HardRedType)

		headerText.TextSize = HeaderFontSize
		// keys
		kernelTextL := canvas.NewText("  Kernel", SoftGreenType)
		kernelTextL.TextSize = TextFontSize
		versionTextL := canvas.NewText("  Version", SoftGreenType)
		versionTextL.TextSize = TextFontSize
		archTextL := canvas.NewText("  Architecture", SoftGreenType)
		archTextL.TextSize = TextFontSize
		osTextL := canvas.NewText("  OS", SoftGreenType)
		osTextL.TextSize = TextFontSize
		// vals
		kernelText := canvas.NewText(sysData["kernel"].(string), GoldType)
		kernelText.TextSize = TextFontSize
		versionText := canvas.NewText(sysData["version"].(string), GoldType)
		versionText.TextSize = TextFontSize
		archText := canvas.NewText(sysData["arch"].(string), GoldType)
		archText.TextSize = TextFontSize
		osText := canvas.NewText(sysData["os"].(string), GoldType)
		osText.TextSize = TextFontSize

		topLevelLayOutContent := container.New(layout.NewVBoxLayout(),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(1), layout.NewSpacer())),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(3), layout.NewSpacer(), headerText, layout.NewSpacer())),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(1), layout.NewSpacer())),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(2), kernelTextL, kernelText)),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(2), versionTextL, versionText)),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(2), archTextL, archText)),
			container.New(layout.NewMaxLayout(),
				container.New(layout.NewGridLayoutWithColumns(2), osTextL, osText)),
		)

		Slot = container.New(layout.NewPaddedLayout(), layout.NewSpacer(), topLevelLayOutContent)
		SetSlot()
	}()
	Slot = container.New(layout.NewVBoxLayout(), layout.NewSpacer())
}

//#endregion

//#region  createHeaderButtons
func createHeaderButtons() *fyne.Container {

	startButton := NewImageButton(icon.Starticon, func() {
		toggleStartStopImportance("start")
	})
	stopButton := NewImageButton(icon.Stopicon, func() {
		toggleStartStopImportance("stop")
	})

	ret := container.New(layout.NewGridLayoutWithColumns(4),
		startButton, layout.NewSpacer(), layout.NewSpacer(), stopButton)
	ret.Resize(fyne.NewSize(float32(600), float32(80)))
	return ret
}

//#endregion
