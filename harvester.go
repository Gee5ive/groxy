package groxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

// ProviderResponse is a struct returned from a proxy provider function, it contains a list of proxies and a possible error
// which in some cases may be a github.com/hashicorp/go-multierror
type ProviderResponse struct {
	Proxies []*Proxy
	Err     error
}

// Provider is any function that can fetch proxies from a remote location, the function returns a provider response
type Provider func() ProviderResponse

// Harvester is a struct that uses a list of provider functions to harvest proxies and stores the harvested proxies in a slice
type Harvester struct {
	providers []Provider
	proxies   []*Proxy
}

// NewHarvester constructs a new harvester struct using the list of provider functions passed in as arguments to harvest proxies
// Use the WithAllProviders function to construct the harvester with all available providers
func NewHarvester(providers ...Provider) *Harvester {
	return &Harvester{providers: providers}
}

// Harvest fetches proxies using the list of providers contained in Harvester's internal providers list
// The results are stored in the proxies list and can be obtained using the Proxies() method
func (h *Harvester) Harvest() {
	for _, provider := range h.providers {
		resp := provider()
		if resp.Err == nil {
			h.proxies = append(h.proxies, resp.Proxies...)
		}
	}
}

// Proxies returns the list of proxies contained in the Harvester struct
func (h *Harvester) Proxies() []*Proxy {
	return h.proxies
}

func getBody(site string) (string, error) {
	var bodyString string
	client := http.Client{}
	resp, err := client.Get(site)
	if err != nil {
		return "", err

	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return "", err
		}
		bodyString = string(bodyBytes)

	}
	defer resp.Body.Close()
	return bodyString, nil
}

// ProxyListDL is a provider which fetches proxies from https://www.proxy-list.download
func ProxyListDL() ProviderResponse {
	var resultErr error
	var proxies []*Proxy
	kinds := []string{"http", "https"}
	respStream := make(chan ProviderResponse, len(kinds))
	get := func(kind string) {
		var list []*Proxy
		site := fmt.Sprintf("https://www.proxy-list.download/api/v1/get?type=%s", kind)

		bodyString, err := getBody(site)
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return
		}

		for _, ip := range strings.Fields(bodyString) {
			proxy := New(ip, "", "")

			list = append(list, proxy)
		}
		respStream <- ProviderResponse{Proxies: list, Err: nil}

	}

	wg := sync.WaitGroup{}
	wg.Add(len(kinds))
	for _, kind := range kinds {
		go func(kind string, wg *sync.WaitGroup) {
			defer wg.Done()
			get(kind)
		}(kind, &wg)
	}
	wg.Wait()

	close(respStream)
	for result := range respStream {
		proxies = append(proxies, result.Proxies...)
		if result.Err != nil {
			resultErr = multierror.Append(resultErr, result.Err)
		}
	}
	return ProviderResponse{Proxies: proxies, Err: resultErr}
}

// FateProxyList is a provider which fetches proxies from https://raw.githubusercontent.com/fate0/proxylist
func FateProxyList() ProviderResponse {
	parseJSON := func(jsonStr string) []string {
		type Resp struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}
		var proxies []string
		list := strings.Split(jsonStr, "\n")
		for _, item := range list {
			var v Resp
			itemBytes := []byte(item)
			if err := json.Unmarshal(itemBytes, &v); err != nil {
				continue
			}
			proxies = append(proxies, v.Host+":"+strconv.Itoa(v.Port))
		}
		return proxies
	}
	respStream := make(chan ProviderResponse)
	get := func() {
		defer close(respStream)
		var list []*Proxy
		resp, err := getBody("https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list")
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return
		}
		proxies := parseJSON(resp)
		for _, proxy := range proxies {
			host := New(proxy, "", "")
			list = append(list, host)
		}
		respStream <- ProviderResponse{Proxies: list, Err: nil}
	}
	go func() {
		get()
	}()
	return <-respStream
}

// ClarkTMProxy is a provider which fetches proxies from https://raw.githubusercontent.com/clarketm/proxy-list
func ClarkTMProxy() ProviderResponse {
	respStream := make(chan ProviderResponse)
	parseProxies := func(proxies string) []string {
		list := strings.Split(proxies, "\n")
		list = list[4 : len(list)-2]
		for i, item := range list {
			list[i] = strings.Fields(item)[0]
		}
		return list
	}
	get := func() {
		var proxies []*Proxy
		defer close(respStream)
		resp, err := getBody("https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list.txt")
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return
		}
		ips := parseProxies(resp)
		for _, value := range ips {
			proxy := New(value, "", "")
			proxies = append(proxies, proxy)

		}
		respStream <- ProviderResponse{Proxies: proxies, Err: nil}
	}
	go func() {
		get()
	}()
	return <-respStream
}

// MultiProxy is a provider which fetches proxies from http://multiproxy.org/
func MultiProxy() ProviderResponse {
	respStream := make(chan ProviderResponse)
	get := func() {
		var proxies []*Proxy
		defer close(respStream)
		resp, err := getBody("http://multiproxy.org/txt_all/proxy.txt")
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return
		}
		for _, value := range strings.Split(resp, "\n") {
			proxy := New(value, "", "")
			proxies = append(proxies, proxy)
		}
		respStream <- ProviderResponse{Proxies: proxies, Err: nil}
	}
	go func() {
		get()
	}()
	return <-respStream
}

// SpysME is a provider which fetches proxies from http://spys.me/
func SpysME() ProviderResponse {
	respStream := make(chan ProviderResponse)
	parseProxies := func(proxies string) []string {
		list := strings.Split(proxies, "\n")
		list = list[4 : len(list)-2]
		for i, item := range list {
			list[i] = strings.Fields(item)[0]
		}
		return list
	}
	get := func() {
		var proxies []*Proxy
		defer close(respStream)
		resp, err := getBody("http://spys.me/proxy.txt")
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return
		}
		ips := parseProxies(resp)
		for _, value := range ips {
			proxy := New(value, "", "")

			proxies = append(proxies, proxy)

		}
		respStream <- ProviderResponse{Proxies: proxies, Err: nil}
	}
	go func() {
		get()
	}()
	return <-respStream
}

// ProxyListNET is a provider which fetches proxies from http://www.proxylists.net/
func ProxyListNET() ProviderResponse {
	var resultErr error
	var proxies []*Proxy
	kinds := []string{"http", "http_highanon"}
	respStream := make(chan ProviderResponse, len(kinds))
	get := func(kind string) {
		var list []*Proxy
		site := fmt.Sprintf("http://www.proxylists.net/%s.txt", kind)

		bodyString, err := getBody(site)
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return
		}

		for _, ip := range strings.Fields(bodyString) {
			proxy := New(ip, "", "")

			list = append(list, proxy)
		}
		respStream <- ProviderResponse{Proxies: list, Err: nil}

	}

	wg := sync.WaitGroup{}
	wg.Add(len(kinds))
	for _, kind := range kinds {
		go func(kind string, wg *sync.WaitGroup) {
			defer wg.Done()
			get(kind)
		}(kind, &wg)
	}
	wg.Wait()

	close(respStream)
	for result := range respStream {
		proxies = append(proxies, result.Proxies...)
		if result.Err != nil {
			resultErr = multierror.Append(resultErr, result.Err)
		}

	}
	return ProviderResponse{Proxies: proxies, Err: resultErr}
}

// WithAllProviders is a simple utility function which is used to pass all provider functions to the NewHarvester constructor
func WithAllProviders() []Provider {
	return []Provider{ProxyListNET, SpysME, MultiProxy, ClarkTMProxy, FateProxyList, ProxyListDL}
}
