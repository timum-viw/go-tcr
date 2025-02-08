package main

import (
	"fmt"
    "log"
    "timum-viw/supercoop/qr"
    "os"
	"encoding/base64"
	"html/template"
	"net/http"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Image</title>
	<style>
		html, body {
			height: 100%;
			margin: 0;
			padding: 0;
		}
		img {
			width: 100%;
			height: 100%;
			position: absolute;
			object-fit: contain;
		}
	</style>
</head>
<body>
	<img src="data:image/png;base64,{{.Base64Image}}" alt="Image" />
</body>
</html>
`

func handleRequest(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[len("/"):]

	if token == "" {
		http.Error(w, "Token is missing", http.StatusBadRequest)
		return
	}

	png, err := qr.Generate(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not generate qr code (%v)", err), http.StatusBadRequest)
		return
	}
	encodedImage := base64.StdEncoding.EncodeToString(png)

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Could not parse HTML template", http.StatusInternalServerError)
		return
	}

	// Create a data structure to pass to the template
	data := struct {
		Base64Image string
	}{
		Base64Image: encodedImage,
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)

	log.Println("Server listening on port 8080...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
		os.Exit(1)
	}
}