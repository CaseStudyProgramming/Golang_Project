<script lang="ts">
  import type { TaskFormData, Task } from '../types/task.types';

  const { 
    show, 
    editingTask, 
    formData, 
    onSubmit, 
    onCancel 
  }: { 
    show: { value: boolean }, 
    editingTask: { value: Task | null }, 
    formData: { value: TaskFormData }, 
    onSubmit: () => void, 
    onCancel: () => void 
  } = $props();
</script>

{#if show.value}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md mx-4">
      <h2 class="text-xl font-bold mb-4">
        {editingTask.value ? 'Edit Task' : 'New Task'}
      </h2>
      
      <form onsubmit={(e) => { e.preventDefault(); onSubmit(); }}>
        <div class="space-y-4">
          <div>
            <label for="task-title" class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
            <input 
              id="task-title"
              type="text" 
              bind:value={formData.value.title}
              class="input"
              required
            />
          </div>
          
          <div>
            <label for="task-subtitle" class="block text-sm font-medium text-gray-700 mb-1">Subtitle</label>
            <input 
              id="task-subtitle"
              type="text" 
              bind:value={formData.value.sub_title}
              class="input"
            />
          </div>
          
          <div>
            <label for="task-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea 
              id="task-description"
              bind:value={formData.value.description}
              class="input"
              rows="3"
            ></textarea>
          </div>
          
          <div>
            <label for="task-due-date" class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
            <input 
              id="task-due-date"
              type="datetime-local" 
              bind:value={formData.value.due_date}
              class="input"
            />
          </div>
          
          <div>
            <label for="task-priority" class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
            <select id="task-priority" bind:value={formData.value.priority} class="input">
              <option value="LOW">Low</option>
              <option value="MEDIUM">Medium</option>
              <option value="HIGH">High</option>
              <option value="URGENT">Urgent</option>
            </select>
          </div>
        </div>
        
        <div class="flex justify-end gap-2 mt-6">
          <button 
            type="button" 
            onclick={onCancel}
            class="btn btn-secondary"
          >
            Cancel
          </button>
          <button type="submit" class="btn btn-primary">
            {editingTask.value ? 'Update' : 'Create'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
