package resources

import "runtime"

func GetTotalCPUsCount() int {
	return runtime.NumCPU()
}
