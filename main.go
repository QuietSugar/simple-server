package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var (
	port = flag.Int("p", getAnAvailablePort(), "Port of the server")
	root = flag.String("w", ".", "Root dir of the server")
)

func init() {
	flag.Parse()
}

func main() {
	a, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	d := http.Dir(a)
	fs := http.FileServer(d)
	ls := &logServer{Next: fs, Logger: log.New(os.Stderr, "", log.LstdFlags)}
	log.Printf("Simple server on :%d...\n", *port)
	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *port), ls))
	}()
	time.Sleep(time.Second)
	url := fmt.Sprintf("http://localhost:%d", *port)
	openBrowser(url)
	// Keep your server running or perform other tasks
	select {}
}
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		fmt.Println("Unsupported operating system:", runtime.GOOS)
		return
	}
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to open browser:", err)
	}
}

func portCheck(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Failed to Close Listener :", err)
		}
	}(l)
	return true
}

func getAnAvailablePort() int {
	startPort := 8080
	endPort := 9080
	for port := startPort; port <= endPort; port++ {
		if portCheck(port) {
			return port
		}
	}
	return startPort
}
