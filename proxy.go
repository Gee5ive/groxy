package groxy

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/url"
	"os"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

// Proxy represents an http proxy used by accounts for accessing the internet
type Proxy struct {
	host         string
	username     string
	password     string
	secure       bool
	responseTime time.Duration
	alive        bool
}

// Scheme returns the scheme used for the proxy, if it has been tested and is secure, the scheme will be https
func (h *Proxy) Scheme() string {
	if h.Secure() {
		return "https://"
	}
	return "http://"
}

// Host portion of the proxy as a string
func (h *Proxy) Host() string {
	return h.Scheme() + h.host
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

func (h *Proxy) Secure() bool {
	return h.secure
}

func (h *Proxy) AsCSV() []string {
	return []string{h.host, h.username, h.password}
}

// New returns a pointer to a proxy, if the info provided cannot be validated it uses the host machine's preferred
// outbound ip Address.
func New(uri string, username string, password string) (*Proxy, error) {
	_, err := url.Parse("http://" + uri)
	if err != nil {
		return nil, err

	}
	proxy := Proxy{host: uri, username: username, password: password}
	return &proxy, nil

}

func SaveToFile(file string, proxies []*Proxy) error {
	var result error
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(bufio.NewWriter(f))
	defer writer.Flush()
	defer f.Close()
	for _, proxy := range proxies {
		if err := writer.Write(proxy.AsCSV()); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

func FromFile(file string) ([]*Proxy, error) {
	var result error
	handleErr := func(err error) {
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	f, err := os.Open(file)
	if err != nil {
		return []*Proxy{}, err
	}
	defer f.Close()
	reader := csv.NewReader(bufio.NewReader(f))
	reader.FieldsPerRecord = -1
	var proxies []*Proxy
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return []*Proxy{}, err
		}
		switch len(line) {
		case 0:
			continue
		case 1:
			proxy, err := New(line[0], "", "")
			if err != nil {
				handleErr(err)
			}
			proxies = append(proxies, proxy)
		case 2:
			proxy, err := New(line[0], line[1], "")
			if err != nil {
				handleErr(err)
			}
			proxies = append(proxies, proxy)

		case 3:
			proxy, err := New(line[0], line[1], line[2])
			if err != nil {
				handleErr(err)
			}
			proxies = append(proxies, proxy)
		}

	}

	return proxies, result
}
