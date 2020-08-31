package chip8

import (
	"fmt"
)

// todo: order cases

const (
	addrMask   uint16 = 0x0FFF
	nibbleMask uint16 = 0x000F
	xMask      uint16 = 0x0F00
	yMask      uint16 = 0x00F0
	byteMask   uint16 = 0x00FF

	// aliases
	nnn = addrMask
	n   = nibbleMask
	x   = xMask
	y   = yMask
	kk  = byteMask

	nibble0Mask uint16 = 0xF000
	nibble1Mask        = xMask
	nibble2Mask        = yMask
	nibble3Mask        = nibbleMask
	byte0Mask   uint16 = 0xFF00
	byte1Mask          = byteMask
)

type Opcode uint8 // internal codes to discern different instructions

// standard Chip-8 instruction codes
const (
	SYSaddr   Opcode = iota + 1 // 0nnn - jump to a machine code routine at nnn
	CLS                         // 00E0 - clear the display
	RET                         // 00EE - return from a subroutine (the interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer)
	JPaddr                      // 1nnn - jump to location nnn (sets the program counter to nnn)
	CALLaddr                    // 2nnn - call subroutine at nnn (increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn)
	SEVxByte                    // 3xkk - skip next instruction if Vx == kk (if they are equal, increments the program counter by 2)
	SNEVxByte                   // 4xkk - skip next instruction if Vx != kk (if they are NOT equal, increments the program counter by 2)
	SEVxVy                      // 5xy0 - skip next instruction if Vx == Vy (if they are equal, increments the program counter by 2)

	LDVxByte  // 6xkk - Vx = kk
	ADDVxByte // 7xkk - Vx = Vx + kk
	LDVxVy    // 8xy0 - Vx = Vy
	ORVxVy    // 8xy1 - Vx = Vx OR Vy
	ANDVxVy   // 8xy2 - Vx = Vx AND Vy
	XORVxVy   // 8xy3 - Vx = Vx XOR Vy
	ADDVxVy   // 8xy4 - Vx = Vx + Vy, set VF = carry (if result > 255 set VF=1, otherwise 0)
	SUBVxVy   // 8xy5 - Vx = Vx - Vy, set VF = NOT borrow (if Vx > Vy then VF = 1, otherwise 0)
	SHRVxVy   // 8xy6 - Vx = Vx SHR 1 (if the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2)
	SUBNVxVy  // 8xy7 - Vx = Vy - Vx, set VF = NOT borrow (analogous to SUBVxVy)
	SHLVxVy   // 8xyE - Vx = Vx SHL 1 (if the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2)

	SNEVxVy  // 9xy0 - skip next instruction if Vx != Vy (the values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2)
	LDIAddr  // Annn - I = nnn
	JPV0Addr // Bnnn - jump to location nnn+V0 (PC=nnn+V0)

	RNDVxByte // Cxkk - Vx = random byte AND kk (generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx)
	/*
		Dxyn - display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision

				The interpreter reads n bytes from memory, starting at the address stored in I.
			These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
			Sprites are XORed onto the existing screen.
			If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.
			If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen.
	*/
	DRWVxVyNibble

	SKPVx  // Ex9E - skip next instruction if key with the value of Vx is pressed (checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2)
	SKNPVx // ExA1- skip next instruction if key with the value of Vx is NOT pressed (checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2)

	LDVxDT // Fx07 - Vx = delay timer value
	LDVxK  // Fx0A - wait for a key press, store the value of the key in Vx (all execution stops until a key is pressed, then the value of that key is stored in Vx
	LDDTVx // Fx15 - set delay timer = Vx
	LDSTVx // Fx18 - set sound timer = Vx

	ADDIVx // Fx1E - I = I + Vx
	LDFVx  // Fx29 - set I = location of sprite for digit Vx (the value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx)
	/*
		Fx33 - store BCD representation of Vx in memory locations I, I+1 and I+2

			Takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2
	*/
	LDBVx
	LDIVx // Fx55 - store registers V0 through Vx in memory starting at location I
	LDVxI // Fx65 - read registers V0 through Vx from memory starting at location I

	SCDNibble // 00Cn - SCD nibble
	SCR       // 00FB
	SCL       // 00FC
	EXIT      // 00FD
	LOW       // 00FE
	HIGH      // 00FF
	DRWVxVy0  // Dxy0 - DRW Vx, Vy, 0
	LDHFVx    // Fx30 - LD HF, Vx
	LDRVx     // Fx75 - LD R, Vx
	LDVxR     // Fx85 - LD Vx, R

)

