<script lang="ts">
	import type { TaskFilters, TaskStatus, TaskPriority } from '../types/task.types';

	let {
		filters = $bindable({}),
		onFilterChange,
		onClearFilters
	}: {
		filters?: TaskFilters;
		onFilterChange?: (filters: TaskFilters) => void;
		onClearFilters?: () => void;
	} = $props();

	const statusOptions: { value: TaskStatus; label: string }[] = [
		{ value: 'todo', label: 'To Do' },
		{ value: 'in_progress', label: 'In Progress' },
		{ value: 'completed', label: 'Completed' },
		{ value: 'cancelled', label: 'Cancelled' }
	];

	const priorityOptions: { value: TaskPriority; label: string }[] = [
		{ value: 'low', label: 'Low' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'high', label: 'High' },
		{ value: 'urgent', label: 'Urgent' }
	];

	/**
	 * Handle filter change
	 */
	function handleFilterChange(key: keyof TaskFilters, event: Event): void {
		const target = event.target as HTMLInputElement;
		const value = target.value;
		const newFilters = { ...filters, [key]: value || undefined };
		filters = newFilters;
		onFilterChange?.(newFilters);
	}

	/**
	 * Clear all filters
	 */
	function handleClearFilters(): void {
		filters = {};
		onClearFilters?.();
	}

	/**
	 * Check if any filters are active
	 */
	const hasActiveFilters = $derived(
		Object.values(filters).some(
			(value) => value !== undefined && value !== '' && (!Array.isArray(value) || value.length > 0)
		)
	);
</script>

<div class="bg-white rounded-lg shadow p-4 mb-4">
	<div class="flex items-center justify-between mb-4">
		<h3 class="text-lg font-medium text-gray-900">Filters</h3>
		{#if hasActiveFilters}
			<button
				onclick={handleClearFilters}
				class="text-sm text-blue-600 hover:text-blue-700"
			>
				Clear all
			</button>
		{/if}
	</div>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
		<div>
			<label for="status" class="block text-sm font-medium text-gray-700 mb-1">Status</label>
			<select
				id="status"
				value={filters.status || ''}
				onchange={(e) => handleFilterChange('status', e)}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			>
				<option value="">All Statuses</option>
				{#each statusOptions as option}
					<option value={option.value}>{option.label}</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="priority" class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
			<select
				id="priority"
				value={filters.priority || ''}
				onchange={(e) => handleFilterChange('priority', e)}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			>
				<option value="">All Priorities</option>
				{#each priorityOptions as option}
					<option value={option.value}>{option.label}</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="category" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
			<input
				id="category"
				type="text"
				value={filters.categoryId || ''}
				oninput={(e) => handleFilterChange('categoryId', e)}
				placeholder="Filter by category ID"
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			/>
		</div>
	</div>
</div>