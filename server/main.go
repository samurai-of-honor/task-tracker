package main

import (
	"html/template"
	"log"
	"net/http"
)

func page(w http.ResponseWriter, _ *http.Request) {
	tmpl, _ := template.ParseFiles("page.html")
	_ = tmpl.Execute(w, "")
}

func handleRequest() {
	http.HandleFunc("/", page)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	handleRequest()
}
