package groxy

import (
	"context"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// TestResult represents the result of running a test on the given proxy
type TestResult struct {
	Err   error
	Proxy *Proxy
}

// Manager struct controls af the methods used for operating on proxy lists, such as checking validity, response time,
// sorting, and filtering
type Manager struct {
	queue    chan *Proxy
	timeout  time.Duration
	queryURL *url.URL
	inputs   []*Proxy
	outputs  []TestResult
	kill     bool
}

// NewManager constructs a new manager struct, maxConn set the number of connections too use at at time for checking proxies
// timeout sets the timeout to be used for connections, queryUrl sets the url to be used for testing proxies
func NewManager(maxConn int, timeout time.Duration, queryUrl string) (*Manager, error) {
	queue := make(chan *Proxy, maxConn)
	target, err := url.Parse(queryUrl)
	if err != nil {
		return nil, err
	}
	return &Manager{queue: queue, timeout: timeout, queryURL: target, kill: false}, nil
}

// Distinct removes all duplicate proxies for a list
func (m *Manager) Distinct(proxies []*Proxy) []*Proxy {
	keys := make(map[string]bool)
	var list []*Proxy
	for _, proxy := range proxies {
		if _, value := keys[proxy.Host()]; !value {
			keys[proxy.Host()] = true
			list = append(list, proxy)
		}
	}
	return list
}

func (m *Manager) doRequest(proxy *Proxy) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxy.ToURL()),
		},
	}
	req, _ := http.NewRequest("HEAD", m.queryURL.String(), nil)
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (m *Manager) checkProxy(proxy *Proxy) TestResult {
	m.getToken(proxy)
	defer m.releaseToken(proxy)
	isSecure := make(chan bool, 1)
	var result TestResult
	var resultProxy *Proxy

	t0 := time.Now()
	resp, err := m.doRequest(proxy)
	if err != nil {

		result = TestResult{Err: err, Proxy: proxy}
	}
	if resp != nil && resp.Status == "200 OK" {
		go func() {
			defer close(isSecure)
			isSecure <- m.checkIsSecure(proxy)
		}()
		resultProxy.responseTime = time.Since(t0)
		resultProxy.alive = true
		resultProxy.username = proxy.Username()
		resultProxy.password = proxy.Password()
		resultProxy.secure = <-isSecure
		resultProxy.host = proxy.Host()
		result = TestResult{Err: nil, Proxy: resultProxy}
	}
	return result
}
func (m *Manager) checkIsSecure(proxy *Proxy) bool {
	m.getToken(proxy)
	defer m.releaseToken(proxy)
	addr := proxy.ToURL()
	addr.Scheme = "https://"
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(addr),
		},
	}
	req, _ := http.NewRequest("HEAD", m.queryURL.String(), nil)
	req = req.WithContext(ctx)
	_, err := client.Do(req)
	return err == nil
}
func (m *Manager) getToken(proxy *Proxy) {
	m.queue <- proxy
}

func (m *Manager) releaseToken(proxy *Proxy) {
	<-m.queue
}

// Add adds a list of proxies to the manager for checking
func (m *Manager) Add(proxies ...*Proxy) {
	m.inputs = m.Distinct(proxies)
}

// Run starts the proxy checking process
func (m *Manager) Run() {
	res := make(chan TestResult, len(m.inputs))
	defer close(res)
	results := res
	for _, inp := range m.inputs {
		if !m.kill {
			go func(proxy *Proxy) {
				results <- m.checkProxy(proxy)
			}(inp)
			m.outputs = append(m.outputs, <-results)
		} else {
			return
		}
	}
}

func (m *Manager) Proxies() {

}

func (m *Manager) Stop() {
	m.kill = true
}

type Predicate int8

const (
	Alive Predicate = iota
	Secure
	ResponseTime
)

func (m *Manager) Filter(predicate Predicate, time ...time.Duration) []*Proxy {
	var retList []*Proxy
	switch predicate {
	case Secure:
		for _, res := range m.outputs {
			if res.Proxy.Secure() {
				retList = append(retList, res.Proxy)
			}
		}
	case Alive:
		for _, res := range m.outputs {
			if res.Proxy.alive {
				retList = append(retList, res.Proxy)
			}
		}
	case ResponseTime:
		for _, res := range m.outputs {
			if res.Proxy.responseTime < time[0] {
				retList = append(retList, res.Proxy)
			}
		}
	}
	return retList
}

func (m *Manager) SortByResponseTime() []*Proxy {
	var proxies []*Proxy
	for _, val := range m.outputs {
		proxies = append(proxies, val.Proxy)
	}
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].responseTime < proxies[j].responseTime
	})
	return proxies
}
