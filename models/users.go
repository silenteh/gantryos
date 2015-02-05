package models

import (
	protobuf "github.com/gogo/protobuf/proto"
	"github.com/silenteh/gantryos/core/proto"
)

type user struct {
	Uid  int
	Gid  int
	Name string
}

func NewUser(uid, gid int, name string) *user {
	return &user{
		Uid:  uid,
		Gid:  gid,
		Name: name,
	}
}

func (u *user) ToProtoBuf() *proto.User {

	user := new(proto.User)
	user.Gid = protobuf.Int32(int32(u.Gid))
	user.Uid = protobuf.Int32(int32(u.Uid))
	user.Name = &u.Name

	return user

}
