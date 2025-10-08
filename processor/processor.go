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
	memory     [4096]byte
	registers  [16]byte
	ir         uint16
	delayTimer byte
	soundTimer byte
	pc         uint16
	sp         byte
	stack      [16]uint16
}

func NewProcessor() *Processor {
	memory := [4096]byte{}

	copy(memory[0x50:], fontset[:])

	return &Processor{
		memory:     memory,
		registers:  [16]byte{},
		ir:         0,
		delayTimer: 0,
		soundTimer: 0,
		pc:         0x200, // start of rom in memory
		sp:         0,
		stack:      [16]uint16{},
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

	x := (instruction >> 8 & 0x000F)
	y := (instruction >> 4 & 0x000F)
	kk := (instruction & 0x00FF)

	switch instruction & 0xF000 {
	case 0x0000:
		switch instruction {
		case 0x00E0:
			// 00E0 - CLS
			// Clear the display.
			display.Clear()
			p.pc += 2
		case 0x00EE:
			// 00EE - RET
			// Return from a subroutine.
			p.sp--
			p.pc = p.stack[p.sp]
			p.pc += 2
		}
	case 0x1000:
		// 1nnn - JP addr
		// Jump to location nnn.
		p.pc = instruction & 0x0FFF
	case 0x2000:
		// 2nnn - CALL addr
		// Call subroutine at nnn.
		p.stack[p.sp] = p.pc
		p.sp++
		p.pc = instruction & 0x0FFF
	case 0x3000:
		// 3xkk - SE Vx, byte
		// Skip next instruction if Vx = kk.
		if p.registers[x] == byte(kk) {
			p.pc += 2
		}
		p.pc += 2
	case 0x4000:
		// 4xkk - SNE Vx, byte
		// Skip next instruction if Vx != kk.
		if p.registers[x] != byte(kk) {
			p.pc += 2
		}
		p.pc += 2
	case 0x5000:
		// 5xy0 - SE Vx, Vy
		// Skip next instruction if Vx = Vy.
		if p.registers[x] == p.registers[y] {
			p.pc += 2
		}
		p.pc += 2
	case 0x6000:
		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		kk := (instruction & 0x00FF)
		p.registers[x] = byte(kk)
		p.pc += 2
	case 0x7000:
		// 7xkk - ADD Vx, byte
		// Set Vx = Vx + kk.
		kk := (instruction & 0x00FF)
		p.registers[x] += byte(kk)
		p.pc += 2
	case 0x8000:
		switch instruction & 0x000F {
		case 0x0000:
			// 8xy0 - LD Vx, Vy
			// Set Vx = Vy.
			p.registers[x] = p.registers[y]
			p.pc += 2
		case 0x0001:
			// 8xy1 - OR Vx, Vy
			// Set Vx = Vx OR Vy.
			p.registers[x] = p.registers[x] | p.registers[y]
			p.pc += 2
		case 0x0002:
			// 8xy2 - AND Vx, Vy
			// Set Vx = Vx AND Vy.
			p.registers[x] = p.registers[x] & p.registers[y]
			p.pc += 2
		case 0x0003:
			// 8xy3 - XOR Vx, Vy
			// Set Vx = Vx XOR Vy.
			p.registers[x] = p.registers[x] ^ p.registers[y]
			p.pc += 2
		case 0x0004:
			// 8xy4 - ADD Vx, Vy
			// Set Vx = Vx + Vy, set VF = carry.
			var value uint16 = uint16(p.registers[x]) + uint16(p.registers[y])
			if value >= 255 {
				p.registers[15] = 1
			} else {
				p.registers[15] = 0
			}
			p.registers[x] = byte(value)
			p.pc += 2
		case 0x0005:
			// 8xy5 - SUB Vx, Vy
			// Set Vx = Vx - Vy, set VF = NOT borrow.
			if p.registers[x] > p.registers[y] {
				p.registers[15] = 1
			} else {
				p.registers[15] = 0
			}
			p.registers[x] = p.registers[x] - p.registers[y]
			p.pc += 2
		case 0x0006:
			// 8xy6 - SHR Vx {, Vy}
			// Set Vx = Vx SHR 1.

			// chops off everything except the smallest bit
			p.registers[15] = p.registers[x] & 0x01

			// same as dividing by 2
			p.registers[x] = p.registers[x] >> 1
			p.pc += 2
		case 0x0007:
			// 8xy7 - SUBN Vx, Vy
			// Set Vx = Vy - Vx, set VF = NOT borrow.
			if p.registers[y] > p.registers[x] {
				p.registers[15] = 1
			} else {
				p.registers[15] = 0
			}
			p.registers[x] = p.registers[y] - p.registers[x]
			p.pc += 2
		case 0x000E:
			// 8xyE - SHL Vx {, Vy}
			// Set Vx = Vx SHL 1.

			// chops off everything except the last bit
			p.registers[15] = p.registers[x] >> 7

			// same as multiplying by 2
			p.registers[x] = p.registers[x] << 1
			p.pc += 2
		}
	case 0x9000:
		// 9xy0 - SNE Vx, Vy
		// Skip next instruction if Vx != Vy.
		if p.registers[x] != p.registers[y] {
			p.pc += 2
		}
		p.pc += 2
	case 0xA000:
		// Annn - LD I, addr
		// Set I = nnn.
		nnn := instruction & 0x0FFF
		p.ir = nnn
		p.pc += 2
	case 0xB000:
		// Bnnn - JP V0, addr
		// Jump to location nnn + V0.
		nnn := instruction & 0x0FFF
		p.pc = nnn + uint16(p.registers[0])
	case 0xC000:
		// Cxkk - RND Vx, byte
		// Set Vx = random byte AND kk.
		kk := instruction & 0x00FF
		value := byte(rand.Intn(256)) & byte(kk)
		p.registers[x] = value
		p.pc += 2
	case 0xD000:
		// Dxyn - DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		n := (instruction & 0x000F)

		vx := p.registers[x] % 64
		vy := p.registers[y] % 32

		p.registers[15] = 0

		for row := range n {
			var spriteByte byte = p.memory[p.ir+row]
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
			// Ex9E - SKP Vx
			// Skip next instruction if key with the value of Vx is pressed.
			if keyboard.IsKeyDown(p.registers[x]) {
				p.pc += 2
			}
			p.pc += 2
		case 0x0001:
			// ExA1 - SKNP Vx
			// Skip next instruction if key with the value of Vx is not pressed.
			if !keyboard.IsKeyDown(p.registers[x]) {
				p.pc += 2
			}
			p.pc += 2
		}
	case 0xF000:
		switch instruction & 0x00FF {
		case 0x0007:
			// Fx07 - LD Vx, DT
			// Set Vx = delay timer value.
			p.registers[x] = p.delayTimer
			p.pc += 2
		case 0x000A:
			// Fx0A - LD Vx, K
			// Wait for a key press, store the value of the key in Vx.
			pressed := keyboard.AnyKeyDown()
			if pressed < 16 {
				p.registers[x] = pressed
				p.pc += 2
			}
		case 0x0015:
			// Fx15 - LD DT, Vx
			// Set delay timer = Vx.
			p.delayTimer = p.registers[x]
			p.pc += 2
		case 0x0018:
			// Fx18 - LD ST, Vx
			// Set sound timer = Vx.
			p.soundTimer = p.registers[x]
			p.pc += 2
		case 0x001E:
			// Fx1E - ADD I, Vx
			// Set I = I + Vx.
			p.ir += uint16(p.registers[x])
			p.pc += 2
		case 0x0029:
			// Fx29 - LD F, Vx
			// Set I = location of sprite for digit Vx.
			p.ir = uint16(p.registers[x])*5 + 0x50
			p.pc += 2
		case 0x0033:
			// Fx33 - LD B, Vx
			// Store BCD representation of Vx in memory locations I, I+1, and I+2.
			vx := p.registers[x]
			p.memory[p.ir] = vx / 100
			p.memory[p.ir+1] = (vx / 10) % 10
			p.memory[p.ir+2] = vx % 10
			p.pc += 2
		case 0x0055:
			// Fx55 - LD [I], Vx
			// Store registers V0 through Vx in memory starting at location I.
			for i := range x + 1 {
				p.memory[p.ir+uint16(i)] = p.registers[i]
			}
			p.pc += 2
		case 0x0065:
			// Fx65 - LD Vx, [I]
			// Read registers V0 through Vx from memory starting at location I.
			for i := range x + 1 {
				p.registers[i] = p.memory[p.ir+uint16(i)]
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
