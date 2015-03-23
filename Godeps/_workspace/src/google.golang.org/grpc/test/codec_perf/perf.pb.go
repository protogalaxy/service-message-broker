// Code generated by protoc-gen-go.
// source: perf.proto
// DO NOT EDIT!

/*
Package codec_perf is a generated protocol buffer package.

It is generated from these files:
	perf.proto

It has these top-level messages:
	Buffer
*/
package codec_perf

import proto "github.com/protogalaxy/service-message-broker/Godeps/_workspace/src/github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// Buffer is a message that contains a body of bytes that is used to exercise
// encoding and decoding overheads.
type Buffer struct {
	Body             []byte `protobuf:"bytes,1,opt,name=body" json:"body,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Buffer) Reset()         { *m = Buffer{} }
func (m *Buffer) String() string { return proto.CompactTextString(m) }
func (*Buffer) ProtoMessage()    {}

func (m *Buffer) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
}
