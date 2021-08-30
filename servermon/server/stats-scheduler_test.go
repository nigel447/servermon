package server

import (
	"fmt"
	"testing"
)

func TestRunScheduler(t *testing.T) {

	done := make(chan struct{})
	RunProfiler()
	<-done
}

func TestChannel(t *testing.T) {
	done := make(chan bool)
	dataChannel := make(chan ProfileData)

	go func() {
		for {
			select {
			case data := <-dataChannel:
				fmt.Println("received data", data)

			}
		}
	}()

	data := ProfileData{}
	dataChannel <- data
	dataChannel <- data

	<-done
}
