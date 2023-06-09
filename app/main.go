package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("index.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(content)
	})

	http.HandleFunc("/api/number", func(w http.ResponseWriter, r *http.Request) {
		number := rand.Intn(100) + 1
		w.Write([]byte(strconv.Itoa(number)))
	})

	listener, err := net.Listen("tcp", ":8080")
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
