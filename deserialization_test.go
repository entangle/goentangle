package goentangle

import (
	"bytes"
	"math"
	"testing"
)

func TestDeserializeString(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected string
	} {
		{"string value", "string value"},
		{"æøå", "æøå"},
	} {
		actual, err := DeserializeString(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		nil,
		[]byte("input"),
		true,
		false,
		0,
	} {
		_, err := DeserializeString(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeBool(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected bool
	} {
		{true,  true},
		{false, false},
	} {
		actual, err := DeserializeBool(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		nil,
		"input",
		[]byte("input"),
		0,
	} {
		_, err := DeserializeBool(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeBinary(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected []byte
	} {
		{"string value", []byte("string value")},
		{"æøå", []byte("æøå")},
		{[]byte("æøå"), []byte("æøå")},
	} {
		actual, err := DeserializeBinary(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if !bytes.Equal(actual, testCase.Expected) {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		nil,
		true,
		false,
		0,
	} {
		_, err := DeserializeBinary(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeInt8(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected int8
	} {
		{int8(math.MinInt8), math.MinInt8},
		{int16(math.MinInt8), math.MinInt8},
		{int32(math.MinInt8), math.MinInt8},
		{int64(math.MinInt8), math.MinInt8},
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt8), math.MaxInt8},
		{int32(math.MaxInt8), math.MaxInt8},
		{int64(math.MaxInt8), math.MaxInt8},
		{uint8(math.MaxInt8), math.MaxInt8},
		{uint16(math.MaxInt8), math.MaxInt8},
		{uint32(math.MaxInt8), math.MaxInt8},
		{uint64(math.MaxInt8), math.MaxInt8},
	} {
		actual, err := DeserializeInt8(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int16(math.MinInt8)-1,
		int32(math.MinInt8)-1,
		int64(math.MinInt8)-1,
		int16(math.MaxInt8)+1,
		int32(math.MaxInt8)+1,
		int64(math.MaxInt8)+1,
		uint8(math.MaxInt8)+1,
		uint16(math.MaxInt8)+1,
		uint32(math.MaxInt8)+1,
		uint64(math.MaxInt8)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeInt8(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeInt16(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected int16
	} {
		{int8(math.MinInt8), math.MinInt8},
		{int16(math.MinInt16), math.MinInt16},
		{int32(math.MinInt16), math.MinInt16},
		{int64(math.MinInt16), math.MinInt16},
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt16), math.MaxInt16},
		{int32(math.MaxInt16), math.MaxInt16},
		{int64(math.MaxInt16), math.MaxInt16},
		{uint8(math.MaxInt8), math.MaxInt8},
		{uint16(math.MaxInt16), math.MaxInt16},
		{uint32(math.MaxInt16), math.MaxInt16},
		{uint64(math.MaxInt16), math.MaxInt16},
	} {
		actual, err := DeserializeInt16(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int32(math.MinInt16)-1,
		int64(math.MinInt16)-1,
		int32(math.MaxInt16)+1,
		int64(math.MaxInt16)+1,
		uint16(math.MaxInt16)+1,
		uint32(math.MaxInt16)+1,
		uint64(math.MaxInt16)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeInt16(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeInt32(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected int32
	} {
		{int8(math.MinInt8), math.MinInt8},
		{int16(math.MinInt16), math.MinInt16},
		{int32(math.MinInt32), math.MinInt32},
		{int64(math.MinInt32), math.MinInt32},
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt16), math.MaxInt16},
		{int32(math.MaxInt32), math.MaxInt32},
		{int64(math.MaxInt32), math.MaxInt32},
		{uint8(math.MaxInt8), math.MaxInt8},
		{uint16(math.MaxInt16), math.MaxInt16},
		{uint32(math.MaxInt32), math.MaxInt32},
		{uint64(math.MaxInt32), math.MaxInt32},
	} {
		actual, err := DeserializeInt32(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int64(math.MinInt32)-1,
		int64(math.MaxInt32)+1,
		uint32(math.MaxInt32)+1,
		uint64(math.MaxInt32)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeInt32(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeInt64(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected int64
	} {
		{int8(math.MinInt8), math.MinInt8},
		{int16(math.MinInt16), math.MinInt16},
		{int32(math.MinInt32), math.MinInt32},
		{int64(math.MinInt64), math.MinInt64},
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt16), math.MaxInt16},
		{int32(math.MaxInt32), math.MaxInt32},
		{int64(math.MaxInt64), math.MaxInt64},
		{uint8(math.MaxInt8), math.MaxInt8},
		{uint16(math.MaxInt16), math.MaxInt16},
		{uint32(math.MaxInt32), math.MaxInt32},
		{uint64(math.MaxInt64), math.MaxInt64},
	} {
		actual, err := DeserializeInt64(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		uint64(math.MaxInt64)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeInt64(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeUint8(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected uint8
	} {
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxUint8), math.MaxUint8},
		{int32(math.MaxUint8), math.MaxUint8},
		{int64(math.MaxUint8), math.MaxUint8},
		{uint8(math.MaxUint8), math.MaxUint8},
		{uint16(math.MaxUint8), math.MaxUint8},
		{uint32(math.MaxUint8), math.MaxUint8},
		{uint64(math.MaxUint8), math.MaxUint8},
	} {
		actual, err := DeserializeUint8(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int8(-1),
		int16(-1),
		int32(-1),
		int64(-1),
		int16(math.MaxUint8)+1,
		int32(math.MaxUint8)+1,
		int64(math.MaxUint8)+1,
		uint16(math.MaxUint8)+1,
		uint32(math.MaxUint8)+1,
		uint64(math.MaxUint8)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeUint8(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeUint16(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected uint16
	} {
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt16), math.MaxInt16},
		{int32(math.MaxUint16), math.MaxUint16},
		{int64(math.MaxUint16), math.MaxUint16},
		{uint8(math.MaxUint8), math.MaxUint8},
		{uint16(math.MaxUint16), math.MaxUint16},
		{uint32(math.MaxUint16), math.MaxUint16},
		{uint64(math.MaxUint16), math.MaxUint16},
	} {
		actual, err := DeserializeUint16(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int8(-1),
		int16(-1),
		int32(-1),
		int64(-1),
		int32(math.MaxUint16)+1,
		int64(math.MaxUint16)+1,
		uint32(math.MaxUint16)+1,
		uint64(math.MaxUint16)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeUint16(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeUint32(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected uint32
	} {
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt16), math.MaxInt16},
		{int32(math.MaxInt32), math.MaxInt32},
		{int64(math.MaxUint32), math.MaxUint32},
		{uint8(math.MaxUint8), math.MaxUint8},
		{uint16(math.MaxUint16), math.MaxUint16},
		{uint32(math.MaxUint32), math.MaxUint32},
		{uint64(math.MaxUint32), math.MaxUint32},
	} {
		actual, err := DeserializeUint32(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int8(-1),
		int16(-1),
		int32(-1),
		int64(-1),
		int64(math.MaxUint32)+1,
		uint64(math.MaxUint32)+1,
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeUint32(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeUint64(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected uint64
	} {
		{int8(0), 0},
		{int16(0), 0},
		{int32(0), 0},
		{int64(0), 0},
		{uint8(0), 0},
		{uint16(0), 0},
		{uint32(0), 0},
		{uint64(0), 0},
		{int8(math.MaxInt8), math.MaxInt8},
		{int16(math.MaxInt16), math.MaxInt16},
		{int32(math.MaxInt32), math.MaxInt32},
		{int64(math.MaxInt64), math.MaxInt64},
		{uint8(math.MaxUint8), math.MaxUint8},
		{uint16(math.MaxUint16), math.MaxUint16},
		{uint32(math.MaxUint32), math.MaxUint32},
		{uint64(math.MaxUint64), math.MaxUint64},
	} {
		actual, err := DeserializeUint64(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		int8(-1),
		int16(-1),
		int32(-1),
		int64(-1),
		nil,
		true,
		false,
		"input",
		[]byte("input"),
	} {
		_, err := DeserializeUint64(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeFloat32(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected float32
	} {
		{float32(0.0), 0.0},
		{float64(0.0), 0.0},
		{float32(123.0), 123.0},
		{float64(123.0), 123.0},
		{float32(-123.0), -123.0},
		{float64(-123.0), -123.0},
	} {
		actual, err := DeserializeFloat32(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		nil,
		"string value",
		[]byte("input"),
		true,
		false,
		0,
	} {
		_, err := DeserializeFloat32(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}

func TestDeserializeFloat64(t *testing.T) {
	// Valid.
	for _, testCase := range []struct{
		Input    interface{}
		Expected float64
	} {
		{float32(0.0), 0.0},
		{float64(0.0), 0.0},
		{float32(123.0), 123.0},
		{float64(123.0), 123.0},
		{float32(-123.0), -123.0},
		{float64(-123.0), -123.0},
	} {
		actual, err := DeserializeFloat64(testCase.Input)
		if err != nil {
			t.Errorf("Unexpected error deserializing %v: %v", testCase.Input, err)
		} else if actual != testCase.Expected {
			t.Errorf("Expected deserialized value to be %v, but it is %v", testCase.Expected, actual)
		}
	}

	// Invalid.
	for _, input := range []interface{} {
		nil,
		"string value",
		[]byte("input"),
		true,
		false,
		0,
	} {
		_, err := DeserializeFloat64(input)
		if err == nil {
			t.Errorf("Expected error while deserializing %v", input)
		} else if err != ErrDeserializationError {
			t.Errorf("Unexpected error while deserializing %v: %v", input, err)
		}
	}
}
