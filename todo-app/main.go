package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

var (
	tasks []string
	mutex sync.Mutex
)

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	mutex.Lock()
	defer mutex.Unlock()
	tmpl.Execute(w, tasks)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		task := r.FormValue("task")
		mutex.Lock()
		tasks = append(tasks, task)
		mutex.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		index := r.FormValue("index")
		var i int
		fmt.Sscanf(index, "%d", &i)

		mutex.Lock()
		if i >= 0 && i < len(tasks) {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
		mutex.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
