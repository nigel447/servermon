package memstats

import (
	"fmt"
)

type (
	Stats struct {
		Total, Used, Buffers, Cached, Free, Available, Active, Inactive,
		SwapTotal, SwapUsed, SwapCached, SwapFree, Mapped, Shmem, Slab,
		PageTables, Committed, VmallocUsed uint64
		MemAvailableEnabled bool
	}

	MemRangeData struct {
		Total  string `json:"totalMem"`
		Used   string `json:"usedMem"`
		Cached string `json:"cachedMem"`
		Free   string `json:"freeMem"`
	}
)

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
