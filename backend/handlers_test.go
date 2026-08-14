package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestFullTaskLifecycleAndValidations(t *testing.T) {
	store := NewTaskStore()
	handler := NewTaskHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.GetAllTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)
	server := httptest.NewServer(enableCORS(mux))
	defer server.Close()

	// 1. Test GET /tasks
	resp, err := http.Get(server.URL + "/tasks")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /tasks falhou: %v, status: %d", err, resp.StatusCode)
	}
	var tasks []Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 3 {
		t.Fatalf("esperava 3 tarefas iniciais, obteve %d", len(tasks))
	}

	// 2. Test POST /tasks - Validação de título obrigatório
	body, _ := json.Marshal(map[string]string{"title": "   ", "status": "A Fazer"})
	resp, _ = http.Post(server.URL+"/tasks", "application/json", bytes.NewReader(body))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("esperava 400 ao criar sem título, obteve %d", resp.StatusCode)
	}

	// 3. Test POST /tasks - Validação de status inválido
	body, _ = json.Marshal(map[string]string{"title": "Nova Tarefa", "status": "Invalido"})
	resp, _ = http.Post(server.URL+"/tasks", "application/json", bytes.NewReader(body))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("esperava 400 com status inválido, obteve %d", resp.StatusCode)
	}

	// 4. Test POST /tasks - Sucesso
	body, _ = json.Marshal(map[string]string{
		"title":       "Minha Nova Tarefa",
		"description": "Descricao da tarefa",
		"status":      "A Fazer",
	})
	resp, _ = http.Post(server.URL+"/tasks", "application/json", bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201 ao criar tarefa válida, obteve %d", resp.StatusCode)
	}
	var created Task
	json.NewDecoder(resp.Body).Decode(&created)
	if created.ID == "" || created.Title != "Minha Nova Tarefa" {
		t.Fatalf("tarefa criada incorretamente: %+v", created)
	}

	// 5. Test PUT /tasks/{id} - Mover status (apenas status no payload)
	moveBody, _ := json.Marshal(map[string]string{"status": "Em Progresso"})
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/tasks/"+created.ID, bytes.NewReader(moveBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT para mover status falhou: %v, status: %d", err, resp.StatusCode)
	}
	var moved Task
	json.NewDecoder(resp.Body).Decode(&moved)
	if moved.Status != "Em Progresso" || moved.Title != "Minha Nova Tarefa" || moved.Description != "Descricao da tarefa" {
		t.Fatalf("mover status alterou outros campos indevidamente: %+v", moved)
	}

	// 6. Test PUT /tasks/{id} - Editar apenas título
	editBody, _ := json.Marshal(map[string]string{"title": "Título Atualizado"})
	req, _ = http.NewRequest(http.MethodPut, server.URL+"/tasks/"+created.ID, bytes.NewReader(editBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT para editar título falhou, status: %d", resp.StatusCode)
	}
	var edited Task
	json.NewDecoder(resp.Body).Decode(&edited)
	if edited.Title != "Título Atualizado" || edited.Status != "Em Progresso" || edited.Description != "Descricao da tarefa" {
		t.Fatalf("editar título alterou status ou descrição indevidamente: %+v", edited)
	}

	// 7. Test PUT /tasks/{id} - Status inválido
	invalidPut, _ := json.Marshal(map[string]string{"status": "StatusDesconhecido"})
	req, _ = http.NewRequest(http.MethodPut, server.URL+"/tasks/"+created.ID, bytes.NewReader(invalidPut))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("esperava 400 ao atualizar com status inválido, obteve %d", resp.StatusCode)
	}

	// 8. Test DELETE /tasks/{id} - Sucesso
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/tasks/"+created.ID, nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 ao excluir tarefa, obteve %d", resp.StatusCode)
	}

	// 9. Test DELETE /tasks/{id} - Não encontrado
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/tasks/"+created.ID, nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("esperava 404 ao excluir tarefa inexistente, obteve %d", resp.StatusCode)
	}

	// 10. Test CORS Preflight OPTIONS
	req, _ = http.NewRequest(http.MethodOptions, server.URL+"/tasks", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 no OPTIONS CORS, obteve %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("cabeçalho CORS Allow-Origin incorreto: %s", resp.Header.Get("Access-Control-Allow-Origin"))
	}

	fmt.Println("TODOS OS 10 TESTES DE INTEGRAÇÃO PASSARAM COM SUCESSO!")
}

func TestConcurrentAccessAndRWMutexSafety(t *testing.T) {
	store := NewTaskStore()
	var wg sync.WaitGroup
	numGoroutines := 50

	// Concorrentemente cria tarefas
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			store.Create(Task{
				Title:  fmt.Sprintf("Tarefa Concorrente %d", idx),
				Status: StatusTodo,
			})
		}(i)
	}

	// Concorrentemente lê todas as tarefas
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.GetAll()
		}()
	}

	wg.Wait()

	tasks := store.GetAll()
	// 3 iniciais + 50 concorrentes = 53
	if len(tasks) != 3+numGoroutines {
		t.Fatalf("esperava %d tarefas após concorrência, obteve %d", 3+numGoroutines, len(tasks))
	}
	fmt.Println("TESTE DE CONCORRÊNCIA E RWMUTEX PASSOU COM SUCESSO!")
}
