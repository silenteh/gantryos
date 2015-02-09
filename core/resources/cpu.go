package resources

import "runtime"

func GetTotalCPUsCount() float64 {
	return float64(runtime.NumCPU())
}
