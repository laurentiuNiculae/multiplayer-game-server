// automatically generated by the FlatBuffers compiler, do not modify

/* eslint-disable @typescript-eslint/no-unused-vars, @typescript-eslint/no-explicit-any, @typescript-eslint/no-non-null-assertion */

import * as flatbuffers from '../../flatbuffers/flatbuffers.js';

import { EventKind } from '../../flatgen/game/event-kind.js';


export class PlayerHelloConfirm {
  bb: flatbuffers.ByteBuffer|null = null;
  bb_pos = 0;
  __init(i:number, bb:flatbuffers.ByteBuffer):PlayerHelloConfirm {
  this.bb_pos = i;
  this.bb = bb;
  return this;
}

static getRootAsPlayerHelloConfirm(bb:flatbuffers.ByteBuffer, obj?:PlayerHelloConfirm):PlayerHelloConfirm {
  return (obj || new PlayerHelloConfirm()).__init(bb.readInt32(bb.position()) + bb.position(), bb);
}

static getSizePrefixedRootAsPlayerHelloConfirm(bb:flatbuffers.ByteBuffer, obj?:PlayerHelloConfirm):PlayerHelloConfirm {
  bb.setPosition(bb.position() + flatbuffers.SIZE_PREFIX_LENGTH);
  return (obj || new PlayerHelloConfirm()).__init(bb.readInt32(bb.position()) + bb.position(), bb);
}

kind():EventKind {
  const offset = this.bb!.__offset(this.bb_pos, 4);
  return offset ? this.bb!.readUint8(this.bb_pos + offset) : EventKind.NilEvent;
}

id():number {
  const offset = this.bb!.__offset(this.bb_pos, 6);
  return offset ? this.bb!.readInt32(this.bb_pos + offset) : 0;
}

static startPlayerHelloConfirm(builder:flatbuffers.Builder) {
  builder.startObject(2);
}

static addKind(builder:flatbuffers.Builder, kind:EventKind) {
  builder.addFieldInt8(0, kind, EventKind.NilEvent);
}

static addId(builder:flatbuffers.Builder, id:number) {
  builder.addFieldInt32(1, id, 0);
}

static endPlayerHelloConfirm(builder:flatbuffers.Builder):flatbuffers.Offset {
  const offset = builder.endObject();
  return offset;
}

static createPlayerHelloConfirm(builder:flatbuffers.Builder, kind:EventKind, id:number):flatbuffers.Offset {
  PlayerHelloConfirm.startPlayerHelloConfirm(builder);
  PlayerHelloConfirm.addKind(builder, kind);
  PlayerHelloConfirm.addId(builder, id);
  return PlayerHelloConfirm.endPlayerHelloConfirm(builder);
}
}
