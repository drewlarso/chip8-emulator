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

func main() {
	rl.InitWindow(int32(ScreenWidth*WindowScale), int32(ScreenHeight*WindowScale), "Chip8 Emulator")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	disp := display.Display{}
	cpu := processor.Processor{}
	fmt.Println(cpu)
	kb := keyboard.NewKeyboard()

	disp.Set(532, true)

	for !rl.WindowShouldClose() {
		kb.Update()
		disp.Draw()
	}
}
