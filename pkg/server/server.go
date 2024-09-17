package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sort"
	"time"

	"test/pkg/log"
	. "test/pkg/types"
	flatgen "test/pkg/types/flatgen/game"
	"test/pkg/types/utils"

	"github.com/coder/websocket"
	flatbuffers "github.com/google/flatbuffers/go"
)

var ServerFPS = 30
var WorldWidth = float64(800 * 2)
var WorldHeight = float64(600 * 2)
var Port = "6969"
var Address = "127.0.0.1:" + Port
var HttpAddress = "http://127.0.0.1:" + Port

type IdGenerator struct {
	idCounter int
}

func (igen *IdGenerator) NewId() int {
	igen.idCounter++

	return igen.idCounter
}

type GameServer struct {
	Players        PlayerStore
	EventQueue     chan Event
	IdGenerator    IdGenerator
	EventCollector *EventCollector
	mux            *http.ServeMux
	log            log.MeloLog
}

func NewGame() GameServer {
	return GameServer{
		Players:        NewPlayerStore(),
		EventQueue:     make(chan Event, 1000),
		IdGenerator:    IdGenerator{},
		EventCollector: NewEventCollector(),
		mux:            http.NewServeMux(),
		log:            log.New(os.Stdout),
	}
}

func (game *GameServer) Start() {
	game.mux.Handle("/", http.FileServer(http.Dir(".")))
	game.mux.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		wcon, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		}

		playerId := game.IdGenerator.NewId()

		defer func() {
			builder := flatbuffers.NewBuilder(128)

			game.EventQueue <- Event{
				PlayerId: playerId,
				Kind:     PlayerQuitKind,
				Conn:     wcon,
				Data:     flatgen.GetRootAsPlayerQuit(GetFlatPlayerQuit(builder, playerId), 0),
			}

			game.log.Infof("Player '%v' diconnected", playerId)
		}()

		game.EventQueue <- Event{
			PlayerId: playerId,
			Kind:     PlayerHelloKind,
			Conn:     wcon,
			Data:     PlayerHello{Kind: PlayerHelloKind, Id: playerId},
		}

		for {
			_, dataBytes, err := wcon.Read(ctx)
			if err != nil {
				return
			}

			kind, data, err := utils.ParseEventBytes(dataBytes)
			if err != nil {
				game.log.Errorf("err: %v\n", err)
				continue
			}

			game.EventQueue <- Event{
				PlayerId: playerId,
				Kind:     kind,
				Data:     data,
				Conn:     wcon,
			}
		}
	})

	go func() {
		WaitServerIsReady(HttpAddress)
		game.log.Info("Listening to server")
	}()

	go game.Tick()

	err := http.ListenAndServe(Address, game.mux)
	if err != nil {
		game.log.Errorf("err: %s\n", err.Error())
	}
}

