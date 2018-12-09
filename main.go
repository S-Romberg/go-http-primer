package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
)

type Home struct {
	Path  string
	Name  string
	Color string
}

func main() {
	port := ":3000"

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", root)
	http.HandleFunc("/makepost", makeyPosty)
	http.HandleFunc("/api/json", jsonResponse)
	fmt.Printf("Starting server on port %s\n", port)
	panic(http.ListenAndServe(port, nil))

}

// root is the handler function, it tells me what to do when a user visits root
func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html; charset=utf-8")
	color := r.URL.Query().Get("color")
	if color == "" {
		color = "black"
	}
	data := Home{
		Path:  html.EscapeString(r.URL.Path),
		Name:  "Spencer",
		Color: color,
	}
	t, _ := template.ParseFiles("templates/home.html")
	t.Execute(w, data)
	return
}

func makeyPosty(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "templates/makeyposty.html")
	case "POST":
		r.ParseForm()
		fmt.Println(r.Form["name"][0])
		fmt.Println(r.Form["email"][0])
		http.ServeFile(w, r, "templates/makeyposty.html")
	default:
		http.NotFound(w, r)
	}
}

func jsonResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	returnObject := struct {
		Name            string
		Age             int
		FavoriteSubmits map[string]int
	}{
		Name: "Savannah Grace Null",
		Age:  24,
		FavoriteSubmits: map[string]int{
			"twitter":                      10,
			"slack":                        45,
			"Denver Post comments section": 120,
		},
	}

	body, err := json.Marshal(returnObject)
	if err != nil {
		log.Panicf("api: could not locate json response: %s", err)
	}
	w.Write(body)
	return
}
