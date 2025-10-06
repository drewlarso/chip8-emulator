package cpu

type CPU struct {
	registers     [16]byte
	indexRegister uint16
	delayTimer    byte
	soundTimer    byte
	pc            uint16
	sp            byte
	stack         [16]uint16
}
