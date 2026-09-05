import { taskAPI } from '../api/task.api';
import type { Task, TaskFormData, TaskQueryParams } from '../types/task.types';

class TaskStore {
  tasks = $state<Task[]>([]);
  loading = $state(true);
  error = $state<string | null>(null);
  
  // Filter state
  searchQuery = $state('');
  filterCompleted = $state<boolean | undefined>(undefined);
  currentPage = $state(1);
  limit = $state(10);

  // Form state
  showForm = $state(false);
  editingTask = $state<Task | null>(null);
  formData = $state<TaskFormData>({
    title: '',
    sub_title: '',
    description: '',
    due_date: '',
    priority: 'MEDIUM'
  });

  async loadTasks() {
    this.loading = true;
    this.error = null;
    
    try {
      const params: TaskQueryParams = {
        page: this.currentPage,
        limit: this.limit,
        search: this.searchQuery || undefined,
        completed: this.filterCompleted
      };

      const response = await taskAPI.getTasks(params);
      if (response.status === 'success') {
        this.tasks = response.data.data;
      }
    } catch (err) {
      this.error = err instanceof Error ? err.message : 'Failed to load tasks';
    } finally {
      this.loading = false;
    }
  }

  async handleSubmit() {
    try {
      if (this.editingTask) {
        await taskAPI.updateTask(this.editingTask.id, this.formData);
      } else {
        await taskAPI.createTask(this.formData);
      }
      this.resetForm();
      await this.loadTasks();
    } catch (err) {
      this.error = err instanceof Error ? err.message : 'Failed to save task';
    }
  }

  async handleDelete(id: number) {
    if (!confirm('Are you sure you want to delete this task?')) return;
    
    try {
      await taskAPI.deleteTask(id);
      await this.loadTasks();
    } catch (err) {
      this.error = err instanceof Error ? err.message : 'Failed to delete task';
    }
  }

  async toggleComplete(task: Task) {
    try {
      if (task.completed) {
        await taskAPI.uncompleteTask(task.id);
      } else {
        await taskAPI.completeTask(task.id);
      }
      await this.loadTasks();
    } catch (err) {
      this.error = err instanceof Error ? err.message : 'Failed to update task';
    }
  }

  editTask(task: Task) {
    this.editingTask = task;
    this.formData = {
      title: task.title,
      sub_title: task.sub_title,
      description: task.description,
      due_date: task.due_date || '',
      priority: task.priority || 'MEDIUM'
    };
    this.showForm = true;
  }

  resetForm() {
    this.editingTask = null;
    this.formData = {
      title: '',
      sub_title: '',
      description: '',
      due_date: '',
      priority: 'MEDIUM'
    };
    this.showForm = false;
  }

  clearError() {
    this.error = null;
  }
}

export const taskStore = new TaskStore();
