// Copyright 2021, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build windows

package windows

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/unicode"
)

func TestBufferReadBytes(t *testing.T) {
	buffer := NewBuffer()
	utf8 := []byte("test")
	utf16, _ := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Bytes(utf8)
	for i, b := range utf16 {
		buffer.buffer[i] = b
	}
	offset := uint32(len(utf16))
	bytes, err := buffer.ReadBytes(offset)
	require.NoError(t, err)
	require.Equal(t, utf8, bytes)
}

func TestBufferReadString(t *testing.T) {
	buffer := NewBuffer()
	utf8 := []byte("test")
	utf16, _ := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Bytes(utf8)
	for i, b := range utf16 {
		buffer.buffer[i] = b
	}
	offset := uint32(len(utf16))
	result, err := buffer.ReadString(offset)
	require.NoError(t, err)
	require.Equal(t, "test", result)
}

func TestBufferUpdateSize(t *testing.T) {
	buffer := NewBuffer()
	buffer.UpdateSize(1)
	require.Equal(t, 1, len(buffer.buffer))
}

func TestBufferSize(t *testing.T) {
	buffer := NewBuffer()
	require.Equal(t, uint32(defaultBufferSize), buffer.Size())
}

func TestBufferFirstByte(t *testing.T) {
	buffer := NewBuffer()
	buffer.buffer[0] = '1'
	require.Equal(t, &buffer.buffer[0], buffer.FirstByte())
}
