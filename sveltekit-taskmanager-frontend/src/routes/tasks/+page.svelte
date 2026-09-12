<script lang="ts">
import { onMount } from 'svelte'
import { goto } from '$app/navigation'
import { taskStore } from '$lib/features/tasks'
import TaskSearch from '$lib/features/tasks/components/TaskSearch.svelte'
import TaskFilters from '$lib/features/tasks/components/TaskFilters.svelte'
import TaskList from '$lib/features/tasks/components/TaskList.svelte'
import Pagination from '$lib/features/tasks/components/Pagination.svelte'
import TaskForm from '$lib/features/tasks/components/TaskForm.svelte'

let showCreateModal = $state(false)
let searchQuery = $state('')

onMount(async () => {
	try {
		await taskStore.fetchTasks()
	} catch (error) {
		console.error('Failed to fetch tasks:', error)
	}
})

/**
 * Handle clear filters
 */
function handleClearFilters(): void {
	taskStore.clearFilters()
	taskStore.fetchTasks()
}

/**
 * Handle create task
 */
async function handleCreateTask(data: {
	categoryId?: string
	description?: string
	dueDate?: string
	priority?: string
	tags?: string[]
	title: string
}): Promise<void> {
	try {
		await taskStore.createTask(data)
		showCreateModal = false
		await taskStore.fetchTasks()
	} catch (error) {
		console.error('Failed to create task:', error)
	}
}

/**
 * Handle delete task
 */
async function handleDeleteTask(task: (typeof taskStore.state.tasks)[0]): Promise<void> {
	if (confirm(`Are you sure you want to delete "${task.title}"?`)) {
		try {
			await taskStore.deleteTask(task.id)
		} catch (error) {
			console.error('Failed to delete task:', error)
		}
	}
}

/**
 * Handle edit task
 */
function handleEditTask(task: (typeof taskStore.state.tasks)[0]): void {
	goto(`/tasks/${task.id}/edit`)
}

/**
 * Handle filter change
 */
function handleFilterChange(filters: typeof taskStore.state.filters): void {
	taskStore.setFilters(filters)
	taskStore.fetchTasks()
}

/**
 * Handle page change
 */
function handlePageChange(page: number): void {
	taskStore.setPage(page)
	taskStore.fetchTasks()
}

/**
 * Handle search
 */
function handleSearch(query: string): void {
	searchQuery = query
	taskStore.setFilters({ search: query || undefined })
	taskStore.fetchTasks()
}

/**
 * Handle view task
 */
function handleViewTask(task: (typeof taskStore.state.tasks)[0]): void {
	goto(`/tasks/${task.id}`)
}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="flex items-center justify-between mb-6">
		<h1 class="text-3xl font-bold text-gray-900">Tasks</h1>
		<button
			onclick={() => (showCreateModal = true)}
			class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			Create Task
		</button>
	</div>

	<div class="mb-6">
		<TaskSearch
			bind:searchQuery={searchQuery}
			onSearch={handleSearch}
			placeholder="Search tasks by title or description..."
		/>
	</div>

	<TaskFilters
		bind:filters={taskStore.state.filters}
		onFilterChange={handleFilterChange}
		onClearFilters={handleClearFilters}
	/>

	<TaskList
		tasks={taskStore.state.tasks}
		isLoading={taskStore.state.isLoading}
		onViewTask={handleViewTask}
		onEditTask={handleEditTask}
		onDeleteTask={handleDeleteTask}
	/>

	{#if taskStore.state.pagination.totalPages > 1}
	<div class="mt-6">
		<Pagination
			bind:currentPage={taskStore.state.pagination.page}
			totalPages={taskStore.state.pagination.totalPages}
			onPageChange={handlePageChange}
		/>
	</div>
	{/if}

	{#if showCreateModal}
		<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
			<div class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
				<div class="p-6">
					<div class="flex items-center justify-between mb-4">
						<h2 class="text-xl font-semibold text-gray-900">Create New Task</h2>
						<button
							onclick={() => (showCreateModal = false)}
							class="text-gray-400 hover:text-gray-600"
							title="Close"
						>
							<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>
					<TaskForm mode="create" onSubmit={handleCreateTask} onCancel={() => (showCreateModal = false)} />
				</div>
			</div>
		</div>
	{/if}
</div>