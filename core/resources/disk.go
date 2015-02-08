package resources

import (
	"fmt"
	"github.com/silenteh/gantryos/utils"
	"strings"
)

type diskLayoutInfo struct {
	Device    string  // /dev/sda1
	Size      float64 // in GB
	Used      float64 // in GB
	Available float64 // in GB
	Usage     string  // this is the percentageof usage
	Mounted   string  // like /
	Ts        int32
}

type diskStat struct {
	DeviceName     string
	ReadsN         uint64 // number of reads
	ReadsSectorN   uint64 // number of sectors read
	ReadsMS        uint64 // number of MS spent in reads
	WritesN        uint64 // number of Writes
	WritesSectorN  uint64 // number of sectors written
	WritesMS       uint64 // number of MS spent in writes
	IOOperationsN  uint64 // number of I/O operations in progress
	IOOperationsMS uint64 // number of MS spent in I/O operations
	IOBacklog      uint64 // complicated ! https://www.kernel.org/doc/Documentation/iostats.txt
	Ts             int32
}

func layout() map[string]diskLayoutInfo {

	layout := make(map[string]diskLayoutInfo)

	detectedOs := detectOS()
	switch detectedOs {
	case BSD:
		dfResult := utils.ExecCommand(false, "df", "-h")
		dfResult = strings.Replace(dfResult, "map ", "", -1)
		dfResult = removeUnits(dfResult)

		output := utils.ParseOutputCommandWithHeader(dfResult, 1)
		dataMapArray, err := utils.CommandOutputToMapArray(output, 8)
		if err != nil {
			return layout
		}

		for _, element := range dataMapArray {
			size := len(element)
			if size >= 9 {
				l := diskLayoutInfo{}
				l.Device = element[0]
				l.Size = utils.StringToFloat64(element[1], false)
				l.Used = utils.StringToFloat64(element[2], false)
				l.Available = utils.StringToFloat64(element[3], false)
				l.Usage = element[4]
				l.Mounted = element[8]
				l.Ts = utils.UTCTimeStamp()
				layout[l.Mounted] = l
			}

		}

		break
	case LINUX:
		dfResult := utils.ExecCommand(false, "df", "-h")
		fmt.Println(dfResult)
		output := utils.ParseOutputCommandWithHeader(dfResult, 1)
		dataMapArray, err := utils.CommandOutputToMapArray(output, 5)
		if err != nil {
			return layout
		}

		for _, element := range dataMapArray {
			size := len(element)
			if size >= 5 {
				l := diskLayoutInfo{}
				l.Device = element[0]
				l.Size = utils.StringToFloat64(element[1], false)
				l.Used = utils.StringToFloat64(element[2], false)
				l.Available = utils.StringToFloat64(element[3], false)
				l.Usage = element[4]
				l.Mounted = element[5]
				l.Ts = utils.UTCTimeStamp()
				layout[l.Mounted] = l
			}
		}
		break
	}

	return layout
}

// TODO: BSD not supported yet
func ioinfo() []diskStat {

	var iostats []diskStat
	var iostatsTemp []diskStat

	detectedOs := detectOS()
	switch detectedOs {
	case BSD:

		resultBytes := utils.ReadFile("linux_diskstats.txt")
		dfResult := string(resultBytes)

		output := utils.ParseOutputCommandWithHeader(dfResult, 1)
		dataMapArray, err := utils.CommandOutputToMapArray(output, 8)
		if err != nil {
			return nil
		}

		added := 0
		for _, element := range dataMapArray {

			t := element[2]
			if strings.Contains(t, "sd") || strings.Contains(t, "hd") || strings.Contains(t, "xv") || strings.Contains(t, "md") {
				l := diskStat{}
				l.DeviceName = t
				l.ReadsN = utils.StringToUINT64(element[3], false)
				l.ReadsSectorN = utils.StringToUINT64(element[5], false)
				l.ReadsMS = utils.StringToUINT64(element[6], false)
				l.WritesN = utils.StringToUINT64(element[7], false)
				l.WritesSectorN = utils.StringToUINT64(element[9], false)
				l.WritesMS = utils.StringToUINT64(element[10], false)
				l.IOOperationsN = utils.StringToUINT64(element[11], false)
				l.IOOperationsMS = utils.StringToUINT64(element[12], false)
				l.IOBacklog = utils.StringToUINT64(element[13], false)
				l.Ts = utils.UTCTimeStamp()

				iostatsTemp[added] = l
				added++
			}

		}

		iostats = make([]diskStat, added)
		copy(iostats, iostatsTemp)

		break
	case LINUX:
		dfResult := utils.ExecCommand(false, "cat", "/proc/diskstats")

		output := utils.ParseOutputCommandWithHeader(dfResult, 1)
		dataMapArray, err := utils.CommandOutputToMapArray(output, 8)
		if err != nil {
			return nil
		}
		iostatsTemp = make([]diskStat, len(dataMapArray))
		added := 0
		for _, element := range dataMapArray {

			t := element[2]
			if strings.Contains(t, "sd") || strings.Contains(t, "hd") || strings.Contains(t, "xv") || strings.Contains(t, "md") {
				l := diskStat{}
				l.DeviceName = t
				l.ReadsN = utils.StringToUINT64(element[3], false)
				l.ReadsSectorN = utils.StringToUINT64(element[5], false)
				l.ReadsMS = utils.StringToUINT64(element[6], false)
				l.WritesN = utils.StringToUINT64(element[7], false)
				l.WritesSectorN = utils.StringToUINT64(element[9], false)
				l.WritesMS = utils.StringToUINT64(element[10], false)
				l.IOOperationsN = utils.StringToUINT64(element[11], false)
				l.IOOperationsMS = utils.StringToUINT64(element[12], false)
				l.IOBacklog = utils.StringToUINT64(element[13], false)
				l.Ts = utils.UTCTimeStamp()

				iostatsTemp[added] = l
				added++
			}

		}

		iostats = make([]diskStat, added)
		copy(iostats, iostatsTemp)
		break
	}

	return iostats

}

func removeUnits(content string) string {
	content = strings.ToLower(content)
	content = strings.Replace(content, "gi", "", -1)
	content = strings.Replace(content, "gb", "", -1)
	content = strings.Replace(content, "mi", "", -1)
	content = strings.Replace(content, "mb", "", -1)
	content = strings.Replace(content, "ki", "", -1)
	content = strings.Replace(content, "kb", "", -1)
	content = strings.Replace(content, "bi", "", -1)
	content = strings.Replace(content, "%", "", -1)
	return content
}
