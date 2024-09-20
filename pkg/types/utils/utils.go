package utils

import (
	"fmt"
	"net/http"
	. "test/pkg/types"
	flatgen "test/pkg/types/flatgen/game"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

func NewFlatEvent(builder *flatbuffers.Builder, kind string, bytes []byte) *flatgen.Event {
	flatKind := builder.CreateByteString([]byte(kind))
	flatBytes := builder.CreateByteVector(bytes)

	flatgen.EventStart(builder)
	flatgen.EventAddKind(builder, flatKind)
	flatgen.EventAddData(builder, flatBytes)
	builder.Finish(flatgen.EventEnd(builder))

	return flatgen.GetRootAsEvent(builder.FinishedBytes(), 0)
}

func NewFlatPlayerHello(builder *flatbuffers.Builder, newPlayer Player) *flatgen.PlayerHello {
	flatgen.PlayerHelloStart(builder)
	flatgen.PlayerHelloAddId(builder, int32(newPlayer.Id))
	flatgen.FinishPlayerHelloBuffer(builder, flatgen.PlayerHelloEnd(builder))

	return flatgen.GetRootAsPlayerHello(builder.FinishedBytes(), 0)
}

func NewFlatPlayerHelloConfirm(builder *flatbuffers.Builder, id int) *flatgen.PlayerHelloConfirm {
	flatgen.PlayerHelloConfirmStart(builder)
	flatgen.PlayerHelloConfirmAddId(builder, int32(id))
	flatgen.FinishPlayerHelloConfirmBuffer(builder, flatgen.PlayerHelloConfirmEnd(builder))

	return flatgen.GetRootAsPlayerHelloConfirm(builder.FinishedBytes(), 0)
}

func NewFlatPlayerQuit(builder *flatbuffers.Builder, playerId int) *flatgen.PlayerQuit {
	flatgen.PlayerQuitStart(builder)
	flatgen.PlayerQuitAddId(builder, int32(playerId))
	flatgen.FinishPlayerQuitBuffer(builder, flatgen.PlayerQuitEnd(builder))

	return flatgen.GetRootAsPlayerQuit(builder.FinishedBytes(), 0)
}

func NewFlatPlayerJoined(builder *flatbuffers.Builder, newPlayer Player) *flatgen.PlayerJoined {
	flatPlayer := NewFlatPlayer(builder, newPlayer)

	flatgen.PlayerJoinedStart(builder)
	flatgen.PlayerJoinedAddPlayer(builder, flatPlayer)
	flatgen.FinishPlayerJoinedBuffer(builder, flatgen.PlayerJoinedEnd(builder))

	return flatgen.GetRootAsPlayerJoined(builder.FinishedBytes(), 0)
}

func NewFlatPlayerMoved(builder *flatbuffers.Builder, newPlayer Player) *flatgen.PlayerMoved {
	flatPlayer := NewFlatPlayer(builder, newPlayer)

	flatgen.PlayerMovedStart(builder)
	flatgen.PlayerMovedAddPlayer(builder, flatPlayer)
	flatgen.FinishPlayerMovedBuffer(builder, flatgen.PlayerMovedEnd(builder))

	return flatgen.GetRootAsPlayerMoved(builder.FinishedBytes(), 0)
}

func NewFlatPlayerMovedList(builder *flatbuffers.Builder, movingPlayers []*flatgen.PlayerMoved) *flatgen.PlayerMovedList {
	flatgen.PlayerMovedListStartPlayersVector(builder, len(movingPlayers))
	for i := range movingPlayers {
		NewFlatPlayerFromFlat(builder, movingPlayers[i].Player(nil))
	}
	movingPlayersVecOffset := builder.EndVector(len(movingPlayers))

	flatgen.PlayerMovedListStart(builder)
	flatgen.PlayerMovedListAddPlayers(builder, movingPlayersVecOffset)
	flatgen.FinishPlayerMovedListBuffer(builder, flatgen.PlayerMovedListEnd(builder))

	return flatgen.GetRootAsPlayerMovedList(builder.FinishedBytes(), 0)
}

func NewFlatPlayerJoinedList(builder *flatbuffers.Builder, joinedPlayers []Player) *flatgen.PlayerJoinedList {
	flatgen.PlayerJoinedListStartPlayersVector(builder, len(joinedPlayers))
	for i := range joinedPlayers {
		NewFlatPlayer(builder, joinedPlayers[i])
	}
	movingPlayersVecOffset := builder.EndVector(len(joinedPlayers))

	flatgen.PlayerJoinedListStart(builder)
	flatgen.PlayerJoinedListAddPlayers(builder, movingPlayersVecOffset)
	flatgen.FinishPlayerJoinedListBuffer(builder, flatgen.PlayerJoinedListEnd(builder))

	return flatgen.GetRootAsPlayerJoinedList(builder.FinishedBytes(), 0)
}

func NewFlatPlayer(builder *flatbuffers.Builder, newPlayer Player) flatbuffers.UOffsetT {
	return flatgen.CreatePlayer(builder,
		int32(newPlayer.Id),
		int32(newPlayer.X),
		int32(newPlayer.Y),
		int32(newPlayer.Speed),
		newPlayer.MovingLeft,
		newPlayer.MovingRight,
		newPlayer.MovingUp,
		newPlayer.MovingDown,
	)
}

func NewFlatPlayerFromFlat(builder *flatbuffers.Builder, newPlayer *flatgen.Player) flatbuffers.UOffsetT {
	return flatgen.CreatePlayer(builder,
		int32(newPlayer.Id()),
		int32(newPlayer.X()),
		int32(newPlayer.Y()),
		int32(newPlayer.Speed()),
		newPlayer.MovingLeft(),
		newPlayer.MovingRight(),
		newPlayer.MovingUp(),
		newPlayer.MovingDown(),
	)
}

func ParseEventBytes(data []byte) (eventKind string, eventData any, err error) {
	flatEvent := flatgen.GetRootAsEvent(data, 0)
	eventKind = string(flatEvent.Kind())

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("was panic, returned panic value '%v'", r)
		}
	}()

	switch eventKind {
	case PlayerHelloKind:
		flatPlayerHello := flatgen.GetRootAsPlayerHello(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerHello, nil
	case PlayerHelloConfirmKind:
		flatPlayerHelloConfirm := flatgen.GetRootAsPlayerHelloConfirm(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerHelloConfirm, nil
	case PlayerQuitKind:
		flatPlayerQuit := flatgen.GetRootAsPlayerQuit(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerQuit, nil
	case PlayerJoinedKind:
		flatPlayerJoined := flatgen.GetRootAsPlayerJoined(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerJoined, nil
	case PlayerJoinedListKind:
		flatPlayerJoinedList := flatgen.GetRootAsPlayerJoinedList(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerJoinedList, nil
	case PlayerMovedKind:
		flatPlayerMoved := flatgen.GetRootAsPlayerMoved(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerMoved, nil
	case PlayerMovedListKind:
		flatPlayerMovedList := flatgen.GetRootAsPlayerMovedList(flatEvent.DataBytes(), 0)

		return eventKind, flatPlayerMovedList, nil
	default:
		return "", nil, fmt.Errorf("ERROR: bogus-amogus kind '%s'", string(flatEvent.Kind()))
	}
}

func GetFlatPlayerHello(builder *flatbuffers.Builder, newPlayer Player) []byte {
	flatgen.PlayerHelloStart(builder)
	flatgen.PlayerHelloAddId(builder, int32(newPlayer.Id))
	flatgen.FinishPlayerHelloBuffer(builder, flatgen.PlayerHelloEnd(builder))

	return builder.FinishedBytes()
}

func GetFlatPlayerQuit(builder *flatbuffers.Builder, playerId int) []byte {
	flatgen.PlayerQuitStart(builder)
	flatgen.PlayerQuitAddId(builder, int32(playerId))
	flatgen.FinishPlayerQuitBuffer(builder, flatgen.PlayerQuitEnd(builder))

	return builder.FinishedBytes()
}

func GetFlatPlayerJoined(builder *flatbuffers.Builder, newPlayer Player) []byte {
	flatPlayer := NewFlatPlayer(builder, newPlayer)

	flatgen.PlayerJoinedStart(builder)
	flatgen.PlayerJoinedAddPlayer(builder, flatPlayer)
	flatgen.FinishPlayerJoinedBuffer(builder, flatgen.PlayerJoinedEnd(builder))

	return builder.FinishedBytes()
}

func GetFlatEvent(builder *flatbuffers.Builder, kind string, bytes []byte) *flatgen.Event {
	flatKind := builder.CreateByteString([]byte(kind))
	flatBytes := builder.CreateByteVector(bytes)

	flatgen.EventStart(builder)
	flatgen.EventAddKind(builder, flatKind)
	flatgen.EventAddData(builder, flatBytes)
	builder.Finish(flatgen.EventEnd(builder))

	return flatgen.GetRootAsEvent(builder.FinishedBytes(), 0)
}

func WaitServerIsReady(url string) {
	for {
		_, err := http.Get(url)
		if err == nil {
			return
		}

		time.Sleep(1 * time.Second)
	}
}
