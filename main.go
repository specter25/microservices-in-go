package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
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

	})
	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")

	})

	http.ListenAndServe(":9090", nil)

}

//step1 -- create a server
//parameter one is the ip:port and the parameter 2 is the handler
//now as we have specified nil in the http handler what it does is it creates a server and uses a default serve mux as the handler
//step 2 create hhtp handler
//what this function does is it creates a http handler and it attaches t to the defaukt servemux
//step 3 create the goodye handler
// now curl at goodbye and logs goodbye world
//anything other than goodbys matches the first route thus it is a greedy matching
// to read info from a req.body -- body is a ioRead enclosure i.e it implements the io.read interafce
//step5
//send the data as res.send
//it implements io.Writer

//step 6
