package gui

import (
	"fmt"
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
	// https://github.com/mitchellh/mapstructure
	RBoxInfo struct {
		Kernel   string `mapstructure:"kernel" json:"kernel"`
		KVersion string `mapstructure:"version" json:"kernel"`
		Arch     string `mapstructure:"arch" json:"kernel"`
		OS       string `mapstructure:"os" json:"kernel"`
	}
)

var (
	Boot      = make(chan [2]int)
	StartStop = make(chan string)
	// a lot of waiting data can come true so large buffer
	DataPipe = make(chan []byte, 20)
	BootPipe = make(chan []byte, 1)
	CpuPipe  = make(chan float64, 10)
	MemPipe  = make(chan float64, 10)
	DiskPipe = make(chan float64, 10)
)

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
