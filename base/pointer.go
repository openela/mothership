// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

func Pointer[T any](v T) *T {
	return &v
}
