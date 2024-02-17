package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Task item
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks = make(map[string]Task)
var nextID = 1
var lock sync.Mutex

func main() {
	tasks["one"] = Task{ID: 1, Description: "The first task", Completed: false}

	log.Printf("Tasks init: %#v", tasks)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /task/{id}", getTask)
	mux.HandleFunc("POST /task/create", createTask)
	mux.HandleFunc("DELETE /task/del/{id}", delTask)

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	log.Printf("Received {id} %s", idStr)

	lock.Lock()
	task, ok := tasks[idStr]
	lock.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	lock.Lock()
	task.ID = nextID
	tasks[fmt.Sprintf("%d", nextID)] = task
	nextID++
	lock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func delTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	lock.Lock()
	delete(tasks, idStr)
	lock.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
