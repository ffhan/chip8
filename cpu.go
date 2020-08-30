package chip8

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
)

const (
	numOfGeneralPurposeRegisters = 16
)

type CPU struct { // todo: delay & sound timers
	registers [numOfGeneralPurposeRegisters]byte // 0-F registers
	iRegister uint16                             // I register
	pc        uint16                             // program counter
	sp        byte                               // stack pointer

	// additional registers
	dt DelayTimer
	st SoundTimer

	// memory
	stack  [16]uint16
	memory Memory

	clock Clock

	// external HW
	display  Display
	speaker  Speaker
	keyboard Keyboard

	random *rand.Rand
}

func NewCPU(display Display, speaker Speaker, keyboard Keyboard, seed int64) *CPU {
	return &CPU{
		dt:       DelayTimer{},
		st:       SoundTimer{},
		memory:   Memory{},
		clock:    Clock{},
		display:  display,
		speaker:  speaker,
		keyboard: keyboard,
		random:   rand.New(rand.NewSource(seed)),
	}
}

func (c *CPU) LoadRom(rom io.ReadCloser) error {
	defer rom.Close()
	reader := bufio.NewReader(rom)
	i := 0
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			return err
		}
		c.memory.Store(uint16(i), readByte)
	}
}

func (c *CPU) Run() {
	for {
		err := c.Step()
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
}

func (c *CPU) Step() error {
	defer c.clock.Step()
	instr := c.memory.ReadWord(c.pc)
	instruction := ParseInstruction(instr)
	if err := c.execute(instruction); err != nil {
		return wrapError(instruction, err)
	}
	return nil
}

func wrapError(instruction Instruction, err error) error {
	return fmt.Errorf("failed executing %s: %w", instruction.String(), err)
}

func (c *CPU) execute(instr Instruction) error {
	incrementPc := true
	defer func() {
		if incrementPc {
			c.pc += 1
		}
	}()
	// todo: SYSaddr implementation
	switch instr.Opcode {
	case SYSaddr:
		return errors.New("unimplemented")
	case CLS:
		c.display.Clear()
	case RET:
		pointer := c.stack[c.sp]
		c.pc = pointer
		c.sp -= 1
		incrementPc = false
	case JPaddr:
		c.pc = instr.nnn
		incrementPc = false
	case CALLaddr:
		c.sp += 1
		c.stack[c.sp] = c.pc
		c.pc = nnn
		incrementPc = false
	case SEVxByte:
		if c.registers[instr.x] == instr.kk {
			c.pc += 1
		}
	case SNEVxByte:
		if c.registers[instr.x] != instr.kk {
			c.pc += 1
		}
	case SEVxVy:
		if c.registers[instr.x] == c.registers[instr.y] {
			c.pc += 1
		}
	case LDVxByte:
		c.registers[instr.x] = instr.kk
	case ADDVxByte:
		c.registers[instr.x] += instr.kk
	case LDVxVy:
		c.registers[instr.x] = c.registers[instr.y]
	case ORVxVy:
		c.registers[instr.x] |= c.registers[instr.y]
	case ANDVxVy:
		c.registers[instr.x] &= c.registers[instr.y]
	case XORVxVy:
		c.registers[instr.x] ^= c.registers[instr.y]
	case ADDVxVy:
		result := uint16(c.registers[instr.x]) + uint16(c.registers[instr.y])
		c.registers[instr.x] = byte(result)
		if result > 255 {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
	case SUBVxVy:
		notBorrow := c.registers[instr.x] > c.registers[instr.y]
		c.registers[instr.x] -= c.registers[instr.y]
		if notBorrow {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
	case SHRVxVy:
		c.registers[0xF] = c.registers[instr.x] & 0x1
		c.registers[instr.x] >>= 1
	case SUBNVxVy:
		notBorrow := c.registers[instr.y] > c.registers[instr.x]
		c.registers[instr.x] = c.registers[instr.y] - c.registers[instr.x]
		if notBorrow {
			c.registers[0xF] = 1
		} else {
			c.registers[0xF] = 0
		}
	case SHLVxVy:
		c.registers[0xF] = (c.registers[instr.x] & 0x80) >> 7
		c.registers[instr.x] <<= 1
	case SNEVxVy:
		if c.registers[instr.x] != c.registers[instr.y] {
			c.pc += 1
		}
	case LDIAddr:
		c.iRegister = instr.nnn
	case JPV0Addr:
		c.pc = instr.nnn + uint16(c.registers[0])
		incrementPc = false
	case RNDVxByte:
		c.registers[instr.x] = byte(c.random.Uint32()) & instr.kk
	case DRWVxVyNibble:
		// todo: xor coords and fill VF, wrap-around
		pointer := c.iRegister
		n := instr.n
		x := instr.x
		y := instr.y
		bytes := c.memory.ReadBytes(pointer, n)
		c.display.Write(x, y, bytes)
	case SKPVx:
		if c.keyboard.IsDown(Key(c.registers[instr.x])) {
			c.pc += 1
		}
	case SKNPVx:
		if !c.keyboard.IsDown(Key(c.registers[instr.x])) {
			c.pc += 1
		}
	case LDVxDT:
		c.registers[instr.x] = c.dt.Get()
	case LDVxK:
		c.registers[instr.x] = byte(c.keyboard.Wait())
	case LDDTVx:
		c.dt.Set(c.registers[instr.x])
	case LDSTVx:
		c.st.Set(c.registers[instr.x])
	case ADDIVx:
		c.iRegister += uint16(c.registers[instr.x])
	case LDFVx:
		return errors.New("unimplemented")
	case LDBVx:
		val := c.registers[instr.x]
		ones := val % 10
		tenths := (val / 10) % 10
		hundredths := val / 100

		pointer := c.iRegister
		c.memory.StoreBytes(pointer, hundredths, tenths, ones)
	case LDIVx:
		c.memory.StoreBytes(c.iRegister, c.registers[:instr.x]...)
	case LDVxI:
		pointer := c.iRegister
		for i := range c.registers {
			c.registers[i] = c.memory.Read(pointer + uint16(i))
		}
	default:
		return errors.New("unimplemented")
	}
	return nil
}
