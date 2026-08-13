import type { ChangeEvent } from "react";
import type { Task } from "../model/Task";
import { TASK_STATUSES, type TaskStatus } from "../model/TaskStatus";

interface CardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, newStatus: TaskStatus) => void;
}

export const Card = ({ task, onEdit, onDelete, onStatusChange }: CardProps) => {
  const selectId = `task-status-${task.id}`;

  const handleStatusChange = (event: ChangeEvent<HTMLSelectElement>) => {
    onStatusChange(task.id, event.target.value as TaskStatus);
  };

  return (
    <div className="card-trello-style">
      <h3 className="card-title">{task.title}</h3>
      {task.description && <p>{task.description}</p>}

      <div className="status-selector">
        <label htmlFor={selectId}>Status</label>
        <select id={selectId} value={task.status} onChange={handleStatusChange}>
          {TASK_STATUSES.map((status) => (
            <option key={status} value={status}>
              {status}
            </option>
          ))}
        </select>
      </div>

      <div className="actions">
        <button type="button" onClick={() => onEdit(task)}>
          Editar
        </button>
        <button type="button" className="danger" onClick={() => onDelete(task.id)}>
          Excluir
        </button>
      </div>
    </div>
  );
};