package forge

import (
	transport_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testForge struct {
	Forge

	callToGetAuthenticator int
}

func (t *testForge) GetAuthenticator() (*Authenticator, error) {
	t.callToGetAuthenticator++
	return &Authenticator{
		AuthMethod: &transport_http.BasicAuth{
			Username: "test",
			Password: "test",
		},
		AuthorName:  "test",
		AuthorEmail: "test@resf.org",
		Expires:     time.Now().Add(time.Minute * 45),
	}, nil
}

func TestNewCacher(t *testing.T) {
	f := &testForge{}
	fAuth, err := f.GetAuthenticator()
	require.Nil(t, err)
	fAuthAuthMethod := fAuth.AuthMethod.(*transport_http.BasicAuth)
	require.Equal(t, "test", fAuthAuthMethod.Username)
	require.Equal(t, "test", fAuthAuthMethod.Password)

	c := NewCacher(f)
	cAuth, err := c.GetAuthenticator()
	require.Nil(t, err)
	cAuthAuthMethod := cAuth.AuthMethod.(*transport_http.BasicAuth)
	require.Equal(t, "test", cAuthAuthMethod.Username)
	require.Equal(t, "test", cAuthAuthMethod.Password)
}

func TestCacher_GetAuthenticator(t *testing.T) {
	f := &testForge{}
	c := NewCacher(f)

	cAuth, err := c.GetAuthenticator()
	require.Nil(t, err)
	cAuthAuthMethod := cAuth.AuthMethod.(*transport_http.BasicAuth)
	require.Equal(t, "test", cAuthAuthMethod.Username)
	require.Equal(t, "test", cAuthAuthMethod.Password)
}

func TestCacher_GetAuthenticator_Cached(t *testing.T) {
	f := &testForge{}
	c := NewCacher(f)

	cAuth, err := c.GetAuthenticator()
	require.Nil(t, err)
	require.Equal(t, "test", cAuth.AuthorName)

	cAuth, err = c.GetAuthenticator()
	require.Nil(t, err)
	require.Equal(t, "test", cAuth.AuthorName)

	require.Equal(t, 1, f.callToGetAuthenticator)
}

func TestCacher_GetAuthenticator_Expired(t *testing.T) {
	f := &testForge{}
	c := NewCacher(f)

	cAuth, err := c.GetAuthenticator()
	require.Nil(t, err)
	require.Equal(t, "test", cAuth.AuthorName)

	cAuth.Expires = time.Now().Add(-time.Minute * 10)

	cAuth, err = c.GetAuthenticator()
	require.Nil(t, err)
	require.Equal(t, "test", cAuth.AuthorName)

	require.Equal(t, 2, f.callToGetAuthenticator)
}
