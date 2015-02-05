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

func (l *label) ToProtoBuf() *proto.Label {

	label := new(proto.Label)
	label.Key = &l.Key
	label.Value = &l.Value
	return label
}

func (ls labels) ToProtoBuf() *proto.Labels {

	labels := new(proto.Labels)

	protoLabels := make([]*proto.Label, len(ls))
	for index, res := range ls {
		protoLabels[index] = res.ToProtoBuf()
	}
	labels.Labels = protoLabels

	return labels

}
