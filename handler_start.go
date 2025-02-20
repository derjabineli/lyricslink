package main

import (
	"fmt"
	"net/http"
)

func handlerStart(w http.ResponseWriter, r *http.Request) {
	response := `
		<html>
		<body>
			<h1>Hello!</h1>
		</body>
		</html>
	`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusAccepted)
    fmt.Fprint(w, response)
}