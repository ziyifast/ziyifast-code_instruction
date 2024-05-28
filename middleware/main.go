package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) responses.
func CORSMiddleware(next http.Handler) http.Handler {
	fmt.Println("cors middleware....")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
			//如果前后端需要传递自定义请求头，需要再Access-Control-Allow-Headers中匹配(Yi-Auth-Token)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Yi-Auth-Token")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept,Yi-Auth-Token")
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware simulates a simple authentication middleware.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware...")
		//store info in ctx
		token := r.Header.Get("Token")
		if len(token) != 0 {
			//TODO 1. check token 2. get userinfo from token
			userID := "1"
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// AuditMiddleware simulates an audit logging middleware.
func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("audit middleware...")
		next.ServeHTTP(w, r)
	})
}

// SmokeHandler returns the current time as a string.
func SmokeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("smoke handle....")
	_, err := w.Write([]byte(time.Now().String()))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewMiddlewareChain creates a new middleware chain with the given middlewares.
func NewMiddlewareChain(middlewares ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

func RunAndServe() error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("err=", e)
		}
	}()

	mux := http.NewServeMux()

	// Create middleware chains for routes.
	authMiddlewareChain := NewMiddlewareChain(CORSMiddleware, AuthMiddleware, AuditMiddleware)
	//noAuthMiddlewareChain := NewMiddlewareChain(CORSMiddleware)

	// Convert the middleware chain result to http.HandlerFunc.
	smokeHandlerWrapped := func(w http.ResponseWriter, r *http.Request) {
		authMiddlewareChain(http.HandlerFunc(SmokeHandler)).ServeHTTP(w, r)
	}

	mux.HandleFunc("/smoke", smokeHandlerWrapped)

	fmt.Printf("listening on http://localhost:%d\n", 9999)
	return http.ListenAndServe(":9999", mux)
}

func main() {
	RunAndServe()
}
