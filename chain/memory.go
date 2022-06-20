// Copyright 2019 The go-tchain Authors
// This file is part of go-tchain.
//
// go-tchain is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-tchain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//

package chain

import (
	"fmt"
)

// Memory implements a simple memory model for the ethereum virtual machine.
type WasmMemory struct {
	Memory []byte
	Pos    int
	Size   map[uint64]int
}

func NewWasmMemory() *WasmMemory {
	return &WasmMemory{
		Size: make(map[uint64]int),
	}
}

// Set sets offset + size to value
func (m *WasmMemory) Set(offset, size uint64, value []byte) {
	// length of Memory may never be less than offset + size.
	// The Memory should be resized PRIOR to setting the memory
	if size > uint64(len(m.Memory)) {
		panic("INVALID memory: Memory empty")
	}

	// It's possible the offset is greater than 0 and size equals 0. This is because
	// the calcMemSize (common.go) could potentially return 0 when size is zero (NO-OP)
	if size > 0 {
		copy(m.Memory[offset:offset+size], value)
		m.Size[offset] = len(value)
		m.Pos = m.Pos + int(size)
	} else {
		m.Size[offset] = 0
		m.Pos = m.Pos + 1
	}
}
func (m *WasmMemory) SetBytes(value []byte) (offset int) {
	offset = m.Len()
	m.Set(uint64(offset), uint64(len(value)), value)
	return
}

// Resize resizes the memory to size
func (m *WasmMemory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.Memory = append(m.Memory, make([]byte, size-uint64(m.Len()))...)
	}
}

// Get returns offset + size as a new slice
func (m *WasmMemory) Get(offset uint64) (cpy []byte) {
	ptr := uint32(offset)
	if int32(ptr) < 0 {
		ptr = uint32(int32(len(m.Memory)) + int32(ptr))
	}
	offset = uint64(ptr)
	var size int
	var ok bool
	if size, ok = m.Size[offset]; ok {
		if size == 0 {
			return nil
		}
	} else {
		return nil
	}

	if len(m.Memory) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.Memory[offset:offset+uint64(size)])
		return
	}

	return
}

func (m *WasmMemory) NormalizeOffset(offset uint32) uint32 {
	if int32(offset) < 0 {
		offset = uint32(int32(m.MemSize()) + int32(offset))
	}
	return offset
}

// GetPtr returns the offset + size
func (m *WasmMemory) GetPtr(offset uint64) []byte {
	ptr := uint32(offset)
	if int32(ptr) < 0 {
		ptr = uint32(int32(len(m.Memory)) + int32(ptr))
	}
	offset = uint64(ptr)
	var size int
	var ok bool
	if size, ok = m.Size[offset]; ok {
		if size == 0 {
			return nil
		}
	} else {
		return nil
	}
	if len(m.Memory) > int(offset) {
		return m.Memory[offset : offset+uint64(size)]
	}
	return nil
}

// Len returns the length of the backing slice
func (m *WasmMemory) Len() int {
	return m.Pos
}

func (m *WasmMemory) MemSize() int {
	return len(m.Memory)
}

// Data returns the backing slice
func (m *WasmMemory) Data() []byte {
	return m.Memory
}

func (m *WasmMemory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.Memory))
	if len(m.Memory) > 0 {
		addr := 0
		for i := 0; i+32 <= len(m.Memory); i += 32 {
			fmt.Printf("%03d: % x\n", addr, m.Memory[i:i+32])
			addr++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("####################")
}
