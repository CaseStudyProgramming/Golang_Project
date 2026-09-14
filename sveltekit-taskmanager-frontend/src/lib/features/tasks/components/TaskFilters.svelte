<script lang="ts">
import { browser } from '$app/environment'
import { categoryStore } from '$lib/features/categories'
import { tagStore } from '$lib/features/tags'
import TagInput from '$lib/features/tags/components/TagInput.svelte'

import type { TaskFilters, TaskPriority, TaskStatus } from '../types/task.types'

let {
	filters = $bindable({}),
	onClearFilters,
	onFilterChange,
}: {
	filters?: TaskFilters
	onClearFilters?: () => void
	onFilterChange?: (filters: TaskFilters) => void
} = $props()

let categories = $derived(browser ? categoryStore.state.categories : [])
let availableTags = $derived(browser ? tagStore.state.tags : [])

/**
 * Initialize categories and tags on mount
 */
$effect(() => {
	if (browser) {
		categoryStore.fetchCategories()
		tagStore.fetchTags()
	}
})

const statusOptions: { label: string; value: TaskStatus }[] = [
	{ label: 'To Do', value: 'todo' },
	{ label: 'In Progress', value: 'in_progress' },
	{ label: 'Completed', value: 'completed' },
	{ label: 'Cancelled', value: 'cancelled' },
]

const priorityOptions: { label: string; value: TaskPriority }[] = [
	{ label: 'Low', value: 'low' },
	{ label: 'Medium', value: 'medium' },
	{ label: 'High', value: 'high' },
	{ label: 'Urgent', value: 'urgent' },
]

/**
 * Handle filter change
 */
function handleFilterChange(key: keyof TaskFilters, event: Event): void {
	const target = event.target as HTMLInputElement
	const value = target.value
	const newFilters = { ...filters, [key]: value || undefined }
	filters = newFilters
	onFilterChange?.(newFilters)
}

/**
 * Handle tag filter change
 */
function handleTagFilterChange(tagIds: string[]): void {
	const newFilters = { ...filters, tags: tagIds.length > 0 ? tagIds : undefined }
	filters = newFilters
	onFilterChange?.(newFilters)
}

/**
 * Create reactive binding for tags filter
 */
let tagFilterIds = $state(filters.tags || [])

/**
 * Sync tag filter changes with filters
 */
$effect(() => {
	if (tagFilterIds.length > 0) {
		handleTagFilterChange(tagFilterIds)
	} else if (filters.tags) {
		handleTagFilterChange([])
	}
})

/**
 * Clear all filters
 */
function handleClearFilters(): void {
	filters = {}
	tagFilterIds = []
	onClearFilters?.()
}

/**
 * Check if any filters are active
 */
const hasActiveFilters = $derived(
	tagFilterIds.length > 0 ||
		Object.values(filters).some(
			(value) => value !== undefined && value !== '' && (!Array.isArray(value) || value.length > 0)
		)
)
</script>

<div class="bg-white rounded-lg shadow p-4 mb-4">
	<div class="flex items-center justify-between mb-4">
		<h3 class="text-lg font-medium text-gray-900">Filters</h3>
		{#if hasActiveFilters}
			<button
				onclick={() => handleClearFilters()}
				class="text-sm text-blue-600 hover:text-blue-700"
			>
				Clear all
			</button>
		{/if}
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
		<div>
			<label for="status" class="block text-sm font-medium text-gray-700 mb-1">Status</label>
			<select
				id="status"
				data-testid="status-filter"
				value={filters.status || ''}
				onchange={(e) => handleFilterChange('status', e)}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			>
				<option value="">All Statuses</option>
				{#each statusOptions as option (option.value)}
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
				{#each priorityOptions as option (option.value)}
					<option value={option.value}>{option.label}</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="category" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
			<select
				id="category"
				value={filters.categoryId || ''}
				onchange={(e) => handleFilterChange('categoryId', e)}
				class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			>
				<option value="">All Categories</option>
				{#each categories as category (category.id)}
					<option value={category.id}>{category.icon ? category.icon + ' ' : ''}{category.name}</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="tags" class="block text-sm font-medium text-gray-700 mb-1">Tags</label>
			<TagInput
				bind:selectedTags={tagFilterIds}
				bind:availableTags={availableTags}
				placeholder="Filter by tags..."
			/>
		</div>
	</div>
</div>