package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

// recoverPanic recovers from a panic, and set the connection to close
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// new limiter which allows an average of 2 requests per second, with a maximum of 4 requests in a single 'burst'
	limiter := rate.NewLimiter(2, 4)

	log.Println(`Number of tokens:`, limiter.Tokens())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(`Number of tokens left:`, limiter.Tokens())
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
