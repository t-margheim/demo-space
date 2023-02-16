package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	insecureServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	http.HandleFunc("/", handler)
	go func() {
		err := insecureServer.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	secureServer := &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	err := secureServer.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Println(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
