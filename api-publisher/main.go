package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	log.SetOutput(os.Stdout)
	log.Println("Starting application")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/health")
		fmt.Fprint(w, "{\"status\":\"UP\"}")
	})

	log.Println("App has started on port 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
