package pbutil

import (
	"github.com/tanyudii/core-go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
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
