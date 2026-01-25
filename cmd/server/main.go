package main

import (
	""
	"net/http"
)

const (
	IP   string = "localhost"
	PORT string = "8080"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello World")) })
	http.ListenAndServe(IP+":"+PORT, mux)
}
