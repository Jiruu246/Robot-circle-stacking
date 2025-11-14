package main

import (
	"errors"
	"fmt"
	"time"
)

type Service struct {
	storage *DataStore
}

func NewService(storage *DataStore) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetState() State {
	s.storage.Mu.Lock()
	defer s.storage.Mu.Unlock()
	return s.storage.State
}

func (s *Service) GetHistory() []MovementHistory {
	s.storage.Mu.Lock()
	defer s.storage.Mu.Unlock()
	return s.storage.History
}

func (s *Service) Move(direction Direction) (State, error) {
	s.storage.Mu.Lock()
	defer s.storage.Mu.Unlock()

	robot := &s.storage.State.Robot

	new_x, new_y := robot.PositionX, robot.PositionY
	switch direction {
	case Up:
		new_y--
	case Down:
		new_y++
	case Left:
		new_x--
	case Right:
		new_x++
	default:
		return s.storage.State, nil
	}

	if outOfBounds(new_x, new_y) {
		return State{}, errors.New("cannot move further in that direction")
	}

	robot.PositionX, robot.PositionY = new_x, new_y
	s.storage.History = append(s.storage.History, MovementHistory{
		Timestamp: time.Now(),
		Moves:     fmt.Sprintf("Moved %s", direction),
	})

	return s.storage.State, nil
}

func (s *Service) Pick() (State, error) {
	s.storage.Mu.Lock()
	defer s.storage.Mu.Unlock()

	robot := &s.storage.State.Robot
	if robot.Holding != nil {
		return State{}, errors.New("already holding a circle")
	}

	stack := s.storage.State.Grid[robot.PositionX][robot.PositionY]
	if len(stack) == 0 {
		return State{}, errors.New("no circles to pick up")
	}

	robot.Holding = &stack[len(stack)-1]
	s.storage.State.Grid[robot.PositionX][robot.PositionY] = stack[:len(stack)-1]
	s.storage.History = append(s.storage.History, MovementHistory{
		Timestamp: time.Now(),
		Moves:     fmt.Sprintf("Picked up a %s circle", *robot.Holding),
	})

	return s.storage.State, nil
}

func (s *Service) Drop() (State, error) {
	s.storage.Mu.Lock()
	defer s.storage.Mu.Unlock()

	robot := &s.storage.State.Robot
	if robot.Holding == nil {
		return State{}, errors.New("not holding any circle to drop")
	}

	stack := s.storage.State.Grid[robot.PositionX][robot.PositionY]

	if !canDropCircle(stack, *robot.Holding) {
		return State{}, errors.New("cannot drop circle here due to stacking rules")
	}

	dropped := *robot.Holding
	s.storage.State.Grid[robot.PositionX][robot.PositionY] = append(stack, dropped)
	robot.Holding = nil
	s.storage.History = append(s.storage.History, MovementHistory{
		Timestamp: time.Now(),
		Moves:     fmt.Sprintf("Dropped a %s circle", dropped),
	})

	return s.storage.State, nil
}

func outOfBounds(x int, y int) bool {
	return x < 0 || x >= GridSize || y < 0 || y >= GridSize
}

func canDropCircle(stack []Circle, circle Circle) bool {
	if len(stack) == 0 {
		return true
	}

	top := stack[len(stack)-1]

	switch top {
	case Red:
		return false
	case Green:
		return true
	case Blue:
		return circle == Red
	}

	return false
}
