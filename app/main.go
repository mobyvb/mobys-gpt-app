package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("index.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(content)
	})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	log.Printf("Starting server on :%d", addr.Port)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal(err)
	}
}
