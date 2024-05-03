// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package state

type PackageState struct {
	// Packages is a map of RPM path to SHA256 hash of the package.
	// The RPM path is the base path of the RPM in the storage. So only NVRA is
	// stored here.
	Packages map[string]string `json:"packages"`
}

type State interface {
	FetchNewPackageState() error
	GetDirtyObjects() []string
	WritePackageState() error
	GetState() *PackageState
}
