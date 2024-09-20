// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package game

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type PlayerJoined struct {
	_tab flatbuffers.Table
}

func GetRootAsPlayerJoined(buf []byte, offset flatbuffers.UOffsetT) *PlayerJoined {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PlayerJoined{}
	x.Init(buf, n+offset)
	return x
}

func FinishPlayerJoinedBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsPlayerJoined(buf []byte, offset flatbuffers.UOffsetT) *PlayerJoined {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &PlayerJoined{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedPlayerJoinedBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *PlayerJoined) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PlayerJoined) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PlayerJoined) Player(obj *Player) *Player {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(Player)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func PlayerJoinedStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func PlayerJoinedAddPlayer(builder *flatbuffers.Builder, player flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(player), 0)
}
func PlayerJoinedEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
