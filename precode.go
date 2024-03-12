package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getAllTasks(w http.ResponseWriter, r *http.Request) { //для получения всех задач
	jsonResp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //статус 200 OK
	w.Write(jsonResp)
}

func getTaskID(w http.ResponseWriter, r *http.Request) { // для получения по id
	taskId := chi.URLParam(r, "id")
	task, ok := tasks[taskId]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest) // если не найдена 400
		return
	}
	jsonResp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //статус 200 OK
	w.Write(jsonResp)
}
func createTasks(w http.ResponseWriter, r *http.Request) { // создание задач
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	// проверка ID
	if _, exists := tasks[task.ID]; exists {
		http.Error(w, "Этот ID занят", http.StatusBadRequest) // 400
		return
	}

	tasks[task.ID] = task             // cоздаем задачу
	w.WriteHeader(http.StatusCreated) // статус 201
}
func deleteTaskID(w http.ResponseWriter, r *http.Request) { // удаление задачи
	taskId := chi.URLParam(r, "id") // значение id
	_, ok := tasks[taskId]
	if !ok {
		http.Error(w, "Задача была удалена", http.StatusBadRequest) // 400
		return
	}
	delete(tasks, taskId) // удалить задачу
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasks)          // обрабочик для получения всех задач
	r.Post("/tasks", createTasks)         // создание задач
	r.Get("/tasks/{id}", getTaskID)       // обработчик доя задачи по ID
	r.Delete("/tasks/{id}", deleteTaskID) // обработчик для удаления задачи

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
