package groxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

type ProviderResponse struct {
	Proxies []*Proxy
	Err     error
}

type Provider interface {
	Provide() ProviderResponse
}
type Harvester struct {
	sources []Provider
}

type FreeProxyListDownload struct {
}

func (p *FreeProxyListDownload) Provide() ProviderResponse {
	var resultErr error
	var proxies []*Proxy
	kinds := []string{"http", "https"}
	respStream := make(chan ProviderResponse, len(kinds))
	get := func(kind string) {
		var list []*Proxy
		client := http.Client{}
		site := fmt.Sprintf("https://www.proxy-list.download/api/v1/get?type=%s", kind)
		resp, err := client.Get(site)
		if err != nil {
			respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
			return

		}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				respStream <- ProviderResponse{Proxies: []*Proxy{}, Err: err}
				return
			}
			bodyString := string(bodyBytes)
			for _, ip := range strings.Fields(bodyString) {
				proxy, err := New(ip, "", "")
				if err != nil {
					continue
				}
				list = append(list, proxy)
			}
			respStream <- ProviderResponse{Proxies: list, Err: nil}
		}
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
