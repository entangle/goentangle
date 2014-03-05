package goentangle

import (
	"fmt"
)

// Opcode.
type Opcode uint8

// Opcodes.
const (
	// Request opcode.
	RequestOpcode Opcode = iota

	// Notification opcode.
	NotificationOpcode

	// Response opcode.
	ResponseOpcode

	// Exception opcode.
	ExceptionOpcode

	// Notification acknowledgement opcode.
	NotificationAcknowledgementOpcode

	// Compressed message opcode.
	CompressedMessageOpcode = 0x7f
)

// Opcode names.
var opcodeNames = map[Opcode]string{
	RequestOpcode:      "request",
	NotificationOpcode: "notification",
	ResponseOpcode:     "response",
	NotificationAcknowledgementOpcode: "notification acknowledgement",
	ExceptionOpcode:    "exception",
}

func (o Opcode) String() string {
	if name, ok := opcodeNames[o]; ok {
		return name
	}

	return fmt.Sprintf("<invalid: %d>", o)
}

// Test if an opcode is valid.
func (o Opcode) Valid() bool {
	return o == CompressedMessageOpcode || o >= RequestOpcode && o <= NotificationAcknowledgementOpcode
}

// Parse an opcode.
func ParseOpcode(input interface{}) (o Opcode, ok bool) {
	raw, err := DeserializeUint8(input)
	return Opcode(raw), err == nil
}
