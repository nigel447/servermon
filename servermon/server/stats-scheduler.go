package server

import (
	// "fmt"
	"github.com/go-co-op/gocron"
	"time"
)

func RunProfiler() {
	s := gocron.NewScheduler(time.UTC)
	// job, _ := s.Every(5).Second().Do(task)
	// go func() {
	// 	for {
	// 		fmt.Println("Run count", job.RunCount())
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	s.Every(5).Second().Do(task)
	s.StartAsync()
}
