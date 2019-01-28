package groxy

import (
	"net/url"
	"time"
)

// Proxy represents an http proxy used by accounts for accessing the internet
type Proxy struct {
	scheme       string
	host         string
	username     string
	password     string
	secure       bool
	responseTime time.Duration
	alive        bool
}

// Host portion of the proxy as a string
func (h *Proxy) Host() string {
	return h.scheme + h.host
}

// Username portion of the proxy if present
func (h *Proxy) Username() string {
	return h.username
}

// Password portion of the proxy if present
func (h *Proxy) Password() string {
	return h.password
}

// ToURL converts the proxy to a *url.URL if possible
func (h *Proxy) ToURL() *url.URL {
	rsp, _ := url.Parse(h.Host())
	if len(h.username) > 0 && len(h.password) > 0 {
		rsp.User = url.UserPassword(h.username, h.password)
	}
	return rsp
}

func (h *Proxy) CheckSecure() bool {
	return true
}

// New returns a pointer to a proxy, if the info provided cannot be validated it uses the host machine's preferred
// outbound ip Address.
func New(uri string, username string, password string) (*Proxy, error) {
	_, err := url.Parse(uri)
	if err != nil {
		return nil, err

	}
	proxy := Proxy{host: uri, username: username, password: password}
	return &proxy, nil

}
