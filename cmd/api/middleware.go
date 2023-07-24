package main

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/tomasen/realip"

	"github.com/blessedmadukoma/greenlight/internal/data"
	"github.com/blessedmadukoma/greenlight/internal/validator"
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
			// ip, _, err := net.SplitHostPort(r.RemoteAddr)

			// if err != nil {
			// 	app.serverErrorResponse(w, r, err)
			// 	return
			// }

			ip := realip.FromRequest(r)

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

// authenticate	- authenticate the user
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization") // response may vary based on the value of the Authorization header in the request.

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// validate the authorization header
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

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

// requireAuthenticatedUser	- check if the user is authenticated
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requireActivatedUser - check if the user is activated after checking if the user is authenticated
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

// requirePermission - check if the user has the required permission
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return app.requireActivatedUser(fn)
}

// enableCORS - enable CORS
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin") // response may vary based on the value of the Origin header in the request.

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)      // allow the origin
					w.Header().Set("Access-Control-Allow-Credentials", "true") // allow credentials

					// check if the request is a preflight request, i.e if the method is OPTIONS and there is an Access-Control-Request-Method header in the request
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, DELETE")    // allow the methods
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type") // allow the headers

						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// metrics - keep track of the number of requests received by the application
func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_µs")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()

		totalRequestsReceived.Add(1)

		// next.ServeHTTP(w, r)
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponsesSent.Add(1)

		// duration := time.Now().Sub(start).Microseconds()

		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)

		// totalProcessingTimeMicroseconds.Add(duration)
	})
}
