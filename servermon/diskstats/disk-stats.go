package diskstats

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	human "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
)

func GetDiskIo() (diskStats []IoStats) {
	// Reference: man 5 proc, Documentation/filesystems/proc.txt in Linux source code
	contents, err := ioutil.ReadFile("/proc/diskstats")
	handleError(err)
	return collectDiskIoStats(contents)
}

func collectDiskIoStats(out []byte) (diskStatsData []IoStats) {
	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)
	var diskStats []IoStats
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		name := fields[2]
		readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			handleError(fmt.Errorf("failed to parse reads completed of %s", name))
		}
		writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			handleError(fmt.Errorf("failed to parse writes completed of %s", name))
		}
		diskStats = append(diskStats, IoStats{
			Name:            name,
			ReadsCompleted:  readsCompleted,
			WritesCompleted: writesCompleted,
		})
	}
	if err := scanner.Err(); err != nil {
		handleError(fmt.Errorf("scan error for /proc/diskstats: %s", err))
	}
	diskStatsData = diskStats
	return
}

func GetDiskSpace() (df []DiskFreeSpace) {

	parts, _ := disk.Partitions(false)
	for _, p := range parts {
		device := p.Mountpoint
		s, _ := disk.Usage(device)
		if s.Total == 0 {
			continue
		}

		if (s.Fstype != "squashfs") && (s.Fstype != "msdos") {
			percent := fmt.Sprintf("%2.f%%", s.UsedPercent)
			df = append(df, DiskFreeSpace{
				FsType:     s.Fstype,
				Total:      human.Bytes(s.Total),
				Used:       human.Bytes(s.Used),
				Free:       human.Bytes(s.Free),
				Percent:    percent,
				Mountpoint: p.Mountpoint,
			})
		}
	}

	return
}
