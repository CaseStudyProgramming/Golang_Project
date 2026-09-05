<script lang="ts">
	import { onMount } from 'svelte';
	import { taskStore, TaskFilters, TaskForm, TaskList } from '$lib/index';
	import type { Task } from '$lib/index';

	onMount(() => {
		taskStore.loadTasks();
	});

	$effect(() => {
		if (taskStore.searchQuery !== undefined || taskStore.filterCompleted !== undefined) {
			taskStore.loadTasks();
		}
	});
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-center">
		<h1 class="text-3xl font-bold text-gray-900">Task Manager</h1>
		<button 
			onclick={() => { taskStore.resetForm(); taskStore.showForm = true; }}
			class="btn btn-primary"
		>
			+ New Task
		</button>
	</div>

	<!-- Filters -->
	<TaskFilters 
		searchQuery={{ value: taskStore.searchQuery }} 
		filterCompleted={{ value: taskStore.filterCompleted }} 
	/>

	<!-- Task Form Modal -->
	<TaskForm 
		show={{ value: taskStore.showForm }}
		editingTask={{ value: taskStore.editingTask }}
		formData={{ value: taskStore.formData }}
		onSubmit={() => taskStore.handleSubmit()}
		onCancel={() => taskStore.resetForm()}
	/>

	<!-- Error Display -->
	{#if taskStore.error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
			{taskStore.error}
			<button onclick={() => taskStore.clearError()} class="ml-2 text-red-800 underline">Dismiss</button>
		</div>
	{/if}

	<!-- Task List -->
	<TaskList 
		tasks={{ value: taskStore.tasks }}
		loading={{ value: taskStore.loading }}
		onToggleComplete={(task: Task) => taskStore.toggleComplete(task)}
		onEdit={(task: Task) => taskStore.editTask(task)}
		onDelete={(id: number) => taskStore.handleDelete(id)}
	/>
</div>
