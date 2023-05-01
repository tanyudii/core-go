package pbutil

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/tanyudii/core-go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func TimeToPbTimestamp(val *time.Time) *timestamppb.Timestamp {
	if val == nil {
		return nil
	}
	return &timestamppb.Timestamp{Seconds: val.Unix()}
}

func PbTimestampToTime(val *timestamppb.Timestamp) *time.Time {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.AsTime())
}

func Float64PointerToDoubleValue(val *float64) *wrappers.DoubleValue {
	if val == nil {
		return nil
	}
	return &wrappers.DoubleValue{
		Value: *val,
	}
}

func PbDoubleValueToFloat64Pointer(val *wrappers.DoubleValue) *float64 {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func Float32PointerToFloatValue(val *float32) *wrappers.FloatValue {
	if val == nil {
		return nil
	}
	return &wrappers.FloatValue{
		Value: *val,
	}
}

func PbFloatValueToFloat32Pointer(val *wrappers.FloatValue) *float32 {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func Int64PointerToInt64Value(val *int64) *wrappers.Int64Value {
	if val == nil {
		return nil
	}
	return &wrappers.Int64Value{
		Value: *val,
	}
}

func PbInt64ValueToInt64Pointer(val *wrappers.Int64Value) *int64 {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func UInt64PointerToUInt64Value(val *uint64) *wrappers.UInt64Value {
	if val == nil {
		return nil
	}
	return &wrappers.UInt64Value{
		Value: *val,
	}
}

func PbUInt64ValueToUInt64Pointer(val *wrappers.UInt64Value) *uint64 {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func Int32PointerToInt32Value(val *int32) *wrappers.Int32Value {
	if val == nil {
		return nil
	}
	return &wrappers.Int32Value{
		Value: *val,
	}
}

func PbInt32ValueToInt32Pointer(val *wrappers.Int32Value) *int32 {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func UInt32PointerToUInt32Value(val *uint32) *wrappers.UInt32Value {
	if val == nil {
		return nil
	}
	return &wrappers.UInt32Value{
		Value: *val,
	}
}

func PbUInt32ValueToUInt32Pointer(val *wrappers.UInt32Value) *uint32 {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func BoolPointerToBoolValue(val *bool) *wrapperspb.BoolValue {
	if val == nil {
		return nil
	}
	return &wrapperspb.BoolValue{
		Value: *val,
	}
}

func PbBoolValueToBoolPointer(val *wrapperspb.BoolValue) *bool {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func StringPointerToStringValue(val *string) *wrappers.StringValue {
	if val == nil {
		return nil
	}
	return &wrappers.StringValue{
		Value: *val,
	}
}

func PbStringValueToStringPointer(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func BytesPointerToBytesValue(val *[]byte) *wrappers.BytesValue {
	if val == nil {
		return nil
	}
	return &wrappers.BytesValue{
		Value: *val,
	}
}

func PbBytesValueToBytesPointer(val *wrappers.BytesValue) *[]byte {
	if val == nil {
		return nil
	}
	return common.PointerVal(val.Value)
}

func NullableStringValue(val string) *wrappers.StringValue {
	if val == "" {
		return nil
	}
	return &wrappers.StringValue{
		Value: val,
	}
}
