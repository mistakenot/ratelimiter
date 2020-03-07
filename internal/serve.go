package internal

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Serve() error {
	http.HandleFunc("/", handler)
	return http.ListenAndServe(":8080", nil)
}
