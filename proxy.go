package groxy

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
)

// ID is a uuid used to distinguish between objects in the system
type ID = uuid.UUID

// NewID returns an ID instance
func NewID() ID {
	return uuid.New()
}

func IDFromString(id string) ID {
	if id == "" {
		return NewID()
	}
	return uuid.MustParse(id)
}

// Proxy represents an http proxy used for accessing the internet anonymously
type Proxy struct {
	id           ID
	host         string
	username     string
	password     string
	secure       bool
	responseTime time.Duration
	alive        bool
}

func (h *Proxy) Id() string {
	return h.id.String()
}

// Scheme returns the scheme used for the proxy, if it has been tested and is secure, the scheme will be https
func (h *Proxy) Scheme() string {
	if h.Secure() {
		return "https://"
	}
	return "http://"
}

// Host returns the host portion of the proxy as a string
func (h *Proxy) Host() string {
	return h.Scheme() + h.host
}

// Username returns the username portion of the proxy if present
func (h *Proxy) Username() string {
	return h.username
}

// Password returns the password portion of the proxy if present
func (h *Proxy) Password() string {
	return h.password
}

// Alive returns whether or not the proxy is dead
func (h *Proxy) Alive() bool {
	return h.alive
}

// ResponseTime returns the proxy response time
func (h *Proxy) ResponseTime() time.Duration {
	return h.responseTime
}

// ToURL converts the proxy to a *url.URL
func (h *Proxy) ToURL() *url.URL {
	rsp, _ := url.Parse(h.Host())
	if len(h.username) > 0 && len(h.password) > 0 {
		rsp.User = url.UserPassword(h.username, h.password)
	}
	return rsp
}

// Secure identifies whether or not the proxy is secure
func (h *Proxy) Secure() bool {
	return h.secure
}

// AsCSV converts the proxy to csv format for saving to disk
func (h *Proxy) AsCSV() []string {
	return []string{h.host, h.username, h.password}
}

// New returns a pointer to a proxy, if the url provided cannot be parsed it returns an error
func New(uri string, username string, password string) (*Proxy, error) {
	_, err := url.Parse("http://" + uri)
	if err != nil {
		return nil, err

	}
	proxy := Proxy{id: NewID(), host: uri, username: username, password: password}
	return &proxy, nil

}

// SaveToFile saves a list of proxies to a CSV file
func SaveToFile(file string, proxies []*Proxy) error {
	var result error
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(bufio.NewWriter(f))

	defer f.Close()
	for _, proxy := range proxies {
		if err := writer.Write(proxy.AsCSV()); err != nil {
			result = multierror.Append(result, err)
		}
	}
	writer.Flush()
	return result
}

// FromFile loads a list of proxies from a file on disk, it returns an error if there is a problem parsing the file
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

// FromExisting returns a proxy from an existing  proxy that was used with this package, normally a database
func FromExisting(
	id string,
	host string,
	username string,
	password string,
	secure bool,
	responseTime time.Duration,
	alive bool) *Proxy {
	return &Proxy{id: IDFromString(id), host: host, username: username, password: password, secure: secure, responseTime: responseTime, alive: alive}

}
