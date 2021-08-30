package memstats

import (
	"fmt"
	"testing"
	"time"
	// "bufio"
	// "os"
)

func TestGetMem(t *testing.T) {

	done := make(chan struct{})

	go func() {
		for {
			dataPoint := GetMemGb()
			fmt.Println("Total: ", dataPoint.Total)
			fmt.Println("Used: ", dataPoint.Used)
			fmt.Println("Cached: ", dataPoint.Cached)
			fmt.Println("Free: ", dataPoint.Free)
			time.Sleep(2 * time.Second)
		}
	}()
	<-done

}
