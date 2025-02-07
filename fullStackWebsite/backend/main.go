package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	tasks  = []Task{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/add", addTaskHandler)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil);

	err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}


func tasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /tasks")
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /add")
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	task.ID = nextID
	nextID++
	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	log.Printf("Task added: %+v", task)
}