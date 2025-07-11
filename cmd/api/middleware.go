package main

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AmiyoKm/green_light/internal/store"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

// Set a custom header and a error response to the client when the server
// recovers from an error
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
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// creates a background goroutine that runs alongside the main go routine and
	// deletes ip addresses from the client hashmap that has not made a request in
	// 3 Minutes
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()

		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				// if client ip not found in the clients HashMap then make a new Limiter instance
				// ip as the key
				// limiter as the value
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			// update the lastSeen with every new request made
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})

}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// response may vary depending on Authorization Header
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r := app.contextSetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		// v := validator.New()
		// if store.ValidateTokenPlaintext(v, token); !v.Valid() {
		// 	app.invalidAuthenticationTokenResponse(w, r)
		// 	return
		// }
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.logger.PrintInfo(jwtToken.Raw, nil)

		claims := jwtToken.Claims

		sub, err := claims.GetSubject()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		userID, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		//user, err := app.store.Users.GetForToken(r.Context(), store.ScopeAuthentication, token)
		user, err := app.store.Users.Get(r.Context(), userID)
		if err != nil {
			switch err {
			case store.ErrorNotFound:
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

func (app *application) requiredAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requiredActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)

	}
	return app.requiredAuthenticatedUser(fn)
}

// Middleware that checks if the authenticated and activated user has the required permission.
// This wraps requiredActivatedUser, which itself wraps requiredAuthenticatedUser.
// The middleware chain is: requirePermission → requiredActivatedUser → requiredAuthenticatedUser.
// Each middleware is executed in order when next.ServeHTTP() is called.
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		permissions, err := app.store.Permissions.GetAllForUser(r.Context(), user.ID)
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
	return app.requiredActivatedUser(fn)
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// add Vary Header because responses will be different due to these Headers
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {

					w.Header().Set("Access-Control-Allow-Origin", origin)

					// checks if it is a preflight request by checking method OPTIONS and Header
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

						// Adds Access-Control-Allow Headers for response to the preflight request
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Access-Control-Allow-Methods and Access-Control-Allow-Headers cached for 15 seconds.
						// means no need for again sending the OPTIONS preflight request for 15s.
						w.Header().Set("Access-Control-Max-Age", "15")

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

type metricsResponseWriter struct {
	wrapper       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapper.Header()
}

// custom WriteHeader for our ResponseWriter for satisfying the interface and collect data
func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapper.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

// custom Write for our ResponseWriter for satisfying the interface and collect data
func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	if !mw.headerWritten {
		mw.statusCode = http.StatusOK
		mw.headerWritten = true
	}

	return mw.wrapper.Write(b)
}

func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapper
}

func (app *application) metrics(next http.Handler) http.Handler {
	var (
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		total_responsesSent             = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_microseconds")
		tokenResponseSentByStatus       = expvar.NewMap("total_response_sent_by_status")
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		totalRequestsReceived.Add(1)

		mw := &metricsResponseWriter{wrapper: w}

		// The next HandlerFunc will use our custom metricsResponseWriter (mw) to write the JSON response and set headers.
		// This allows us to track metrics (like status code and bytes written) in our middleware.
		next.ServeHTTP(mw, r)

		total_responsesSent.Add(1)

		tokenResponseSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)

		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}
