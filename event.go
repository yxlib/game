// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import "errors"

type PlayerOnlineEvent struct {
	PlayerId uint64
}

func (e *PlayerOnlineEvent) Exec(s Game) error {
	s.SetPlayerOnline(e.PlayerId)
	return nil
}

type PlayerOfflineEvent struct {
	PlayerId uint64
}

func (e *PlayerOfflineEvent) Exec(s Game) error {
	s.SetPlayerOffline(e.PlayerId)
	return nil
}

type PlayerJoinEvent struct {
	player Player
}

func (e *PlayerJoinEvent) Exec(s Game) error {
	playerId := e.player.GetPlayerId()
	existGame, ok := GameMgr.GetGameByPlayerId(playerId)
	if !ok {
		return errors.New("player not join game")
	}

	if existGame != s {
		return errors.New("not same game")
	}

	err := s.AddPlayers([]Player{e.player})
	return err
}

type PlayerRemoveEvent struct {
	PlayerId uint64
}

func (e *PlayerRemoveEvent) Exec(s Game) error {
	s.RemovePlayer(e.PlayerId)
	return nil
}
