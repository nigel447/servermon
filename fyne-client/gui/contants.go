package gui

// https://github.com/mitchellh/mapstructure
import (
	"image/color"

	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/size"
)

type (
	CpuRangeData struct {
		UserSpace string `json:"userMem"`
		System    string `json:"sysMem"`
		Idle      string `json:"idleMem"`
	}

	DiskFreeSpace struct {
		Total      string `json:"totalDisk"`
		Used       string `json:"usedDisk"`
		Free       string `json:"freeDisk"`
		FsType     string `json:"typeDisk"`
		Percent    string `json:"percentDisk"`
		Mountpoint string `json:"mountDisk"`
	}

	MemRangeData struct {
		Total  string `json:"totalMem"`
		Used   string `json:"usedMem"`
		Cached string `json:"cachedMem"`
		Free   string `json:"freeMem"`
	}

	ProfileData struct {
		CpuPoint  CpuRangeData    `json:"cpu"`
		MemPoint  MemRangeData    `json:"mem"`
		DiskPoint []DiskFreeSpace `json:"disk"`
	}
	EventHandler func()
	ImageButton  struct {
		widget.Card
		Handler EventHandler
	}
)

var (
	Boot      = make(chan [2]int)
	StartStop = make(chan string)
	// a lot of waiting data can come true so large buffer
	DataPipe = make(chan map[string]interface{}, 20)
	BootPipe = make(chan map[string]interface{}, 1)
	CpuPipe  = make(chan float64, 10)
	MemPipe  = make(chan float64, 10)
	DiskPipe = make(chan float64, 10)

	HardRedType   = &color.NRGBA{R: 0xfa, A: 0xff}
	SoftGreenType = &color.NRGBA{G: 0xfa, A: 0x99}
	GoldType      = &color.NRGBA{0xff, 0xc1, 0x07, 0xff}

	// purple = &color.NRGBA{R: 128, G: 0, B: 128, A: 255}
	// orange = &color.NRGBA{R: 198, G: 123, B: 0, A: 255}
	// grey   = &color.Gray{Y: 123}

	HeaderFontSize = float32(24)
	TextFontSize   = float32(18)
)

func NewImageButton(res fyne.Resource, f EventHandler) *ImageButton {
	card := &ImageButton{Handler: f}
	card.ExtendBaseWidget(card)
	card.SetImage(canvas.NewImageFromResource(res))

	return card
}

func (t *ImageButton) Tapped(_ *fyne.PointEvent) {
	t.Handler()
}

func (t *ImageButton) TappedSecondary(_ *fyne.PointEvent) {
}

func ScreenDims() {
	driver.Main(func(s screen.Screen) {
		wS, _ := s.NewWindow(&screen.NewWindowOptions{
			Title: "GoLang Server system metrics",
		})

		for {
			e := wS.NextEvent()
			switch e := e.(type) {
			case size.Event:
				dims := e.Size()
				wS.Release()
				var a [2]int
				a[0] = dims.X
				a[1] = dims.Y
				Boot <- a
			}

		}
	})
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
