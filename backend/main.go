package main

import (
	"fmt"
	"log"
	"net/http"
)

// enableCORS configura os cabeçalhos de Cross-Origin Resource Sharing
// Permitindo a comunicação segura com qualquer origem do frontend em desenvolvimento
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Responde antecipadamente às requisições preflight (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware exibe no console as requisições recebidas para fácil rastreamento
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s -> %s", r.Method, r.RemoteAddr, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Inicialização da camada de dados em memória e dos handlers
	store := NewTaskStore()
	handler := NewTaskHandler(store)

	// Roteador nativo do Go (Go 1.22+ com suporte a métodos e path parameters)
	mux := http.NewServeMux()

	// Mapeamento das rotas RESTful
	mux.HandleFunc("GET /tasks", handler.GetAllTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	// Aplicação da cadeia de middlewares (CORS + Logging)
	handlerWithMiddlewares := loggingMiddleware(enableCORS(mux))

	port := ":3000"
	fmt.Println("==================================================")
	fmt.Println("🚀 Mini Kanban Backend (Veritas) - Servidor Iniciado")
	fmt.Printf("📍 Escutando em: http://localhost%s\n", port)
	fmt.Println("🌐 CORS habilitado para origens do frontend")
	fmt.Println("📋 Endpoints disponíveis:")
	fmt.Println("   • GET    /tasks       -> Listar todas as tarefas")
	fmt.Println("   • POST   /tasks       -> Criar nova tarefa")
	fmt.Println("   • PUT    /tasks/{id}  -> Atualizar tarefa existente")
	fmt.Println("   • DELETE /tasks/{id}  -> Remover tarefa por ID")
	fmt.Println("==================================================")

	if err := http.ListenAndServe(port, handlerWithMiddlewares); err != nil {
		log.Fatalf("Erro fatal ao iniciar o servidor: %v", err)
	}
}
