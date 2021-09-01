package cpustats

import (
	"fmt"
)

type (
	Stats struct {
		User, Nice, System, Idle, Iowait, Irq, Softirq, Steal, Guest, GuestNice, Total uint64
		CPUCount, StatCount                                                            int
	}

	cpuStat struct {
		name string
		ptr  *uint64
	}

	CpuRangeData struct {
		RType     string `json:"type"`
		UserSpace string `json:"user"`
		System    string `json:"sys"`
		Idle      string `json:"idle"`
	}
)

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
