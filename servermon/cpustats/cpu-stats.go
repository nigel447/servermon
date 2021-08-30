package cpustats

// Portions of this file are derived from (https://en.mackerel.io/).
// https://github.com/mackerelio/go-osstat/blob/master/cpu/cpu_linux.go

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	// "os"
	// "io"
	"bufio"
	"unicode"
)

// Get cpu statistics
func GetCpuStats() (cpuStats *Stats) {
	contents, err := ioutil.ReadFile("/proc/stat")
	handleError(err)
	cpuStats = collectCPUStats(contents)
	return
}

func collectCPUStats(out []byte) (cpuData *Stats) {
	reader := bytes.NewReader(out)
	var cpu Stats
	scanner := bufio.NewScanner(reader)

	cpuStats := []cpuStat{
		{"user", &cpu.User},
		{"nice", &cpu.Nice},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"iowait", &cpu.Iowait},
		{"irq", &cpu.Irq},
		{"softirq", &cpu.Softirq},
		{"steal", &cpu.Steal},
		{"guest", &cpu.Guest},
		{"guest_nice", &cpu.GuestNice},
	}
	// need to scan to get a line
	if !scanner.Scan() {
		handleError(fmt.Errorf("failed to scan /proc/stat"))
	}
	valStrs := strings.Fields(scanner.Text())[1:]
	cpu.StatCount = len(valStrs)
	for i, valStr := range valStrs {
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			handleError(fmt.Errorf("failed to scan %s from /proc/stat", cpuStats[i].name))
		}
		*cpuStats[i].ptr = val
		cpu.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
	cpu.Total -= cpu.Guest
	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	cpu.Total -= cpu.GuestNice

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			cpu.CPUCount++
		}
	}

	if err := scanner.Err(); err != nil {
		handleError(fmt.Errorf("scan error for /proc/stat: %s", err))
	}
	cpuData = &cpu
	return
}

func GetCpuPercent() (dataPoint *CpuRangeData) {
	before := GetCpuStats()
	time.Sleep(time.Duration(1) * time.Second)
	after := GetCpuStats()
	total := float64(after.Total - before.Total)
	dataPoint = &CpuRangeData{
		UserSpace: getFloat64PercentDiff(after.User, before.User, total),
		System:    getFloat64PercentDiff(after.System, before.System, total),
		Idle:      getFloat64PercentDiff(after.Idle, before.Idle, total),
	}
	return
}

func getFloat64PercentDiff(end uint64, start uint64, total float64) string {
	return fmt.Sprintf("%f", float64(end-start)/total*100)
}
