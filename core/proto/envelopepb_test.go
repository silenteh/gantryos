// Code generated by protoc-gen-gogo.
// source: envelope.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	envelope.proto
	gantryos.proto
	messages.proto

It has these top-level messages:
	Envelope
*/
package proto

import testing14 "testing"
import math_rand14 "math/rand"
import time14 "time"
import github_com_gogo_protobuf_proto8 "github.com/gogo/protobuf/proto"
import testing15 "testing"
import math_rand15 "math/rand"
import time15 "time"
import encoding_json2 "encoding/json"
import testing16 "testing"
import math_rand16 "math/rand"
import time16 "time"
import github_com_gogo_protobuf_proto9 "github.com/gogo/protobuf/proto"
import math_rand17 "math/rand"
import time17 "time"
import testing17 "testing"
import fmt4 "fmt"
import math_rand18 "math/rand"
import time18 "time"
import testing18 "testing"
import github_com_gogo_protobuf_proto10 "github.com/gogo/protobuf/proto"
import math_rand19 "math/rand"
import time19 "time"
import testing19 "testing"
import fmt5 "fmt"
import go_parser2 "go/parser"
import math_rand20 "math/rand"
import time20 "time"
import testing20 "testing"
import github_com_gogo_protobuf_proto11 "github.com/gogo/protobuf/proto"

func TestEnvelopeProto(t *testing14.T) {
	popr := math_rand14.New(math_rand14.NewSource(time14.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, false)
	data, err := github_com_gogo_protobuf_proto8.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &Envelope{}
	if err := github_com_gogo_protobuf_proto8.Unmarshal(data, msg); err != nil {
		panic(err)
	}
	for i := range data {
		data[i] = byte(popr.Intn(256))
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("%#v !Proto %#v", msg, p)
	}
}

func TestEnvelopeMarshalTo(t *testing14.T) {
	popr := math_rand14.New(math_rand14.NewSource(time14.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, false)
	size := p.Size()
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(popr.Intn(256))
	}
	_, err := p.MarshalTo(data)
	if err != nil {
		panic(err)
	}
	msg := &Envelope{}
	if err := github_com_gogo_protobuf_proto8.Unmarshal(data, msg); err != nil {
		panic(err)
	}
	for i := range data {
		data[i] = byte(popr.Intn(256))
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("%#v !Proto %#v", msg, p)
	}
}

func BenchmarkEnvelopeProtoMarshal(b *testing14.B) {
	popr := math_rand14.New(math_rand14.NewSource(616))
	total := 0
	pops := make([]*Envelope, 10000)
	for i := 0; i < 10000; i++ {
		pops[i] = NewPopulatedEnvelope(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := github_com_gogo_protobuf_proto8.Marshal(pops[i%10000])
		if err != nil {
			panic(err)
		}
		total += len(data)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkEnvelopeProtoUnmarshal(b *testing14.B) {
	popr := math_rand14.New(math_rand14.NewSource(616))
	total := 0
	datas := make([][]byte, 10000)
	for i := 0; i < 10000; i++ {
		data, err := github_com_gogo_protobuf_proto8.Marshal(NewPopulatedEnvelope(popr, false))
		if err != nil {
			panic(err)
		}
		datas[i] = data
	}
	msg := &Envelope{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += len(datas[i%10000])
		if err := github_com_gogo_protobuf_proto8.Unmarshal(datas[i%10000], msg); err != nil {
			panic(err)
		}
	}
	b.SetBytes(int64(total / b.N))
}

func TestEnvelopeJSON(t *testing15.T) {
	popr := math_rand15.New(math_rand15.NewSource(time15.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, true)
	jsondata, err := encoding_json2.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &Envelope{}
	err = encoding_json2.Unmarshal(jsondata, msg)
	if err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("%#v !Json Equal %#v", msg, p)
	}
}
func TestEnvelopeProtoText(t *testing16.T) {
	popr := math_rand16.New(math_rand16.NewSource(time16.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, true)
	data := github_com_gogo_protobuf_proto9.MarshalTextString(p)
	msg := &Envelope{}
	if err := github_com_gogo_protobuf_proto9.UnmarshalText(data, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("%#v !Proto %#v", msg, p)
	}
}

func TestEnvelopeProtoCompactText(t *testing16.T) {
	popr := math_rand16.New(math_rand16.NewSource(time16.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, true)
	data := github_com_gogo_protobuf_proto9.CompactTextString(p)
	msg := &Envelope{}
	if err := github_com_gogo_protobuf_proto9.UnmarshalText(data, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("%#v !Proto %#v", msg, p)
	}
}

func TestEnvelopeStringer(t *testing17.T) {
	popr := math_rand17.New(math_rand17.NewSource(time17.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, false)
	s1 := p.String()
	s2 := fmt4.Sprintf("%v", p)
	if s1 != s2 {
		t.Fatalf("String want %v got %v", s1, s2)
	}
}
func TestEnvelopeSize(t *testing18.T) {
	popr := math_rand18.New(math_rand18.NewSource(time18.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, true)
	size2 := github_com_gogo_protobuf_proto10.Size(p)
	data, err := github_com_gogo_protobuf_proto10.Marshal(p)
	if err != nil {
		panic(err)
	}
	size := p.Size()
	if len(data) != size {
		t.Fatalf("size %v != marshalled size %v", size, len(data))
	}
	if size2 != size {
		t.Fatalf("size %v != before marshal proto.Size %v", size, size2)
	}
	size3 := github_com_gogo_protobuf_proto10.Size(p)
	if size3 != size {
		t.Fatalf("size %v != after marshal proto.Size %v", size, size3)
	}
}

func BenchmarkEnvelopeSize(b *testing18.B) {
	popr := math_rand18.New(math_rand18.NewSource(616))
	total := 0
	pops := make([]*Envelope, 1000)
	for i := 0; i < 1000; i++ {
		pops[i] = NewPopulatedEnvelope(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += pops[i%1000].Size()
	}
	b.SetBytes(int64(total / b.N))
}

func TestEnvelopeGoString(t *testing19.T) {
	popr := math_rand19.New(math_rand19.NewSource(time19.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, false)
	s1 := p.GoString()
	s2 := fmt5.Sprintf("%#v", p)
	if s1 != s2 {
		t.Fatalf("GoString want %v got %v", s1, s2)
	}
	_, err := go_parser2.ParseExpr(s1)
	if err != nil {
		panic(err)
	}
}
func TestEnvelopeVerboseEqual(t *testing20.T) {
	popr := math_rand20.New(math_rand20.NewSource(time20.Now().UnixNano()))
	p := NewPopulatedEnvelope(popr, false)
	data, err := github_com_gogo_protobuf_proto11.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &Envelope{}
	if err := github_com_gogo_protobuf_proto11.Unmarshal(data, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseEqual %#v, since %v", msg, p, err)
	}
}

//These tests are generated by github.com/gogo/protobuf/plugin/testgen
