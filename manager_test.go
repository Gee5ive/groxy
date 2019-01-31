package groxy

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

var manager, _ = NewManager(4, time.Second*1, "http://google.com")

func TestManager_Distinct(t *testing.T) {
	makeDuplicates := func() []*Proxy {
		var proxies []*Proxy
		hosts := []string{"5.135.164.72:3128", "178.128.21.47:3128", "46.105.190.37:3128", "157.230.44.89:3128", "138.201.223.250:31288",
			"5.135.164.72:3128", "178.128.21.47:3128", "178.128.21.47:3128", "157.230.44.89:3128"}
		for _, host := range hosts {
			proxy, _ := New(host, "", "")
			proxies = append(proxies, proxy)
		}
		return proxies
	}
	type args struct {
		proxies []*Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want []*Proxy
	}{
		{name: "Should not return a list with any duplicate proxies", args: args{proxies: makeDuplicates()}, m: manager, want: makeProxies()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Distinct(tt.args.proxies); !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("Manager.Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_checkProxy(t *testing.T) {
	badProxy, _ := New("118.175.93.68:38176", "", "")
	myErr := errors.New("error")
	type args struct {
		proxy *Proxy
	}
	tests := []struct {
		name string
		m    *Manager
		args args
		want TestResult
	}{
		{name: "Dead proxies should not be marked as alive", m: manager, args: args{proxy: badProxy}, want: TestResult{Err: myErr, Proxy: badProxy}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkProxy(tt.args.proxy); !reflect.DeepEqual(got.Proxy.Alive(), tt.want.Proxy.Alive()) {
				t.Errorf("Manager.checkProxy() = %v, want %v", got, tt.want)
			}
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
