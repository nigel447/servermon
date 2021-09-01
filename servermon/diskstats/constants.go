package diskstats

import (
	"fmt"
)

/**
The /proc/diskstats file displays the I/O statistics
of block devices. Each line contains the following 14
fields:
1 - major number
2 - minor mumber
3 - device name
4 - reads completed successfully
5 - reads merged
6 - sectors read
7 - time spent reading (ms)
8 - writes completed
9 - writes merged
10 - sectors written
11 - time spent writing (ms)
12 - I/Os currently in progress
13 - time spent doing I/Os (ms)
14 - weighted time spent doing I/Os (ms)
*/

type (
	IoStats struct {
		Name            string // device name; like "hda"
		ReadsCompleted  uint64 // total number of reads completed successfully
		WritesCompleted uint64 // total number of writes completed successfully
	}

	DiskFreeSpace struct {
		Total      string `json:"total"`
		Used       string `json:"used"`
		Free       string `json:"free"`
		FsType     string `json:"fs"`
		Percent    string `json:"pct"`
		Mountpoint string `json:"mt"`
	}

	DiskRet struct {
		RType string        `json:"type"`
		Data  DiskFreeSpace `json:"data"`
	}
)

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
