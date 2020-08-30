package chip8

type Memory struct {
	space [0xFFF]byte
}

func (m *Memory) Read(address uint16) byte {
	return m.space[address]
}

func (m *Memory) ReadWord(address uint16) uint16 {
	return (uint16(m.space[address]) << 8) | uint16(m.space[address+1])
}

func (m *Memory) ReadBytes(address uint16, n byte) []byte {
	return m.space[address : address+uint16(n)]
}

func (m *Memory) Store(address uint16, value byte) {
	m.space[address] = value
}

func (m *Memory) StoreBytes(address uint16, values ...byte) {
	for i, val := range values {
		m.space[address+uint16(i)] = val
	}
}
