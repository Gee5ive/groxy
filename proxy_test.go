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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Scheme(); got != tt.want {
				t.Errorf("Proxy.Scheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_Host(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Host(); got != tt.want {
				t.Errorf("Proxy.Host() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_Username(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Username(); got != tt.want {
				t.Errorf("Proxy.Username() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_Password(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Password(); got != tt.want {
				t.Errorf("Proxy.Password() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_ToURL(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want *url.URL
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ToURL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxy.ToURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_Secure(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Secure(); got != tt.want {
				t.Errorf("Proxy.Secure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxy_AsCSV(t *testing.T) {
	tests := []struct {
		name string
		h    *Proxy
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.AsCSV(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxy.AsCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		uri      string
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *Proxy
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.uri, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
