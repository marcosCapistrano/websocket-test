package main

type InputKey int

const (
	LEFT InputKey = iota
	RIGHT
	UP
	DOWN
)

type KeyState int

const (
	UNPRESSED KeyState = iota
	PRESSED
)

type InputState map[InputKey]KeyState

func NewInputState() *InputState {
	inputState := make(InputState)
	inputState[LEFT] = UNPRESSED
	inputState[RIGHT] = UNPRESSED
	inputState[UP] = UNPRESSED
	inputState[DOWN] = UNPRESSED

	return &inputState
}

func (input *InputState) Update(key InputKey, state KeyState) {
	(*input)[key] = state
}
