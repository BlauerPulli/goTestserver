package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func availableImages(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		cmd := exec.Command("docker", "images")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	} else {
		err := "Wrong request method!"
		fmt.Printf("Request failed with %s\n", err)
	}
}

func runningContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cmd := exec.Command("docker", "ps")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	} else {
		err := "Wrong request method!"
		fmt.Printf("Request failed with %s\n", err)
	}

}

/*
 func routeTwo (w http.ResponseWriter, r *http.Request) {
 	fmt.Fprintf(w, "In route Zwei -> gewählte Route: %s", r.URL.Path[1:])
}
*/

func main() {
	http.HandleFunc("/testserver/getAvailableImages", availableImages)
	http.HandleFunc("/testserver/getRunningContainers", runningContainers)
	// http.HandleFunc("/routeTwo", routeTwo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
