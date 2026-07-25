package web

import (
	"strconv"
	"strings"
	"testing"
)

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		name            string
		resp            Response
		wantLine        string
		wantBody        string
		wantHeader      string
		wantHeaderValue string
	}{
		{
			name:     "200 OK with body",
			resp:     Response{StatusCode: 200, Status: "OK", Body: []byte("hello")},
			wantLine: "HTTP/1.1 200 OK\r\n",
			wantBody: "hello",
		},
		{
			name:     "404 Not Found empty body",
			resp:     Response{StatusCode: 404, Status: "Not Found"},
			wantLine: "HTTP/1.1 404 Not Found\r\n",
			wantBody: "",
		},
		{
			name:            "302 redirect keeps custom header",
			resp:            Response{StatusCode: 302, Status: "Found", Headers: map[string]string{"Location": "https://example.com"}},
			wantLine:        "HTTP/1.1 302 Found\r\n",
			wantBody:        "",
			wantHeader:      "Location",
			wantHeaderValue: "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := string(WriteResponse(tt.resp))

			if !strings.HasPrefix(data, tt.wantLine) {
				t.Errorf("status line: got %q, want prefix %q", data, tt.wantLine)
			}
			if !strings.HasSuffix(data, tt.wantBody) {
				t.Errorf("body: got %q, want suffix %q", data, tt.wantBody)
			}

			wantCL := "Content-Length: " + strconv.Itoa(len(tt.resp.Body))
			if !strings.Contains(data, wantCL) {
				t.Errorf("expected %q in response, got %q", wantCL, data)
			}

			if tt.wantHeader != "" {
				wantHeaderLine := tt.wantHeader + ": " + tt.wantHeaderValue
				if !strings.Contains(data, wantHeaderLine) {
					t.Errorf("expected header %q in response, got %q", wantHeaderLine, data)
				}
			}
		})
	}
}
