package main

type CommandRequest struct {
	Action    Action    `json:"action"`
	Direction Direction `json:"direction,omitempty"`
}

type StateResponse struct {
	PositionX int                          `json:"position_x"`
	PositionY int                          `json:"position_y"`
	Holding   *Circle                      `json:"holding,omitempty"`
	Grid      [GridSize][GridSize][]Circle `json:"grid"`
	Won       bool                         `json:"won"`
}
