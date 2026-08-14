import { useCallback, useEffect, useMemo, useState } from 'react';
import { Column } from './components/Column';
import type { Task } from './model/Task';
import { TASK_STATUSES, type TaskStatus } from './model/TaskStatus';
import { taskService } from './services/taskService';
import './App.css';

const COLUMN_STYLES: Record<TaskStatus, 'todo' | 'doing' | 'done'> = {
  'A Fazer': 'todo',
  'Em Progresso': 'doing',
  'Concluídas': 'done',
};

export const App = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTasks = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const apiTasks = await taskService.getAll();
      setTasks(apiTasks);
    } catch {
      setError('Erro ao carregar tarefas. Verifique se o backend está em execução na porta 3000.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const tasksByStatus = useMemo<Record<TaskStatus, Task[]>>(
    () => ({
      'A Fazer': tasks.filter((task) => task.status === 'A Fazer'),
      'Em Progresso': tasks.filter((task) => task.status === 'Em Progresso'),
      'Concluídas': tasks.filter((task) => task.status === 'Concluídas'),
    }),
    [tasks],
  );

  const handleEditTask = useCallback((task: Task) => {
    const editedTitle = window.prompt('Novo título da tarefa:', task.title);

    if (!editedTitle || editedTitle.trim().length === 0) {
      return;
    }

    const editedDescriptionInput = window.prompt(
      'Nova descrição da tarefa (opcional):',
      task.description ?? '',
    );

    if (editedDescriptionInput === null) {
      return;
    }

    void (async () => {
      try {
        const updatedTask = await taskService.update(task.id, {
          title: editedTitle.trim(),
          description: editedDescriptionInput.trim() || undefined,
          status: task.status,
        });

        setError(null);
        setTasks((previousTasks) =>
          previousTasks.map((currentTask) =>
            currentTask.id === task.id ? updatedTask : currentTask,
          ),
        );
      } catch {
        setError('Erro ao editar tarefa.');
      }
    })();
  }, []);

  const handleCreateTask = useCallback((status: TaskStatus) => {
    const title = window.prompt('Título da tarefa:');

    if (!title || title.trim().length === 0) {
      return;
    }

    const description = window.prompt('Descrição da tarefa (opcional):') ?? undefined;

    void (async () => {
      try {
        const createdTask = await taskService.create({
          title: title.trim(),
          description: description?.trim() || undefined,
          status,
        });

        setError(null);
        setTasks((previousTasks) => [...previousTasks, createdTask]);
      } catch {
        setError('Erro ao criar tarefa.');
      }
    })();
  }, []);

  const handleDeleteTask = useCallback((id: string) => {
    void (async () => {
      try {
        await taskService.remove(id);
        setError(null);
        setTasks((previousTasks) => previousTasks.filter((task) => task.id !== id));
      } catch {
        setError('Erro ao excluir tarefa.');
      }
    })();
  }, []);

  const handleMoveTask = useCallback((id: string, newStatus: TaskStatus) => {
    void (async () => {
      try {
        const updatedTask = await taskService.update(id, { status: newStatus });
        setError(null);
        setTasks((previousTasks) =>
          previousTasks.map((task) => (task.id === id ? updatedTask : task)),
        );
      } catch {
        setError('Erro ao atualizar status da tarefa.');
      }
    })();
  }, []);

  if (loading) {
    return (
      <div className="board-state">
        <div className="loading-spinner" />
        <p>Carregando Kanban...</p>
      </div>
    );
  }

  if (error && tasks.length === 0) {
    return (
      <div className="board-state error-message">
        <p>{error}</p>
        <button type="button" className="retry-button" onClick={fetchTasks}>
          Tentar novamente
        </button>
      </div>
    );
  }

  return (
    <div className="kanban-shell">
      <header className="board-header">
        <h1>Veritas Kanban</h1>
        <p>Quadro com 3 colunas fixas para o CRUD de tarefas.</p>
      </header>

      {error && (
        <div className="error-banner">
          <span>{error}</span>
          <button type="button" className="error-banner-close" onClick={() => setError(null)}>
            ×
          </button>
        </div>
      )}

      <div className="kanban-board">
        {TASK_STATUSES.map((status) => (
          <Column
            key={status}
            title={status}
            tone={COLUMN_STYLES[status]}
            tasks={tasksByStatus[status]}
            onCreateTask={handleCreateTask}
            onEditTask={handleEditTask}
            onDeleteTask={handleDeleteTask}
            onMoveTask={handleMoveTask}
          />
        ))}
      </div>
    </div>
  );
};
