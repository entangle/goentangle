package goentangle

import (
	"errors"
	"math"
)

var (
	ErrDeserializationError = errors.New("deserialization error")
)

func DeserializeString(input interface{}) (string, error) {
	switch input.(type) {
	case string:
		return input.(string), nil
	}

	return "", ErrDeserializationError
}

func DeserializeBool(input interface{}) (bool, error) {
	switch input.(type) {
	case bool:
		return input.(bool), nil
	}

	return false, ErrDeserializationError
}

func DeserializeBinary(input interface{}) ([]byte, error) {
	switch input.(type) {
	case []byte:
		if result := input.([]byte); result != nil {
			return result, nil
		}

	case string:
		return []byte(input.(string)), nil
	}

	return nil, ErrDeserializationError
}

func DeserializeInt8(input interface{}) (int8, error) {
	switch input.(type) {
	case int8:
		return input.(int8), nil

	case int16:
		if input.(int16) >= math.MinInt8 && input.(int16) <= math.MaxInt8 {
			return int8(input.(int16)), nil
		}

	case int32:
		if input.(int32) >= math.MinInt8 && input.(int32) <= math.MaxInt8 {
			return int8(input.(int32)), nil
		}

	case int64:
		if input.(int64) >= math.MinInt8 && input.(int64) <= math.MaxInt8 {
			return int8(input.(int64)), nil
		}

	case int:
		if input.(int) >= math.MinInt8 && input.(int) <= math.MaxInt8 {
			return int8(input.(int)), nil
		}

	case uint8:
		if input.(uint8) <= math.MaxInt8 {
			return int8(input.(uint8)), nil
		}

	case uint16:
		if input.(uint16) <= math.MaxInt8 {
			return int8(input.(uint16)), nil
		}

	case uint32:
		if input.(uint32) <= math.MaxInt8 {
			return int8(input.(uint32)), nil
		}

	case uint64:
		if input.(uint64) <= math.MaxInt8 {
			return int8(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxInt8 {
			return int8(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeInt16(input interface{}) (int16, error) {
	switch input.(type) {
	case int8:
		return int16(input.(int8)), nil

	case int16:
		return input.(int16), nil

	case int32:
		if input.(int32) >= math.MinInt16 && input.(int32) <= math.MaxInt16 {
			return int16(input.(int32)), nil
		}

	case int64:
		if input.(int64) >= math.MinInt16 && input.(int64) <= math.MaxInt16 {
			return int16(input.(int64)), nil
		}

	case int:
		if input.(int) >= math.MinInt16 && input.(int) <= math.MaxInt16 {
			return int16(input.(int)), nil
		}

	case uint8:
		return int16(input.(uint8)), nil

	case uint16:
		if input.(uint16) <= math.MaxInt16 {
			return int16(input.(uint16)), nil
		}

	case uint32:
		if input.(uint32) <= math.MaxInt16 {
			return int16(input.(uint32)), nil
		}

	case uint64:
		if input.(uint64) <= math.MaxInt16 {
			return int16(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxInt16 {
			return int16(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeInt32(input interface{}) (int32, error) {
	switch input.(type) {
	case int8:
		return int32(input.(int8)), nil

	case int16:
		return int32(input.(int16)), nil

	case int32:
		return input.(int32), nil

	case int64:
		if input.(int64) >= math.MinInt32 && input.(int64) <= math.MaxInt32 {
			return int32(input.(int64)), nil
		}

	case int:
		if input.(int) >= math.MinInt32 && input.(int) <= math.MaxInt32 {
			return int32(input.(int)), nil
		}

	case uint8:
		return int32(input.(uint8)), nil

	case uint16:
		return int32(input.(uint16)), nil

	case uint32:
		if input.(uint32) <= math.MaxInt32 {
			return int32(input.(uint32)), nil
		}

	case uint64:
		if input.(uint64) <= math.MaxInt32 {
			return int32(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxInt32 {
			return int32(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeInt64(input interface{}) (int64, error) {
	switch input.(type) {
	case int8:
		return int64(input.(int8)), nil

	case int16:
		return int64(input.(int16)), nil

	case int32:
		return int64(input.(int32)), nil

	case int64:
		return input.(int64), nil

	case int:
		if input.(int) >= math.MinInt64 && input.(int) <= math.MaxInt64 {
			return int64(input.(int)), nil
		}

	case uint8:
		return int64(input.(uint64)), nil

	case uint16:
		return int64(input.(uint16)), nil

	case uint32:
		return int64(input.(uint32)), nil

	case uint64:
		if input.(uint64) <= math.MaxInt64 {
			return int64(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxInt64 {
			return int64(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeUint8(input interface{}) (uint8, error) {
	switch input.(type) {
	case int8:
		if input.(int8) >= 0 {
			return uint8(input.(int8)), nil
		}

	case int16:
		if input.(int16) >= 0 && input.(int16) <= math.MaxUint8 {
			return uint8(input.(int16)), nil
		}

	case int32:
		if input.(int32) >= 0 && input.(int32) <= math.MaxUint8 {
			return uint8(input.(int32)), nil
		}

	case int64:
		if input.(int64) >= 0 && input.(int64) <= math.MaxUint8 {
			return uint8(input.(int64)), nil
		}

	case int:
		if input.(int) >= 0 && input.(int) <= math.MaxUint8 {
			return uint8(input.(int)), nil
		}

	case uint8:
		if input.(uint8) <= math.MaxUint8 {
			return uint8(input.(uint8)), nil
		}

	case uint16:
		if input.(uint16) <= math.MaxUint8 {
			return uint8(input.(uint16)), nil
		}

	case uint32:
		if input.(uint32) <= math.MaxUint8 {
			return uint8(input.(uint32)), nil
		}

	case uint64:
		if input.(uint64) <= math.MaxUint8 {
			return uint8(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxUint8 {
			return uint8(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeUint16(input interface{}) (uint16, error) {
	switch input.(type) {
	case int8:
		if input.(int8) >= 0 {
			return uint16(input.(int8)), nil
		}

	case int16:
		if input.(int16) >= 0 {
			return uint16(input.(int16)), nil
		}

	case int32:
		if input.(int32) >= 0 && input.(int32) <= math.MaxUint16 {
			return uint16(input.(int32)), nil
		}

	case int64:
		if input.(int64) >= 0 && input.(int64) <= math.MaxUint16 {
			return uint16(input.(int64)), nil
		}

	case int:
		if input.(int) >= 0 && input.(int) <= math.MaxUint16 {
			return uint16(input.(int)), nil
		}

	case uint8:
		return uint16(input.(uint8)), nil

	case uint16:
		return uint16(input.(uint16)), nil

	case uint32:
		if input.(uint32) <= math.MaxUint16 {
			return uint16(input.(uint32)), nil
		}

	case uint64:
		if input.(uint64) <= math.MaxUint16 {
			return uint16(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxUint16 {
			return uint16(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeUint32(input interface{}) (uint32, error) {
	switch input.(type) {
	case int8:
		if input.(int8) >= 0 {
			return uint32(input.(int8)), nil
		}

	case int16:
		if input.(int16) >= 0 {
			return uint32(input.(int16)), nil
		}

	case int32:
		if input.(int32) >= 0 {
			return uint32(input.(int32)), nil
		}

	case int64:
		if input.(int64) >= 0 && input.(int64) <= math.MaxUint32 {
			return uint32(input.(int64)), nil
		}

	case int:
		if input.(int) >= 0 && input.(int) <= math.MaxUint32 {
			return uint32(input.(int)), nil
		}

	case uint8:
		return uint32(input.(uint8)), nil

	case uint16:
		return uint32(input.(uint16)), nil

	case uint32:
		return uint32(input.(uint32)), nil

	case uint64:
		if input.(uint64) <= math.MaxUint32 {
			return uint32(input.(uint64)), nil
		}

	case uint:
		if input.(uint) <= math.MaxUint32 {
			return uint32(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeUint64(input interface{}) (uint64, error) {
	switch input.(type) {
	case int8:
		if input.(int8) >= 0 {
			return uint64(input.(int8)), nil
		}

	case int16:
		if input.(int16) >= 0 {
			return uint64(input.(int16)), nil
		}

	case int32:
		if input.(int32) >= 0 {
			return uint64(input.(int32)), nil
		}

	case int64:
		if input.(int64) >= 0 {
			return uint64(input.(int64)), nil
		}

	case int:
		if input.(int) >= 0 {
			return uint64(input.(int)), nil
		}

	case uint8:
		return uint64(input.(uint8)), nil

	case uint16:
		return uint64(input.(uint16)), nil

	case uint32:
		return uint64(input.(uint32)), nil

	case uint64:
		return uint64(input.(uint64)), nil

	case uint:
		if input.(uint) <= math.MaxUint64 {
			return uint64(input.(uint)), nil
		}
	}

	return 0, ErrDeserializationError
}

func DeserializeFloat64(input interface{}) (float64, error) {
	switch input.(type) {
	case float32:
		return float64(input.(float32)), nil

	case float64:
		return input.(float64), nil
	}

	return 0, ErrDeserializationError
}

func DeserializeFloat32(input interface{}) (float32, error) {
	switch input.(type) {
	case float32:
		return input.(float32), nil

	case float64:
		return float32(input.(float64)), nil
	}

	return 0, ErrDeserializationError
}
