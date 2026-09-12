<script lang="ts">
import { onMount } from 'svelte'
import { goto } from '$app/navigation'
import { page } from '$app/stores'
import { categoryStore } from '$lib/features/categories'
import { tagStore } from '$lib/features/tags'
import { activityStore, taskStore } from '$lib/features/tasks'
import ActivityLog from '$lib/features/tasks/components/ActivityLog.svelte'
import SubtaskList from '$lib/features/tasks/components/SubtaskList.svelte'
import TaskForm from '$lib/features/tasks/components/TaskForm.svelte'

let taskId = $derived($page.params.id || '')
let isEditing = $state(false)
let showDeleteConfirm = $state(false)

let categories = $derived(categoryStore.state.categories)
let tags = $derived(tagStore.state.tags)

onMount(async () => {
	if (!taskId) return
	try {
		await taskStore.fetchTaskById(taskId)
		await activityStore.fetchActivities({ taskId })
		await categoryStore.fetchCategories()
		await tagStore.fetchTags()
	} catch (error) {
		console.error('Failed to fetch task:', error)
	}
})

/**
 * Format date for display
 */
function formatDate(dateString?: string): string {
	if (!dateString) return 'Not set'
	return new Date(dateString).toLocaleDateString('en-US', {
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
		month: 'long',
		year: 'numeric',
	})
}

/**
 * Get category object by ID
 */
function getCategory(categoryId?: string) {
	if (!categoryId) return null
	return categories.find((c) => c.id === categoryId) || null
}

/**
 * Get priority color class
 */
function getPriorityColor(priority: string): string {
	const colors = {
		high: 'bg-orange-100 text-orange-800',
		low: 'bg-green-100 text-green-800',
		medium: 'bg-yellow-100 text-yellow-800',
		urgent: 'bg-red-100 text-red-800',
	}
	return colors[priority as keyof typeof colors] || 'bg-gray-100 text-gray-800'
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
		todo: 'bg-gray-100 text-gray-800',
	}
	return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800'
}

/**
 * Get tag objects by IDs
 */
function getTagObjects(tagIds?: string[]) {
	if (!tagIds || tagIds.length === 0) return []
	return tagIds
		.map((id) => tags.find((t) => t.id === id))
		.filter((tag): tag is (typeof tags)[0] => tag !== undefined)
}

/**
 * Handle add subtask
 */
async function handleAddSubtask(title: string): Promise<void> {
	if (!taskId) return
	try {
		await taskStore.addSubtask(taskId, title)
		await taskStore.fetchTaskById(taskId)
	} catch (error) {
		console.error('Failed to add subtask:', error)
	}
}

/**
 * Handle bulk complete subtasks
 */
async function handleBulkComplete(subtaskIds: string[]): Promise<void> {
	if (!taskId || subtaskIds.length === 0) return
	try {
		await taskStore.bulkCompleteSubtasks(taskId, subtaskIds)
		await taskStore.fetchTaskById(taskId)
	} catch (error) {
		console.error('Failed to bulk complete subtasks:', error)
	}
}

/**
 * Handle bulk delete subtasks
 */
async function handleBulkDelete(subtaskIds: string[]): Promise<void> {
	if (!taskId || subtaskIds.length === 0) return
	try {
		await taskStore.bulkDeleteSubtasks(taskId, subtaskIds)
		await taskStore.fetchTaskById(taskId)
	} catch (error) {
		console.error('Failed to bulk delete subtasks:', error)
	}
}

/**
 * Handle delete subtask
 */
async function handleDeleteSubtask(subtaskId: string): Promise<void> {
	if (!taskId) return
	try {
		await taskStore.deleteSubtask(taskId, subtaskId)
		await taskStore.fetchTaskById(taskId)
	} catch (error) {
		console.error('Failed to delete subtask:', error)
	}
}

/**
 * Handle delete task
 */
async function handleDeleteTask(): Promise<void> {
	if (!taskId) return
	try {
		await taskStore.deleteTask(taskId)
		goto('/tasks')
	} catch (error) {
		console.error('Failed to delete task:', error)
	}
}

/**
 * Handle toggle subtask
 */
async function handleToggleSubtask(subtaskId: string): Promise<void> {
	if (!taskId) return
	try {
		await taskStore.toggleSubtask(taskId, subtaskId)
		await taskStore.fetchTaskById(taskId)
	} catch (error) {
		console.error('Failed to toggle subtask:', error)
	}
}

/**
 * Handle update task
 */
async function handleUpdateTask(data: {
	categoryId?: string
	description?: string
	dueDate?: string
	priority?: string
	tags?: string[]
	title: string
}): Promise<void> {
	if (!taskId) return
	try {
		await taskStore.updateTask(taskId, data)
		isEditing = false
		await taskStore.fetchTaskById(taskId)
	} catch (error) {
		console.error('Failed to update task:', error)
	}
}
</script>

