<script lang="ts">
	import { taskStore, TaskForm } from '$lib/features/tasks';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	let taskId = $derived($page.params.id || '');

	onMount(async () => {
		if (!taskId) return;
		try {
			await taskStore.fetchTaskById(taskId);
		} catch (error) {
			console.error('Failed to fetch task:', error);
		}
	});

	/**
	 * Handle update task
	 */
	async function handleUpdateTask(data: any): Promise<void> {
		if (!taskId) return;
		try {
			await taskStore.updateTask(taskId, data);
			goto(`/tasks/${taskId}`);
		} catch (error) {
			console.error('Failed to update task:', error);
		}
	}

	/**
	 * Handle cancel
	 */
	function handleCancel(): void {
		goto(`/tasks/${taskId}`);
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="flex items-center justify-between mb-6">
			<button
				onclick={handleCancel}
				class="text-blue-600 hover:text-blue-700 flex items-center"
			>
				<svg class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
				</svg>
				Back to Task
			</button>
			<h1 class="text-2xl font-bold text-gray-900">Edit Task</h1>
		</div>

		{#if taskStore.state.isLoading}
			<div class="text-center py-12">
				<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
				<p class="mt-4 text-gray-500">Loading task...</p>
			</div>
		{:else if taskStore.state.currentTask}
			<div class="bg-white rounded-lg shadow p-6">
				<TaskForm
					mode="edit"
					initialData={taskStore.state.currentTask}
					onSubmit={handleUpdateTask}
					onCancel={handleCancel}
				/>
			</div>
		{:else}
			<div class="text-center py-12">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900">Task not found</h3>
				<p class="mt-1 text-sm text-gray-500">The task you're looking for doesn't exist.</p>
				<button
					onclick={() => goto('/tasks')}
					class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				>
					Back to Tasks
				</button>
			</div>
		{/if}
	</div>
</div>