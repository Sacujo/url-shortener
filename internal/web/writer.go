package web

import (
	"strconv"
	"strings"
)

var statusTexts = map[int]string{
	200: "OK",
	201: "Created",
	302: "Found",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

func WriteResponse(resp Response) []byte {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}

	status := statusTexts[resp.StatusCode]

	statusLine := "HTTP/1.1 " +
		strconv.Itoa(resp.StatusCode) +
		" " +
		status +
		"\r\n"

	if _, ok := resp.Headers["Content-Length"]; !ok {
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
