package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Constantes com os três status estritamente permitidos pelo edital
const (
	StatusTodo       = "A Fazer"
	StatusInProgress = "Em Progresso"
	StatusDone       = "Concluídas"
)

// Task representa o modelo de dados de uma tarefa no Kanban
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TaskStore gerencia o armazenamento em memória de forma thread-safe usando sync.RWMutex
type TaskStore struct {
	mu    sync.RWMutex
	tasks []Task
}

// NewTaskStore inicializa o repositório em memória com tarefas de exemplo
func NewTaskStore() *TaskStore {
	now := time.Now()
	return &TaskStore{
		tasks: []Task{
			{
				ID:          uuid.NewString(),
				Title:       "Configurar ambiente",
				Description: "Instalar dependências e preparar o workspace do projeto",
				Status:      StatusDone,
				CreatedAt:   now.Add(-2 * time.Hour),
			},
			{
				ID:          uuid.NewString(),
				Title:       "Desenvolver backend em Go",
				Description: "Implementar API RESTful com rotas, validações e CORS",
				Status:      StatusInProgress,
				CreatedAt:   now.Add(-1 * time.Hour),
			},
			{
				ID:          uuid.NewString(),
				Title:       "Integrar com Frontend React",
				Description: "Conectar as colunas do Kanban com os endpoints da API",
				Status:      StatusTodo,
				CreatedAt:   now,
			},
		},
	}
}

// IsValidStatus valida se o status fornecido é um dos três valores permitidos
func IsValidStatus(status string) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

// Validate executa as validações de integridade da Task
func (t *Task) Validate() error {
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return errors.New("o campo 'title' é obrigatório")
	}

	// Se o status não for informado na criação, assume o padrão inicial 'A Fazer'
	if t.Status == "" {
		t.Status = StatusTodo
	}

	if !IsValidStatus(t.Status) {
		return fmt.Errorf("o campo 'status' deve ser '%s', '%s' ou '%s'", StatusTodo, StatusInProgress, StatusDone)
	}

	return nil
}

// GetAll retorna uma cópia segura de todas as tarefas cadastradas (Read Lock)
func (s *TaskStore) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

// GetByID busca uma tarefa específica pelo seu ID (Read Lock)
func (s *TaskStore) GetByID(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		if task.ID == id {
			return task, true
		}
	}
	return Task{}, false
}

// Create insere uma nova tarefa gerando ID com google/uuid e CreatedAt automaticamente (Write Lock)
func (s *TaskStore) Create(task Task) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = uuid.NewString()
	task.CreatedAt = time.Now().UTC()
	s.tasks = append(s.tasks, task)
	return task
}

// UpdateTaskInput define os campos aceitos na atualização de uma tarefa
type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

// Update atualiza os campos de uma tarefa existente preservando ID, CreatedAt e campos não alterados (Write Lock)
func (s *TaskStore) Update(id string, input UpdateTaskInput) (Task, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			if input.Title != nil {
				trimmedTitle := strings.TrimSpace(*input.Title)
				if trimmedTitle == "" {
					return Task{}, true, errors.New("o campo 'title' não pode ser vazio")
				}
				task.Title = trimmedTitle
			}

			if input.Description != nil {
				task.Description = strings.TrimSpace(*input.Description)
			}

			if input.Status != nil {
				if !IsValidStatus(*input.Status) {
					return Task{}, true, fmt.Errorf("o campo 'status' deve ser '%s', '%s' ou '%s'", StatusTodo, StatusInProgress, StatusDone)
				}
				task.Status = *input.Status
			}

			s.tasks[i] = task
			return task, true, nil
		}
	}
	return Task{}, false, nil
}

// Delete remove uma tarefa pelo ID (Write Lock)
func (s *TaskStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}
	return false
}
