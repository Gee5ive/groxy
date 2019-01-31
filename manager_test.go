package groxy

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	type args struct {
		maxConn  int
		timeout  time.Duration
		queryURL string
	}
	tests := []struct {
		name    string
		args    args
		want    *Manager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManager(tt.args.maxConn, tt.args.timeout, tt.args.queryURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Distinct(t *testing.T) {
	type args struct {
		proxies []*Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []*Proxy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Distinct(tt.args.proxies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_doRequest(t *testing.T) {
	type args struct {
		proxy *Proxy
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.doRequest(tt.args.proxy)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.doRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_checkProxy(t *testing.T) {
	type args struct {
		proxy *Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want TestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkProxy(tt.args.proxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.checkProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_checkIsSecure(t *testing.T) {
	type args struct {
		proxy *Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkIsSecure(tt.args.proxy); got != tt.want {
				t.Errorf("Manager.checkIsSecure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getToken(t *testing.T) {
	type args struct {
		proxy *Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.getToken(tt.args.proxy)
		})
	}
}

func TestManager_releaseToken(t *testing.T) {
	type args struct {
		proxy *Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.releaseToken(tt.args.proxy)
		})
	}
}

func TestManager_Add(t *testing.T) {
	type args struct {
		proxies []*Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Add(tt.args.proxies...)
		})
	}
}

func TestManager_Run(t *testing.T) {
	tests := []struct {
		name string
		m    *Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Run()
		})
	}
}

func TestManager_Proxies(t *testing.T) {
	tests := []struct {
		name string
		m    *Manager
		want []*Proxy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Proxies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.Proxies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Stop(t *testing.T) {
	tests := []struct {
		name string
		m    *Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Stop()
		})
	}
}

func TestManager_Filter(t *testing.T) {
	type args struct {
		predicate Predicate
		time      []time.Duration
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []*Proxy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Filter(tt.args.predicate, tt.args.time...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_SortByResponseTime(t *testing.T) {
	tests := []struct {
		name string
		m    *Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SortByResponseTime()
		})
	}
}
