package web

import (
	"strconv"
	"strings"
)

func WriteResponse(resp Response) []byte {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}

	statusLine := "HTTP/1.1 " + strconv.Itoa(resp.StatusCode) + " " + resp.Status + "\r\n"

	_, ok := resp.Headers["Content-Length"]
	if !ok {
		resp.Headers["Content-Length"] = strconv.Itoa(len(resp.Body))
	}

	var b strings.Builder

	b.WriteString(statusLine)
	for key, value := range resp.Headers {
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(value)
		b.WriteString("\r\n")
	}

	b.WriteString("\r\n")
	b.Write(resp.Body)

	return []byte(b.String())
}
