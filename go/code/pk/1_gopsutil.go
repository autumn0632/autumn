package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func main() {
	totalCpuPercent, _ := cpu.Percent(3*time.Second, false)
	mem, _ := mem.VirtualMemory()

	fmt.Printf("total cpu percent:%.4f, total mem percent: %.4f",
				totalCpuPercent,
				mem.UsedPercent)

	var total, used uint64
	infos, _ := disk.Partitions(false)
	for _, info := range infos {
		totalDisk, _ := disk.Usage(info.Mountpoint)
		usedDisk, _ := disk.Usage(info.Mountpoint)

		total += totalDisk.Total
		used += usedDisk.Used
	}

	fmt.Printf("total disk percent:%0.4f:", float64(used) / float64(total))

}
