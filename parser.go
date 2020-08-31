package chip8

import (
	"bufio"
	"encoding/hex"
	"io"
)

func ParseCodes(reader io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	result := make([]byte, 0, 0xFFF)

	for scanner.Scan() {
		instr := scanner.Text()
		h, err := hex.DecodeString(instr)
		if err != nil {
			return nil, err
		}
		result = append(result, h...)
	}
	return result, nil
}
