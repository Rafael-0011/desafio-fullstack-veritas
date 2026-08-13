import type { Task } from "./Task";

export const initialTasks: Task[] = [
  {
    id: '1',
    title: 'Configurar Monorepo',
    description: 'Estruturar pastas backend e frontend conforme o edital.',
    status: 'A Fazer'
  },
  {
    id: '2',
    title: 'Desenhar User Flow',
    description: 'Criar o diagrama obrigatório para a pasta /docs.',
    status: 'Em Progresso'
  },
  {
    id: '3',
    title: 'Criar README.md',
    description: 'Documentar instruções de execução e decisões técnicas.',
    status: 'Concluídas'
  }
];