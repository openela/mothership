package forge

import "time"

type Cacher struct {
	Forge

	authenticator *Authenticator
}

func NewCacher(f Forge) *Cacher {
	return &Cacher{
		Forge: f,
	}
}

func (c *Cacher) GetAuthenticator() (*Authenticator, error) {
	if c.authenticator != nil {
		// Check if the token is expired, if not return it
		if !c.authenticator.Expires.Before(time.Now()) {
			return c.authenticator, nil
		}
	}

	// Otherwise, get a new token
	c.authenticator = nil
	a, err := c.Forge.GetAuthenticator()
	if err != nil {
		return nil, err
	}

	c.authenticator = a

	return a, nil
}
