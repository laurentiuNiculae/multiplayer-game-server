// automatically generated by the FlatBuffers compiler, do not modify
/* eslint-disable @typescript-eslint/no-unused-vars, @typescript-eslint/no-explicit-any, @typescript-eslint/no-non-null-assertion */
import * as flatbuffers from '../../flatbuffers/flatbuffers.js';
export class Player {
    bb = null;
    bb_pos = 0;
    __init(i, bb) {
        this.bb_pos = i;
        this.bb = bb;
        return this;
    }
    static getRootAsPlayer(bb, obj) {
        return (obj || new Player()).__init(bb.readInt32(bb.position()) + bb.position(), bb);
    }
    static getSizePrefixedRootAsPlayer(bb, obj) {
        bb.setPosition(bb.position() + flatbuffers.SIZE_PREFIX_LENGTH);
        return (obj || new Player()).__init(bb.readInt32(bb.position()) + bb.position(), bb);
    }
    id() {
        const offset = this.bb.__offset(this.bb_pos, 4);
        return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
    }
    x() {
        const offset = this.bb.__offset(this.bb_pos, 6);
        return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
    }
    y() {
        const offset = this.bb.__offset(this.bb_pos, 8);
        return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
    }
    speed() {
        const offset = this.bb.__offset(this.bb_pos, 10);
        return offset ? this.bb.readInt32(this.bb_pos + offset) : 0;
    }
    movingLeft() {
        const offset = this.bb.__offset(this.bb_pos, 12);
        return offset ? !!this.bb.readInt8(this.bb_pos + offset) : false;
    }
    movingRight() {
        const offset = this.bb.__offset(this.bb_pos, 14);
        return offset ? !!this.bb.readInt8(this.bb_pos + offset) : false;
    }
    movingUp() {
        const offset = this.bb.__offset(this.bb_pos, 16);
        return offset ? !!this.bb.readInt8(this.bb_pos + offset) : false;
    }
    movingDown() {
        const offset = this.bb.__offset(this.bb_pos, 18);
        return offset ? !!this.bb.readInt8(this.bb_pos + offset) : false;
    }
    static startPlayer(builder) {
        builder.startObject(8);
    }
    static addId(builder, id) {
        builder.addFieldInt32(0, id, 0);
    }
    static addX(builder, x) {
        builder.addFieldInt32(1, x, 0);
    }
    static addY(builder, y) {
        builder.addFieldInt32(2, y, 0);
    }
    static addSpeed(builder, speed) {
        builder.addFieldInt32(3, speed, 0);
    }
    static addMovingLeft(builder, movingLeft) {
        builder.addFieldInt8(4, +movingLeft, +false);
    }
    static addMovingRight(builder, movingRight) {
        builder.addFieldInt8(5, +movingRight, +false);
    }
    static addMovingUp(builder, movingUp) {
        builder.addFieldInt8(6, +movingUp, +false);
    }
    static addMovingDown(builder, movingDown) {
        builder.addFieldInt8(7, +movingDown, +false);
    }
    static endPlayer(builder) {
        const offset = builder.endObject();
        return offset;
    }
    static createPlayer(builder, id, x, y, speed, movingLeft, movingRight, movingUp, movingDown) {
        Player.startPlayer(builder);
        Player.addId(builder, id);
        Player.addX(builder, x);
        Player.addY(builder, y);
        Player.addSpeed(builder, speed);
        Player.addMovingLeft(builder, movingLeft);
        Player.addMovingRight(builder, movingRight);
        Player.addMovingUp(builder, movingUp);
        Player.addMovingDown(builder, movingDown);
        return Player.endPlayer(builder);
    }
}
