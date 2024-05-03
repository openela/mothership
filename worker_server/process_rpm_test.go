// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package mothership_worker_server

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWorker_VerifyResourceExists(t *testing.T) {
	require.Nil(t, testW.VerifyResourceExists("memory://efi-rpm-macros-3-3.el8.src.rpm"))
}

func TestWorker_VerifyResourceExists_NotFound(t *testing.T) {
	err := testW.VerifyResourceExists("memory://not-found.rpm")
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "resource does not exist")
}

func TestWorker_VerifyResourceExists_CannotRead(t *testing.T) {
	err := testW.VerifyResourceExists("bad-protocol://not-found.rpm")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "client submitted a resource URI that cannot be read by server")
}

func TestWorker_ImportRPM(t *testing.T) {
	require.False(t, inmf.repos["efi-rpm-macros"])

	res, err := testW.ImportRPM(
		"memory://efi-rpm-macros-3-3.el8.src.rpm",
		"518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		"Rocky Linux release 8.8 (Green Obsidian)",
	)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, "efi-rpm-macros", res.Pkg)
	require.Equal(t, "efi-rpm-macros-0:3-3.el8.noarch.rpm", res.Nevra)

	require.True(t, inmf.repos["efi-rpm-macros"])
}

func TestWorker_ImportRPM_Existing(t *testing.T) {
	require.False(t, inmf.repos["basesystem"])

	res, err := testW.ImportRPM(
		"memory://basesystem-11-5.el8.src.rpm",
		"6beff4cbfd5425e2c193312a9a184969a27d6bbd2d4cc29d7ce72dbe3d9f6416",
		"Rocky Linux release 8.8 (Green Obsidian)",
	)
	require.Nil(t, err)
	require.NotNil(t, res)

	require.True(t, inmf.repos["basesystem"])

	res, err = testW.ImportRPM(
		"memory://basesystem-11-5.el8.src.rpm",
		"6beff4cbfd5425e2c193312a9a184969a27d6bbd2d4cc29d7ce72dbe3d9f6416",
		"Rocky Linux release 8.8 (Green Obsidian)",
	)
	require.Nil(t, err)
	require.NotNil(t, res)

	require.True(t, inmf.repos["basesystem"])

	remote := testW.forge.GetRemote("basesystem")

	repo, err := getRepo(remote, nil)
	require.Nil(t, err)

	commitIter, err := repo.CommitObjects()
	require.Nil(t, err)
	c, err := commitIter.Next()
	require.Nil(t, err)
	require.NotNil(t, c)
	require.Equal(t, "import basesystem-11-5.el8", c.Message)
	c, err = commitIter.Next()
	require.Nil(t, err)
	require.NotNil(t, c)
	require.Equal(t, "import basesystem-11-5.el8", c.Message)
}

func TestWorker_ImportRPM_ChecksumDoesntMatch(t *testing.T) {
	res, err := testW.ImportRPM(
		"memory://efi-rpm-macros-3-3.el8.src.rpm",
		"518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d27",
		"Rocky Linux release 8.8 (Green Obsidian)",
	)
	require.NotNil(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "checksum does not match")
}

func TestWorker_ImportRPM_AuthError(t *testing.T) {
	inmf.noAuthMethod = true

	res, err := testW.ImportRPM(
		"memory://efi-rpm-macros-3-3.el8.src.rpm",
		"518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		"Rocky Linux release 8.8 (Green Obsidian)",
	)
	require.NotNil(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "auth failed")

	inmf.noAuthMethod = false
}

func TestWorker_ImportRPM_InvalidCredentials(t *testing.T) {
	inmf.invalidUsernamePass = true

	res, err := testW.ImportRPM(
		"memory://efi-rpm-macros-3-3.el8.src.rpm",
		"518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		"Rocky Linux release 8.8 (Green Obsidian)",
	)
	require.NotNil(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "username or password incorrect")

	inmf.invalidUsernamePass = false
}
