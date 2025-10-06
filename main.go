package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(640, 320, "Chip8 Emulator")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	screenBuffer := [64 * 32]bool{}

	screenBuffer[531] = true

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		for i := range len(screenBuffer) {
			x := i % 64
			y := i / 64
			if screenBuffer[i] {
				rl.DrawRectangle(int32(x*10), int32(y*10), 10, 10, rl.White)
			}
		}

		rl.EndDrawing()
	}
}
