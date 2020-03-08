package serve

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mistakenot/ratelimiter/pkg/limiter"
)

func Serve(limiter limiter.RateLimiter, port int) error {
	r := mux.NewRouter()

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		err := limiter.Healthcheck()

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.Write([]byte("Ratelimiter Ok."))
	}

	// Index and healthcheck url
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/healthz", indexHandler)

	// Rate limit url
	r.HandleFunc("/token/{tokenId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenId := vars["tokenId"]

		if tokenId == "" {
			w.WriteHeader(404)
			return
		}

		result, err := limiter.TakeToken(tokenId)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		err = json.NewEncoder(w).Encode(result)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
	})

	addr := fmt.Sprintf(":%v", port)

	fmt.Printf("Listening on %v\n", addr)

	// TODO is this the best place to put this?
	// Makes sense to link its lifetime with the
	//  lifetime of the server.
	defer limiter.Close()

	return http.ListenAndServe(addr, r)
}
