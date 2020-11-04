package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// http handler is an interface which implements a fuction so what simply we have to do is tht
type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ooops", http.StatusBadRequest)
		return
		//this is the 3 tasks that http.Error accomplishes
		// w.WriteHeader(http.StatusBadRequest) // similar to res.status
		// w.Write([]byte("Ooops"))
		// return
	}
	log.Printf("Data %s= \n", d)
	fmt.Fprintf(w, "Hello %s", d)
}
