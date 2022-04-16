package main

import (
	"fmt"
	"net/http"
	"os"

)

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello\n")
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))

	FRONTEND_PORT := os.Getenv("FRONTEND_PORT")
	if FRONTEND_PORT == "" {
		FRONTEND_PORT = "8081"
	}

	fmt.Printf("Starting frontend test server on port %s\n", FRONTEND_PORT)
	// use DefaultServerMux when handler is nil
	panic(http.ListenAndServe(fmt.Sprintf(":%s",FRONTEND_PORT), nil))
}
