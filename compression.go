package goentangle

import (
	"github.com/golang/snappy/snappy"
	"fmt"
)

// Compression method.
type CompressionMethod uint8

// Compression methods.
const (
	// Snappy.
	SnappyCompression CompressionMethod = iota
)

// Compression method names.
var compressionMethodNames = map[CompressionMethod]string{
	SnappyCompression: "Snappy",
}

func (m CompressionMethod) String() string {
	if name, ok := compressionMethodNames[m]; ok {
		return name
	}

	return fmt.Sprintf("<invalid: %d>", m)
}

// Test if a compression method is valid.
func (m CompressionMethod) Valid() bool {
	return m == SnappyCompression
}

// Deserialize a compression method.
func DeserializeCompressionMethod(input interface{}) (m CompressionMethod, ok bool) {
	raw, err := DeserializeUint8(input)
	return CompressionMethod(raw), err == nil
}

// Compress.
func (m CompressionMethod) Compress(input []byte) (output []byte, err error) {
	switch m {
	case SnappyCompression:
		output = make([]byte, snappy.MaxEncodedLen(len(input)))
		output, err = snappy.Encode(output, input)

	default:
		panic("invalid compression method")
	}

	return
}

// Decompress.
func (m CompressionMethod) Decompress(input []byte) (output []byte, err error) {
	switch m {
	case SnappyCompression:
		var s int
		if s, err = snappy.DecodedLen(input); err != nil {
			return
		}
		output = make([]byte, s)
		output, err = snappy.Decode(output, input)

	default:
		panic("invalid compression method")
	}

	return
}
