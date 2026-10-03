package http

import "net/http"

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Bytes      int
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		StatusCode:     http.StatusOK, // Default status code is 200 OK
	}
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *ResponseRecorder) Write(b []byte) (int, error) {
	n, err := rr.ResponseWriter.Write(b)
	rr.Bytes += n
	return n, err
}

func (rr *ResponseRecorder) Flush() {
	if flusher, ok := rr.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
