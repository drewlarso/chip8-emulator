package display

import rl "github.com/gen2brain/raylib-go/raylib"

type Display struct {
	buffer [64 * 32]bool
}

func (d *Display) Clear() {
	d.buffer = [64 * 32]bool{}
}

func (d *Display) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	for i := range len(d.buffer) {
		x := i % 64
		y := i / 64
		if d.buffer[i] {
			rl.DrawRectangle(int32(x*10), int32(y*10), 10, 10, rl.White)
		}
	}
	rl.EndDrawing()
}

func (d *Display) GetBuffer() *[64 * 32]bool {
	return &d.buffer
}

func (d *Display) Set(index int, value bool) {
	d.buffer[index] = value
}
