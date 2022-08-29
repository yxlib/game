// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

type GameInput struct {
	PlayerId uint64
	Cmd      string
	Params   []interface{}
}

type GameInputProcessor interface {
	HandleGameInput(playerId uint64, cmd string, params []interface{})
}

type GameEvent interface {
	Exec(s Game) error
}

type GameEventProcessor interface {
	HandleGameEvent(e GameEvent)
}

type Player interface {
	GetPlayerId() uint64
}

type PlayerMgr interface {
	AddPlayers(players []Player) error
	GetPlayerIds() []uint64
	IsExistPlayer(playerId uint64) bool
	SetPlayerOnline(playerId uint64)
	SetPlayerOffline(playerId uint64)
	RemovePlayer(playerId uint64)
}

type Game interface {
	PlayerMgr

	SetGameID(gameId uint16)
	GetGameID() uint16
	GetTemplateID() uint16

	AddInput(playerId uint64, cmd string, params ...interface{}) error
	HandleInputs()
	AddEvent(e GameEvent) error
	HandleEvents()

	Init() error
	Start(now int64) error
	Update(dt int64)
	IsGameOver() bool
	Stop()
	IsAutoRestart() bool
	Destroy()
}
