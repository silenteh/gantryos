package models

import "github.com/silenteh/gantryos/core/proto"

type labels []*label

type label struct {
	Key   string
	Value string
}

func NewLabel(key, value string) *label {
	return &label{
		Key:   key,
		Value: value,
	}
}

func (l *label) toProtoBuf() *proto.Label {

	label := new(proto.Label)
	label.Key = &l.Key
	label.Value = &l.Value
	return label
}

func (ls labels) toProtoBuf() *proto.Labels {

	labels := new(proto.Labels)

	protoLabels := make([]*proto.Label, len(ls))
	for index, res := range ls {
		protoLabels[index] = res.toProtoBuf()
	}
	labels.Labels = protoLabels

	return labels

}
