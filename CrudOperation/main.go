package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// Task represents a single to-do item
type Task struct {
	ID   int
	Name string
}

var (
	tasks      = []Task{}
	nextTaskID = 1
	mutex      sync.Mutex
)

var tmpl = template.Must(template.New("tasks").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>To-Do List</title>
</head>
<body>
	<h1>To-Do List</h1>
	<form action="/add" method="POST">
		<input type="text" name="task" placeholder="New Task" required>
		<button type="submit">Add Task</button>
	</form>
	<ul>
		{{range .}}
			<li>{{.Name}} <a href="/delete?id={{.ID}}">Delete</a></li>
		{{end}}
	</ul>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", viewTasks)
	http.HandleFunc("/add", addTask)
	http.HandleFunc("/delete", deleteTask)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func viewTasks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	tmpl.Execute(w, tasks)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	taskName := r.FormValue("task")
	if taskName == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mutex.Lock()
	tasks = append(tasks, Task{ID: nextTaskID, Name: taskName})
	nextTaskID++
	mutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	mutex.Lock()
	defer mutex.Unlock()

	for i, task := range tasks {
		if fmt.Sprintf("%d", task.ID) == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}