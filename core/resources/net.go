package resources

import (
	"github.com/silenteh/gantryos/utils"
	"strings"
)

type netStat struct {
	InterfaceName string
	RXBytes       uint64 //
	RXPackets     uint64 //
	RXErrors      uint64 //
	RXDrops       uint64 //
	RXFifo        uint64 //
	RXFrame       uint64 //
	RXCompressed  uint64 //
	RXMulticast   uint64 //
	TXBytes       uint64 //
	TXPackets     uint64 //
	TXErrors      uint64 //
	TXDrops       uint64 //
	TXFifo        uint64 //
	TXCollisions  uint64 //
	TXCarrier     uint64 //
	TXCompressed  uint64 //
	Ts            int32
}

func netStats() map[string]netStat {

	stats := make(map[string]netStat)

	detectedOs := detectOS()

	switch detectedOs {
	case BSD:

		resultBytes := utils.ReadFile("netinfo.txt")
		// this is for the tests
		// TODO: fix this hack !
		if len(resultBytes) == 0 {
			resultBytes = utils.ReadFile("../../netinfo.txt")
		}
		dfResult := string(resultBytes)

		itemsArray := utils.ParseOutputCommandWithHeader(dfResult, 2)

		for _, element := range itemsArray {

			size := len(element)
			if size >= 17 {
				l := netStat{}
				l.InterfaceName = strings.Replace(element[0], ":", "", -1)
				l.RXBytes = utils.StringToUINT64(element[1], false)
				l.RXPackets = utils.StringToUINT64(element[2], false)
				l.RXErrors = utils.StringToUINT64(element[3], false)
				l.RXDrops = utils.StringToUINT64(element[4], false)
				l.RXFifo = utils.StringToUINT64(element[5], false)
				l.RXFrame = utils.StringToUINT64(element[6], false)
				l.RXCompressed = utils.StringToUINT64(element[7], false)
				l.RXMulticast = utils.StringToUINT64(element[8], false)

				l.TXBytes = utils.StringToUINT64(element[9], false)
				l.TXPackets = utils.StringToUINT64(element[10], false)
				l.TXErrors = utils.StringToUINT64(element[11], false)
				l.TXDrops = utils.StringToUINT64(element[12], false)
				l.TXFifo = utils.StringToUINT64(element[13], false)
				l.TXCollisions = utils.StringToUINT64(element[14], false)
				l.TXCarrier = utils.StringToUINT64(element[15], false)
				l.TXCompressed = utils.StringToUINT64(element[16], false)
				l.Ts = utils.UTCTimeStamp()
				stats[l.InterfaceName] = l
			}
		}

		break
	case LINUX:
		output := utils.ExecCommand(false, "cat", "/proc/net/dev")

		//fmt.Println(output)

		itemsArray := utils.ParseOutputCommandWithHeader(output, 2)

		for _, element := range itemsArray {

			size := len(element)
			if size >= 17 {
				l := netStat{}
				l.InterfaceName = strings.Replace(element[0], ":", "", -1)
				l.RXBytes = utils.StringToUINT64(element[1], false)
				l.RXPackets = utils.StringToUINT64(element[2], false)
				l.RXErrors = utils.StringToUINT64(element[3], false)
				l.RXDrops = utils.StringToUINT64(element[4], false)
				l.RXFifo = utils.StringToUINT64(element[5], false)
				l.RXFrame = utils.StringToUINT64(element[6], false)
				l.RXCompressed = utils.StringToUINT64(element[7], false)
				l.RXMulticast = utils.StringToUINT64(element[8], false)

				l.TXBytes = utils.StringToUINT64(element[9], false)
				l.TXPackets = utils.StringToUINT64(element[10], false)
				l.TXErrors = utils.StringToUINT64(element[11], false)
				l.TXDrops = utils.StringToUINT64(element[12], false)
				l.TXFifo = utils.StringToUINT64(element[13], false)
				l.TXCollisions = utils.StringToUINT64(element[14], false)
				l.TXCarrier = utils.StringToUINT64(element[15], false)
				l.TXCompressed = utils.StringToUINT64(element[16], false)
				l.Ts = utils.UTCTimeStamp()
				stats[l.InterfaceName] = l
			}
		}
		break
	}
	return stats
}
