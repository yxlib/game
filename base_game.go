// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import (
	"errors"

	"github.com/yxlib/yx"
)

type BaseGame struct {
	gameId     uint16
	templateId uint16
	stage      int
	startTime  int64
	curTime    int64
	endTime    int64

	queInputs yx.Queue
	queEvents yx.Queue

	inputProc GameInputProcessor
	eventProc GameEventProcessor
}

func NewBaseGame(templateId uint16, inputProc GameInputProcessor, eventProc GameEventProcessor) *BaseGame {
	return &BaseGame{
		gameId:     0,
		templateId: templateId,
		stage:      0,
		startTime:  0,
		curTime:    0,
		endTime:    0,
		queInputs:  yx.NewSyncLinkedQueue(),
		queEvents:  yx.NewSyncLinkedQueue(),
		inputProc:  inputProc,
		eventProc:  eventProc,
	}
}

func (s *BaseGame) SetGameID(gameId uint16) {
	s.gameId = gameId
}

func (s *BaseGame) GetGameID() uint16 {
	return s.gameId
}

func (s *BaseGame) GetTemplateID() uint16 {
	return s.templateId
}

func (s *BaseGame) AddInput(playerId uint64, cmd string, params ...interface{}) error {
	input := &GameInput{
		PlayerId: playerId,
		Cmd:      cmd,
		Params:   params,
	}

	s.queInputs.Enqueue(input)
	return nil
}

func (s *BaseGame) HandleInputs() {
	for {
		if s.queInputs.GetSize() == 0 {
			break
		}

		item, err := s.queInputs.Dequeue()
		if err != nil {
			continue
		}

		input := item.(*GameInput)
		s.inputProc.HandleGameInput(input.PlayerId, input.Cmd, input.Params)
	}
}

func (s *BaseGame) AddEvent(e GameEvent) error {
	if e == nil {
		return errors.New("event is nil")
	}

	s.queEvents.Enqueue(e)
	return nil
}

func (s *BaseGame) HandleEvents() {
	for {
		if s.queEvents.GetSize() == 0 {
			break
		}

		item, err := s.queEvents.Dequeue()
		if err != nil {
			continue
		}

		evt := item.(GameEvent)
		s.eventProc.HandleGameEvent(evt)
	}
}

func (s *BaseGame) SetStage(stage int) {
	s.stage = stage
}

func (s *BaseGame) GetStage() int {
	return s.stage
}

func (s *BaseGame) SetStartTime(startTime int64) {
	s.startTime = startTime
}

func (s *BaseGame) GetStartTime() int64 {
	return s.startTime
}

func (s *BaseGame) SetCurTime(curTime int64) {
	s.curTime = curTime
}

func (s *BaseGame) AddCurTime(dt int64) int64 {
	s.curTime += dt
	return s.curTime
}

func (s *BaseGame) GetCurTime() int64 {
	return s.curTime
}

func (s *BaseGame) SetEndTime(endTime int64) {
	s.endTime = endTime
}

func (s *BaseGame) GetEndTime() int64 {
	return s.endTime
}
