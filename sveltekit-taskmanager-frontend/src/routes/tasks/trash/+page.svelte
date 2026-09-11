<script lang="ts">
	import { goto } from '$app/navigation';
	import { taskStore } from '$lib/features/tasks';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			// Filter to show only deleted tasks
			taskStore.setFilters({ status: 'deleted' });
			await taskStore.fetchTasks();
		} catch (error) {
			console.error('Failed to fetch deleted tasks:', error);
		}
	});

	/**
	 * Format date for display
	 */
	function formatDate(dateString?: string): string {
		if (!dateString) return 'No due date';
		return new Date(dateString).toLocaleDateString('en-US', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}

	/**
	 * Get status color class
	 */
	function getStatusColor(status: string): string {
		const colors = {
			cancelled: 'bg-red-100 text-red-800',
			completed: 'bg-green-100 text-green-800',
			deleted: 'bg-gray-300 text-gray-600',
			in_progress: 'bg-blue-100 text-blue-800',
			todo: 'bg-gray-100 text-gray-800'
		};
		return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
	}

	/**
	 * Handle permanent delete
	 */
	async function handlePermanentDelete(task: typeof taskStore.state.tasks[0]): Promise<void> {
		if (confirm(`Are you sure you want to permanently delete "${task.title}"? This action cannot be undone.`)) {
			try {
				await taskStore.permanentDeleteTask(task.id);
				await taskStore.fetchTasks();
			} catch (error) {
				console.error('Failed to permanently delete task:', error);
			}
		}
	}

	/**
	 * Handle restore task
	 */
	async function handleRestoreTask(task: typeof taskStore.state.tasks[0]): Promise<void> {
		try {
			await taskStore.restoreTask(task.id);
			await taskStore.fetchTasks();
		} catch (error) {
			console.error('Failed to restore task:', error);
		}
	}

	/**
	 * Handle view task
	 */
	function handleViewTask(task: typeof taskStore.state.tasks[0]): void {
		goto(`/tasks/${task.id}`);
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Trash</h1>
			<p class="text-gray-600 mt-1">Manage your deleted tasks</p>
		</div>
		<button
			onclick={() => goto('/tasks')}
			class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			Back to Tasks
		</button>
	</div>

	{#if taskStore.state.isLoading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-500">Loading trash...</p>
		</div>
	{:else if taskStore.state.tasks.length === 0}
		<div class="text-center py-12 bg-white rounded-lg shadow">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
			</svg>
			<h3 class="mt-2 text-sm font-medium text-gray-900">Trash is empty</h3>
			<p class="mt-1 text-sm text-gray-500">No deleted tasks to restore.</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each taskStore.state.tasks as task (task.id)}
				<div class="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-4">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2 mb-2">
								<h3 class="text-lg font-medium text-gray-900 truncate">{task.title}</h3>
								<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusColor(task.status)}">
									{task.status.replace('_', ' ')}
								</span>
							</div>
							{#if task.description}
								<p class="text-sm text-gray-600 mb-2 line-clamp-2">{task.description}</p>
							{/if}
							<div class="flex items-center gap-4 text-sm text-gray-500">
								<div class="flex items-center">
									<svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
									</svg>
									{formatDate(task.dueDate)}
								</div>
								<div class="flex items-center">
									<svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
									</svg>
									Deleted {new Date(task.updatedAt).toLocaleDateString()}
								</div>
							</div>
						</div>
						<div class="flex items-center gap-2 ml-4">
							<button
								onclick={() => handleViewTask(task)}
								class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
								title="View task"
							>
								<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
								</svg>
							</button>
							<button
								onclick={() => handleRestoreTask(task)}
								class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
								title="Restore task"
							>
								<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
								</svg>
							</button>
							<button
								onclick={() => handlePermanentDelete(task)}
								class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
								title="Permanently delete"
							>
								<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
								</svg>
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>