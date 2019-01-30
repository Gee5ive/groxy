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

type ProviderResponse struct {
	Proxies []*Proxy
	Err     error
}
type Provider func() ProviderResponse

type Harvester struct {
	providers []Provider
	proxies   []*Proxy
}

func NewHarvester(providers ...Provider) *Harvester {
	return &Harvester{providers: providers}
}

func (h *Harvester) Harvest() {
	for _, provider := range h.providers {
		resp := provider()
		if resp.Err == nil {
			h.proxies = append(h.proxies, resp.Proxies...)
		}
	}
}

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

func FreeProxyListDL() ProviderResponse {
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
			proxy, err := New(ip, "", "")
			if err != nil {
				continue
			}
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
		resultErr = multierror.Append(resultErr, result.Err)
	}
	return ProviderResponse{Proxies: proxies, Err: resultErr}
}

func FateProxyList() ProviderResponse {
	parseJson := func(jsonStr string) []string {
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
		proxies := parseJson(resp)
		for _, proxy := range proxies {
			host, err := New(proxy, "", "")
			if err != nil {
				continue
			}
			list = append(list, host)
		}
		respStream <- ProviderResponse{Proxies: list, Err: nil}
	}
	go func() {
		get()
	}()
	return <-respStream
}

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
			proxy, err := New(value, "", "")
			if err != nil {
				continue
			}
			proxies = append(proxies, proxy)

		}
		respStream <- ProviderResponse{Proxies: proxies, Err: nil}
	}
	go func() {
		get()
	}()
	return <-respStream
}
