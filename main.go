package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Picture struct {
	Id   int
	Name string
}

var (
	mu           sync.Mutex
	currentMaxId = 3
)

func getNextId() int {
	mu.Lock()
	defer mu.Unlock()
	currentMaxId++
	return currentMaxId
}

func main() {
	fmt.Println("Hello dog")

	handleRoot := func(w http.ResponseWriter, r *http.Request) {
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

	handleCreatePicture := func(w http.ResponseWriter, r *http.Request) {
		log.Print("HTMX request received")
		log.Print(r.Header.Get("HX-Request"))

		// extract the POST data from request
		name := r.PostFormValue("name")

		// get next available id
		id := getNextId()

		// create new picture
		// newPicture := Picture{
		// 	Id: id,
		// 	Name: name,
		// }

		// append new picture to list
		// mu.Lock()
		// defer mu.Unlock()
		// pictures

		htmlStr := fmt.Sprintf("<li class='flex items-center py-4 px-6'><span class='text-gray-700 text-lg font-medium mr-4'>%d</span><span class='text-gray-800 text-lg font medium'>%s</span></li>", id, name)
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/add-picture/", handleCreatePicture)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
