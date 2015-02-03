//go:generate protoc --proto_path=${GOPATH}/src:${GOPATH}/src/github.com/gogo/protobuf/protobuf:. --gogo_out=. *.proto

package proto
