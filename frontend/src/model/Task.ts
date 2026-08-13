import type { TaskStatus } from "./TaskStatus";

export interface Task {
  id: string;                // Identificador único (necessário para editar e excluir)
  title: string;             // Campo OBRIGATÓRIO conforme o edital [3, 4]
  description?: string;      // Campo OPCIONAL conforme o edital [3, 4]
  status: TaskStatus;        // Deve ser um dos três status fixos [3, 5]
  createdAt?: Date;          // Opcional: útil para ordenação ou controle interno
}