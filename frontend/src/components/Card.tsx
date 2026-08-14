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
    event.stopPropagation();
    onStatusChange(task.id, event.target.value as TaskStatus);
  };

  return (
    <div
      className="card-trello-style"
      role="button"
      tabIndex={0}
      onClick={() => onEdit(task)}
      onKeyDown={(event) => {
        if (event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          onEdit(task);
        }
      }}
    >
      <h3 className="card-title">{task.title}</h3>
      {task.description && <p>{task.description}</p>}

      <div className="status-selector" onClick={(e) => e.stopPropagation()}>
        <label htmlFor={selectId}>Status</label>
        <select
          id={selectId}
          value={task.status}
          onChange={handleStatusChange}
          onClick={(e) => e.stopPropagation()}
        >
          {TASK_STATUSES.map((status) => (
            <option key={status} value={status}>
              {status}
            </option>
          ))}
        </select>
      </div>

      <div className="actions" onClick={(e) => e.stopPropagation()}>
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            onEdit(task);
          }}
        >
          Editar
        </button>
        <button
          type="button"
          className="danger"
          onClick={(e) => {
            e.stopPropagation();
            onDelete(task.id);
          }}
        >
          Excluir
        </button>
      </div>
    </div>
  );
};