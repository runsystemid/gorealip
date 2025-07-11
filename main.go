package main

import (
	"log"
	"net/http"

	"github.com/realclientip/realclientip-go"
)

func realIPMiddleware(next http.Handler) http.Handler {
	strat, err := realclientip.NewRightmostTrustedCountStrategy("X-Forwarded-For", 1)
	if err != nil {
		log.Fatalf("strategy creation error: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := strat.ClientIP(r.Header, r.RemoteAddr)
		if realIP == "" {
			log.Println("WARNING: real IP was empty")
		}
		chain := r.Header.Get("X-Forwarded-For")
		log.Printf("RealIP: %s | ProxyChain: \"%s\" | %s %s %s",
			realIP, chain, r.Method, r.URL.Path, r.Proto)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", realIPMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello\n"))
	})))
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
