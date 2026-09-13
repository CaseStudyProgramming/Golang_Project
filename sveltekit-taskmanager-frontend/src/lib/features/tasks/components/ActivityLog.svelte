<script lang="ts">
import { activityStore } from '../stores/activity.store.svelte.ts'
import type { Activity, ActivityFilters, ActivityType } from '../types/task.types'

let {
	activities = $bindable(activityStore.state.activities),
	isLoading = $bindable(activityStore.state.isLoading),
	taskId,
}: {
	activities?: Activity[]
	isLoading?: boolean
	taskId?: string
} = $props()

let filters = $state<ActivityFilters>({})
let selectedType = $state<ActivityType | undefined>(undefined)

const activityTypes: { label: string; value: ActivityType }[] = [
	{ label: 'Task Created', value: 'task_created' },
	{ label: 'Task Updated', value: 'task_updated' },
	{ label: 'Task Deleted', value: 'task_deleted' },
	{ label: 'Task Completed', value: 'task_completed' },
	{ label: 'Subtask Added', value: 'subtask_added' },
	{ label: 'Subtask Completed', value: 'subtask_completed' },
	{ label: 'Subtask Deleted', value: 'subtask_deleted' },
	{ label: 'Category Assigned', value: 'category_assigned' },
	{ label: 'Tag Added', value: 'tag_added' },
	{ label: 'Tag Removed', value: 'tag_removed' },
	{ label: 'Status Changed', value: 'status_changed' },
	{ label: 'Priority Changed', value: 'priority_changed' },
]

/**
 * Get filtered activities
 */
let filteredActivities = $derived(() => {
	let filtered = activities

	if (taskId) {
		filtered = filtered.filter((a) => a.taskId === taskId)
	}

	if (selectedType) {
		filtered = filtered.filter((a) => a.type === selectedType)
	}

	if (filters.dateFrom) {
		filtered = filtered.filter((a) => new Date(a.createdAt) >= new Date(filters.dateFrom!))
	}

	if (filters.dateTo) {
		filtered = filtered.filter((a) => new Date(a.createdAt) <= new Date(filters.dateTo!))
	}

	return filtered
})

/**
 * Clear all filters
 */
function clearFilters(): void {
	selectedType = undefined
	filters = {}
}

/**
 * Format activity timestamp
 */
function formatTimestamp(timestamp: string): string {
	const date = new Date(timestamp)
	const now = new Date()
	const diffMs = now.getTime() - date.getTime()
	const diffMins = Math.floor(diffMs / 60000)
	const diffHours = Math.floor(diffMs / 3600000)
	const diffDays = Math.floor(diffMs / 86400000)

	if (diffMins < 1) return 'Just now'
	if (diffMins < 60) return `${diffMins}m ago`
	if (diffHours < 24) return `${diffHours}h ago`
	if (diffDays < 7) return `${diffDays}d ago`
	return date.toLocaleDateString()
}

/**
 * Get activity icon based on type
 */
function getActivityIcon(type: ActivityType): string {
	const icons: Record<ActivityType, string> = {
		category_assigned: '📁',
		priority_changed: '⚡',
		status_changed: '🔄',
		subtask_added: '➕',
		subtask_completed: '☑️',
		subtask_deleted: '❌',
		tag_added: '🏷️',
		tag_removed: '🏷️',
		task_completed: '✅',
		task_created: '✨',
		task_deleted: '🗑️',
		task_updated: '✏️',
	}
	return icons[type] || '📝'
}

/**
 * Handle type filter change
 */
function handleTypeFilterChange(event: Event): void {
	const target = event.target as HTMLSelectElement
	const value = target.value as ActivityType | undefined
	selectedType = value
	filters = { ...filters, type: value }
}
</script>

<div class="bg-white rounded-lg shadow p-4">
	<div class="flex items-center justify-between mb-4">
		<h3 class="text-lg font-medium text-gray-900">Activity Log</h3>
		<div class="flex items-center gap-2">
			{#if selectedType || filters.dateFrom || filters.dateTo}
				<button
					onclick={() => clearFilters()}
					class="text-sm text-blue-600 hover:text-blue-700"
				>
					Clear filters
				</button>
			{/if}
		</div>
	</div>

	<!-- Filters -->
	<div class="flex gap-2 mb-4">
		<select
			bind:value={selectedType}
			onchange={(e) => handleTypeFilterChange(e)}
			class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
		>
			<option value={undefined}>All Types</option>
			{#each activityTypes as type (type.value)}
				<option value={type.value}>{type.label}</option>
			{/each}
		</select>
		<input
			type="date"
			bind:value={filters.dateFrom}
			class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>
		<input
			type="date"
			bind:value={filters.dateTo}
			class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>
	</div>

	{#if isLoading}
		<div class="text-center py-8">
			<div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
			<p class="mt-2 text-sm text-gray-500">Loading activities...</p>
		</div>
	{:else if filteredActivities().length === 0}
		<div class="text-center py-8 text-gray-500">
			<p>No activity found</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each filteredActivities() as activity (activity.id)}
				<div class="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
					<div class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-white rounded-full shadow">
						<span class="text-lg">{getActivityIcon(activity.type)}</span>
					</div>
					<div class="flex-1 min-w-0">
						<p class="text-sm text-gray-900">{activity.description}</p>
						<div class="flex items-center gap-2 mt-1">
							<span class="text-xs text-gray-500">{activity.userName || 'User'}</span>
							<span class="text-xs text-gray-400">•</span>
							<span class="text-xs text-gray-500">{formatTimestamp(activity.createdAt)}</span>
						</div>
						{#if activity.changes && Object.keys(activity.changes).length > 0}
							<div class="mt-2 text-xs">
								{#each Object.entries(activity.changes) as [field, change] (field)}
									<div class="flex items-center gap-2">
										<span class="text-gray-600">{field}:</span>
										<span class="text-red-600 line-through">{String((change as { old: unknown }).old)}</span>
										<span class="text-gray-400">→</span>
										<span class="text-green-600">{String((change as { new: unknown }).new)}</span>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>