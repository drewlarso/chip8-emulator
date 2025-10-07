package main

import (
	"fmt"

	"github.com/drewlarso/chip8-emulator/display"
	"github.com/drewlarso/chip8-emulator/keyboard"
	"github.com/drewlarso/chip8-emulator/processor"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const ScreenWidth int = 64
const ScreenHeight int = 32
const WindowScale int = 10
const CyclesPerFrame int = 10

func main() {
	rl.InitWindow(int32(ScreenWidth*WindowScale), int32(ScreenHeight*WindowScale), "Chip8 Emulator")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	disp := display.Display{} // disp is the actual value
	cpu := processor.NewProcessor()
	fmt.Println(cpu)
	kb := keyboard.NewKeyboard() // already a pointer to my keyboard struct

	disp.Set(532, true)

	for !rl.WindowShouldClose() {
		kb.Update()

		for range CyclesPerFrame {
			cpu.Cycle(&disp, kb)
		}

		cpu.UpdateTimers()

		disp.Draw()
	}
}
