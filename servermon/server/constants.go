package server

import (
	"encoding/json"
	"time"
	//"fmt"
	"net"
	"srvmon/cpustats"
	"srvmon/diskstats"
	"srvmon/memstats"

	"golang.org/x/crypto/ssh"
	//"sync"
	"sync/atomic"

	"github.com/go-co-op/gocron"
)

type (
	Server struct {
		Addr           string
		ClienPublicKey ssh.PublicKey
		config         *ssh.ServerConfig
	}
)

var (
	lc net.ListenConfig
	//formatter = "%-14s %7s %7s %7s %4s %s\n"
	sheduer             *gocron.Scheduler
	sheduerCount        uint64
	mainServerLoopCount uint64
	//wg            sync.WaitGroup
	ProfileDataCh = make(chan []byte, 1)

	task = func() {
		CpuPoint := *cpustats.GetCpuPercent()
		MemPoint := *memstats.GetMemGb()
		DiskPoint := diskstats.GetDiskSpace()
		CpuPoint.RType = "cpu"
		jBytes, err := json.Marshal(CpuPoint)
		handleError(err)
		ProfileDataCh <- jBytes
		time.Sleep(200 * time.Millisecond)
		MemPoint.RType = "mem"
		jBytes, err = json.Marshal(MemPoint)
		handleError(err)
		ProfileDataCh <- jBytes
		time.Sleep(200 * time.Millisecond)
		for _, dsk := range DiskPoint {
			d := diskstats.DiskRet{
				RType: "disk",
				Data:  dsk,
			}
			jBytes, err = json.Marshal(d)
			handleError(err)
			ProfileDataCh <- jBytes
			time.Sleep(200 * time.Millisecond)
		}
		atomic.AddUint64(&sheduerCount, 1)
	}
)
