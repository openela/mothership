// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
