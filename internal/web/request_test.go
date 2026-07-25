package web

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    Request
		wantErr bool
	}{
		{
			name: "valid GET request",
			line: "GET /links HTTP/1.1",
			want: Request{Method: "GET", Path: "/links", Version: "HTTP/1.1"},
		},
		{
			name: "valid POST request",
			line: "POST / HTTP/1.1",
			want: Request{Method: "POST", Path: "/", Version: "HTTP/1.1"},
		},
		{
			name:    "missing version",
			line:    "GET /links",
			wantErr: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "too many parts",
			line:    "GET /links HTTP/1.1 extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRequestLine(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Method != tt.want.Method || got.Path != tt.want.Path || got.Version != tt.want.Version {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReadHeaders(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "single header",
			input: "Content-Length: 5\r\n\r\n",
			want:  map[string]string{"Content-Length": "5"},
		},
		{
			name:  "multiple headers",
			input: "Host: localhost\r\nContent-Type: text/plain\r\n\r\n",
			want:  map[string]string{"Host": "localhost", "Content-Type": "text/plain"},
		},
		{
			name:  "no headers",
			input: "\r\n",
			want:  map[string]string{},
		},
		{
			name:    "invalid header without colon",
			input:   "InvalidHeader\r\n\r\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got, err := readHeaders(reader)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("header %q: got %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
