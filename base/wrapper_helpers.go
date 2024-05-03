// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func SqlNullTime(t sql.NullTime) *timestamppb.Timestamp {
	if !t.Valid {
		return nil
	}
	return timestamppb.New(t.Time)
}

func SqlNullString(s sql.NullString) *wrapperspb.StringValue {
	if !s.Valid {
		return nil
	}
	return wrapperspb.String(s.String)
}
