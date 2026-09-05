<script lang="ts">
	import { onMount } from 'svelte';
	import { apiService } from '$lib/services/api';
	import type { Task, TaskFormData, TaskQueryParams } from '$lib/types';

	let tasks: Task[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	
	// Form state
	let showForm = $state(false);
	let editingTask = $state<Task | null>(null);
	let formData = $state<TaskFormData>({
		title: '',
		sub_title: '',
		description: '',
		due_date: '',
		priority: 'MEDIUM'
	});

	// Filter state
	let searchQuery = $state('');
	let filterCompleted = $state<boolean | undefined>(undefined);
	let currentPage = $state(1);
	let limit = $state(10);

	async function loadTasks() {
		loading = true;
		error = null;
		
		try {
			const params: TaskQueryParams = {
				page: currentPage,
				limit,
				search: searchQuery || undefined,
				completed: filterCompleted
			};

			const response = await apiService.getTasks(params);
			if (response.status === 'success') {
				tasks = response.data.data;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load tasks';
		} finally {
			loading = false;
		}
	}

	async function handleSubmit() {
		try {
			if (editingTask) {
				await apiService.updateTask(editingTask.id, formData);
			} else {
				await apiService.createTask(formData);
			}
			resetForm();
			await loadTasks();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save task';
		}
	}

	async function handleDelete(id: number) {
		if (!confirm('Are you sure you want to delete this task?')) return;
		
		try {
			await apiService.deleteTask(id);
			await loadTasks();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete task';
		}
	}

	async function toggleComplete(task: Task) {
		try {
			if (task.completed) {
				await apiService.uncompleteTask(task.id);
			} else {
				await apiService.completeTask(task.id);
			}
			await loadTasks();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update task';
		}
	}

	function editTask(task: Task) {
		editingTask = task;
		formData = {
			title: task.title,
			sub_title: task.sub_title,
			description: task.description,
			due_date: task.due_date || '',
			priority: task.priority || 'MEDIUM'
		};
		showForm = true;
	}

	function resetForm() {
		editingTask = null;
		formData = {
			title: '',
			sub_title: '',
			description: '',
			due_date: '',
			priority: 'MEDIUM'
		};
		showForm = false;
	}

	onMount(() => {
		loadTasks();
	});

	$: if (searchQuery !== undefined || filterCompleted !== undefined) {
		loadTasks();
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-center">
		<h1 class="text-3xl font-bold text-gray-900">Task Manager</h1>
		<button 
			onclick={() => { resetForm(); showForm = true; }}
			class="btn btn-primary"
		>
			+ New Task
		</button>
	</div>

	<!-- Filters -->
	<div class="bg-white rounded-lg shadow p-4">
		<div class="flex flex-col sm:flex-row gap-4">
			<div class="flex-1">
				<input 
					type="text" 
					placeholder="Search tasks..." 
					bind:value={searchQuery}
					class="input"
				/>
			</div>
			<div class="flex gap-2">
				<button 
					class="btn btn-secondary"
					class:btn-primary={filterCompleted === undefined}
					onclick={() => filterCompleted = undefined}
				>
					All
				</button>
				<button 
					class="btn btn-secondary"
					class:btn-primary={filterCompleted === false}
					onclick={() => filterCompleted = false}
				>
					Active
				</button>
				<button 
					class="btn btn-secondary"
					class:btn-primary={filterCompleted === true}
					onclick={() => filterCompleted = true}
				>
					Completed
				</button>
			</div>
		</div>
	</div>

	<!-- Task Form Modal -->
	{#if showForm}
		<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
			<div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md mx-4">
				<h2 class="text-xl font-bold mb-4">
					{editingTask ? 'Edit Task' : 'New Task'}
				</h2>
				
				<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
					<div class="space-y-4">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
							<input 
								type="text" 
								bind:value={formData.title}
								class="input"
								required
							/>
						</div>
						
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Subtitle</label>
							<input 
								type="text" 
								bind:value={formData.sub_title}
								class="input"
							/>
						</div>
						
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
							<textarea 
								bind:value={formData.description}
								class="input"
								rows="3"
							/>
						</div>
						
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
							<input 
								type="datetime-local" 
								bind:value={formData.due_date}
								class="input"
							/>
						</div>
						
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
							<select bind:value={formData.priority} class="input">
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
							onclick={resetForm}
							class="btn btn-secondary"
						>
							Cancel
						</button>
						<button type="submit" class="btn btn-primary">
							{editingTask ? 'Update' : 'Create'}
						</button>
					</div>
				</form>
			</div>
		</div>
	{/if}

	<!-- Error Display -->
	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
			{error}
		</div>
	{/if}

	<!-- Task List -->
	{#if loading}
		<div class="text-center py-8">
			<p class="text-gray-500">Loading tasks...</p>
		</div>
	{:else if tasks.length === 0}
		<div class="text-center py-8">
			<p class="text-gray-500">No tasks found. Create your first task!</p>
		</div>
	{:else}
		<div class="bg-white rounded-lg shadow overflow-hidden">
			<ul class="divide-y divide-gray-200">
				{#each tasks as task}
					<li class="p-4 hover:bg-gray-50 transition-colors">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-4 flex-1">
								<input 
									type="checkbox" 
									checked={task.completed}
									onchange={() => toggleComplete(task)}
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
									onclick={() => editTask(task)}
									class="text-blue-600 hover:text-blue-800 px-2 py-1"
								>
									Edit
								</button>
								<button 
									onclick={() => handleDelete(task.id)}
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
</div>
