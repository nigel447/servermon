package server

import (
	"encoding/json"
	//"fmt"
	"net"
	"srvmon/cpustats"
	"srvmon/diskstats"
	"srvmon/memstats"

	"golang.org/x/crypto/ssh"
	//"sync"
	"github.com/go-co-op/gocron"
	"sync/atomic"
)

type (
	Server struct {
		Addr           string
		ClienPublicKey ssh.PublicKey
		config         *ssh.ServerConfig
	}

	ProfileData struct {
		CpuPoint  cpustats.CpuRangeData     `json:"cpu"`
		MemPoint  memstats.MemRangeData     `json:"mem"`
		DiskPoint []diskstats.DiskFreeSpace `json:"disk"`
	}
)

var (
	lc net.ListenConfig
	//formatter = "%-14s %7s %7s %7s %4s %s\n"
	sheduer             *gocron.Scheduler
	sheduerCount        uint64
	mainServerLoopCount uint64
	//wg            sync.WaitGroup
	ProfileDataCh = make(chan string, 1)

	task = func() {
		data := ProfileData{}
		data.CpuPoint = *cpustats.GetCpuPercent()
		data.MemPoint = *memstats.GetMemGb()
		data.DiskPoint = diskstats.GetDiskSpace()
		jBytes, err := json.Marshal(data)
		handleError(err)
		ProfileDataCh <- string(jBytes)
		atomic.AddUint64(&sheduerCount, 1)
	}
)