func PrintMemUsage(log log.MeloLog) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Debugf("Alloc = %v MiB\n\tTotalAlloc = %v MiB\n\tSys = %v MiB\n", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (game *GameServer) Tick() {
	WaitServerIsReady(HttpAddress)

	tickTimeArr := make([]float64, 30)
	timeI := 0

	ticker := time.NewTicker(1 * time.Second / time.Duration(ServerFPS))
	previousTime, delta := time.Now(), time.Duration(0)

	playerMovedBuilder := flatbuffers.NewBuilder(512)
	playerMovedList := []*flatgen.PlayerMoved{}

	<-ticker.C

	for range ticker.C {
		start := time.Now()
		ctx := context.Background()

		for range len(game.EventQueue) {
			event := <-game.EventQueue

			switch event.Kind {
			case PlayerHelloKind:
				playerHello := event.Data.(PlayerHello)

				newPlayer := PlayerWithSocket{
					Conn: event.Conn,
					Player: Player{
						Id:    playerHello.Id,
						Speed: rand.Float64()*100 + 200,
						X:     rand.Float64()*float64(WorldWidth)/4 + float64(WorldWidth)/2,
						Y:     rand.Float64()*float64(WorldHeight)/4 + +float64(WorldHeight)/2,
					},
				}

				game.Players.Set(newPlayer.Id, newPlayer)

				game.log.Infof("Player connected: '%v'", playerHello.Id)

				builder := flatbuffers.NewBuilder(512)

				eventData := GetFlatPlayerHello(builder, newPlayer.Player)
				eventBytes := GetFlatEvent(builder, PlayerHelloKind, eventData).Table().Bytes

				err := newPlayer.Conn.Write(ctx, websocket.MessageBinary, eventBytes)
				if err != nil {
					game.log.Errorf("err: %s\n", err.Error())
				}
			case PlayerHelloConfirmKind:
				helloResponse := event.Data.(*flatgen.PlayerHelloConfirm)

				if helloResponse.Id() == int32(event.PlayerId) {
					game.log.Debug("HELLO CONFIRMED BY PLAYER")
				} else {
					game.log.Debugf("player ID doesn't match expected:'%d', given:'%d'", event.PlayerId, helloResponse.Id())
				}

				builder := flatbuffers.NewBuilder(512)

				newPlayer, _ := game.Players.Get(event.PlayerId)

				flatNewPlayerJoined := GetFlatPlayerJoined(builder, newPlayer.Player)
				flatNewPlayerJoinedEvent := GetFlatEvent(builder, PlayerJoinedKind, flatNewPlayerJoined)

				for _, otherPlayer := range game.Players.All() {
					builder2 := flatbuffers.NewBuilder(512)
					game.EventCollector.AddEvent(otherPlayer.Id, flatNewPlayerJoinedEvent)
					// otherPlayer.Conn.Write(ctx, websocket.MessageBinary, flatNewPlayerJoinedEvent)

					flatOtherPlayerJoined := GetFlatPlayerJoined(builder2, otherPlayer.Player)
					flatOtherPlayerJoinedEvent := GetFlatEvent(builder2, PlayerJoinedKind, flatOtherPlayerJoined)
					if otherPlayer.Id != newPlayer.Id {
						game.EventCollector.AddEvent(newPlayer.Id, flatOtherPlayerJoinedEvent)
						// newPlayer.Conn.Write(ctx, websocket.MessageBinary, flatOtherPlayerJoinedEvent)
					}
				}
			case PlayerQuitKind:
				playerQuit := event.Data.(*flatgen.PlayerQuit)

				if playerQuit.Id() != int32(event.PlayerId) {
					event.Conn.CloseNow()
					game.log.Errorf("player '%s' tried to cheat", event.PlayerId)
				}

				playerQuitEvent := GetFlatEvent(flatbuffers.NewBuilder(512), PlayerQuitKind,
					playerQuit.Table().Bytes)

				game.Players.Delete(event.PlayerId)
				game.EventCollector.RemovePlayer(event.PlayerId)

				for _, player := range game.Players.All() {
					game.EventCollector.AddEvent(player.Id, playerQuitEvent)
				}

			case PlayerMovedKind:
				playerMoved := event.Data.(*flatgen.PlayerMoved)
				newPlayerInfo := playerMoved.Player(nil)

				if newPlayerInfo.Id() != int32(event.PlayerId) {
					event.Conn.CloseNow()
					game.log.Errorf("player '%s' tried to cheat", event.PlayerId)
				}

				player, _ := game.Players.Get(int(newPlayerInfo.Id())) // TODO _
				player.MovingLeft = newPlayerInfo.MovingLeft()
				player.MovingRight = newPlayerInfo.MovingRight()
				player.MovingUp = newPlayerInfo.MovingUp()
				player.MovingDown = newPlayerInfo.MovingDown()

				playerMoved.Player(nil).MutateX(int32(player.X))
				playerMoved.Player(nil).MutateY(int32(player.Y))

				game.Players.Set(int(newPlayerInfo.Id()), player)

				// playerMovedEvent := GetFlatEvent(flatbuffers.NewBuilder(512), PlayerMovedKind,
				// 	playerMoved.Table().Bytes)

				playerMovedList = append(playerMovedList, playerMoved)

				// for _, player := range game.Players.All() {
				// 	game.EventCollector.AddEvent(player.Id, playerMovedEvent)
				// }
			}
		}

		// calculate all players that moved event and send it.
		if len(playerMovedList) > 0 {
			flatPlayerMovedList := utils.NewFlatPlayerMovedList(playerMovedBuilder, playerMovedList)
			game.EventCollector.AddGeneralEvent(utils.NewFlatEvent(playerMovedBuilder, PlayerMovedListKind,
				flatPlayerMovedList.Table().Bytes))
		}

		// collect events here then send them.
		for id, player := range game.Players.All() {
			eventList := game.EventCollector.GetPlayerEventList(id)

			if eventList != nil {
				player.Conn.Write(ctx, websocket.MessageBinary, eventList.Table().Bytes)
			}
		}

		game.EventCollector.Reset()
		clear(playerMovedList)
		playerMovedList = playerMovedList[:0]
		playerMovedBuilder.Reset()

		delta, previousTime = time.Since(previousTime), time.Now()

		for i, player := range game.Players.All() {
			movedDelta := delta.Seconds() * player.Speed

			if player.MovingLeft && player.X-movedDelta >= 0 {
				player.X = player.X - movedDelta
			}
			if player.MovingRight && player.X+movedDelta < WorldWidth {
				player.X = player.X + movedDelta
			}
			if player.MovingUp && player.Y-movedDelta >= 0 {
				player.Y = player.Y - movedDelta
			}
			if player.MovingDown && player.Y+movedDelta < WorldHeight {
				player.Y = player.Y + movedDelta
			}

			game.Players.Set(i, player)
		}

		if timeI == 29 {
			sort.Slice(tickTimeArr, func(i, j int) bool {
				return tickTimeArr[i] < tickTimeArr[j]
			})

			game.log.Debugf("%06f", tickTimeArr[30/2])
			PrintMemUsage(game.log)

			timeI = 0
		}

		tickTimeArr[timeI] = time.Since(start).Seconds()
		timeI++
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
	flatPlayer := GetFlatPlayer(builder, newPlayer)

	flatgen.PlayerJoinedStart(builder)
	flatgen.PlayerJoinedAddPlayer(builder, flatPlayer)
	flatgen.FinishPlayerJoinedBuffer(builder, flatgen.PlayerJoinedEnd(builder))

	return builder.FinishedBytes()
}

func GetFlatPlayer(builder *flatbuffers.Builder, newPlayer Player) flatbuffers.UOffsetT {
	flatgen.PlayerStart(builder)
	flatgen.PlayerAddId(builder, int32(newPlayer.Id))
	flatgen.PlayerAddX(builder, int32(newPlayer.X))
	flatgen.PlayerAddY(builder, int32(newPlayer.Y))
	flatgen.PlayerAddSpeed(builder, int32(newPlayer.Speed))
	flatgen.PlayerAddMovingDown(builder, newPlayer.MovingDown)
	flatgen.PlayerAddMovingLeft(builder, newPlayer.MovingLeft)
	flatgen.PlayerAddMovingRight(builder, newPlayer.MovingRight)
	flatgen.PlayerAddMovingUp(builder, newPlayer.MovingUp)

	return flatgen.PlayerEnd(builder)
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

func (game *GameServer) NotifyAll(msg []byte) {
	for _, player := range game.Players.All() {
		err := player.Conn.Write(context.Background(), websocket.MessageBinary, msg)
		if err != nil {
			continue
		}
	}
}

func (game *GameServer) NotifyAllElse(msg []byte, except ...int) {
	for _, player := range game.Players.All() {
		if !slices.Contains(except, player.Id) {
			err := player.Conn.Write(context.Background(), websocket.MessageBinary, msg)
			if err != nil {
				continue
			}
		}
	}
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
