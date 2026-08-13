import type { TaskStatus } from "./TaskStatus";

export interface Task {
  id: string;                
  title: string;             
  description?: string;     
  status: TaskStatus;        
  createdAt?: Date;      
}