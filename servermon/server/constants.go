package server

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"strings"
	"time"
	//"fmt"

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

	ServerSysData struct {
		Type     string `mapstructure:"type" json:"type"`
		Kernel   string `mapstructure:"kernel" json:"kernel"`
		KVersion string `mapstructure:"version" json:"version"`
		Arch     string `mapstructure:"arch" json:"arch"`
		OS       string `mapstructure:"os" json:"os"`
	}
)

var (
	lc net.ListenConfig
	//formatter = "%-14s %7s %7s %7s %4s %s\n"
	sheduer             *gocron.Scheduler
	sheduerCount        uint64
	mainServerLoopCount uint64
	//wg            sync.WaitGroup
	ProfileDataCh = make(chan []byte, 6)

	task = func() {
		CpuPoint := *cpustats.GetCpuPercent()
		MemPoint := *memstats.GetMemGb()
		DiskPoint := diskstats.GetDiskSpace()
		CpuPoint.RType = "cpu"
		jBytes, err := json.Marshal(CpuPoint)
		handleError(err)
		ProfileDataCh <- jBytes
		time.Sleep(500 * time.Millisecond)
		MemPoint.RType = "mem"
		jBytes, err = json.Marshal(MemPoint)
		handleError(err)
		ProfileDataCh <- jBytes
		time.Sleep(500 * time.Millisecond)
		for _, dsk := range DiskPoint {
			d := diskstats.DiskRet{
				RType: "disk",
				Data:  dsk,
			}
			jBytes, err = json.Marshal(d)
			handleError(err)
			ProfileDataCh <- jBytes
			time.Sleep(500 * time.Millisecond)
		}
		atomic.AddUint64(&sheduerCount, 1)
	}
)

func GetServerSysData() (data ServerSysData) {
	cmd := exec.Command("uname", "-srio")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	handleError(err)
	ret := strings.Split(out.String(), " ")
	data = ServerSysData{
		Type:     "sys",
		Kernel:   ret[0],
		KVersion: ret[1],
		Arch:     ret[2],
		OS:       ret[3],
	}
	return
}
