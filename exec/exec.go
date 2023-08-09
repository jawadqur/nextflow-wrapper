package exec

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

// Create function to create a temporary nextflow.config file
func createConfig(username string) string {

	// dynamic file name so it doesn't crash with multiple users
	fileName := "tmp/" + username + ".config"

	// Create file
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Print("Creating config file")
	// Write to file
	_, err = f.WriteString("process.executor = 'local'\n")
	if err != nil {
		log.Fatal(err)
	}
	return fileName
}

func ExecHandler(w http.ResponseWriter, r *http.Request) {
	// Path to the binary
	binary := "/nextflow/nextflow"
	flusher, ok := w.(http.Flusher)
	if !ok {
		// streaming not supported
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// read the workflow name from the url params
	vars := mux.Vars(r)
	workflow := vars["workflow"]

	// if workflow is empty fallback to hello
	if workflow == "" {
		workflow = "hello"
	}

	//these two headers are needed to get the http chunk incremently
	//this header had no effect
	w.Header().Set("Connection", "Keep-Alive")

	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	log.Print("Executing command")

	user := "test"
	fileName := createConfig(user)
	// create a list of commands
	commands := []string{"-log", "logs/" + user + "nextflow.log", "run", workflow, "-c", fileName}
	// "-r", "3.12.0"}

	w.Write([]byte(fmt.Sprint(commands) + "\n"))
	flusher.Flush()
	log.Printf("Running command and waiting for it to finish...", commands)
	// Create a command with arguments
	cmd := exec.Command(binary, commands...)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start command
	cmd.Start()

	// log command started
	log.Printf("Command started with PID %d", cmd.Process.Pid)

	// Create scanner for stdout
	scanner := bufio.NewScanner(stdout)

	// stderr scanner
	stderrScanner := bufio.NewScanner(stderr)

	// Copy stderr to response
	go func() {
		for stderrScanner.Scan() {
			w.Write([]byte(stderrScanner.Text() + "\n"))
			fmt.Println(stderrScanner.Text())
			flusher.Flush()
		}
	}()

	// Stream output to response writer
	for scanner.Scan() {
		w.Write([]byte(scanner.Text() + "\n"))
		fmt.Println(scanner.Text())
		flusher.Flush()
	}

	// Check errors
	if err := scanner.Err(); err != nil {
		http.Error(w, "Error streaming response", 500)
		return
	}

	// Wait for command to finish
	cmd.Wait()

}
