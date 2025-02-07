package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type URLMapping struct {
	ShortURL string
	LongURL  string
}

var (
	urlMappings = map[string]string{}
	mutex       sync.Mutex
	chars       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>URL Shortener</title>
</head>
<body>
	<h1>URL Shortener</h1>
	<form action="/shorten" method="POST">
		<input type="url" name="url" placeholder="Enter URL" required>
		<button type="submit">Shorten</button>
	</form>
	{{if .}}
		<p>Shortened URL: <a href="/{{.}}">http://localhost:8080/{{.}}</a></p>
	{{end}}
</body>
</html>
`))

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/redirect/", redirectHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	shortURL := generateShortURL()

	mutex.Lock()
	urlMappings[shortURL] = longURL
	mutex.Unlock()

	tmpl.Execute(w, shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/redirect/"):]

	mutex.Lock()
	longURL, exists := urlMappings[shortURL]
	mutex.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func generateShortURL() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}