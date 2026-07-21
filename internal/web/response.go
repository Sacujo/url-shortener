package web

type Response struct {
	StatusCode int
	Status     string

	Headers map[string]string

	Body []byte
}
