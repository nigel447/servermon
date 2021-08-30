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
		UserSpace string `json:"userMem"`
		System    string `json:"sysMem"`
		Idle      string `json:"idleMem"`
	}
)

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
