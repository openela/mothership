package mothership_worker_server

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetObjectPath_Path_S3(t *testing.T) {
	object, err := getObjectPath("s3://mship/test.rpm")
	require.Nil(t, err)
	require.Equal(t, "test.rpm", object)
}

func TestGetObjectPath_Host_Memory(t *testing.T) {
	object, err := getObjectPath("memory://test.rpm")
	require.Nil(t, err)
	require.Equal(t, "test.rpm", object)
}

func TestGetObjectPath_InvalidURI(t *testing.T) {
	_, err := getObjectPath("test://test:test/")
	require.NotNil(t, err)
}
