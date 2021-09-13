package server

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRunScheduler(t *testing.T) {

	done := make(chan struct{})
	RunProfiler()
	<-done
}

func TestServerData(t *testing.T) {
	data := GetServerSysData()
	fmt.Println(data)
	respBytes, err := json.Marshal(data)
	handleError(err)
	fmt.Println(string(respBytes))

}
