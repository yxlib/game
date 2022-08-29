// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import (
	"errors"
	"sync"
	"time"

	"github.com/yxlib/yx"
)

const MAX_GAME_ID = 1024

type gameMgr struct {
	mapId2Game map[uint16]Game
	lckGame    *sync.RWMutex

	mapPlayerId2GameId map[uint64]uint16
	lckPlayer          *sync.RWMutex

	idGen *yx.IdGenerator
	ec    *yx.ErrCatcher
}

var GameMgr = &gameMgr{
	mapId2Game: make(map[uint16]Game),
	lckGame:    &sync.RWMutex{},

	mapPlayerId2GameId: make(map[uint64]uint16),
	lckPlayer:          &sync.RWMutex{},

	idGen: yx.NewIdGenerator(1, MAX_GAME_ID),
	ec:    yx.NewErrCatcher("GameMgr"),
}

func (m *gameMgr) RunGame(s Game, updateDuration int64) {
	go m.startGame(s, updateDuration)
}

func (m *gameMgr) AddGame(s Game) uint16 {
	m.lckGame.Lock()
	defer m.lckGame.Unlock()

	id, _ := m.idGen.GetId()
	gameId := uint16(id)
	s.SetGameID(gameId)

	// id := s.GetGameID()
	m.mapId2Game[gameId] = s
	return gameId
}

func (m *gameMgr) RemoveGame(gameId uint16) {
	m.lckGame.Lock()
	defer m.lckGame.Unlock()

	_, ok := m.mapId2Game[gameId]
	if ok {
		delete(m.mapId2Game, gameId)
	}
}

func (m *gameMgr) GetGame(gameId uint16) (Game, bool) {
	m.lckGame.RLock()
	defer m.lckGame.RUnlock()

	s, ok := m.mapId2Game[gameId]
	return s, ok
}

func (m *gameMgr) AddMatchPlayers(players []Player, gameId uint16) error {
	playerIds := make([]uint64, 0, len(players))
	for _, player := range players {
		playerId := player.GetPlayerId()
		_, ok := m.GetGameByPlayerId(playerId)
		if ok {
			return errors.New("player has join game")
		}

		playerIds = append(playerIds, playerId)
	}

	s, ok := m.GetGame(gameId)
	if !ok {
		return errors.New("game not exist")
	}

	err := s.AddPlayers(players)
	if err != nil {
		return err
	}

	m.addPlayersImpl(playerIds, gameId)
	return nil
}

func (m *gameMgr) JoinPlayer(p Player, gameId uint16) error {
	playerId := p.GetPlayerId()
	_, ok := m.GetGameByPlayerId(playerId)
	if ok {
		return errors.New("player has join game")
	}

	s, ok := m.GetGame(gameId)
	if !ok {
		return errors.New("game not exist")
	}

	e := &PlayerJoinEvent{
		player: p,
	}

	err := s.AddEvent(e)
	if err != nil {
		return err
	}

	m.addPlayerImpl(playerId, gameId)
	return nil
}

func (m *gameMgr) RemovePlayer(playerId uint64) {
	m.removePlayerImpl(playerId)

	s, ok := m.GetGameByPlayerId(playerId)
	if ok {
		e := &PlayerRemoveEvent{
			PlayerId: playerId,
		}

		s.AddEvent(e)
	}
}

func (m *gameMgr) GetGameByPlayerId(playerId uint64) (Game, bool) {
	m.lckPlayer.Lock()
	defer m.lckPlayer.Unlock()

	gameId, ok := m.mapPlayerId2GameId[playerId]
	if !ok {
		return nil, false
	}

	return m.GetGame(gameId)
}

func (m *gameMgr) SetPlayerOnline(playerId uint64) {
	s, ok := m.GetGameByPlayerId(playerId)
	if ok {
		e := &PlayerOnlineEvent{
			PlayerId: playerId,
		}

		s.AddEvent(e)
	}
}

func (m *gameMgr) SetPlayerOffline(playerId uint64) {
	s, ok := m.GetGameByPlayerId(playerId)
	if ok {
		e := &PlayerOfflineEvent{
			PlayerId: playerId,
		}

		s.AddEvent(e)
	}
}

func (m *gameMgr) startGame(s Game, updateDuration int64) {
	err := s.Init()
	if err != nil {
		m.ec.Catch("startGame", &err)
		return
	}

	for {
		m.gameLoop(s, updateDuration)
		if !s.IsAutoRestart() {
			break
		}
	}

	playerIds := s.GetPlayerIds()
	m.removePlayersImpl(playerIds)

	s.Destroy()
}

func (m *gameMgr) gameLoop(s Game, updateDuration int64) {
	var err error = nil
	defer m.ec.Catch("gameLoop", &err)

	lastUpdate := time.Now().UnixNano()
	err = s.Start(lastUpdate)
	if err != nil {
		return
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(updateDuration))

	for {
		<-ticker.C

		s.HandleEvents()
		s.HandleInputs()

		now := time.Now().UnixNano()
		dt := now - lastUpdate
		lastUpdate = now
		s.Update(dt)

		if s.IsGameOver() {
			break
		}
	}

	ticker.Stop()

	s.Stop()
}

func (m *gameMgr) addPlayerImpl(playerId uint64, gameId uint16) {
	m.lckPlayer.Lock()
	defer m.lckPlayer.Unlock()

	m.mapPlayerId2GameId[playerId] = gameId
}

func (m *gameMgr) addPlayersImpl(playerIds []uint64, gameId uint16) {
	m.lckPlayer.Lock()
	defer m.lckPlayer.Unlock()

	for _, playerId := range playerIds {
		m.mapPlayerId2GameId[playerId] = gameId
	}
}

func (m *gameMgr) removePlayerImpl(playerId uint64) {
	m.lckPlayer.Lock()
	defer m.lckPlayer.Unlock()

	_, ok := m.mapPlayerId2GameId[playerId]
	if ok {
		delete(m.mapPlayerId2GameId, playerId)
	}
}

func (m *gameMgr) removePlayersImpl(playerIds []uint64) {
	m.lckPlayer.Lock()
	defer m.lckPlayer.Unlock()

	for _, playerId := range playerIds {
		_, ok := m.mapPlayerId2GameId[playerId]
		if ok {
			delete(m.mapPlayerId2GameId, playerId)
		}
	}
}
