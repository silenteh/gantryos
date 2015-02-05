package models

import (
	"errors"
	"github.com/silenteh/gantryos/core/proto"
)

const (
	CONTAINER_VOLUME_RO = proto.Volume_RO
	CONTAINER_VOLUME_RW = proto.Volume_RW
)

type containerVolumes []*containerVolume

type containerVolume struct {
	ContainerPath string
	HostPath      string
	Persistent    bool
	Permission    proto.Volume_Mode // CONTAINER_VOLUME_RO || CONTAINER_VOLUME_RW
}

func (cv containerVolume) ToProtoBuf() *proto.Volume {
	cvProto := new(proto.Volume)
	cvProto.HostPath = &cv.HostPath
	cvProto.ContainerPath = &cv.ContainerPath
	cvProto.Persistent = &cv.Persistent
	cvProto.Mode = &cv.Permission
	return cvProto
}

func (cvs containerVolumes) ToProtoBuf() []*proto.Volume {
	hcsProto := make([]*proto.Volume, len(cvs))

	for index, el := range cvs {
		hcsProto[index] = el.ToProtoBuf()
	}

	return hcsProto
}

func NewContainerVolume(containerPath, hostPath string, persistent bool, permission proto.Volume_Mode) (*containerVolume, error) {
	v := &containerVolume{}
	if containerPath == "" {
		return v, errors.New("A container volume must have the container path")
	}

	if hostPath == "" {
		return v, errors.New("A container volume must have the host path")
	}

	v.ContainerPath = containerPath
	v.HostPath = hostPath
	v.Persistent = persistent

	// TODO: check the permission input value
	v.Permission = permission

	return v, nil

}
