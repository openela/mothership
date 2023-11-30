package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/openela/mothership/base"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func (s *server) authLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Generate a random state
	state := ""
	for i := 0; i < 16; i++ {
		state += strconv.Itoa(rand.Intn(10))
	}

	// Generate the auth url
	authURL := s.oauth2Config.AuthCodeURL(state)

	// Set the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "auth_state",
		Value: state,
		Path:  "/",
		// expires in 2 minutes
		MaxAge: 120,
		// secure if self is https
		Secure: strings.HasPrefix(s.frontendInfo.Self, "https://"),
	})

	// Redirect to the auth url
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (s *server) authCallbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the state cookie
	stateCookie, err := r.Cookie("auth_state")
	if err != nil {
		s.errorView("No state cookie")(w, r, nil)
		return
	}

	// Get the state query param
	stateQueryParam := r.URL.Query().Get("state")
	if stateQueryParam == "" {
		s.errorView("No state query param")(w, r, nil)
		return
	}

	// Check if the state cookie and state query param match
	if stateCookie.Value != stateQueryParam {
		s.errorView("State cookie and state query param do not match")(w, r, nil)
		return
	}

	// Exchange the code for a token
	token, err := s.oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		s.errorView(fmt.Sprintf("Failed to exchange code for token: %v", err))(w, r, nil)
		return
	}

	// Verify the token
	accessToken := token.AccessToken
	userInfo, err := s.oidcProvider.UserInfo(r.Context(), oauth2.StaticTokenSource(token))
	if err != nil {
		s.errorView(fmt.Sprintf("Failed to get user info: %v", err))(w, r, nil)
		return
	}

	// Check if the user is in the group
	if s.frontendInfo.OIDCGroup != "" {
		var claims base.OidcClaims
		err := userInfo.Claims(&claims)
		if err != nil {
			s.errorView(fmt.Sprintf("Failed to get user claims: %v", err))(w, r, nil)
			return
		}

		groups := claims.Groups

		found := false
		for _, group := range groups {
			if group == s.frontendInfo.OIDCGroup {
				found = true
				break
			}
		}

		if !found {
			s.errorView(fmt.Sprintf("User is not in group %s", s.frontendInfo.OIDCGroup))(w, r, nil)
			return
		}
	}

	// Set the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:  base.FrontendAuthCookieKey,
		Value: accessToken,
		Path:  "/",
		// expires in 2 hours
		MaxAge: 7200,
		// secure if self is https
		Secure: strings.HasPrefix(s.frontendInfo.Self, "https://"),
	})

	// Redirect to self, this is due to the "root" not being / for all apps
	http.Redirect(w, r, s.frontendInfo.Self, http.StatusFound)
}

func (s *server) authLogoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Delete the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:   base.FrontendAuthCookieKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		// secure if self is https
		Secure: strings.HasPrefix(s.frontendInfo.Self, "https://"),
	})

	// Redirect to self, this is due to the "root" not being / for all apps
	http.Redirect(w, r, s.frontendInfo.Self, http.StatusFound)
}

func (s *server) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excludedPaths := []string{
			"/auth/login",
			"/auth/callback",
			"/auth/logout",
		}

		// Check if the path is excluded
		for _, path := range excludedPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				h.ServeHTTP(w, r)
				return
			}
		}

		ctx := r.Context()

		// get auth cookie
		authCookie, err := r.Cookie(base.FrontendAuthCookieKey)
		if err == nil {
			// verify the token
			userInfo, err := s.oidcProvider.UserInfo(r.Context(), oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: authCookie.Value,
				TokenType:   "Bearer",
			}))
			if err != nil {
				// redirect to login
				http.Redirect(w, r, s.frontendInfo.Self+"/auth/login", http.StatusFound)
				return
			}

			// Check if the user is in the group
			var claims base.OidcClaims
			err = userInfo.Claims(&claims)
			if err != nil {
				// redirect to login
				http.Redirect(w, r, s.frontendInfo.Self+"/auth/login", http.StatusFound)
				return
			}

			groups := claims.Groups
			if s.frontendInfo.OIDCGroup != "" {
				if !base.Contains(groups, s.frontendInfo.OIDCGroup) {
					// show unauthenticated page
					s.errorView(fmt.Sprintf("User is not in group %s", s.frontendInfo.OIDCGroup))(w, r, nil)
					return
				}
			}

			// Add the user to the context
			ctx = context.WithValue(ctx, base.UserContextKey, userInfo)

			// Add the token to the context
			ctx = context.WithValue(ctx, base.TokenContextKey, authCookie.Value)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
