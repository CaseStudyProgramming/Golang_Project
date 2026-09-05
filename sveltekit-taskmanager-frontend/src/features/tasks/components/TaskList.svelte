<script lang="ts">
  import type { Task } from '../types/task.types';

  let { 
    tasks, 
    loading, 
    onToggleComplete, 
    onEdit, 
    onDelete 
  }: { 
    tasks: { value: Task[] }, 
    loading: { value: boolean }, 
    onToggleComplete: (task: Task) => void, 
    onEdit: (task: Task) => void, 
    onDelete: (id: number) => void 
  } = $props();
</script>

{#if loading.value}
  <div class="text-center py-8">
    <p class="text-gray-500">Loading tasks...</p>
  </div>
{:else if tasks.value.length === 0}
  <div class="text-center py-8">
    <p class="text-gray-500">No tasks found. Create your first task!</p>
  </div>
{:else}
  <div class="bg-white rounded-lg shadow overflow-hidden">
    <ul class="divide-y divide-gray-200">
      {#each tasks.value as task}
        <li class="p-4 hover:bg-gray-50 transition-colors">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4 flex-1">
              <input 
                type="checkbox" 
                checked={task.completed}
                onchange={() => onToggleComplete(task)}
                class="h-5 w-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
              />
              
              <div class="flex-1">
                <h3 class="text-lg font-medium {task.completed ? 'line-through text-gray-400' : 'text-gray-900'}">
                  {task.title}
                </h3>
                {#if task.sub_title}
                  <p class="text-sm text-gray-500">{task.sub_title}</p>
                {/if}
                {#if task.description}
                  <p class="text-sm text-gray-600 mt-1">{task.description}</p>
                {/if}
                <div class="flex gap-4 mt-2 text-xs text-gray-500">
                  <span>Created: {new Date(task.created_at).toLocaleDateString()}</span>
                  {#if task.due_date}
                    <span>Due: {new Date(task.due_date).toLocaleDateString()}</span>
                  {/if}
                  {#if task.priority}
                    <span class="px-2 py-1 rounded-full text-xs font-medium
                      {task.priority === 'URGENT' ? 'bg-red-100 text-red-800' : 
                       task.priority === 'HIGH' ? 'bg-orange-100 text-orange-800' :
                       task.priority === 'MEDIUM' ? 'bg-yellow-100 text-yellow-800' :
                       'bg-green-100 text-green-800'}">
                      {task.priority}
                    </span>
                  {/if}
                </div>
              </div>
            </div>
            
            <div class="flex gap-2">
              <button 
                onclick={() => onEdit(task)}
                class="text-blue-600 hover:text-blue-800 px-2 py-1"
              >
                Edit
              </button>
              <button 
                onclick={() => onDelete(task.id)}
                class="text-red-600 hover:text-red-800 px-2 py-1"
              >
                Delete
              </button>
            </div>
          </div>
        </li>
      {/each}
    </ul>
  </div>
{/if}