type Instruction struct {
	Opcode   Opcode
	Original uint16
	nnn      uint16
	n        byte
	x        byte
	y        byte
	kk       byte
}

func (i Instruction) String() string {
	return fmt.Sprintf("{Opcode: %d, Original: %X}", i.Opcode, i.Original)
}

func ParseInstruction(instr uint16) Instruction {
	return Instruction{
		Opcode:   getOpcode(instr),
		Original: instr,
		nnn:      instr & nnn,
		n:        byte(instr & n),
		x:        byte((instr & x) >> 8),
		y:        byte((instr & y) >> 4),
		kk:       byte(instr & kk),
	}
}

func getOpcode(instr uint16) Opcode {
	switch (instr & nibble0Mask) >> 12 {
	case 0:
		return parse0(instr)
	case 1:
		return JPaddr
	case 2:
		return CALLaddr
	case 3:
		return SEVxByte
	case 4:
		return SNEVxByte
	case 5:
		return SEVxVy
	case 6:
		return LDVxByte
	case 7:
		return ADDVxByte
	case 8:
		return parse8(instr)
	case 9:
		return SNEVxVy
	case 0xA:
		return LDIAddr
	case 0xB:
		return JPV0Addr
	case 0xC:
		return RNDVxByte
	case 0xD:
		return parseD(instr)
	case 0xE:
		return parseE(instr)
	case 0xF:
		return parseF(instr)
	}
	panic(fmt.Errorf("invalid instruction %X\n", instr))
}

func parseD(instr uint16) Opcode {
	// these are the same instruction, but super chip 48 has the extra instruction
	switch instr & nibbleMask {
	case 0:
		return DRWVxVy0
	default:
		return DRWVxVyNibble
	}
}

func parseF(instr uint16) Opcode {
	switch instr & byteMask {
	case 0x07:
		return LDVxDT
	case 0x0A:
		return LDVxK
	case 0x15:
		return LDDTVx
	case 0x18:
		return LDSTVx
	case 0x1E:
		return ADDIVx
	case 0x29:
		return LDFVx
	case 0x33:
		return LDBVx
	case 0x55:
		return LDIVx
	case 0x65:
		return LDVxI
	case 0x30:
		return LDHFVx
	case 0x75:
		return LDRVx
	case 0x85:
		return LDVxR
	}
	panic(fmt.Errorf("invalid instruction %X\n", instr))
}

func parseE(instr uint16) Opcode {
	switch instr & byteMask {
	case 0x9E:
		return SKPVx
	case 0xA1:
		return SKNPVx
	}
	panic(fmt.Errorf("invalid instruction %X\n", instr))
}

func parse8(instr uint16) Opcode {
	const (
		LDVxVyMask   = 0x0
		ORVxVyMask   = 0x1
		ANDVxVyMask  = 0x2
		XORVxVyMask  = 0x3
		ADDVxVyMask  = 0x4
		SUBVxVyMask  = 0x5
		SHRVxVyMask  = 0x6
		SUBNVxVyMask = 0x7
		SHLVxVyMask  = 0xE
	)
	switch instr & nibbleMask {
	case LDVxVyMask:
		return LDVxVy
	case ORVxVyMask:
		return ORVxVy
	case ANDVxVyMask:
		return ANDVxVy
	case XORVxVyMask:
		return XORVxVy
	case ADDVxVyMask:
		return ADDVxVy
	case SUBVxVyMask:
		return SUBVxVy
	case SHRVxVyMask:
		return SHRVxVy
	case SUBNVxVyMask:
		return SUBNVxVy
	case SHLVxVyMask:
		return SHLVxVy
	}
	panic(fmt.Errorf("invalid instruction %X\n", instr))
}

func parse0(instr uint16) Opcode {
	const (
		clsMask       = 0xE0
		retMask       = 0xEE
		scdNibbleMask = 0xC0
		scr           = 0xFB
		sclMask       = 0xFC
		exitMask      = 0xFD
		lowMask       = 0xFE
		highMask      = 0xFF
	)
	if clsMask == instr {
		return CLS
	} else if retMask == instr {
		return RET
	} else if instr == scr {
		return SCR
	} else if instr == sclMask {
		return SCL
	} else if instr == exitMask {
		return EXIT
	} else if instr == lowMask {
		return LOW
	} else if instr == highMask {
		return HIGH
	} else if instr&scdNibbleMask == scdNibbleMask {
		return SCDNibble
	}
	return SYSaddr
}
