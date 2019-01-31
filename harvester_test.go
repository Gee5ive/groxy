package groxy

import (
	"reflect"
	"testing"
)

func TestNewHarvester(t *testing.T) {
	type args struct {
		providers []Provider
	}
	tests := []struct {
		name string
		args args
		want *Harvester
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHarvester(tt.args.providers...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHarvester() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHarvester_Harvest(t *testing.T) {
	tests := []struct {
		name string
		h    *Harvester
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Harvest()
		})
	}
}

func TestHarvester_Proxies(t *testing.T) {
	tests := []struct {
		name string
		h    *Harvester
		want []*Proxy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Proxies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Harvester.Proxies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBody(t *testing.T) {
	type args struct {
		site string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBody(tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyListDL(t *testing.T) {
	tests := []struct {
		name string
		want ProviderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProxyListDL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProxyListDL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFateProxyList(t *testing.T) {
	tests := []struct {
		name string
		want ProviderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FateProxyList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FateProxyList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClarkTMProxy(t *testing.T) {
	tests := []struct {
		name string
		want ProviderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClarkTMProxy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClarkTMProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiProxy(t *testing.T) {
	tests := []struct {
		name string
		want ProviderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MultiProxy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MultiProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpysME(t *testing.T) {
	tests := []struct {
		name string
		want ProviderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpysME(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpysME() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyListNET(t *testing.T) {
	tests := []struct {
		name string
		want ProviderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProxyListNET(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProxyListNET() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAllProviders(t *testing.T) {
	tests := []struct {
		name string
		want []Provider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAllProviders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAllProviders() = %v, want %v", got, tt.want)
			}
		})
	}
}
