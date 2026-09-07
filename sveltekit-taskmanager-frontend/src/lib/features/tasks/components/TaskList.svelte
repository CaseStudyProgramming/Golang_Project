<script lang="ts">
	import { taskStore } from '../stores/task.store';
	import type { Task } from '../types/task.types';

	let { 
		tasks = $bindable(taskStore.state.tasks),
		isLoading = $bindable(taskStore.state.isLoading),
		onViewTask,
		onEditTask,
		onDeleteTask
	}: {
		tasks?: Task[];
		isLoading?: boolean;
		onViewTask?: (task: Task) => void;
		onEditTask?: (task: Task) => void;
		onDeleteTask?: (task: Task) => void;
	} = $props();

	/**
	 * Get priority color class
	 */
	function getPriorityColor(priority: string): string {
		const colors = {
			low: 'bg-green-100 text-green-800',
			medium: 'bg-yellow-100 text-yellow-800',
			high: 'bg-orange-100 text-orange-800',
			urgent: 'bg-red-100 text-red-800'
		};
		return colors[priority as keyof typeof colors] || 'bg-gray-100 text-gray-800';
	}

	/**
	 * Get status color class
	 */
	function getStatusColor(status: string): string {
		const colors = {
			todo: 'bg-gray-100 text-gray-800',
			in_progress: 'bg-blue-100 text-blue-800',
			completed: 'bg-green-100 text-green-800',
			cancelled: 'bg-red-100 text-red-800',
			deleted: 'bg-gray-300 text-gray-600'
		};
		return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
	}

	/**
	 * Format date for display
	 */
	function formatDate(dateString?: string): string {
		if (!dateString) return 'No due date';
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
</script>

<div class="space-y-4">
	{#if isLoading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-500">Loading tasks...</p>
		</div>
	{:else if tasks.length === 0}
		<div class="text-center py-12 bg-white rounded-lg shadow">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
			</svg>
			<h3 class="mt-2 text-sm font-medium text-gray-900">No tasks</h3>
			<p class="mt-1 text-sm text-gray-500">Get started by creating a new task.</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each tasks as task}
				<div class="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-4">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2 mb-2">
								<h3 class="text-lg font-medium text-gray-900 truncate">{task.title}</h3>
								<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getPriorityColor(task.priority)}">
									{task.priority}
								</span>
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
								{#if task.categoryId}
									<div class="flex items-center">
										<svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
										</svg>
										{task.categoryId}
									</div>
								{/if}
							</div>
						</div>
						<div class="flex items-center gap-2 ml-4">
							{#if onViewTask}
								<button
									onclick={() => onViewTask(task)}
									class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
									title="View task"
								>
									<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
									</svg>
								</button>
							{/if}
							{#if onEditTask}
								<button
									onclick={() => onEditTask(task)}
									class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
									title="Edit task"
								>
									<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
									</svg>
								</button>
							{/if}
							{#if onDeleteTask}
								<button
									onclick={() => onDeleteTask(task)}
									class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
									title="Delete task"
								>
									<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
									</svg>
								</button>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>