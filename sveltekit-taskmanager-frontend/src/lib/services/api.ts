import type { 
  Task, 
  TaskFormData, 
  TaskQueryParams, 
  TaskListAPIResponse, 
  APIResponse 
} from '$lib/types';

// Use proxy during development, direct URL in production
const API_BASE_URL = import.meta.env.DEV ? '/api' : 'http://localhost:8080';

class APIService {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'An error occurred' }));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Task operations
  async getTasks(params?: TaskQueryParams): Promise<TaskListAPIResponse> {
    const queryParams = new URLSearchParams();
    
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.completed !== undefined) queryParams.append('completed', params.completed.toString());
    if (params?.search) queryParams.append('search', params.search);
    if (params?.priority) queryParams.append('priority', params.priority);
    if (params?.category_id) queryParams.append('category_id', params.category_id.toString());
    if (params?.sort_by) queryParams.append('sort_by', params.sort_by);
    if (params?.order) queryParams.append('order', params.order);

    const queryString = queryParams.toString();
    const endpoint = `/tasks${queryString ? `?${queryString}` : ''}`;
    
    return this.request<TaskListAPIResponse>(endpoint);
  }

  async getTask(id: number): Promise<APIResponse<Task>> {
    return this.request<APIResponse<Task>>(`/tasks/${id}`);
  }

  async createTask(data: TaskFormData): Promise<APIResponse<Task>> {
    return this.request<APIResponse<Task>>('/tasks', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateTask(id: number, data: Partial<TaskFormData>): Promise<APIResponse<Task>> {
    return this.request<APIResponse<Task>>(`/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteTask(id: number): Promise<APIResponse<void>> {
    return this.request<APIResponse<void>>(`/tasks/${id}`, {
      method: 'DELETE',
    });
  }

  async completeTask(id: number): Promise<APIResponse<Task>> {
    return this.request<APIResponse<Task>>(`/tasks/${id}/complete`, {
      method: 'PATCH',
    });
  }

  async uncompleteTask(id: number): Promise<APIResponse<Task>> {
    return this.request<APIResponse<Task>>(`/tasks/${id}/uncomplete`, {
      method: 'PATCH',
    });
  }

  async restoreTask(id: number): Promise<APIResponse<Task>> {
    return this.request<APIResponse<Task>>(`/tasks/${id}/restore`, {
      method: 'PATCH',
    });
  }

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    return this.request<{ status: string }>('/health');
  }
}

export const apiService = new APIService();