package models

import "github.com/silenteh/gantryos/core/proto"

type environmentVariables []*environmentVariable

type environmentVariable struct {
	Name  string
	Value string
}

func (e *environmentVariable) ToProtoBuf() *proto.Environment_Variable {

	env := new(proto.Environment_Variable)
	env.Name = &e.Name
	env.Value = &e.Value
	return env
}

func (ev environmentVariables) ToProtoBuf() *proto.Environment {

	envs := new(proto.Environment)

	protoEnvs := make([]*proto.Environment_Variable, len(ev))
	for index, res := range ev {
		protoEnvs[index] = res.ToProtoBuf()
	}
	envs.Variables = protoEnvs

	return envs

}

func NewEnvironmentVariable(name, value string) *environmentVariable {
	return &environmentVariable{
		Name:  name,
		Value: value,
	}
}

func NewEnvironmentVariables(envs ...*environmentVariable) environmentVariables {
	var allEnvs environmentVariables
	allEnvs = envs
	return allEnvs
}
