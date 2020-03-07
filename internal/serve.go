package internal

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func Serve(port int) error {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)

	// Health check url. Need to double check this.
	r.HandleFunc("/zhealth", indexHandler)

	addr := fmt.Sprintf(":%v", port)

	fmt.Printf("Listening on %v\n", addr)

	return http.ListenAndServe(addr, r)
}
