package main

import (
	"sync"
	"time"
)

type Circle string

const (
	Red   Circle = "red"
	Green Circle = "green"
	Blue  Circle = "blue"
)

type Action string

const (
	PickUp Action = "pick_up"
	Drop   Action = "drop"
	Move   Action = "move"
)

type Direction string

const (
	Up    Direction = "up"
	Down  Direction = "down"
	Left  Direction = "left"
	Right Direction = "right"
)

type Robot struct {
	PositionX int
	PositionY int
	Holding   *Circle
}

type State struct {
	Robot Robot
	Grid  [GridSize][GridSize][]Circle
}

type MovementHistory struct {
	Timestamp time.Time
	Moves     string
}

type DataStore struct {
	Mu      sync.Mutex
	State   State
	History []MovementHistory
}

func NewDataStore() *DataStore {
	grid := [GridSize][GridSize][]Circle{
		{{Red}, {Green}, {Green}},
		{{Blue}, {Red}, {Blue}},
		{{Green}, {Blue}, {Red}},
	}

	return &DataStore{
		State: State{
			Robot: Robot{
				PositionX: 0,
				PositionY: 0,
				Holding:   nil,
			},
			Grid: grid,
		},
		History: []MovementHistory{},
	}
}
