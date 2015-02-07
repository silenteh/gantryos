package resources

import "os"

func GetHostname() string {

	if hostname, err := os.Hostname(); err != nil {
		return "unknown"
	} else {
		return hostname
	}
}
