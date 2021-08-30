package memstats

// Portions of this file are derived from (https://en.mackerel.io/).
// https://github.com/mackerelio/go-osstat/blob/master/cpu/cpu_linux.go

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetMem() (memStats *Stats) {
	// Reference: man 5 proc, Documentation/filesystems/proc.txt in Linux source code
	contents, err := ioutil.ReadFile("/proc/meminfo")
	handleError(err)
	return collectMemoryStats(contents)
}

func GetMemGb() (dataPoint *MemRangeData) {
	memory := GetMem()
	dataPoint = &MemRangeData{
		Total:  fmt.Sprintf("%.2fGb", float64(memory.Total/(1048576))),
		Used:   fmt.Sprintf("%.2fGb", float64(memory.Used/(1048576))),
		Cached: fmt.Sprintf("%.2fGb", float64(memory.Cached/(1048576))),
		Free:   fmt.Sprintf("%.2fGb", float64(memory.Free/(1048576))),
	}
	return
}

// https://stackoverflow.com/questions/4938612/how-do-i-print-the-pointer-value-of-a-go-object-what-does-the-pointer-value-mea
func collectMemoryStats(out []byte) (memStatData *Stats) {
	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)
	var memory Stats
	memStats := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemFree":      &memory.Free,
		"MemAvailable": &memory.Available,
		"Buffers":      &memory.Buffers,
		"Cached":       &memory.Cached,
		"Active":       &memory.Active,
		"Inactive":     &memory.Inactive,
		"SwapCached":   &memory.SwapCached,
		"SwapTotal":    &memory.SwapTotal,
		"SwapFree":     &memory.SwapFree,
		"Mapped":       &memory.Mapped,
		"Shmem":        &memory.Shmem,
		"Slab":         &memory.Slab,
		"PageTables":   &memory.PageTables,
		"Committed_AS": &memory.Committed,
		"VmallocUsed":  &memory.VmallocUsed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		fld := line[:i]
		if ptr := memStats[fld]; ptr != nil {

			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v
			}
			// if fld =="MemFree" {
			// 	fmt.Printf("test memfree val %v\n",*ptr)
			// }
			if fld == "MemAvailable" {
				memory.MemAvailableEnabled = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		handleError(fmt.Errorf("scan error for /proc/meminfo: %s", err))
	}

	memory.SwapUsed = memory.SwapTotal - memory.SwapFree

	if memory.MemAvailableEnabled {
		memory.Used = memory.Total - memory.Available
	} else {
		memory.Used = memory.Total - memory.Free - memory.Buffers - memory.Cached
	}

	memStatData = &memory
	return
}
