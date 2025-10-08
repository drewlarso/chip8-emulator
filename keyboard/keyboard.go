package keyboard

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Keyboard struct {
	// keys map[rune]bool
	keys [16]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{keys: [16]bool{}}
}

func (k *Keyboard) Update() {

	// 1 2 3 C
	k.keys[1] = rl.IsKeyDown(rl.KeyOne)   // 1
	k.keys[2] = rl.IsKeyDown(rl.KeyTwo)   // 2
	k.keys[3] = rl.IsKeyDown(rl.KeyThree) // 3
	k.keys[12] = rl.IsKeyDown(rl.KeyFour) // C

	// 4 5 6 D
	k.keys[4] = rl.IsKeyDown(rl.KeyQ)  // 4
	k.keys[5] = rl.IsKeyDown(rl.KeyW)  // 5
	k.keys[6] = rl.IsKeyDown(rl.KeyE)  // 6
	k.keys[13] = rl.IsKeyDown(rl.KeyR) // D

	// 7 8 9 E
	k.keys[7] = rl.IsKeyDown(rl.KeyA)  // 7
	k.keys[8] = rl.IsKeyDown(rl.KeyS)  // 8
	k.keys[9] = rl.IsKeyDown(rl.KeyD)  // 9
	k.keys[14] = rl.IsKeyDown(rl.KeyF) // E

	// A 0 B F
	k.keys[10] = rl.IsKeyDown(rl.KeyZ) // A
	k.keys[0] = rl.IsKeyDown(rl.KeyX)  // 0
	k.keys[11] = rl.IsKeyDown(rl.KeyC) // B
	k.keys[15] = rl.IsKeyDown(rl.KeyV) // F
}

func (k *Keyboard) IsKeyDown(index byte) bool {
	return k.keys[index]
}

func (k *Keyboard) AnyKeyDown() byte {
	for i := range 16 {
		if k.keys[i] {
			return byte(i)
		}
	}
	return 255
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
