// Task Types based on backend API response structure
export interface Task {
  id: number;
  title: string;
  completed: boolean;
  created_at: string;
  updated_at?: string;
  deleted_at: string | null;
  sub_title?: string;
  description?: string;
  due_date?: string;
  priority?: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
  category_id?: number;
}

export interface TaskFormData {
  title: string;
  sub_title?: string;
  description?: string;
  due_date?: string;
  priority?: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
  category_id?: number;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total_data: number;
  total_page: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface TaskListResponse {
  data: Task[];
  meta: PaginationMeta;
}

export interface APIResponse<T> {
  status: 'success' | 'error';
  message: string;
  data: T;
}

export interface TaskListAPIResponse {
  status: 'success' | 'error';
  message: string;
  data: TaskListResponse;
}

export interface TaskQueryParams {
  page?: number;
  limit?: number;
  completed?: boolean;
  search?: string;
  priority?: string;
  category_id?: number;
  sort_by?: string;
  order?: 'asc' | 'desc';
}

export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

export interface Category {
  id: number;
  user_id: number;
  name: string;
  color_hex: string;
  created_at: string;
}

export interface TaskSummary {
  total_tasks: number;
  completed_tasks: number;
  pending_tasks: number;
  overdue_tasks: number;
  completion_rate: number;
  priority_distribution: {
    LOW: number;
    MEDIUM: number;
    HIGH: number;
    URGENT: number;
  };
}