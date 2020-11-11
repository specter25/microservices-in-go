package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

//GzipHandler something
type GzipHandler struct {
}

//GzipMiddleware somwthing
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			//create a 	gzipped response
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(wrw, r)
			defer wrw.Flush()
			return
		}

		next.ServeHTTP(rw, r)

	})
}

//WrappedresponseWriter something
type WrappedresponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

//NewWrappedResponseWriter something
func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedresponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrappedresponseWriter{rw: rw, gw: gw}
}

//Header something
func (wr *WrappedresponseWriter) Header() http.Header {
	return wr.rw.Header()
}

//Write someting
func (wr *WrappedresponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)

}

//WriteHEadersomething
func (wr *WrappedresponseWriter) WriteHeader(statuscode int) {
	wr.rw.WriteHeader(statuscode)
}

//Flush something
func (wr *WrappedresponseWriter) Flush() {
	wr.gw.Flush()
}
