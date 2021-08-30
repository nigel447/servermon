package diskstats

import (
	"fmt"
	"testing"
	// human "github.com/dustin/go-humanize"
	// "github.com/shirou/gopsutil/disk"
)

func TestGetDisk(t *testing.T) {
	diskStats := GetDiskIo()

	for _, disk := range diskStats {
		fmt.Println("disk."+disk.Name+".reads", disk.ReadsCompleted, "-")
		fmt.Println("disk."+disk.Name+".writes", disk.WritesCompleted, "-")
	}

}

func TestDiskSpace(t *testing.T) {
	formatter := "%-14s %7s %7s %7s %4s %s\n"
	fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")
	dFArr := GetDiskSpace()
	for _, df := range dFArr {
		fmt.Printf(formatter,
			df.FsType,
			df.Total,
			df.Used,
			df.Free,
			df.Percent,
			df.Mountpoint,
		)
	}
}

func TestDiskSpaceExtd(t *testing.T) {

}

// func TestDiskSpaceDev(t *testing.T) {
// 	formatter := "%-14s %7s %7s %7s %4s %s\n"
//     fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

//     parts, _ := disk.Partitions(false)
//     for _, p := range parts {
//         device := p.Mountpoint
//         s, _ := disk.Usage(device)

//         if s.Total == 0 {
//             continue
//         }

//         percent := fmt.Sprintf("%2.f%%", s.UsedPercent)
//         if (s.Fstype != "squashfs" ) &&  (s.Fstype != "msdos" )  {
//             fmt.Printf(formatter,
//                 s.Fstype,
//                 human.Bytes(s.Total),
//                 human.Bytes(s.Used),
//                 human.Bytes(s.Free),
//                 percent,
//                 p.Mountpoint,
//             )
//         }
//     }
// }
