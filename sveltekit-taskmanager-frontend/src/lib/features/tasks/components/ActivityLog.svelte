<script lang="ts">
	import { activityStore } from '../stores/activity.store';
	import type { Activity, ActivityFilters, ActivityType } from '../types/task.types';

	let {
		taskId,
		activities = $bindable(activityStore.state.activities),
		isLoading = $bindable(activityStore.state.isLoading),
		showTimeline = true
	}: {
		taskId?: string;
		activities?: Activity[];
		isLoading?: boolean;
		showTimeline?: boolean;
	} = $props();

	let filters = $state<ActivityFilters>({});
	let selectedType = $state<ActivityType | undefined>(undefined);

	const activityTypes: { value: ActivityType; label: string }[] = [
		{ value: 'task_created', label: 'Task Created' },
		{ value: 'task_updated', label: 'Task Updated' },
		{ value: 'task_deleted', label: 'Task Deleted' },
		{ value: 'task_completed', label: 'Task Completed' },
		{ value: 'subtask_added', label: 'Subtask Added' },
		{ value: 'subtask_completed', label: 'Subtask Completed' },
		{ value: 'subtask_deleted', label: 'Subtask Deleted' },
		{ value: 'category_assigned', label: 'Category Assigned' },
		{ value: 'tag_added', label: 'Tag Added' },
		{ value: 'tag_removed', label: 'Tag Removed' },
		{ value: 'status_changed', label: 'Status Changed' },
		{ value: 'priority_changed', label: 'Priority Changed' }
	];

	/**
	 * Get filtered activities
	 */
	const filteredActivities: Activity[] = $derived(() => {
		let filtered = activities;

		if (taskId) {
			filtered = filtered.filter((a) => a.taskId === taskId);
		}

		if (selectedType) {
			filtered = filtered.filter((a) => a.type === selectedType);
		}

		if (filters.dateFrom) {
			filtered = filtered.filter((a) => new Date(a.createdAt) >= new Date(filters.dateFrom!));
		}

		if (filters.dateTo) {
			filtered = filtered.filter((a) => new Date(a.createdAt) <= new Date(filters.dateTo!));
		}

		return filtered;
	});

	/**
	 * Get activity icon based on type
	 */
	function getActivityIcon(type: ActivityType): string {
		const icons: Record<ActivityType, string> = {
			task_created: '✨',
			task_updated: '✏️',
			task_deleted: '🗑️',
			task_completed: '✅',
			subtask_added: '➕',
			subtask_completed: '☑️',
			subtask_deleted: '❌',
			category_assigned: '📁',
			tag_added: '🏷️',
			tag_removed: '🏷️',
			status_changed: '🔄',
			priority_changed: '⚡'
		};
		return icons[type] || '📝';
	}

	/**
	 * Format activity timestamp
	 */
	function formatTimestamp(timestamp: string): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		const diffHours = Math.floor(diffMs / 3600000);
		const diffDays = Math.floor(diffMs / 86400000);

		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		if (diffDays < 7) return `${diffDays}d ago`;
		return date.toLocaleDateString();
	}

	/**
	 * Handle type filter change
	 */
	function handleTypeFilterChange(event: Event): void {
		const target = event.target as HTMLSelectElement;
		const value = target.value as ActivityType | undefined;
		selectedType = value;
		filters = { ...filters, type: value };
	}

	/**
	 * Clear all filters
	 */
	function clearFilters(): void {
		selectedType = undefined;
		filters = {};
	}

	/**
	 * Load activities
	 */
	async function loadActivities(): Promise<void> {
		if (taskId) {
			await activityStore.fetchActivities({ taskId });
		} else {
			await activityStore.fetchActivities(filters);
		}
	}
</script>

<div class="bg-white rounded-lg shadow p-4">
	<div class="flex items-center justify-between mb-4">
		<h3 class="text-lg font-medium text-gray-900">Activity Log</h3>
		<div class="flex items-center gap-2">
			{#if selectedType || filters.dateFrom || filters.dateTo}
				<button
					onclick={clearFilters}
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
			onchange={handleTypeFilterChange}
			class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
		>
			<option value={undefined}>All Types</option>
			{#each activityTypes as type}
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
	{:else if filteredActivities.length === 0}
		<div class="text-center py-8 text-gray-500">
			<p>No activity found</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each filteredActivities as activity}
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
								{#each Object.entries(activity.changes) as [field, change]}
									<div class="flex items-center gap-2">
										<span class="text-gray-600">{field}:</span>
										<span class="text-red-600 line-through">{String(change.old)}</span>
										<span class="text-gray-400">→</span>
										<span class="text-green-600">{String(change.new)}</span>
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