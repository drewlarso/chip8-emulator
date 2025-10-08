package display

import (
	"github.com/drewlarso/chip8-emulator/constants"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Display struct {
	buffer [64 * 32]byte
}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) Clear() {
	d.buffer = [64 * 32]byte{}
}

func (d *Display) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(0x0201127ff))
	for i := range len(d.buffer) {
		x := i % 64
		y := i / 64
		if d.buffer[i] != 0 {
			rl.DrawRectangle(
				int32(x*constants.Scale),
				int32(y*constants.Scale),
				int32(constants.Scale),
				int32(constants.Scale),
				rl.GetColor(0xffeb99ff),
			)
		}
	}
	rl.EndDrawing()
}

func (d *Display) GetBuffer() *[64 * 32]byte {
	return &d.buffer
}

func (d *Display) GetPixel(index int) byte {
	return d.buffer[index]
}

func (d *Display) Set(index int, value byte) {
	d.buffer[index] = value
}
