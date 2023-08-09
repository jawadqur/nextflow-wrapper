package main

import (
	"gen3-nextflow/exec"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// Handle route to execute command
	router.HandleFunc("/exec", exec.ExecHandler)
	router.HandleFunc("/exec/{workflow}", exec.ExecHandler)

	// port variable either from env var or default
	// port := os.Getenv("PORT", "8001")

	// Start server
	log.Print("Server started on localhost:8000\n")

	log.Fatal(http.ListenAndServe("0.0.0.0:8001", router))
}
