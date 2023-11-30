package github_forge

import (
	transport_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCrCn8kH3AjWG36
BH20WgmZBBpOqHTVjMQtH0WXl+7+qhAldpYeoodbKR0KhraEh0EPez0wAnlx5dFN
7EwdajoznaMZkzmjoqa1Bmd4IscUtDdZPEHDMadYe5Go0PvzgDNBviPVEi0E82lE
zpqMClKuUBS8X+lypZGUqQXdCxbWATihS7C6/Opmw8pNi7GZw4SrVYGk6uKFMH/W
eB3SMpifqAS3nQAoTjNj7m8RNhICjov9Xtn3GOCvukqO7FKNirGyqBHq4w/d+ETu
Ni58mkZHDNnDvBJXDXGfqt3G/d4OL2fQ+ZpF7cJCjGtAmidImgv36Hbs8uKM35ep
9xkMFVABAgMBAAECggEAC/JQPhQqD3XyPIfKxemSCQuF0N+oRXAvFZ29DSESEtyL
AyrcwmgEv0PIYP9WyTvvOecYN328wM1WCLTL/jP4u7kzdqpXWMwYC8XWPUhkklgi
E4wHZdxWfXIoEtwB9RiLu/hNZWj/FvzvadxqZASmkMFMKXojgpv45qLFo5JONoVC
y0GX4cApG0ALP1drIkluOq5yYZmjpbrylOeQngBzJ3wPzeEHDwAKk5N/pqCjrtfY
ReW3FVW+WgxydCv2S1zn5cwTSBflqz+42ob0hUnous1FA05E7/Wzsq66vQ/3nJD1
56yi75uTE6PQ2FgxW4tv7IxHCrsR0hpD44UFS9+pvwKBgQDw50WRpVyXfmbGge+y
dbiYNUtUrCts26STikbRVvYn/WAoDXxzyuw2GEsQoWrGR1BfwRoYCOP8tMdwHq0v
+iwFnjmvufGeIud2S+eQLAJqhXXNBrKK92Jlpz0K8HHqQ21KHvAIEPFSEclADNwY
8rEJm1J8tzCcmiYVL6y7Wnm5cwKBgQC1wnJxG7hEn6bH7urMIAMtVCOP/ExLchdF
wlPU6JEvyqzKsBCOmapE5tsWD88xeascDzFvsFzIefebWPx5NuKrl9SHCWvvVOxk
zwCgmBkhNEQJ03nCFTklK+J9MxB1ZeRma1hO8VC8Z75N1tHXp6g4lR2QdDi/vqcq
oWM8OmeDuwKBgETcBqG8H7xZ8Cy7xXVAexRe33qDgCIsol1eACIkdlY18b9hI3rB
vUU1KnfFfAzTI6FLRBcsq2Z3ki51RlHZc63jbV/SicMG/RxuU/F88u/Z2DNTv8ND
NUgTRrqSwi0ROvMd5sSXezNXTCxXwK4M6Rfy4uAtSOLqmQojR3+CPBsLAoGAJkzt
NKx0rfE+gc70p0LvqHOcctDleth10vtaEvlW7s00kBl9w67Z1F8ZN5LpRDGxPt5s
um5dftlEtfWQbjKEnUgHPtVbazln/u4n4a9rTDXpSHDJrX4vZofS2DMUesiX0oU4
PJpZOvpZfamQ2nK33gR+EFyNQMp6C1+qu5xLB9UCgYEAmfR5E1gDRS2oqtnWlWuI
eukgvnao0uSDWK7K4I9y6ZH6MJPAMsSHS2ppNgeiWwLneXz2bI4wxKQCMlj3ZMZb
MAX21piZSr0bp8evE9AHFIXU8Tfj8II6fY5KcWm+3AIae3NQXTtxR+k8wwf5PPDn
KDqvQovQtfR2TdH4t56YHc0=
-----END PRIVATE KEY-----
`

func TestNew(t *testing.T) {
	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	require.Equal(t, "test-org", forge.organization)
	require.Equal(t, "123", forge.appId)
	require.NotNil(t, forge.appPrivateKey)
	require.False(t, forge.shouldMakeRepoPublic)

	require.Equal(t, 65537, forge.appPrivateKey.PublicKey.E)
	require.Equal(t, "21591926237731527940797012101171439523542979145099072456746913347900286948930956360285480936616597570038038179727935785740932577232980060401308780474432320730685823408025458821807894246437469943261837548491954465021772289266754949738031497460967955016818560488471010997875191087486755196596143982294619114076489100464465659028521281331191374862147882697445608174211739135375236895595372412671912637593012110128170709873050090012617276785952825931935301790521557448256516848668806692526033364651010065962357562265998842334259601221952975920205467703506376745867644890126749157744786044647764818435588205721145469915137", forge.appPrivateKey.PublicKey.N.String())
}

func TestGetInstallationToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/orgs/test-org/installation",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":       123456,
			"app_slug": "test_app",
		}))

	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/123456/access_tokens",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"token": "test_token",
		}))

	it, err := forge.getInstallationToken("test-jwt")
	require.Nil(t, err)
	require.Equal(t, "test_token", it.Token)
	require.Equal(t, "test_app", it.AppSlug)
}

func TestGetInstallationTokenError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/orgs/test-org/installation",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":       123456,
			"app_slug": "test_app",
		}))

	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/123456/access_tokens",
		httpmock.NewJsonResponderOrPanic(500, map[string]interface{}{
			"message": "test error",
		}))

	it, err := forge.getInstallationToken("test-jwt")
	require.NotNil(t, err)
	require.Nil(t, it)
	require.Equal(t, "token not found in response", err.Error())
}

func TestGetMeID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/users/test_app[bot]",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id": 123456,
		}))

	id, err := forge.GetMeID("test_token", "test_app")
	require.Nil(t, err)
	require.Equal(t, "123456", id)
}

func TestGetMeIDError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/users/test_app[bot]",
		httpmock.NewJsonResponderOrPanic(500, map[string]interface{}{
			"message": "test error",
		}))

	id, err := forge.GetMeID("test_token", "test_app")
	require.NotNil(t, err)
	require.Equal(t, "", id)
	require.Equal(t, "id not found in response", err.Error())
}

func TestGetRemote(t *testing.T) {
	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	remote := forge.GetRemote("test")
	require.Equal(t, "https://github.com/test-org/test", remote)
}

func TestGetCommitViewerURL(t *testing.T) {
	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	url := forge.GetCommitViewerURL("test", "123456")
	require.Equal(t, "https://github.com/test-org/test/commit/123456", url)
}

func TestGetAuthenticator(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/orgs/test-org/installation",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":       123456,
			"app_slug": "test_app",
		}))

	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/123456/access_tokens",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"token": "test_token",
		}))

	httpmock.RegisterResponder("GET", "https://api.github.com/users/test_app[bot]",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id": 123456,
		}))

	auth, err := forge.GetAuthenticator()
	require.Nil(t, err)
	require.Equal(t, "test_app[bot]", auth.AuthorName)
	require.Equal(t, "123456+test_app[bot]@users.noreply.github.com", auth.AuthorEmail)

	// Ensure expires is at least 44 minutes in the future
	require.True(t, auth.Expires.After(time.Now().Add(44*time.Minute)))

	// Cast AuthMethod to BasicAuth
	basic := auth.AuthMethod.(*transport_http.BasicAuth)
	require.Equal(t, "123", basic.Username)
	require.Equal(t, "test_token", basic.Password)
}

func TestGetAuthenticatorError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/orgs/test-org/installation",
		httpmock.NewJsonResponderOrPanic(500, map[string]interface{}{
			"message": "test error",
		}))

	auth, err := forge.GetAuthenticator()
	require.NotNil(t, err)
	require.Nil(t, auth)
	require.Equal(t, "id not found in response", err.Error())
}

func TestEnsureRepositoryExists_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/orgs/test-org/installation",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":       123456,
			"app_slug": "test_app",
		}))

	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/123456/access_tokens",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"token": "test_token",
		}))

	httpmock.RegisterResponder("GET", "https://api.github.com/users/test_app[bot]",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id": 123456,
		}))

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/test-org/test",
		httpmock.NewJsonResponderOrPanic(404, map[string]interface{}{
			"message": "Not Found",
		}))

	httpmock.RegisterResponder("POST", "https://api.github.com/orgs/test-org/repos",
		httpmock.NewJsonResponderOrPanic(201, map[string]interface{}{
			"name": "test",
		}))

	auth, err := forge.GetAuthenticator()
	require.Nil(t, err)

	err = forge.EnsureRepositoryExists(auth, "test")
	require.Nil(t, err)

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	require.Equal(t, 1, info["POST https://api.github.com/orgs/test-org/repos"])
}

func TestEnsureRepositoryExists_AlreadyExists(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	forge, err := New("test-org", "123", []byte(testPrivateKey), false)
	require.Nil(t, err)

	httpmock.RegisterResponder("GET", "https://api.github.com/orgs/test-org/installation",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":       123456,
			"app_slug": "test_app",
		}))

	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/123456/access_tokens",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"token": "test_token",
		}))

	httpmock.RegisterResponder("GET", "https://api.github.com/users/test_app[bot]",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id": 123456,
		}))

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/test-org/test",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"name": "test",
		}))

	auth, err := forge.GetAuthenticator()
	require.Nil(t, err)

	err = forge.EnsureRepositoryExists(auth, "test")
	require.Nil(t, err)

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	require.Equal(t, 0, info["POST https://api.github.com/orgs/test-org/repos"])
}
