package base

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	UserContextKey        = "user"
	TokenContextKey       = "token"
	FrontendAuthCookieKey = "auth_bearer"
)

// OidcInterceptorDetails contains the details for the OIDC interceptor
type OidcInterceptorDetails struct {
	Provider             OidcProvider
	Group                string
	AllowUnauthenticated bool
}

// OidcClaims contains the claims for the OIDC token
// At least the ones we care for at the moment
type OidcClaims struct {
	Groups []string
}

// OidcProvider is the interface for OIDC providers
type OidcProvider interface {
	UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (UserInfo, error)
}

// UserInfo is the interface for user info
type UserInfo interface {
	Subject() string
	Email() string
	Claims(v interface{}) error
}

// OidcProviderImpl is the implementation of OidcProvider
// This is main usage in "real" applications
// Tests should use the TestOidcProvider
type OidcProviderImpl struct {
	*oidc.Provider
}

// UserInfo gets the user info from the OIDC provider
func (o *OidcProviderImpl) UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (UserInfo, error) {
	userInfo, err := o.Provider.UserInfo(ctx, tokenSource)
	if err != nil {
		return nil, err
	}

	return &OidcUserInfo{userInfo}, nil
}

// OidcUserInfo is the implementation of UserInfo
type OidcUserInfo struct {
	UserInfo *oidc.UserInfo
}

// Subject gets the subject from the user info
func (o *OidcUserInfo) Subject() string {
	return o.UserInfo.Subject
}

// Email gets the email from the user info
func (o *OidcUserInfo) Email() string {
	return o.UserInfo.Email
}

// Claims gets the claims from the user info
func (o *OidcUserInfo) Claims(v interface{}) error {
	return o.UserInfo.Claims(v)
}

// OidcGrpcInterceptor creates a new OIDC interceptor
// This enforces authentication and authorization
// Authorization is as simple as checking if the user is in a group
// If the group is empty, no authorization is enforced
// Authentication enforcement can be disabled by setting AllowUnauthenticated to true
func OidcGrpcInterceptor(details *OidcInterceptorDetails) (grpc.UnaryServerInterceptor, error) {
	provider := details.Provider

	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			if details.AllowUnauthenticated {
				return handler(ctx, req)
			}

			// let's check if there is a cookie
			cookie := md.Get("cookie")
			if len(cookie) == 0 {
				if details.AllowUnauthenticated {
					return handler(ctx, req)
				}
				return nil, status.Error(codes.Unauthenticated, "missing auth token")
			}

			// parse the cookie
			header := http.Header{}
			header.Add("Cookie", cookie[0])
			req := http.Request{Header: header}

			// verify the token
			cookieToken, err := req.Cookie(FrontendAuthCookieKey)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}

			token = cookieToken.Value
		}

		// verify the token
		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
			TokenType:   "bearer",
		}))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// check if the user is in the group
		if details.Group != "" {
			var claims OidcClaims
			if err := userInfo.Claims(&claims); err != nil {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}

			if !Contains[string](claims.Groups, details.Group) {
				return nil, status.Error(codes.PermissionDenied, "user not in group")
			}
		}

		// add user to context
		ctx = context.WithValue(ctx, UserContextKey, userInfo)

		// add token to context
		ctx = context.WithValue(ctx, TokenContextKey, token)

		return handler(ctx, req)
	}

	return interceptor, nil
}

func UserFromContext(ctx context.Context) (UserInfo, error) {
	user, ok := ctx.Value(UserContextKey).(UserInfo)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user info")
	}

	return user, nil
}

// TestOidcProvider is a test implementation of OidcProvider
type TestOidcProvider struct {
	// This interface is a pointer on purpose, so we can point it to
	// a value in main_test and change it in the tests
	userInfo *UserInfo
}

// NewTestOidcProvider creates a new TestOidcProvider
func NewTestOidcProvider(userInfo *UserInfo) *TestOidcProvider {
	return &TestOidcProvider{
		userInfo: userInfo,
	}
}

// UserInfo gets the user info from the OIDC provider
func (t *TestOidcProvider) UserInfo(_ context.Context, _ oauth2.TokenSource) (UserInfo, error) {
	if t.userInfo == nil {
		return nil, errors.New("no user info")
	}
	return *t.userInfo, nil
}

// TestUserInfo is a test implementation of UserInfo
type TestUserInfo struct {
	subject string
	email   string
	claims  map[string]any
}

// NewTestUserInfo creates a new TestUserInfo
func NewTestUserInfo(subject string, email string, claims map[string]any) *TestUserInfo {
	return &TestUserInfo{
		subject: subject,
		email:   email,
		claims:  claims,
	}
}

// Subject gets the subject from the user info
func (t *TestUserInfo) Subject() string {
	return t.subject
}

// Email gets the email from the user info
func (t *TestUserInfo) Email() string {
	return t.email
}

// Claims gets the claims from the user info
func (t *TestUserInfo) Claims(v *any) error {
	*v = t.claims
	return nil
}
