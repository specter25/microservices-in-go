package handlers

import (
	"log"
	"net/http"
)

// http handler is an interface which implements a fuction so what simply we have to do is tht
type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (h *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Byee"))
}
