package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Println("inventory service listening on :8083")

	if err := http.ListenAndServe(":8083", mux); err != nil {
		log.Fatal(err)
	}
}