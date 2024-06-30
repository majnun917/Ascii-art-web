package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type FormData struct {
	Text     string
	Banner   string
	ASCIIArt string
}

var asciiArt string

// The main function that handles the HTTP methods (GET & POST)
func MainPage(w http.ResponseWriter, r *http.Request) {
	// Clear the global ASCII art variable
	asciiArt = ""

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := FormData{}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// Handler to process form submission
func SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	if banner == "" {
		banner = "standard"
	}

	art, err := PrintAsciiArt(text, banner) // Call the function to generate ASCII art
	if err != nil {
		if err.Error() == "Error" {
			http.Error(w, "Notice: only printable ASCII characters are allowed", http.StatusBadRequest)
			return
		}
	}

	asciiArt = art // update global asciiArt variable with the generated art

	// Generate a unique URL to prevent iframe caching issues
	uniqueURL := fmt.Sprintf("/ascii_art?timestamp=%d", time.Now().UnixNano())

	// Redirect to ASCII art display page
	http.Redirect(w, r, uniqueURL, http.StatusSeeOther)
}

// Handler to serve the ASCII art
func AsciiArtPage(w http.ResponseWriter, r *http.Request) {
	data := FormData{
		ASCIIArt: asciiArt,
	}
	tmpl, err := template.ParseFiles("./templates/ascii_art.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tmpl.Execute(w, data)
}

// Handler to download the ASCII art
func DownloadArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if asciiArt == "" {
		http.Error(w, "Notice: nothing to download, you didn't submit anything", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(asciiArt))
}
