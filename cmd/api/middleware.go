package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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

// rateLimit - IP-based rate limiting
func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background goroutine to remove old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				// check if the client hasn't been seen for the past 3 minutes
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only carry out the rate limiting if enabled
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// lock the mutex to prevent concurrent execution
			mu.Lock()

			// check if the IP exisits in the map, if it doesn't, initialize a new rate limiter and add the IP address and limiter to the map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			// update the client's last seen
			clients[ip].lastSeen = time.Now()

			// if the request is not allowed, unlock the mutex and send 429 error
			if !clients[ip].limiter.Allow() {
				// fmt.Println("IP:", ip, "\nLast seen:", clients[ip].lastSeen.String(), "\nTokens:", clients[ip].limiter.Tokens(), "\n...")
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// Very Important: unlock the mutex before calling the next handler in the chain.
			mu.Unlock()
		}

		next.ServeHTTP(w, r)

	})
}

// rateLimit - Global rate limiting
// func (app *application) rateLimit(next http.Handler) http.Handler {
// 	// new limiter which allows an average of 2 requests per second, with a maximum of 4 requests in a single 'burst'
// 	limiter := rate.NewLimiter(2, 4)

// 	// log.Println(`Number of tokens:`, limiter.Tokens())

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// log.Println(`Number of tokens left:`, limiter.Tokens())
// 		if !limiter.Allow() {
// 			app.rateLimitExceededResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
