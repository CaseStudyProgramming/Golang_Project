// Placeholder TypeScript file for type definitions
// This file will be expanded when the project is properly initialized

export interface Task {
  id: number;
  title: string;
  completed: boolean;
  created_at: string;
  deleted_at: string | null;
}

export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
}
