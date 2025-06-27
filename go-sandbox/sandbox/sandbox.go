package main

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"
)

// server
type Server struct {
	Router *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		Router: http.NewServeMux(),
	}
	//route
	s.Router.HandleFunc("/run", runHandler)
	return s
}

// runhandler
func runHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST Method Allowed", http.StatusMethodNotAllowed)
		return
	}
	output, err := runCodeInDocker()
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %s\nOutput: %s", err, output), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

// push to  docker
func runCodeInDocker() ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second) //adjust time if timeout error
	defer cancel()

	snippetPath := "/PATH/TO/FILE/snippet.go" // replace with path
	absPath := filepath.ToSlash(snippetPath)

	//docker command
	cmd := exec.CommandContext(ctx, "docker", "run",
		"--rm",
		"--memory=512m",
		"--cpus=0.5",
		"--pids-limit=100",
		"--read-only",
		"--tmpfs", "/tmp:exec",
		"--cap-drop=ALL",
		"--network=none",
		"-e", "GOCACHE=/tmp/gocache",
		"-v", absPath+":/sandbox/snippet.go:ro",
		"golang:latest",
		"go", "run", "/sandbox/snippet.go",
	)

	//output
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {

		return nil, fmt.Errorf("execution timed out")
	} else if err != nil {
		return output, fmt.Errorf("docker Go run failed")
	}

	fmt.Printf("Docker Go run succeeded:\nOutput: %s\n", output)
	return output, nil

}

func main() {
	server := NewServer()
	err := http.ListenAndServe(":8080", server.Router)
	if err != nil {
		fmt.Println("Error COnnection to server")
	}

}
