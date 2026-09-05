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
  priority?: "LOW" | "MEDIUM" | "HIGH" | "URGENT";
  category_id?: number;
}

export interface TaskFormData {
  title: string;
  sub_title?: string;
  description?: string;
  due_date?: string;
  priority?: "LOW" | "MEDIUM" | "HIGH" | "URGENT";
  category_id?: number;
}

export interface TaskQueryParams {
  page?: number;
  limit?: number;
  completed?: boolean;
  search?: string;
  priority?: string;
  category_id?: number;
  sort_by?: string;
  order?: "asc" | "desc";
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
  status: "success" | "error";
  message: string;
  data: T;
}

export interface TaskListAPIResponse {
  status: "success" | "error";
  message: string;
  data: TaskListResponse;
}
