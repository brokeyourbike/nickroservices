package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hi!")
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Oops", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Hello %s", d)
	})
	http.HandleFunc("/bye", func(http.ResponseWriter, *http.Request) {
		log.Println("Bye!")
	})

	http.ListenAndServe("127.0.0.1:9090", nil)
}
