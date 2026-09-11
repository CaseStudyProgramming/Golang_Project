<script lang="ts">
	import { categoryStore } from '$lib/features/categories';
	import { tagStore } from '$lib/features/tags';
	import { EmptyState, LoadingSpinner, ProgressBar } from '$lib/shared/components';
	import { confirmStore } from '$lib/shared/stores';

	import type { Task } from '../types/task.types';

	import { taskStore } from '../stores/task.store';

	let { 
		isLoading = $bindable(taskStore.state.isLoading),
		onDeleteTask,
		onEditTask,
		onViewTask,
		tasks = $bindable(taskStore.state.tasks)
	}: {
		isLoading?: boolean;
		onDeleteTask?: (task: Task) => void;
		onEditTask?: (task: Task) => void;
		onViewTask?: (task: Task) => void;
		tasks?: Task[];
	} = $props();

	let categories = $derived(categoryStore.state.categories);
	let tags = $derived(tagStore.state.tags);

	/**
	 * Initialize categories and tags on mount
	 */
	$effect(() => {
		categoryStore.fetchCategories();
		tagStore.fetchTags();
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
	 * Get category by ID
	 */
	function getCategory(categoryId?: string) {
		if (!categoryId) return null;
		return categories.find((c) => c.id === categoryId) || null;
	}

	/**
	 * Get priority color class
	 */
	function getPriorityColor(priority: string): string {
		const colors = {
			high: 'bg-orange-100 text-orange-800',
			low: 'bg-green-100 text-green-800',
			medium: 'bg-yellow-100 text-yellow-800',
			urgent: 'bg-red-100 text-red-800'
		};
		return colors[priority as keyof typeof colors] || 'bg-gray-100 text-gray-800';
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
	 * Get tag objects by IDs
	 */
	function getTagObjects(tagIds?: string[]) {
		if (!tagIds || tagIds.length === 0) return [];
		return tagIds
			.map((id) => tags.find((t) => t.id === id))
			.filter((tag): tag is typeof tags[0] => tag !== undefined);
	}

	/**
	 * Handle delete task with confirmation
	 */
	async function handleDeleteTask(task: Task) {
		const confirmed = await confirmStore.showConfirm({
			cancelText: 'Cancel',
			confirmText: 'Delete',
			message: `Are you sure you want to delete "${task.title}"? This action cannot be undone.`,
			title: 'Delete Task',
			type: 'danger'
		});

		if (confirmed && onDeleteTask) {
			onDeleteTask(task);
		}
	}
</script>

<div class="space-y-3 sm:space-y-4">
	{#if isLoading}
		<div class="text-center py-8 sm:py-12">
			<LoadingSpinner text="Loading tasks..." />
		</div>
	{:else if tasks.length === 0}
		<div class="bg-white rounded-lg shadow">
			<EmptyState 
				icon='<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>'
				title="No tasks"
				description="Get started by creating a new task to organize your work."
			/>
		</div>
	{:else}
		<div class="space-y-3">
			{#each tasks as task (task.id)}
				<div class="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-3 sm:p-4">
					<div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
						<div class="flex-1 min-w-0">
							<div class="flex flex-wrap items-center gap-2 mb-2">
								<h3 class="text-base sm:text-lg font-medium text-gray-900 truncate">{task.title}</h3>
								<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {getPriorityColor(task.priority)}">
									{task.priority}
								</span>
								<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {getStatusColor(task.status)}">
									{task.status.replace('_', ' ')}
								</span>
							</div>
							{#if task.description}
								<p class="text-sm text-gray-600 mb-2 line-clamp-2">{task.description}</p>
							{/if}
							<div class="flex flex-wrap items-center gap-3 sm:gap-4 text-sm text-gray-500">
								<div class="flex items-center">
									<svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
									</svg>
									{formatDate(task.dueDate)}
								</div>
								{#if getCategory(task.categoryId)}
									<div class="flex items-center">
										<svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
										</svg>
										{getCategory(task.categoryId)?.icon} {getCategory(task.categoryId)?.name}
									</div>
								{/if}
							</div>
							{#if task.progress !== undefined && task.progress > 0}
								<div class="mt-2">
									<ProgressBar progress={task.progress} size="sm" showLabel={true} />
								</div>
							{/if}
							{#if task.tags && task.tags.length > 0}
								<div class="flex flex-wrap gap-1 mt-2">
									{#each getTagObjects(task.tags) as tag (tag.id)}
										<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" style="background-color: {tag.color || '#3B82F6'}20; color: {tag.color || '#3B82F6'}">
											{tag.name}
										</span>
									{/each}
								</div>
							{/if}
						</div>
						<div class="flex items-center gap-2 sm:ml-4">
							{#if onViewTask}
								<button
									onclick={() => onViewTask(task)}
									class="p-3 sm:p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 active:bg-blue-100 rounded-lg transition-colors min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0"
									title="View task"
									aria-label="View task"
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
									class="p-3 sm:p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 active:bg-green-100 rounded-lg transition-colors min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0"
									title="Edit task"
									aria-label="Edit task"
								>
									<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
									</svg>
								</button>
							{/if}
							{#if onDeleteTask}
								<button
									onclick={() => handleDeleteTask(task)}
									class="p-3 sm:p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 active:bg-red-100 rounded-lg transition-colors min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0"
									title="Delete task"
									aria-label="Delete task"
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