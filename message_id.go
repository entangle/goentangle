package goentangle

// Message ID.
type MessageId uint32

// Parse a message ID.
func ParseMessageId(input interface{}) (id MessageId, ok bool) {
	raw, err := DeserializeUint32(input)
	return MessageId(raw), err == nil
}
