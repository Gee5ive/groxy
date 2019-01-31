package groxy

import (
	"net/url"
	"reflect"
	"testing"
)

func TestProxy_Scheme(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want string
	}{
		{name: "Proxy that is not secure should have http:// scheme", h: &Proxy{host: "159.89.46.56:3128"}, want: "http://"},
		{name: "Proxy that is not secure should have https:// scheme", h: &Proxy{host: "159.89.46.56:3128", secure: true}, want: "https://"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Scheme(); got != tt.want {
				t.Errorf("Proxy.Scheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_ToURL(t *testing.T) {
	expected, _ := url.Parse("http://username:password@159.89.46.56:3128")
	tests := []struct {
		name string
		h    *Proxy
		want *url.URL
	}{
		{name: "Should Return properly formatted url", h: &Proxy{host: "159.89.46.56:3128", password: "password", username: "username"},
			want: expected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ToURL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxy.ToURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeProxies() []*Proxy {
	var proxies []*Proxy
	hosts := []string{"5.135.164.72:3128", "178.128.21.47:3128", "46.105.190.37:3128", "157.230.44.89:3128", "138.201.223.250:31288"}
	for _, host := range hosts {
		proxy, _ := New(host, "", "")
		proxies = append(proxies, proxy)
	}
	return proxies
}
func TestSaveToFile(t *testing.T) {
	type args struct {
		file    string
		proxies []*Proxy
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "SaveToFile Should save the given list of proxies to a csv file", args: args{file: "test.csv", proxies: makeProxies()}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveToFile(tt.args.file, tt.args.proxies); (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Proxy
		wantErr bool
	}{
		{name: "FromFile Should return the given list of proxies from said file", args: args{file: "test.csv"}, want: makeProxies(), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
