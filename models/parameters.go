package models

import "github.com/silenteh/gantryos/core/proto"

type parameters []*parameter

type parameter struct {
	Key   string
	Value string
}

func (p *parameter) ToProtoBuf() *proto.Parameter {

	param := new(proto.Parameter)
	param.Key = &p.Key
	param.Value = &p.Value
	return param
}

func (ps parameters) ToProtoBuf() *proto.Parameters {

	params := new(proto.Parameters)

	protoParams := make([]*proto.Parameter, len(ps))
	for index, res := range ps {
		protoParams[index] = res.ToProtoBuf()
	}
	params.Parameter = protoParams

	return params

}

func NewParameter(key, value string) *parameter {
	return &parameter{
		Key:   key,
		Value: value,
	}
}
