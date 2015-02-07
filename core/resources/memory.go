package resources

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/silenteh/gantryos/utils"
	//"strings"
)

type memoryUsage struct {
	Total         float64
	Available     float64 // amount of RAM available
	App           float64 // amount of RAM used by Apps
	OS            float64 // amount of RAM used by ther OS ! cannot be borrowed
	Used          float64
	TotalSwap     float64
	AvailableSwap float64
	UsedSwap      float64
	Ts            int32
}

func totalRam() float64 {
	detectedOs := detectOS()
	availableRam := float64(0)
	switch detectedOs {
	case BSD:
		availableRamString := utils.ExecCommand(true, "/usr/sbin/sysctl", "-n", "hw.memsize")
		availableRam = convertStringBytesToGB(availableRamString)
		break
	case LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/meminfo") // linux
		fmt.Println(output)
		data := utils.ParseOutputCommand(output)

		m, err := utils.CommandOutputToMap(data, 0, 1)
		if err != nil {
			log.Error(err)
			break
		}
		availableRam = convertKiloBytesToGB(m["MemTotal"])
		break
	default:
		fmt.Println("OS not recognized:", detectOS())
	}
	return availableRam
}

func usedRam() memoryUsage {

	memUsage := memoryUsage{}
	memUsage.Total = totalRam()

	detectedOs := detectOS()
	switch detectedOs {
	case BSD:
		output := utils.ExecCommand(false, "vm_stat")
		memUsage = realTimeMemoryFromCommandOutput(output, detectedOs)
		memUsage.Ts = utils.UTCTimeStamp()
		break
	case LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/meminfo")
		memUsage = realTimeMemoryFromCommandOutput(output, detectedOs)
		memUsage.Ts = utils.UTCTimeStamp()
		break
	}
	return memUsage
}

func realTimeMemoryFromCommandOutput(output string, detectedOs string) memoryUsage {
	memUsage := memoryUsage{}

	switch detectedOs {
	case BSD:
		pageSize := pageSize()

		data := utils.ParseOutputCommand(output)

		m, err := utils.CommandOutputToMap(data, 0, 1)
		if err != nil {
			log.Error(err)
			break
		}

		pagesOccupiedByCompressor := (float64(m["Pages occupied by compressor"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//pagesStoredInCompressor := (float64(m["Pages stored in compressor"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		fileCache := (float64(m["File-backed pages"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesFree := (float64(m["Pages free"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//pagesActive := (float64(m["Pages active"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesInactive := (float64(m["Pages inactive"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesWired := (float64(m["Pages wired down"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesSpeculative := (float64(m["Pages speculative"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//anonymousPages := (float64(m["Anonymous pages"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		//pagesPurgeable := (float64(m["Pages purgeable"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))
		pagesReactivated := (float64(m["Pages reactivated"]*pageSize) / float64(1024.00) / float64(1024.00) / float64(1024.00))

		//fmt.Printf("File Cache: %s GB\n", strconv.FormatFloat(fileCache, 'f', 2, 64))
		//fmt.Printf("Wired Memory: %s GB\n", strconv.FormatFloat(pagesWired, 'f', 2, 64))
		//fmt.Printf("Compressed: %s GB\n", strconv.FormatFloat(pagesOccupiedByCompressor, 'f', 3, 64))

		memUsage.Available = pagesFree
		memUsage.App = pagesInactive + pagesSpeculative + pagesReactivated //- pagesInactive - pagesPurgeable
		memUsage.OS = pagesWired
		//memUsage.Used = pagesActive + pagesInactive + pagesSpeculative + pagesWired + pagesOccupiedByCompressor
		memUsage.Used = memUsage.App + fileCache + pagesWired + pagesOccupiedByCompressor

		break
	case LINUX:
		//m := utils.ParseOutputLabelValueFormat(output, false, ":")

		data := utils.ParseOutputCommand(output)

		m, err := utils.CommandOutputToMap(data, 0, 1)
		if err != nil {
			log.Error(err)
			break
		}

		totalMem := convertKiloBytesToGB(m["MemTotal"])
		freeMem := convertKiloBytesToGB(m["MemFree"])

		buffersMem := convertKiloBytesToGB(m["Buffers"])
		cachedMem := convertKiloBytesToGB(m["Cached"])

		totalSwap := convertKiloBytesToGB(m["SwapTotal"])
		freeSwap := convertKiloBytesToGB(m["SwapFree"])

		memUsage.OS = buffersMem + cachedMem

		memUsage.Available = freeMem + memUsage.OS
		memUsage.Used = totalMem - freeMem

		memUsage.App = memUsage.Used - memUsage.OS

		memUsage.TotalSwap = totalSwap
		memUsage.AvailableSwap = freeSwap
		memUsage.UsedSwap = totalSwap - freeSwap

		break
	}

	return memUsage
}

func convertBytesToGB(value uint64) float64 {
	return (float64(value) / float64(1024.00) / float64(1024.00) / float64(1024.00))
}

func convertStringBytesToGB(value string) float64 {
	i := utils.StringToFloat64(value, false)
	floatValue := (i / float64(1024.00) / float64(1024.00) / float64(1024.00))
	return floatValue
}

func convertKiloBytesToGB(value uint64) float64 {
	return (float64(value) / float64(1024.00) / float64(1024.00))
}
