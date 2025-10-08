package processor

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/drewlarso/chip8-emulator/display"
	"github.com/drewlarso/chip8-emulator/keyboard"
)

var fontset [80]byte = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Processor struct {
	memory        [4096]byte
	registers     [16]byte
	indexRegister uint16
	delayTimer    byte
	soundTimer    byte
	pc            uint16
	sp            byte
	stack         [16]uint16
}

func NewProcessor() *Processor {
	memory := [4096]byte{}

	copy(memory[0x50:], fontset[:])

	return &Processor{
		memory:        memory,
		registers:     [16]byte{},
		indexRegister: 0,
		delayTimer:    0,
		soundTimer:    0,
		pc:            0x200, // start of rom in memory
		sp:            0,
		stack:         [16]uint16{},
	}
}

func (p *Processor) LoadROM(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	copy(p.memory[0x200:], content)
}

func (p *Processor) Cycle(display *display.Display, keyboard *keyboard.Keyboard) {
	// instructions are 2 bytes long
	// combine memory[pc] and memory[pc+1] into a 16 bit value
	// shift the first one left 8 to add empty space to the right of the values
	// OR with the next value to combine the two
	// 0xff0000 OR 0x0000ff = 0xff00ff
	instruction := uint16(p.memory[p.pc])<<8 | uint16(p.memory[p.pc+1])

	switch instruction & 0xF000 {
	case 0x0000:
		switch instruction {
		case 0x00E0:
			// clear the screen
			display.Clear()
			p.pc += 2
		case 0x00EE:
			// return from a subroutine
			p.sp--
			p.pc = p.stack[p.sp]
			p.pc += 2
		}
	case 0x1000:
		// 1nnn
		// jump to addr nnn
		// this AND operation basically chops off the first part of our instruction
		p.pc = instruction & 0x0FFF
	case 0x2000:
		// 2nnn
		// call addr
		p.stack[p.sp] = p.pc
		p.sp++
		p.pc = instruction & 0x0FFF
	case 0x3000:
		// 3xkk
		// skips the next instruction if Vx == kk
		x := (instruction >> 8 & 0x000F)
		kk := (instruction & 0x00FF)
		if p.registers[x] == byte(kk) {
			p.pc += 2
		}
		p.pc += 2
	case 0x4000:
		// 4xkk
		// skips the next instruction if Vx != kk
		x := (instruction >> 8 & 0x000F)
		kk := (instruction & 0x00FF)
		if p.registers[x] != byte(kk) {
			p.pc += 2
		}
		p.pc += 2
	case 0x5000:
		// 5xy0
		// skips the next instruction if Vx == Vy
		x := (instruction >> 8 & 0x000F)
		y := (instruction >> 4 & 0x000F)
		if p.registers[x] == p.registers[y] {
			p.pc += 2
		}
		p.pc += 2
	case 0x6000:
		// 6xkk
		// sets Vx to kk
		x := (instruction >> 8 & 0x000F)
		kk := (instruction & 0x00FF)
		p.registers[x] = byte(kk)
		p.pc += 2
	case 0x7000:
		// 7xkk
		// adds kk to Vx
		x := (instruction >> 8 & 0x000F)
		kk := (instruction & 0x00FF)
		p.registers[x] += byte(kk)
		p.pc += 2
	case 0x8000:
		switch instruction & 0x000F {
		case 0x0000:
			// 8xy0
			// set Vx = Vy
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			p.registers[x] = p.registers[y]
			p.pc += 2
		case 0x0001:
			// 8xy1
			// set Vx = Vx OR Vy
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			p.registers[x] = p.registers[x] | p.registers[y]
			p.pc += 2
		case 0x0002:
			// 8xy2
			// set Vx = Vx AND Vy
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			p.registers[x] = p.registers[x] & p.registers[y]
			p.pc += 2
		case 0x0003:
			// 8xy3
			// set Vx = Vx XOR Vy
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			p.registers[x] = p.registers[x] ^ p.registers[y]
			p.pc += 2
		case 0x0004:
			// 8xy4
			// set Vx = Vx + Vy, set Vf = carry
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			var value uint16 = uint16(p.registers[x]) + uint16(p.registers[y])
			if value >= 255 {
				// set Vf to 1
				p.registers[15] = 1
			} else {
				p.registers[15] = 0
			}
			p.registers[x] = byte(value)
			p.pc += 2
		case 0x0005:
			// 8xy5
			// set Vx = Vx - Vy, set VF = NOT borrow
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			if p.registers[x] > p.registers[y] {
				p.registers[15] = 1
			} else {
				p.registers[15] = 0
			}
			p.registers[x] = p.registers[x] - p.registers[y]
			p.pc += 2
		case 0x0006:
			// 8xy6
			// set Vx = Vx SHR 1
			x := (instruction >> 8 & 0x000F)

			// chops off everything except the smallest bit
			p.registers[15] = p.registers[x] & 0x01

			// same as dividing by 2
			p.registers[x] = p.registers[x] >> 1
			p.pc += 2
		case 0x0007:
			// 8xy7
			// set Vx = Vy - Vx, set VF = NOT borrow
			x := (instruction >> 8 & 0x000F)
			y := (instruction >> 4 & 0x000F)
			if p.registers[y] > p.registers[x] {
				p.registers[15] = 1
			} else {
				p.registers[15] = 0
			}
			p.registers[x] = p.registers[y] - p.registers[x]
			p.pc += 2
		case 0x000E:
			// 8xyE
			// set Vx = Vx SHL 1
			x := (instruction >> 8 & 0x000F)

			// chops off everything except the last bit
			p.registers[15] = p.registers[x] >> 7

			// same as multiplying by 2
			p.registers[x] = p.registers[x] << 1
			p.pc += 2
		}
	case 0x9000:
		// 9xy0
		// skip the next instruction if Vx != Vy
		x := (instruction >> 8 & 0x000F)
		y := (instruction >> 4 & 0x000F)
		if p.registers[x] != p.registers[y] {
			p.pc += 2
		}
		p.pc += 2
	case 0xA000:
		// Annn
		// set I = nnn
		nnn := instruction & 0x0FFF
		p.indexRegister = nnn
		p.pc += 2
	case 0xB000:
		// Bnnn
		// jump to location nnn + V0
		nnn := instruction & 0x0FFF
		p.pc = nnn + uint16(p.registers[0])
	case 0xC000:
		// Cxkk
		// set Vx = random byte AND kk
		x := (instruction >> 8 & 0x000F)
		kk := instruction & 0x00FF
		value := byte(rand.Intn(256)) & byte(kk)
		p.registers[x] = value
		p.pc += 2
	case 0xD000:
		// Dxyn
		// display n byte sprite at (Vx, Vy), set VF = collision
		// read n bytes from memory starting at i
		x := (instruction >> 8 & 0x000F)
		y := (instruction >> 4 & 0x000F)
		n := (instruction & 0x000F)

		vx := p.registers[x] % 64
		vy := p.registers[y] % 32

		p.registers[15] = 0

		for row := range n {
			var spriteByte byte = p.memory[p.indexRegister+row]
			for col := range 8 {
				// get data from left to right
				// col = 0 means we shift left 7, col = 1 means shift 6
				// we & with 0x01 to isolate the rightmost bit
				bit := spriteByte >> (7 - col) & 0x01
				pixelX := (vx + byte(col)) % 64
				pixelY := (vy + byte(row)) % 32

				index := int(pixelY)*64 + int(pixelX)
				pixelValue := display.GetPixel(index)
				if pixelValue != 0 {
					p.registers[15] = 1
				}
				display.Set(index, pixelValue^bit)
			}
		}
		p.pc += 2
	case 0xE000:
		switch instruction & 0x000F {
		case 0x000E:
			// Ex9E
			// skip next instruction if key with the value of Vx is pressed
			x := (instruction >> 8 & 0x000F)
			if keyboard.IsKeyDown(p.registers[x]) {
				p.pc += 2
			}
			p.pc += 2
		case 0x0001:
			// ExA1
			// skip next instruction if key with the value of Vx is not pressed
			x := (instruction >> 8 & 0x000F)
			if !keyboard.IsKeyDown(p.registers[x]) {
				p.pc += 2
			}
			p.pc += 2
		}
	case 0xF000:
		switch instruction & 0x00FF {
		case 0x0007:
			// Fx07
			// set Vx = delay timer value
			x := (instruction >> 8 & 0x000F)
			p.registers[x] = p.delayTimer
			p.pc += 2
		case 0x000A:
			// Fx0A
			// wait for a key press, store the value of key in Vx
			x := (instruction >> 8 & 0x000F)
			pressed := keyboard.AnyKeyDown()
			if pressed < 16 {
				p.registers[x] = pressed
				p.pc += 2
			}
		case 0x0015:
			// Fx15
			// set delay timer = Vx
			x := (instruction >> 8 & 0x000F)
			p.delayTimer = p.registers[x]
			p.pc += 2
		case 0x0018:
			// Fx18
			// set sound timer = Vx
			x := (instruction >> 8 & 0x000F)
			p.soundTimer = p.registers[x]
			p.pc += 2
		case 0x001E:
			// Fx1E
			// set I = Vx + I
			x := (instruction >> 8 & 0x000F)
			p.indexRegister += uint16(p.registers[x])
			p.pc += 2
		case 0x0029:
			// Fx29
			// set I = location of sprite for digit Vx
			x := (instruction >> 8 & 0x000F)
			p.indexRegister = uint16(p.registers[x])*5 + 0x50
			p.pc += 2
		case 0x0033:
			// Fx33
			// store BCD representation of Vx in memory locations I, I+1, and I+2
			x := (instruction >> 8 & 0x000F)
			vx := p.registers[x]
			p.memory[p.indexRegister] = vx / 100
			p.memory[p.indexRegister+1] = (vx / 10) % 10
			p.memory[p.indexRegister+2] = vx % 10
			p.pc += 2
		case 0x0055:
			// Fx55
			// store registers V0 through Vx in memory starting at location I
			x := (instruction >> 8 & 0x000F)
			for i := range x + 1 {
				p.memory[p.indexRegister+uint16(i)] = p.registers[i]
			}
			p.pc += 2
		case 0x0065:
			// Fx65
			// read registers V0 through Vx from memory starting at location I
			x := (instruction >> 8 & 0x000F)
			for i := range x + 1 {
				p.registers[i] = p.memory[p.indexRegister+uint16(i)]
			}
			p.pc += 2
		}
	default:
		fmt.Printf("%x is not implemented!\n", instruction)
	}
}

func (p *Processor) UpdateTimers() {
	if p.delayTimer != 0 {
		p.delayTimer--
	}
	if p.soundTimer != 0 {
		p.soundTimer--
	}
}
