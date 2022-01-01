package prototime

import (
	"time"

	//nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// ShouldTimestamp
func ShouldTimestamp(t time.Time) *timestamp.Timestamp {
	if t.IsZero() {
		return nil
	}

	//nolint:staticcheck
	rt, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return rt
}

// ShouldTime
func ShouldTime(timestamp *timestamp.Timestamp) time.Time {
	if timestamp == nil {
		return time.Time{}
	}

	//nolint:staticcheck
	t, err := ptypes.Timestamp(timestamp)
	if err != nil {
		panic(err)
	}
	return t
}
