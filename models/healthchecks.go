package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
	"strings"
)

type healthchecks []*healthcheck

type healthcheck struct {
	Port                 int    //uint32
	Path                 string //string
	ExpectedStatuses     []int  //[]uint32
	DelaySeconds         int    //float64 // Amount of time to wait until starting the health checks. DEFAULT 15
	IntervalSeconds      int    //float64 // Interval between health checks. DEFAULT 20
	TimeoutSeconds       int    //float64 // Amount of time to wait for the health check to complete. DEFAULT 10
	ConsecuritveFailures int    //uint32  // Number of consecutive failures until considered unhealthy. DEFAULT 3
	GracePeriodSeconds   int    //float64 // Amount of time to allow failed health checks since launch. DEFAULT 10
}

func (hc *healthcheck) toProtoBuf() *proto.HealthCheck {

	healthCheckProto := new(proto.HealthCheck)

	httpHealthCheck := new(proto.HealthCheck_HTTP)
	httpHealthCheck.Port = protobuf.Uint32(uint32(hc.Port))
	httpHealthCheck.Path = &hc.Path
	httpHealthCheck.Statuses = convertArrayIntToUint32(hc.ExpectedStatuses)
	healthCheckProto.Http = httpHealthCheck
	healthCheckProto.DelaySeconds = protobuf.Float64(float64(hc.DelaySeconds))
	healthCheckProto.IntervalSeconds = protobuf.Float64(float64(hc.IntervalSeconds))
	healthCheckProto.TimeoutSeconds = protobuf.Float64(float64(hc.TimeoutSeconds))
	healthCheckProto.GracePeriodSeconds = protobuf.Float64(float64(hc.GracePeriodSeconds))
	healthCheckProto.Failures = protobuf.Uint32(uint32(hc.ConsecuritveFailures))

	return healthCheckProto
}

func (hcs healthchecks) toProtoBuf() []*proto.HealthCheck {
	hcsProto := make([]*proto.HealthCheck, len(hcs))

	for index, el := range hcs {
		hcsProto[index] = el.toProtoBuf()
	}

	return hcsProto
}

func convertArrayIntToUint32(original []int) []uint32 {
	newArray := make([]uint32, len(original))
	for index, el := range original {
		newArray[index] = uint32(el)
	}

	return newArray
}

func NewHealthCheck(port int, path string, intervalSec, timeoutSec, delaySec, gracePeriodSec, consecutiveFailures int, expectedStatuses []int) *healthcheck {

	hc := healthcheck{}

	// PATH some checks - need to be extended probably
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	hc.Path = path

	// port
	if port <= 0 || port > 65535 {
		port = 80
	}
	hc.Port = port

	// health check interval
	if intervalSec <= 0 {
		intervalSec = 10
	}
	hc.IntervalSeconds = intervalSec

	// expected statuses
	total := 0
	convertedStatuses := make([]int, len(expectedStatuses))
	for _, status := range expectedStatuses {
		if status < 100 || status > 600 {
			continue
		}
		convertedStatuses[total] = status
		total++
	}
	hc.ExpectedStatuses = convertedStatuses[:total]

	// timeout seconds
	if timeoutSec <= 0 {
		timeoutSec = 5
	}
	hc.TimeoutSeconds = timeoutSec

	// consecutive failures
	if consecutiveFailures <= 0 {
		consecutiveFailures = 1
	}
	hc.ConsecuritveFailures = consecutiveFailures

	// delay seconds
	if delaySec <= 0 {
		delaySec = 5
	}
	hc.DelaySeconds = delaySec

	// grace period
	if gracePeriodSec <= 0 {
		gracePeriodSec = 5
	}
	hc.GracePeriodSeconds = gracePeriodSec

	return &hc

}
