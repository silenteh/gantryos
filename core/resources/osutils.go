package resources

import (
	"github.com/silenteh/gantryos/utils"
	"strings"
)

const (
	LINUX   = "linux"
	BSD     = "bsd"
	WINDOWS = "windows"
	UNKNOWN = "unknown"
)

// TODO
// make sure the command uname exists (windows for example)
func detectOS() string {
	if utils.FileExists(`/usr/bin/uname`) || utils.FileExists(`/usr/sbin/uname`) {
		out := utils.ExecCommand(true, "uname")
		ostype := strings.ToLower(out)

		switch ostype {
		case "darwin", "freebsd", "openbsd":
			return BSD
		case "linux":
			return LINUX
		default:
			return UNKNOWN
		}
	}

	if utils.FileExists("/usr/bin/lsb_release") {
		return LINUX
	}

	return UNKNOWN
}

func pageSize() uint64 {
	detectedOs := detectOS()
	pageSize := uint64(1) // this will get multiplied somwehere else !!!! DO NOT PUT IT TO ZERO !!!
	switch detectedOs {
	case BSD:
		pageSizeString := utils.ExecCommand(true, "pagesize")
		pageSize = utils.StringToUINT64(pageSizeString, true)
		break
	case LINUX:
		//getconf PAGESIZE
		pageSizeString := utils.ExecCommand(true, "getconf", "PAGESIZE")
		pageSize = utils.StringToUINT64(pageSizeString, true)
		break
	}
	return pageSize
}
