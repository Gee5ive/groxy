package groxy

import (
	ctx "context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gammazero/workerpool"
)

const ipCheckAddr = `https://ip4.seeip.org`
const proxyAnonCheckAddr = `http://bot.whatismyipaddress.com`

// TestResult represents the result of running a test on the given proxy
type TestResult struct {
	Err   error
	Proxy *Proxy
}

// Manager struct controls af the methods used for operating on proxy lists, such as checking validity, response time,
// sorting, and filtering
type Manager struct {
	pool     *workerpool.WorkerPool
	timeout  time.Duration
	queryURL *url.URL
	inputs   []*Proxy
	ctx      ctx.Context
	done     ctx.CancelFunc
	realIPs  []string
}

// NewManager constructs a new manager struct, maxConn set the number of connections too use at at time for checking proxies
// timeout sets the timeout to be used for connections, queryUrl sets the url to be used for testing proxies
func NewManager(maxConn int, timeout time.Duration, queryURL string) (*Manager, error) {
	target, err := url.Parse(queryURL)
	if err != nil {
		return nil, err
	}
	proxyCtx, cancel := ctx.WithCancel(ctx.Background())
	return &Manager{pool: workerpool.New(maxConn), timeout: timeout, queryURL: target, inputs: []*Proxy{}, ctx: proxyCtx, done: cancel, realIPs: GetIPs()}, nil
}

// Distinct removes all duplicate proxies from a list
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

	client := &http.Client{
		Timeout: m.timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxy.ToURL()),
		},
	}
	req, _ := http.NewRequest("GET", m.queryURL.String(), nil)
	req.Close = true
	resp, err := client.Do(req)
	return resp, err
}

func getIPV4() string {
	var bodyString string
	resp, err := http.Get(ipCheckAddr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		bodyString = string(bodyBytes)
	}
	return bodyString
}

func getIPV6() string {
	var bodyString string
	resp, err := http.Get(ipCheckAddr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		bodyString = string(bodyBytes)
	}
	return bodyString
}
func GetIPs() []string {
	var r []string

	ips := []string{getIPV4(), getIPV6()}
	for _, ip := range ips {
		if ip != "" {
			r = append(r, ip)
		}
	}
	return r
}
func (m *Manager) isAnon(proxy *Proxy) bool {

	client := &http.Client{
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxy.ToURL()),
		},
	}
	req, _ := http.NewRequest("GET", proxyAnonCheckAddr, nil)
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		bodyString = string(bodyBytes)
	}

	for _, ip := range m.realIPs {
		if ip == bodyString {
			return false
		}
	}
	return true

}
func (m *Manager) checkProxy(proxy *Proxy) TestResult {
	var result TestResult
	resultProxy := &Proxy{}

	t0 := time.Now()
	resp, err := m.doRequest(proxy)
	if err != nil {
		return TestResult{Err: err, Proxy: proxy}

	}
	if resp != nil && resp.Status == "200 OK" {
		anon := m.isAnon(proxy)
		resultProxy.responseTime = time.Since(t0)
		resultProxy.alive = true
		resultProxy.url = proxy.ToURL()
		resultProxy.transParent = anon
		resultProxy.url.User = url.UserPassword(proxy.Username(), proxy.Password())
		result = TestResult{Err: nil, Proxy: resultProxy}
	}
	return result
}

// Add adds a list of proxies to the manager for checking
func (m *Manager) Add(proxies ...*Proxy) {
	m.inputs = m.Distinct(proxies)
}

// Run starts the proxy checking process
func (m *Manager) Run() <-chan TestResult {
	results := make(chan TestResult)
	go func() {
		defer close(results)
		for _, proxy := range m.inputs {
			check := func(prox *Proxy) func() {
				return func() {
					results <- m.checkProxy(prox)
				}
			}
			m.pool.Submit(check(proxy))
		}
		m.pool.StopWait()
	}()

	return results
}

// Stop forces the proxy checking process abort
func (m *Manager) Stop() {
	m.done()
}

func (m *Manager) Inputs() []*Proxy {
	return m.inputs
}
