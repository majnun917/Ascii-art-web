package main

import (
	"flag"
	"fmt"
	"net/http"

	f "ascii-web/handler"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle main page requests
	http.HandleFunc("/", f.MainPage)

	// Handle ASCII art page requests
	http.HandleFunc("/ascii_art", f.AsciiArtPage)

	// Handle form submission and result download
	http.HandleFunc("/submit", f.SubmitForm)
	http.HandleFunc("/download", f.DownloadArt)

	// Parse the server port from the command line flags
	serverPort := flag.String("port", "8080", "port to serve on")
	flag.Parse()

	// Start the server
	fmt.Printf("http://localhost:%s - Server started on port\n", *serverPort)
	http.ListenAndServe(":"+*serverPort, nil)
}
