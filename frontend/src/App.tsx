import { useCallback, useEffect, useMemo, useState } from 'react';
import { Column } from './components/Column';
import { initialTasks } from './model/initialTasks';
import type { Task } from './model/Task';
import { TASK_STATUSES, type TaskStatus } from './model/TaskStatus';
import './App.css';

const API_DELAY_IN_MS = 1200;

const COLUMN_STYLES: Record<TaskStatus, 'todo' | 'doing' | 'done'> = {
  'A Fazer': 'todo',
  'Em Progresso': 'doing',
  'Concluídas': 'done',
};

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));


export const App = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadTasks = async () => {
      try {
        setLoading(true);
        setError(null);

        // Simula o tempo de resposta do "GET /tasks" do backend em Go.
        await sleep(API_DELAY_IN_MS);

        if (!isMounted) {
          return;
        }

        setTasks(initialTasks);
      } catch {
        if (isMounted) {
          setError('Erro ao carregar tarefas.');
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    loadTasks();

    return () => {
      isMounted = false;
    };
  }, []);

  const tasksByStatus = useMemo<Record<TaskStatus, Task[]>>(
    () => ({
      'A Fazer': tasks.filter((task) => task.status === 'A Fazer'),
      'Em Progresso': tasks.filter((task) => task.status === 'Em Progresso'),
      'Concluídas': tasks.filter((task) => task.status === 'Concluídas'),
    }),
    [tasks],
  );

  const handleEditTask = useCallback((task: Task) => {
    console.log('Editar tarefa:', task);
  }, []);

  const handleDeleteTask = useCallback((id: string) => {
    console.log('Excluir tarefa com ID:', id);
  }, []);

  const handleMoveTask = useCallback((id: string, newStatus: TaskStatus) => {
    setTasks((previousTasks) =>
      previousTasks.map((task) => (task.id === id ? { ...task, status: newStatus } : task)),
    );
  }, []);

  if (loading) {
    return <div className="board-state">Carregando Kanban...</div>;
  }

  if (error) {
    return <div className="board-state error-message">{error}</div>;
  }

  return (
    <div className="kanban-shell">
      <header className="board-header">
        <h1>Veritas Kanban</h1>
        <p>Quadro com 3 colunas fixas para o CRUD de tarefas.</p>
      </header>

      <div className="kanban-board">
        {TASK_STATUSES.map((status) => (
          <Column
            key={status}
            title={status}
            tone={COLUMN_STYLES[status]}
            tasks={tasksByStatus[status]}
            onEditTask={handleEditTask}
            onDeleteTask={handleDeleteTask}
            onMoveTask={handleMoveTask}
          />
        ))}
      </div>
    </div>
  );
};
