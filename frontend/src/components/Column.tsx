import type { Task } from "../model/Task";
import type { TaskStatus } from "../model/TaskStatus";
import { Card } from "./Card";

interface ColumnProps {
  title: TaskStatus;
  tone: "todo" | "doing" | "done";
  tasks: Task[];
  onCreateTask: (status: TaskStatus) => void;
  onEditTask: (task: Task) => void;
  onDeleteTask: (id: string) => void;
  onMoveTask: (id: string, newStatus: TaskStatus) => void;
}

export const Column = ({
  title,
  tone,
  tasks,
  onCreateTask,
  onEditTask,
  onDeleteTask,
  onMoveTask,
}: ColumnProps) => {
  return (
    <section className={`kanban-column kanban-column-${tone}`} aria-label={`Coluna ${title}`}>
      <h2 className="column-title">{title}</h2>
      <button type="button" className="column-create-button" onClick={() => onCreateTask(title)}>
        + Criar tarefa
      </button>
      <div className="task-list">
        {tasks.length === 0 && <p className="empty-column">Sem tarefas nesta etapa.</p>}
        {tasks.map((task) => (
          <Card
            key={task.id}
            task={task}
            onEdit={onEditTask}
            onDelete={onDeleteTask} 
            onStatusChange={onMoveTask}
          />
        ))}
      </div>
    </section>
  );
};