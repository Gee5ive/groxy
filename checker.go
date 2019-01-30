package groxy

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type TestResult struct {
	Err   error
	Proxy *Proxy
}

type Checker struct {
	queue    chan *Proxy
	timeout  time.Duration
	queryURL *url.URL
	inputs   []*Proxy
	results  chan TestResult
	kill     bool
}

func NewProxyChecker(maxConn int, timeout time.Duration, queryUrl string) (*Checker, error) {
	queue := make(chan *Proxy, maxConn)
	target, err := url.Parse(queryUrl)
	if err != nil {
		return nil, err
	}
	return &Checker{queue: queue, timeout: timeout, queryURL: target, kill: false}, nil
}

func (c *Checker) dedup(proxies []*Proxy) []*Proxy {
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

func (c *Checker) checkProxy(proxy *Proxy) TestResult {
	c.getToken(proxy)
	defer c.releaseToken(proxy)
	isSecure := make(chan bool, 1)
	var result TestResult
	var resultProxy *Proxy
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxy.ToURL()),
		},
	}
	req, _ := http.NewRequest("GET", c.queryURL.String(), nil)
	t0 := time.Now()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {

		result = TestResult{Err: err, Proxy: proxy}
	}
	if resp != nil && resp.Status == "200 OK" {
		go func() {
			defer close(isSecure)
			isSecure <- c.checkIsSecure(proxy)
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
func (c *Checker) checkIsSecure(proxy *Proxy) bool {
	c.getToken(proxy)
	defer c.releaseToken(proxy)
	addr := proxy.ToURL()
	addr.Scheme = "https://"
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(addr),
		},
	}
	req, _ := http.NewRequest("GET", c.queryURL.String(), nil)
	req = req.WithContext(ctx)
	_, err := client.Do(req)
	return err == nil
}
func (c *Checker) getToken(proxy *Proxy) {
	c.queue <- proxy
}

func (c *Checker) releaseToken(proxy *Proxy) {
	<-c.queue
}

func (c *Checker) Add(proxies ...*Proxy) {
	c.inputs = c.dedup(proxies)
}
func (c *Checker) Run() {
	res := make(chan TestResult, len(c.inputs))
	c.results = res
	for _, inp := range c.inputs {
		if !c.kill {
			go func(proxy *Proxy) {
				c.results <- c.checkProxy(proxy)
			}(inp)
		}
	}
}

func (c *Checker) Stop() {
	c.kill = true
}
