import type {
  Task,
  TaskFormData,
  TaskQueryParams,
  TaskListAPIResponse,
  APIResponse,
} from "../types/task.types";
import { apiUtils } from "../../../shared/utils/api.utils";

class TaskAPI {

  async getTasks(params?: TaskQueryParams): Promise<TaskListAPIResponse> {
    const queryParams = new URLSearchParams();

    if (params?.page) queryParams.append("page", params.page.toString());
    if (params?.limit) queryParams.append("limit", params.limit.toString());
    if (params?.completed !== undefined)
      queryParams.append("completed", params.completed.toString());
    if (params?.search) queryParams.append("search", params.search);
    if (params?.priority) queryParams.append("priority", params.priority);
    if (params?.category_id) queryParams.append("category_id", params.category_id.toString());
    if (params?.sort_by) queryParams.append("sort_by", params.sort_by);
    if (params?.order) queryParams.append("order", params.order);

    const queryString = queryParams.toString();
    const endpoint = `/tasks${queryString ? `?${queryString}` : ""}`;

    return apiUtils.request<TaskListAPIResponse>(endpoint);
  }

  async getTask(id: number): Promise<APIResponse<Task>> {
    return apiUtils.request<APIResponse<Task>>(`/tasks/${id}`);
  }

  async createTask(data: TaskFormData): Promise<APIResponse<Task>> {
    return apiUtils.request<APIResponse<Task>>("/tasks", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateTask(id: number, data: Partial<TaskFormData>): Promise<APIResponse<Task>> {
    return apiUtils.request<APIResponse<Task>>(`/tasks/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteTask(id: number): Promise<APIResponse<void>> {
    return apiUtils.request<APIResponse<void>>(`/tasks/${id}`, {
      method: "DELETE",
    });
  }

  async completeTask(id: number): Promise<APIResponse<Task>> {
    return apiUtils.request<APIResponse<Task>>(`/tasks/${id}/complete`, {
      method: "PATCH",
    });
  }

  async uncompleteTask(id: number): Promise<APIResponse<Task>> {
    return apiUtils.request<APIResponse<Task>>(`/tasks/${id}/uncomplete`, {
      method: "PATCH",
    });
  }

  async restoreTask(id: number): Promise<APIResponse<Task>> {
    return apiUtils.request<APIResponse<Task>>(`/tasks/${id}/restore`, {
      method: "PATCH",
    });
  }
}

export const taskAPI = new TaskAPI();
