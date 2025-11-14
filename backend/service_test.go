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
	tests := []struct {
		name         string
		setupFunc    func(*DataStore)
		expectError  bool
		errorMessage string
	}{
		{
			name: "drop on empty stack",
			setupFunc: func(ds *DataStore) {
				ds.State.Robot.PositionX = 1
				ds.State.Robot.PositionY = 1
				circle := Red
				ds.State.Robot.Holding = &circle
				ds.State.Grid[1][1] = []Circle{}
			},
			expectError: false,
		},
		{
			name: "drop red on green stack",
			setupFunc: func(ds *DataStore) {
				// Pick red from 0,0 and move to 0,1 (green stack)
				svc.Pick()
				svc.Move(Down)
			},
			expectError: false,
		},
		{
			name: "drop red on red stack - should fail",
			setupFunc: func(ds *DataStore) {
				// Pick red from 0,0 and move to 1,1 (red stack)
				svc.Pick()
				svc.Move(Right)
				svc.Move(Down)
			},
			expectError: true,
		},
		{
			name: "drop without holding anything",
			setupFunc: func(ds *DataStore) {
				// Don't pick anything
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			svc := NewService(ds)

			tt.setupFunc(ds)

			state, err := svc.Drop()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if state.Robot.Holding != nil {
				t.Fatalf("expected robot to hold nothing after drop, got %v", state.Robot.Holding)
			}
		})
	}
}

func TestService_GetState(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(*Service)
		validateFunc func(*testing.T, State)
	}{
		{
			name: "initial state",
			setupFunc: func(svc *Service) {
				// No setup needed
			},
			validateFunc: func(t *testing.T, state State) {
				if state.Robot.PositionX != 0 || state.Robot.PositionY != 0 {
					t.Fatalf("expected robot at (0,0), got (%d,%d)",
						state.Robot.PositionX, state.Robot.PositionY)
				}
				if state.Robot.Holding != nil {
					t.Fatalf("expected robot to hold nothing initially, got %v", state.Robot.Holding)
				}
			},
		},
		{
			name: "after move",
			setupFunc: func(svc *Service) {
				svc.Move(Right)
			},
			validateFunc: func(t *testing.T, state State) {
				if state.Robot.PositionX != 1 || state.Robot.PositionY != 0 {
					t.Fatalf("expected robot at (1,0) after move right, got (%d,%d)",
						state.Robot.PositionX, state.Robot.PositionY)
				}
			},
		},
		{
			name: "after pick",
			setupFunc: func(svc *Service) {
				svc.Pick()
			},
			validateFunc: func(t *testing.T, state State) {
				if state.Robot.Holding == nil || *state.Robot.Holding != Red {
					t.Fatalf("expected robot to hold Red after pick, got %v", state.Robot.Holding)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDataStore()
			svc := NewService(ds)

			tt.setupFunc(svc)

			state := svc.GetState()
			tt.validateFunc(t, state)
		})
	}
}

func TestService_GetHistory(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(*Service)
		expectedLength int
		validateFunc   func(*testing.T, []MovementHistory)
	}{
		{
			name: "empty history initially",
			setupFunc: func(svc *Service) {
				// No actions
			},
			expectedLength: 0,
			validateFunc: func(t *testing.T, history []MovementHistory) {
				// Nothing to validate for empty history
			},
		},
		{
			name: "history after one move",
			setupFunc: func(svc *Service) {
				svc.Move(Right)
			},
			expectedLength: 1,
			validateFunc: func(t *testing.T, history []MovementHistory) {
				if history[0].Moves != "Moved right" {
					t.Fatalf("expected 'Moved right', got '%s'", history[0].Moves)
				}
			},
		},
		{
			name: "history after pick and drop",
			setupFunc: func(svc *Service) {
				svc.Pick()
				svc.Move(Down)
				svc.Drop()
			},
			expectedLength: 3,
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
		{
			name: "history after multiple moves",
			setupFunc: func(svc *Service) {
				svc.Move(Right)
				svc.Move(Down)
				svc.Move(Left)
			},
			expectedLength: 3,
			validateFunc: func(t *testing.T, history []MovementHistory) {
				expected := []string{
					"Moved right",
					"Moved down",
					"Moved left",
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

			if len(history) != tt.expectedLength {
				t.Fatalf("expected history length %d, got %d", tt.expectedLength, len(history))
			}

			tt.validateFunc(t, history)
		})
	}
}
