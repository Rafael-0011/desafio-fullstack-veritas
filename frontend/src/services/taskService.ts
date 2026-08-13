import type { Task } from '../model/Task';
import { TASK_STATUSES, type TaskStatus } from '../model/TaskStatus';

const API_URL = (import.meta.env.VITE_API_URL || 'http://localhost:3000').replace(/\/$/, '');
const TASKS_ENDPOINT = `${API_URL}/tasks`;

const VALID_STATUSES = new Set<string>(TASK_STATUSES);

type TaskDTO = {
	id: string;
	title: string;
	description?: string | null;
	status: string;
	createdAt?: string | null;
};

export type CreateTaskInput = {
	title: string;
	description?: string;
	status: TaskStatus;
};

export type UpdateTaskInput = {
	title?: string;
	description?: string;
	status?: TaskStatus;
};

const toTask = (dto: TaskDTO): Task => {
	if (!VALID_STATUSES.has(dto.status)) {
		throw new Error(`Status invalido recebido da API: ${dto.status}`);
	}

	return {
		id: dto.id,
		title: dto.title,
		description: dto.description ?? undefined,
		status: dto.status as TaskStatus,
		createdAt: dto.createdAt ? new Date(dto.createdAt) : undefined,
	};
};

const request = async <T>(url: string, init?: RequestInit): Promise<T> => {
	const response = await fetch(url, {
		headers: { 'Content-Type': 'application/json' },
		...init,
	});

	if (!response.ok) {
		const message = await response.text();
		throw new Error(message || `Erro HTTP ${response.status}`);
	}

	if (response.status === 204) {
		return undefined as T;
	}

	return (await response.json()) as T;
};

export const taskService = {
	async getAll(): Promise<Task[]> {
		const data = await request<TaskDTO[]>(TASKS_ENDPOINT);
		return data.map(toTask);
	},

	async getById(id: string): Promise<Task> {
		const data = await request<TaskDTO>(`${TASKS_ENDPOINT}/${id}`);
		return toTask(data);
	},

	async create(input: CreateTaskInput): Promise<Task> {
		const data = await request<TaskDTO>(TASKS_ENDPOINT, {
			method: 'POST',
			body: JSON.stringify(input),
		});
		return toTask(data);
	},

	async update(id: string, input: UpdateTaskInput): Promise<Task> {
		const data = await request<TaskDTO>(`${TASKS_ENDPOINT}/${id}`, {
			method: 'PUT',
			body: JSON.stringify(input),
		});
		return toTask(data);
	},

	async remove(id: string): Promise<void> {
		await request<void>(`${TASKS_ENDPOINT}/${id}`, { method: 'DELETE' });
	},
};
