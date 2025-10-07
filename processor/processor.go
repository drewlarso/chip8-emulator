package processor

import (
	"github.com/drewlarso/chip8-emulator/display"
	"github.com/drewlarso/chip8-emulator/keyboard"
)

type Processor struct {
	registers     [16]byte
	indexRegister uint16
	delayTimer    byte
	soundTimer    byte
	pc            uint16
	sp            byte
	stack         [16]uint16
}

func NewProcessor() *Processor {
	return &Processor{
		registers:     [16]byte{},
		indexRegister: 0,
		delayTimer:    255,
		soundTimer:    255,
		pc:            0,
		sp:            0,
		stack:         [16]uint16{},
	}
}

func (p *Processor) Cycle(display *display.Display, keyboard *keyboard.Keyboard) {}

func (p *Processor) UpdateTimers() {
	if p.delayTimer != 0 {
		p.delayTimer--
	}
	if p.soundTimer != 0 {
		p.soundTimer--
	}
}
