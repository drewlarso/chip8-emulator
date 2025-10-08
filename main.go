package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/drewlarso/chip8-emulator/constants"
	"github.com/drewlarso/chip8-emulator/display"
	"github.com/drewlarso/chip8-emulator/keyboard"
	"github.com/drewlarso/chip8-emulator/processor"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	constants.CyclesPerFrame = 20
	rom := "roms/logo.ch8"
	if len(os.Args) == 2 {
		rom = os.Args[1]
	} else if len(os.Args) >= 3 {
		cycles, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		constants.CyclesPerFrame = cycles
		rom = os.Args[2]
	}
	fmt.Println(os.Args)
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

	cpu.LoadROM(rom)

	for !rl.WindowShouldClose() {
		kb.Update()

		for range constants.CyclesPerFrame {
			cpu.Cycle(disp, kb)
		}

		cpu.UpdateTimers()

		disp.Draw()
	}
}
