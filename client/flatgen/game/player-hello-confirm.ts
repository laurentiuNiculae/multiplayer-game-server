// automatically generated by the FlatBuffers compiler, do not modify

/* eslint-disable @typescript-eslint/no-unused-vars, @typescript-eslint/no-explicit-any, @typescript-eslint/no-non-null-assertion */

import * as flatbuffers from '../../flatbuffers/flatbuffers.js';

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

id():number {
  const offset = this.bb!.__offset(this.bb_pos, 4);
  return offset ? this.bb!.readInt32(this.bb_pos + offset) : 0;
}

static startPlayerHelloConfirm(builder:flatbuffers.Builder) {
  builder.startObject(1);
}

static addId(builder:flatbuffers.Builder, id:number) {
  builder.addFieldInt32(0, id, 0);
}

static endPlayerHelloConfirm(builder:flatbuffers.Builder):flatbuffers.Offset {
  const offset = builder.endObject();
  return offset;
}

static createPlayerHelloConfirm(builder:flatbuffers.Builder, id:number):flatbuffers.Offset {
  PlayerHelloConfirm.startPlayerHelloConfirm(builder);
  PlayerHelloConfirm.addId(builder, id);
  return PlayerHelloConfirm.endPlayerHelloConfirm(builder);
}
}
