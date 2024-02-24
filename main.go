package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Picture struct {
	Id   int
	Name string
}

func main() {
	fmt.Println("Hello dog")

	handleroot := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("view/index.html"))
		pictures := map[string][]Picture{
			"Pictures": {
				{Id: 1, Name: "flat tire"},
				{Id: 2, Name: "truck tire"},
				{Id: 3, Name: "motor tire"},
			},
		}
		tmpl.Execute(w, pictures)
	}
	http.HandleFunc("/", handleroot)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
