package common

import "time"

func PointerInt8(val int8) *int8 {
	return &val
}

func PointerInt16(val int16) *int16 {
	return &val
}

func PointerInt32(val int32) *int32 {
	return &val
}

func PointerInt64(val int64) *int64 {
	return &val
}

func PointerUint8(val uint8) *uint8 {
	return &val
}

func PointerUint16(val uint16) *uint16 {
	return &val
}

func PointerUint32(val uint32) *uint32 {
	return &val
}

func PointerUint64(val uint64) *uint64 {
	return &val
}

func PointerFloat32(val float32) *float32 {
	return &val
}

func PointerFloat64(val float64) *float64 {
	return &val
}

func PointerString(val string) *string {
	return &val
}

func PointerBool(val bool) *bool {
	return &val
}

func PointerTime(val time.Time) *time.Time {
	return &val
}
