package main

import (
	"github.com/coder/websocket"
)

type Player struct {
	Id          int
	X, Y        int
	MovingLeft  bool
	MovingRight bool
	MovingUp    bool
	MovingDown  bool
}

type PlayerWithSocket struct {
	Player
	Conn *websocket.Conn
}

type Event struct {
	Kind string
	Conn *websocket.Conn
	Data any
}

// ----

const (
	PlayerHelloKind  = "PlayerHello"
	PlayerQuitKind   = "PlayerQuit"
	PlayerJoinedKind = "PlayerJoined"
	PlayerMovedKind  = "PlayerMoved"
)

type KindHolder struct {
	Kind string `json:"Kind"`
}

type PlayerQuit struct {
	Kind string
	Id   int
}

type PlayerJoined struct {
	Kind   string
	Player Player
}

type PlayerHello struct {
	Kind string
	Id   int
}

type PlayerMoved struct {
	Kind        string
	Player      Player
	MovingLeft  bool
	MovingRight bool
	MovingUp    bool
	MovingDown  bool
}
