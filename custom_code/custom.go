package main

import (
	"bytes"
	"chip8"
	"strings"
)

func main() {
	reader := strings.NewReader("6010 6110 FA29 D015 00FD")
	codes, err := chip8.ParseCodes(reader)
	if err != nil {
		panic(err)
	}
	chip8.Run(bytes.NewBuffer(codes), chip8.SuperChip48)
}
