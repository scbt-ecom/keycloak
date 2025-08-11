package persistent

import (
	"fmt"
	"sync"
	"time"

	"github.com/scbt-ecom/keycloak/v2"
)

type Session struct {
	accessToken string
	err         error

	mu sync.RWMutex
}

func (s *Session) GetAccessToken() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.err != nil {
		return "", s.err
	}

	return s.accessToken, nil
}

func (cl *Client) NewSession(creds keycloak.Credentials) *Session {
	s := &Session{}

	var ticker *time.Ticker

	data, err := cl.authWithCredentials(creds)

	s.mu.Lock()
	if err != nil {
		ticker = time.NewTicker(time.Second * 10)
		s.err = fmt.Errorf("auth with credentials: %s", err.Error())
	} else {
		ticker = time.NewTicker(time.Duration(float64(data.ExpiresIn)*0.95) * time.Second)
		s.accessToken = data.AccessToken
	}
	s.mu.Unlock()

	go func() {
		for {
			select {
			case <-ticker.C:
				data, err = cl.authWithCredentials(creds)
				s.mu.Lock()
				if err != nil {
					s.accessToken = ""
					s.err = fmt.Errorf("auth with credentials: %s", err.Error())
				} else {
					s.err = nil
					s.accessToken = data.AccessToken
				}
				s.mu.Unlock()
			}
		}

	}()

	return s
}
