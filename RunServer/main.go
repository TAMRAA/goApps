package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var quotes = []string{
	"Do not take life too seriously. You will never get out of it alive. - Elbert Hubbard",
	"To succeed in life, you need two things: ignorance and confidence. - Mark Twain",
	"The only mystery in life is why the kamikaze pilots wore helmets. - Al McGuire",
	"Life is short, and it is up to you to make it sweet. - Sarah Louise Delany",
	"In three words I can sum up everything I've learned about life: it goes on. - Robert Frost",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		randomIndex := rand.Intn(len(quotes))
		randomQuote := quotes[randomIndex]

		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Random Quote</title>
			</head>
			<body>
				<h1>Random Quote</h1>
				<p>%s</p>
			</body>
			</html>
		`, randomQuote)

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}