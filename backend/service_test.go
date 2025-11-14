package main

import (
	"testing"
)

func TestService_Move(t *testing.T) {
	tests := []struct {
		name         string
		direction    Direction
		initialX     int
		initialY     int
		expectedX    int
		expectedY    int
		expectError  bool
		errorMessage string
	}{
		{
			name:        "move up",
			direction:   Up,
			initialX:    1,
			initialY:    1,
			expectedX:   1,
			expectedY:   0,
			expectError: false,
		},
		{
			name:        "move down",
			direction:   Down,
			initialX:    1,
			initialY:    1,
			expectedX:   1,
			expectedY:   2,
			expectError: false,
		},
		{
			name:        "move left",
			direction:   Left,
			initialX:    1,
			initialY:    1,
			expectedX:   0,
			expectedY:   1,
			expectError: false,
		},
		{
			name:        "move right",
			direction:   Right,
			initialX:    1,
			initialY:    1,
			expectedX:   2,
			expectedY:   1,
			expectError: false,
		},
		{
			name:         "move out of bounds up",
			direction:    Up,
			initialX:     1,
			initialY:     0,
			expectedX:    1,
			expectedY:    0,
			expectError:  true,
			errorMessage: "cannot move further in that direction",
		},
		{
			name:         "move out of bounds down",
			direction:    Down,
			initialX:     1,
			initialY:     2,
			expectedX:    1,
			expectedY:    2,
			expectError:  true,
			errorMessage: "cannot move further in that direction",
		},
		{
			name:         "move out of bounds left",
			direction:    Left,
			initialX:     0,
			initialY:     1,
			expectedX:    0,
			expectedY:    1,
			expectError:  true,
			errorMessage: "cannot move further in that direction",
		},
		{
			name:         "move out of bounds right",
			direction:    Right,
			initialX:     2,
			initialY:     1,
			expectedX:    2,
			expectedY:    1,
			expectError:  true,
			errorMessage: "cannot move further in that direction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			ds.State.Robot.PositionX = tt.initialX
			ds.State.Robot.PositionY = tt.initialY

			svc := NewService(ds)
			state, err := svc.Move(tt.direction)

			if tt.expectError {
				if err.Error() != tt.errorMessage {
					t.Fatalf("expected error '%s', got '%v'", tt.errorMessage, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if state.Robot.PositionX != tt.expectedX || state.Robot.PositionY != tt.expectedY {
				t.Fatalf("expected position (%d,%d), got (%d,%d)",
					tt.expectedX, tt.expectedY, state.Robot.PositionX, state.Robot.PositionY)
			}
		})
	}
}

func TestService_Pick(t *testing.T) {
	redCircle := Red

	tests := []struct {
		name           string
		setupFunc      func(*DataStore)
		expectedCircle *Circle
		expectError    bool
		errorMessage   string
	}{
		{
			name: "pick circle from position",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 2
				ds.State.Robot.PositionY = 2
				ds.State.Grid[2][2] = []Circle{redCircle}
			},
			expectedCircle: &redCircle,
			expectError:    false,
		},
		{
			name: "pick from empty stack",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 2
				ds.State.Robot.PositionY = 2
				ds.State.Grid[2][2] = []Circle{}
			},
			expectedCircle: nil,
			expectError:    true,
			errorMessage:   "no circles to pick up",
		},
		{
			name: "pick when already holding",
			setupFunc: func(ds *DataStore) {
				circle := Red
				ds.State.Robot.Holding = &circle
			},
			expectedCircle: nil,
			expectError:    true,
			errorMessage:   "already holding a circle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			tt.setupFunc(ds)

			svc := NewService(ds)

			state, err := svc.Pick()

			if tt.expectError {
				if err.Error() != tt.errorMessage {
					t.Fatalf("expected error '%s', got '%v'", tt.errorMessage, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedCircle == nil {
				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing, got %v", state.Robot.Holding)
				}
			} else {
				if *state.Robot.Holding != *tt.expectedCircle {
					t.Fatalf("expected robot to hold %v, got %v", *tt.expectedCircle, *state.Robot.Holding)
				}
			}
		})
	}
}

func TestService_Drop(t *testing.T) {
	redCircle := Red

	tests := []struct {
		name         string
		setupFunc    func(*DataStore)
		assertState  func(*testing.T, State)
		expectError  bool
		errorMessage string
	}{
		{
			name: "drop on empty stack",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &redCircle
				ds.State.Grid[1][1] = []Circle{}
			},
			assertState: func(t *testing.T, state State) {
				if len(state.Grid[1][1]) != 1 || state.Grid[1][1][0] != Red {
					t.Fatalf("expected red circle on stack, got %v", state.Grid[1][1])
				}

				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing after drop, got %v", state.Robot.Holding)
				}
			},
			expectError: false,
		},
		{
			name: "drop when not holding",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = nil
			},
			assertState:  func(t *testing.T, state State) {},
			expectError:  true,
			errorMessage: "not holding any circle to drop",
		},
		{
			name: "drop blue on red circle",
			setupFunc: func(ds *DataStore) {
				blueCircle := Blue
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &blueCircle
				ds.State.Grid[1][1] = []Circle{Red}
			},
			assertState:  func(t *testing.T, state State) {},
			expectError:  true,
			errorMessage: "cannot drop circle here due to stacking rules",
		},
		{
			name: "drop green on red circle",
			setupFunc: func(ds *DataStore) {
				greenCircle := Green
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &greenCircle
				ds.State.Grid[1][1] = []Circle{Red}
			},
			assertState:  func(t *testing.T, state State) {},
			expectError:  true,
			errorMessage: "cannot drop circle here due to stacking rules",
		},
		{
			name: "drop red on red circle",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &redCircle
				ds.State.Grid[1][1] = []Circle{Red}
			},
			assertState:  func(t *testing.T, state State) {},
			expectError:  true,
			errorMessage: "cannot drop circle here due to stacking rules",
		},
		{
			name: "drop green on blue circle",
			setupFunc: func(ds *DataStore) {
				greenCircle := Green
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &greenCircle
				ds.State.Grid[1][1] = []Circle{Blue}
			},
			assertState:  func(t *testing.T, state State) {},
			expectError:  true,
			errorMessage: "cannot drop circle here due to stacking rules",
		},
		{
			name: "drop red on blue circle",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &redCircle
				ds.State.Grid[1][1] = []Circle{Blue}
			},
			assertState: func(t *testing.T, state State) {
				if len(state.Grid[1][1]) != 2 || state.Grid[1][1][1] != Red {
					t.Fatalf("expected red circle on top of blue, got %v", state.Grid[1][1])
				}

				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing after drop, got %v", state.Robot.Holding)
				}
			},
			expectError: false,
		},
		{
			name: "drop blue on blue circle",
			setupFunc: func(ds *DataStore) {
				blueCircle := Blue
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &blueCircle
				ds.State.Grid[1][1] = []Circle{Blue}
			},
			assertState:  func(t *testing.T, state State) {},
			expectError:  true,
			errorMessage: "cannot drop circle here due to stacking rules",
		},
		{
			name: "drop blue on green circle",
			setupFunc: func(ds *DataStore) {
				blueCircle := Blue
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &blueCircle
				ds.State.Grid[1][1] = []Circle{Green}
			},
			assertState: func(t *testing.T, state State) {
				if len(state.Grid[1][1]) != 2 || state.Grid[1][1][1] != Blue {
					t.Fatalf("expected blue circle on top, got %v", state.Grid[1][1])
				}

				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing after drop, got %v", state.Robot.Holding)
				}
			},
			expectError: false,
		},
		{
			name: "drop green on green circle",
			setupFunc: func(ds *DataStore) {
				greenCircle := Green
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &greenCircle
				ds.State.Grid[1][1] = []Circle{Green}
			},
			assertState: func(t *testing.T, state State) {
				if len(state.Grid[1][1]) != 2 || state.Grid[1][1][1] != Green {
					t.Fatalf("expected green circle on top, got %v", state.Grid[1][1])
				}

				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing after drop, got %v", state.Robot.Holding)
				}
			},
			expectError: false,
		},
		{
			name: "drop red on green circle",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				ds.State.Robot.Holding = &redCircle
				ds.State.Grid[1][1] = []Circle{Green}
			},
			assertState: func(t *testing.T, state State) {
				if len(state.Grid[1][1]) != 2 || state.Grid[1][1][1] != Red {
					t.Fatalf("expected red circle on top, got %v", state.Grid[1][1])
				}

				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing after drop, got %v", state.Robot.Holding)
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			tt.setupFunc(ds)

			svc := NewService(ds)

			state, err := svc.Drop()

			tt.assertState(t, state)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestService_GetState(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(*testing.T, State)
	}{
		{
			name: "get state successfully",
			validateFunc: func(t *testing.T, state State) {
				if state.Robot.PositionX != 0 || state.Robot.PositionY != 0 {
					t.Fatalf("expected robot at (0,0), got (%d,%d)",
						state.Robot.PositionX, state.Robot.PositionY)
				}
				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing initially, got %v", state.Robot.Holding)
				}

				expectedGrid := [GridSize][GridSize][]Circle{
					{{Red}, {Green}, {Green}},
					{{Blue}, {Red}, {Blue}},
					{{Green}, {Blue}, {Red}},
				}
				for x := 0; x < GridSize; x++ {
					for y := 0; y < GridSize; y++ {
						expectedStack := expectedGrid[x][y]
						actualStack := state.Grid[x][y]
						if len(expectedStack) != len(actualStack) {
							t.Fatalf("expected stack length %d at (%d,%d), got %d",
								len(expectedStack), x, y, len(actualStack))
						}
						for i := range expectedStack {
							if expectedStack[i] != actualStack[i] {
								t.Fatalf("expected circle %v at index %d at (%d,%d), got %v",
									expectedStack[i], i, x, y, actualStack[i])
							}
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			svc := NewService(ds)
			state := svc.GetState()
			tt.validateFunc(t, state)
		})
	}
}

func TestService_GetHistory(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(*Service)
		validateFunc func(*testing.T, []MovementHistory)
	}{
		{
			name:      "empty history initially",
			setupFunc: func(svc *Service) {},
			validateFunc: func(t *testing.T, history []MovementHistory) {
				if len(history) != 0 {
					t.Fatalf("expected empty history, got length %d", len(history))
				}
			},
		},
		{
			name: "history after one move",
			setupFunc: func(svc *Service) {
				svc.Move(Right)
			},
			validateFunc: func(t *testing.T, history []MovementHistory) {
				if history[0].Moves != "Moved right" {
					t.Fatalf("expected 'Moved right', got '%s'", history[0].Moves)
				}
			},
		},
		{
			name: "history after pick, move and drop",
			setupFunc: func(svc *Service) {
				svc.Pick()
				svc.Move(Down)
				svc.Drop()
			},
			validateFunc: func(t *testing.T, history []MovementHistory) {
				expected := []string{
					"Picked up a red circle",
					"Moved down",
					"Dropped a red circle",
				}
				for i, exp := range expected {
					if history[i].Moves != exp {
						t.Fatalf("expected '%s' at index %d, got '%s'", exp, i, history[i].Moves)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			svc := NewService(ds)

			tt.setupFunc(svc)

			history := svc.GetHistory()

			tt.validateFunc(t, history)
		})
	}
}
