package main

import (
	"encoding/json"
	"net/http"
)

// TaskHandler gerencia as requisições HTTP da API de tarefas
type TaskHandler struct {
	store *TaskStore
}

// NewTaskHandler cria uma nova instância de TaskHandler injetando o TaskStore
func NewTaskHandler(store *TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

// writeJSON padroniza o retorno das respostas em JSON com status code adequado
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// writeError padroniza as respostas de erro em formato JSON
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// GetAllTasks responde a GET /tasks
// Retorna a lista completa de tarefas em memória
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.GetAll()
	// Garante retorno de array vazio [] em vez de null caso não existam tarefas
	if tasks == nil {
		tasks = []Task{}
	}
	writeJSON(w, http.StatusOK, tasks)
}

// CreateTask responde a POST /tasks
// Cria uma nova tarefa com validação dos campos obrigatórios e geração automática de ID e CreatedAt
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	if err := task.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created := h.store.Create(task)
	writeJSON(w, http.StatusCreated, created)
}

// UpdateTask responde a PUT /tasks/{id}
// Atualiza os dados de uma tarefa existente (permitindo a transição fluida entre colunas no Kanban)
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID da tarefa é obrigatório")
		return
	}

	var input UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	updated, found, err := h.store.Update(id, input)
	if !found {
		writeError(w, http.StatusNotFound, "Tarefa não encontrada")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// DeleteTask responde a DELETE /tasks/{id}
// Remove uma tarefa existente pelo ID
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "ID da tarefa é obrigatório")
		return
	}

	if !h.store.Delete(id) {
		writeError(w, http.StatusNotFound, "Tarefa não encontrada")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Tarefa excluída com sucesso",
		"id":      id,
	})
}
