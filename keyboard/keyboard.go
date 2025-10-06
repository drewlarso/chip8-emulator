package keyboard

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Keyboard struct {
	keymap map[rune]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{keymap: make(map[rune]bool)}
}

func (k *Keyboard) Update() {
	k.keymap['1'] = rl.IsKeyDown('1')
	k.keymap['2'] = rl.IsKeyDown('2')
	k.keymap['3'] = rl.IsKeyDown('3')
	k.keymap['C'] = rl.IsKeyDown('4')
	k.keymap['4'] = rl.IsKeyDown('Q')
	k.keymap['5'] = rl.IsKeyDown('W')
	k.keymap['6'] = rl.IsKeyDown('E')
	k.keymap['D'] = rl.IsKeyDown('R')
	k.keymap['7'] = rl.IsKeyDown('A')
	k.keymap['8'] = rl.IsKeyDown('S')
	k.keymap['9'] = rl.IsKeyDown('D')
	k.keymap['E'] = rl.IsKeyDown('F')
	k.keymap['A'] = rl.IsKeyDown('Z')
	k.keymap['0'] = rl.IsKeyDown('X')
	k.keymap['B'] = rl.IsKeyDown('C')
	k.keymap['F'] = rl.IsKeyDown('V')
}

func (k *Keyboard) IsKeyDown(key rune) bool {
	return k.keymap[key]
}

/*
|1 2 3 C|
|4 5 6 D|
|7 8 9 E|
|A 0 B F|

translates to

|1 2 3 4|
|Q W E R|
|A S D F|
|Z X C V|
*/
