package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/DanielKnourek/CTC-tasks/trunk/task03/backend/product"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello\n")
}

func CORSheaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		next.ServeHTTP(w, r)
	})
}

// go get -u github.com/gorilla/mux
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", hello).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/product", product.List).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/product", product.Put).Methods(http.MethodPut, http.MethodOptions)

	router_prod := router.PathPrefix("/product/").Subrouter()
	router_prod.HandleFunc("/{id:[[:xdigit:]]{24}}", product.Get).Methods(http.MethodGet, http.MethodOptions)
	router_prod.HandleFunc("/{id:[[:xdigit:]]{24}}", product.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router_prod.HandleFunc("/{id:[[:xdigit:]]{24}}", product.Update).Methods(http.MethodPatch, http.MethodOptions)
	router_prod.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		product.List(w, r)
	})

	router.Use(mux.CORSMethodMiddleware(router), CORSheaderMiddleware)
	http.Handle("/", router)

	BACKEND_PORT := os.Getenv("BACKEND_PORT")
	if BACKEND_PORT == "" {
		BACKEND_PORT = "8080"
	}

	fmt.Printf("Starting backend server on port %s\n", BACKEND_PORT)
	// use DefaultServerMux when handler is nil
	panic(http.ListenAndServe(fmt.Sprintf(":%s", BACKEND_PORT), nil))
}
