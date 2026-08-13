export const TASK_STATUSES = ['A Fazer', 'Em Progresso', 'Concluídas'] as const;

export type TaskStatus = (typeof TASK_STATUSES)[number];
