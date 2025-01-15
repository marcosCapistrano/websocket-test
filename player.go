package main

type Player struct {
	ID           string
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2
	input        *InputState
}

func (p *Player) step() {
	if (*p.input)[UP] == PRESSED {
		p.Acceleration.Y++
	} else if (*p.input)[DOWN] == PRESSED {
		p.Acceleration.Y--
	} else if (*p.input)[RIGHT] == PRESSED {
		p.Acceleration.X++
	} else if (*p.input)[LEFT] == PRESSED {
		p.Acceleration.X--
	}

	p.Velocity.X += p.Acceleration.X
	p.Velocity.Y += p.Acceleration.Y

	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y

	p.Acceleration.X = 0
	p.Acceleration.Y = 0
}
