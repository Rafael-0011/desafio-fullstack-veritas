# 🚀 Desafio Fullstack - Mini Kanban Veritas

Aplicação Fullstack desenvolvida para o Desafio Técnico da **Veritas**, implementando um quadro Kanban minimalista, de alto desempenho e orquestrado via containers Docker.

---

## 📋 Sumário
- [Instruções de Execução](#-instruções-de-execução)
- [Arquitetura e Decisões Técnicas](#-arquitetura-e-decisões-técnicas)
  - [Backend (Go)](#backend-go)
  - [Frontend (React + TypeScript)](#frontend-react--typescript)
- [Fluxo do Usuário (User Flow)](#-fluxo-do-usuário-user-flow)
- [Endpoints da API](#-endpoints-da-api)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Limitações e Melhorias Futuras](#-limitações-e-melhorias-futuras)

---

## ⚡ Instruções de Execução

### Opção 1: Execução com Docker Compose (Recomendado)

Na raiz do projeto, execute o comando único:

```bash
docker compose up --build
```

Após a inicialização dos containers:
- **Frontend (Web)**: Acesse em [http://localhost:3000](http://localhost:3000)
- **Backend (API)**: Disponível em [http://localhost:8080](http://localhost:8080)

Para parar os containers:
```bash
docker compose down
```

---

### Opção 2: Execução Local (Sem Docker)

#### 1. Backend (Go)
Requisitos: Go 1.22+
```bash
cd backend
go run .
```
> O servidor iniciará na porta `8080` (configurável via variável de ambiente `PORT`).

#### 2. Frontend (React)
Requisitos: Node.js 20+
```bash
cd frontend
npm install
npm run dev
```
> O servidor de desenvolvimento iniciará na porta `5173` ou `3000`.

---

## 🏗️ Arquitetura e Decisões Técnicas

O projeto foi concebido seguindo os princípios de simplicidade, código limpo (*Clean Code*) e separação estrita de responsabilidades, atendendo 100% aos requisitos do MVP sem sobrecarga de dependências.

### Backend (Go)
- **Linguagem & Roteador Nativo**: Desenvolvido em **Go**, aproveitando o roteamento nativo do `net/http` (Go 1.22+) com suporte a métodos HTTP e *path parameters* sem a necessidade de frameworks pesados de terceiros.
- **Concorrência e Thread-Safety**: Conforme exigido pelo edital, o armazenamento em memória utiliza `sync.RWMutex` para proteger o slice de tarefas contra *race conditions*. As leituras utilizam `RLock()` com retorno de cópia defensiva, e as mutações (criação, edição, exclusão) utilizam `Lock()`.
- **Validações Estritas**:
  - `title` obrigatório e validado contra strings vazias/espaços.
  - `status` rigorosamente restrito aos 3 valores válidos: `'A Fazer'`, `'Em Progresso'` ou `'Concluídas'`.
  - Tratamento de atualizações parciais via `UpdateTaskInput` com ponteiros, garantindo integridade dos dados existentes ao mover cartões.
- **CORS Middleware**: Habilitado nativamente com suporte a requisições *preflight* (`OPTIONS`) para comunicação segura com o frontend.
- **Estrutura Modular**: Organizado de forma canônica em `main.go`, `handlers.go` e `models.go`.

### Frontend (React + TypeScript)
- **TypeScript & React 19**: Tipagem estática rigorosa para tarefas, status e DTOs da API.
- **Camada de Serviços (`/services`)**: O arquivo `taskService.ts` centraliza todas as chamadas HTTP (Fetch API), mapeando com segurança DTOs da rede para o modelo de domínio.
- **Otimização de Performance**:
  - `useMemo`: Agrupa as tarefas por coluna de forma reativa e otimizada.
  - `useCallback`: Evita recriação desnecessária de handlers nas renderizações dos cartões e colunas.
- **Isolamento de Eventos (`stopPropagation`)**: O clique no card abre a edição, enquanto o seletor de status e os botões de ação ("Excluir" e "Editar") contam com `stopPropagation` para evitar disparos concorrentes indesejados.
- **Resiliência e UX**:
  - Indicador visual de carregamento (*spinner*).
  - Banner de erro não-obstrutivo com botão de fechamento para operações de CRUD.
  - Botão de retry em caso de falha de conexão inicial.
- **Estética Trello**: Interface limpa, responsiva, com cartões escuros contrastantes, realces cromáticos por coluna e cursores interativos.

---

## 🔄 Fluxo do Usuário (User Flow)

```
[ Acesso à Aplicação ]
          │
          ▼
┌────────────────────────────────────────────────────────┐
│ Visualização do Quadro com 3 Colunas Fixas             │
│ • 'A Fazer'   • 'Em Progresso'   • 'Concluídas'        │
└────────────────────────────────────────────────────────┘
          │
          ├─────────────────────────┬─────────────────────────┐
          ▼                         ▼                         ▼
┌──────────────────┐      ┌──────────────────┐      ┌──────────────────┐
│  Criar Tarefa    │      │  Mover Tarefa    │      │  Editar / Excluir│
│  Botão '+' na    │      │  Alterar status  │      │  • Clique no     │
│  coluna desejada │      │  no dropdown do  │      │    card p/ editar│
│  • Título obrig. │      │  card (troca de  │      │  • Botão excluir │
│  • Descrição opc.│      │  coluna fluida)  │      │    remove o card │
└──────────────────┘      └──────────────────┘      └──────────────────┘
```

1. **Abertura do Quadro**: O usuário visualiza as 3 colunas fixas com a listagem sincronizada do backend.
2. **Criação de Tarefa**: O usuário clica em `+ Criar tarefa` no topo de qualquer coluna, inserindo o título (obrigatório) e a descrição (opcional). A tarefa é criada imediatamente na etapa selecionada.
3. **Movimentação entre Etapas**: O usuário altera o status diretamente no seletor do card (`A Fazer` $\rightarrow$ `Em Progresso` $\rightarrow$ `Concluídas`), movendo o cartão dinamicamente entre as colunas sem abrir a tela de edição.
4. **Edição**: O usuário clica no corpo do card para editar o título e a descrição.
5. **Exclusão**: O usuário clica em `Excluir` para remover o card com segurança.

> O diagrama visual do fluxo está disponível em `/docs/user-flow.png`.

---

## 📡 Endpoints da API

| Método | Endpoint | Descrição | Códigos de Retorno |
| :--- | :--- | :--- | :--- |
| `GET` | `/tasks` | Retorna todas as tarefas cadastradas | `200 OK` |
| `POST` | `/tasks` | Cria uma nova tarefa | `201 Created`, `400 Bad Request` |
| `PUT` | `/tasks/{id}` | Atualiza título, descrição ou status da tarefa | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/tasks/{id}` | Remove uma tarefa existente | `200 OK`, `404 Not Found` |

---

## 📁 Estrutura do Projeto

```text
desafio-fullstack-veritas/
├── backend/
│   ├── Dockerfile           # Imagem Go Alpine minimalista
│   ├── go.mod               # Módulo Go
│   ├── go.sum               # Checksums das dependências
│   ├── handlers.go          # Handlers HTTP RESTful
│   ├── handlers_test.go     # Testes unitários, de concorrência e integração
│   ├── main.go              # Inicialização do servidor, rotas e CORS
│   └── models.go            # Modelos, validações e TaskStore com sync.RWMutex
├── frontend/
│   ├── Dockerfile           # Multi-stage build (Node + Nginx)
│   ├── nginx.conf           # Configuração do Nginx na porta 3000 com suporte a SPA
│   ├── package.json         # Dependências do frontend
│   ├── src/
│   │   ├── components/      # Componentes visuais (Card.tsx, Column.tsx)
│   │   ├── model/           # Tipagens de domínio (Task.ts, TaskStatus.ts)
│   │   ├── services/        # Integração HTTP com a API (taskService.ts)
│   │   ├── App.tsx          # Estado global do Kanban e orquestração de eventos
│   │   ├── App.css          # Estilização estilo Trello
│   │   └── main.tsx         # Ponto de entrada do React
│   └── vite.config.ts       # Configuração do Vite
├── docs/
│   └── README.md            # Documentações e diagramas (user-flow.png)
├── docker-compose.yml       # Orquestração do ambiente completo (Backend + Frontend)
└── README.md                # Documentação técnica do projeto
```

---

## 🚀 Limitações e Melhorias Futuras

### Limitações do MVP
- **Persistência Volátil**: As tarefas são mantidas em memória (*slices* em Go), sendo reiniciadas caso o container do backend seja finalizado.

### Evoluções Futuras (Itens Bônus)
1. **Banco de Dados Persistente**: Integração com PostgreSQL ou SQLite via migrations (Golang-migrate) e ORM/SQLc.
2. **Interface Drag and Drop (DnD)**: Implementação de arrastar e soltar cartões entre colunas utilizando `@hello-pangea/dnd` ou HTML5 Drag and Drop API.
3. **Autenticação & Multi-tenancy**: Suporte a múltiplos quadros por usuário com autenticação JWT.
4. **Testes End-to-End (E2E)**: Implementação de testes automatizados com Playwright ou Cypress cobrindo o fluxo completo no navegador.
