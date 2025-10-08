package main

import (
	"github.com/drewlarso/chip8-emulator/constants"
	"github.com/drewlarso/chip8-emulator/display"
	"github.com/drewlarso/chip8-emulator/keyboard"
	"github.com/drewlarso/chip8-emulator/processor"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(
		int32(constants.Width),
		int32(constants.Height),
		"Chip8 Emulator",
	)
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	disp := display.NewDisplay()
	cpu := processor.NewProcessor()
	kb := keyboard.NewKeyboard()

	cpu.LoadROM("roms/pong.ch8")

	for !rl.WindowShouldClose() {
		kb.Update()

		for range constants.CyclesPerFrame {
			cpu.Cycle(disp, kb)
		}

		cpu.UpdateTimers()

		disp.Draw()
	}
}
