package cpustats

import (
	"fmt"
	"testing"
	"time"
	// "bufio"
	// "os"
)

func TestGetCPU(t *testing.T) {
	cpuStats := GetCpuStats()
	fmt.Printf("%+v\n", cpuStats.User)
	fmt.Printf("%+v\n", cpuStats.System)
	fmt.Printf("%+v\n", cpuStats.Idle)
	fmt.Printf("%+v\n", cpuStats.Total)
}

func TestGetCPUPercent(t *testing.T) {
	done := make(chan struct{})

	go func() {
		for {
			dataPoint := GetCpuPercent()
			fmt.Printf("cpu userspace: %+v%%\n", dataPoint.UserSpace)
			fmt.Printf("cpu system: %+v%%\n", dataPoint.System)
			fmt.Printf("cpu idle: %+v%%\n", dataPoint.Idle)
			time.Sleep(2 * time.Second)
		}
	}()
	<-done
}