<div class="container mx-auto px-4 py-8">
	{#if taskStore.state.isLoading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-500">Loading task...</p>
		</div>
	{:else if taskStore.state.currentTask}
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center justify-between mb-6">
				<button
					onclick={() => goto('/tasks')}
					class="text-blue-600 hover:text-blue-700 flex items-center"
				>
					<svg class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
					</svg>
					Back to Tasks
				</button>
				<div class="flex gap-2">
					{#if !isEditing}
						<button
							onclick={() => (isEditing = true)}
							class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
						>
							Edit Task
						</button>
					{/if}
					<button
						onclick={() => (showDeleteConfirm = true)}
						class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
					>
						Delete Task
					</button>
				</div>
			</div>

			{#if isEditing}
				<div class="bg-white rounded-lg shadow p-6">
					<h2 class="text-xl font-semibold text-gray-900 mb-4">Edit Task</h2>
					<TaskForm
						mode="edit"
						initialData={taskStore.state.currentTask}
						onSubmit={handleUpdateTask}
						onCancel={() => (isEditing = false)}
					/>
				</div>
			{:else}
				<div class="bg-white rounded-lg shadow p-6">
					<div class="flex items-start justify-between mb-4">
						<h1 class="text-2xl font-bold text-gray-900">{taskStore.state.currentTask.title}</h1>
						<div class="flex gap-2">
							<span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium {getPriorityColor(taskStore.state.currentTask.priority)}">
								{taskStore.state.currentTask.priority}
							</span>
							<span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium {getStatusColor(taskStore.state.currentTask.status)}">
								{taskStore.state.currentTask.status.replace('_', ' ')}
							</span>
						</div>
					</div>

					{#if taskStore.state.currentTask.description}
						<div class="mb-6">
							<h2 class="text-sm font-medium text-gray-500 mb-2">Description</h2>
							<p class="text-gray-900">{taskStore.state.currentTask.description}</p>
						</div>
					{/if}

					<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
						<div>
							<h2 class="text-sm font-medium text-gray-500 mb-2">Due Date</h2>
							<p class="text-gray-900">{formatDate(taskStore.state.currentTask.dueDate)}</p>
						</div>
						<div>
							<h2 class="text-sm font-medium text-gray-500 mb-2">Category</h2>
							{#if getCategory(taskStore.state.currentTask.categoryId)}
								<p class="text-gray-900">
									{getCategory(taskStore.state.currentTask.categoryId)?.icon}
									{getCategory(taskStore.state.currentTask.categoryId)?.name}
								</p>
							{:else}
								<p class="text-gray-900">None</p>
							{/if}
						</div>
						<div>
							<h2 class="text-sm font-medium text-gray-500 mb-2">Created</h2>
							<p class="text-gray-900">{formatDate(taskStore.state.currentTask.createdAt)}</p>
						</div>
						<div>
							<h2 class="text-sm font-medium text-gray-500 mb-2">Last Updated</h2>
							<p class="text-gray-900">{formatDate(taskStore.state.currentTask.updatedAt)}</p>
						</div>
					</div>

					{#if taskStore.state.currentTask.tags && taskStore.state.currentTask.tags.length > 0}
						<div class="mb-6">
							<h2 class="text-sm font-medium text-gray-500 mb-2">Tags</h2>
							<div class="flex flex-wrap gap-2">
								{#each getTagObjects(taskStore.state.currentTask.tags) as tag (tag.id)}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" style="background-color: {tag.color || '#3B82F6'}20; color: {tag.color || '#3B82F6'}">
										{tag.name}
									</span>
								{/each}
							</div>
						</div>
					{/if}

					{#if taskStore.state.currentTask.completedAt}
						<div>
							<h2 class="text-sm font-medium text-gray-500 mb-2">Completed At</h2>
							<p class="text-gray-900">{formatDate(taskStore.state.currentTask.completedAt)}</p>
						</div>
					{/if}

					{#if taskStore.state.currentTask.progress !== undefined}
						<div>
							<h2 class="text-sm font-medium text-gray-500 mb-2">Progress</h2>
							<div class="flex items-center gap-3">
								<div class="flex-1 bg-gray-200 rounded-full h-2">
									<div
										class="bg-blue-600 h-2 rounded-full transition-all duration-300"
										style="width: {taskStore.state.currentTask.progress}%"
									></div>
								</div>
								<span class="text-sm text-gray-600">{taskStore.state.currentTask.progress}%</span>
							</div>
						</div>
					{/if}
				</div>

				{#if taskStore.state.currentTask.subtasks}
					<div class="mt-6">
						<SubtaskList
							bind:subtasks={taskStore.state.currentTask.subtasks}
							onAddSubtask={handleAddSubtask}
							onToggleSubtask={handleToggleSubtask}
							onDeleteSubtask={handleDeleteSubtask}
							onBulkComplete={handleBulkComplete}
							onBulkDelete={handleBulkDelete}
						/>
					</div>
				{/if}

				<div class="mt-6">
					<ActivityLog taskId={taskId} />
				</div>
			{/if}
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

	{#if showDeleteConfirm}
		<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
			<div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
				<h2 class="text-xl font-semibold text-gray-900 mb-4">Delete Task</h2>
				<p class="text-gray-600 mb-6">
					Are you sure you want to delete "{taskStore.state.currentTask?.title}"? This action can be undone from the trash.
				</p>
				<div class="flex justify-end gap-3">
					<button
						onclick={() => (showDeleteConfirm = false)}
						class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
					>
						Cancel
					</button>
					<button
						onclick={() => handleDeleteTask()}
						class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
					>
						Delete
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>