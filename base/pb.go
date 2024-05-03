// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

type internalGetPb[T any] interface {
	ToPB() T
}

func SliceToPB[T any, U internalGetPb[T]](slice []U) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[i] = v.ToPB()
	}
	return result
}
