<script lang="ts">
import { toastStore } from '$lib/shared/stores'

import type { CreateCategoryPayload, UpdateCategoryPayload } from '../schemas/category.schemas'
import { categoryStore } from '../stores/category.store.svelte.ts'
import type { Category } from '../types/category.types'

let {
	initialData,
	onCancel,
	onSubmit,
}: {
	initialData?: Category
	onCancel?: () => void
	onSubmit?: (data: CreateCategoryPayload | UpdateCategoryPayload) => Promise<void>
} = $props()

let name = $state('')
let description = $state('')
let color = $state('#3B82F6')
let icon = $state('')
let isSubmitting = $state(false)
let error = $state('')

/**
 * Sync form state when initialData reference changes
 */
$effect(() => {
	if (initialData) {
		name = initialData.name
		description = initialData.description || ''
		color = initialData.color || '#3B82F6'
		icon = initialData.icon || ''
	} else {
		name = ''
		description = ''
		color = '#3B82F6'
		icon = ''
	}
})

const predefinedColors = [
	'#3B82F6', // blue
	'#10B981', // green
	'#F59E0B', // yellow
	'#EF4444', // red
	'#8B5CF6', // purple
	'#EC4899', // pink
	'#6366F1', // indigo
	'#14B8A6', // teal
]

/**
 * Handle cancel
 */
function handleCancel() {
	if (onCancel) {
		onCancel()
	}
}

/**
 * Handle form submission
 */
async function handleSubmit() {
	if (!name.trim()) {
		error = 'Name is required'
		return
	}

	isSubmitting = true
	error = ''

	try {
		const payload = {
			color: color || undefined,
			description: description.trim() || undefined,
			icon: icon.trim() || undefined,
			name: name.trim(),
		}

		if (onSubmit) {
			await onSubmit(payload)
		} else {
			if (initialData) {
				await categoryStore.updateCategory(initialData.id, payload)
			} else {
				await categoryStore.createCategory(payload)
			}
		}

		// Show success toast
		toastStore.success(
			initialData ? 'Category updated' : 'Category created',
			initialData
				? 'Your category has been updated successfully'
				: 'Your category has been created successfully'
		)

		// Reset form if creating new category
		if (!initialData) {
			name = ''
			description = ''
			color = '#3B82F6'
			icon = ''
		}
	} catch (err) {
		error = err instanceof Error ? err.message : 'Failed to save category'
		toastStore.error(initialData ? 'Failed to update category' : 'Failed to create category', error)
	} finally {
		isSubmitting = false
	}
}
</script>

<div class="bg-white rounded-lg shadow p-4 sm:p-6">
	<h2 class="text-base sm:text-lg font-medium text-gray-900 mb-4">
		{initialData ? 'Edit Category' : 'Create Category'}
	</h2>

	{#if error}
		<div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
			<p class="text-sm text-red-600">{error}</p>
		</div>
	{/if}

	<form onsubmit={(e) => e.preventDefault()} class="space-y-4">
		<div>
			<label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name *</label>
			<input
				id="name"
				type="text"
				bind:value={name}
				placeholder="Category name"
				class="w-full px-3 py-2 sm:px-4 sm:py-2.5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm sm:text-base"
				disabled={isSubmitting}
			/>
		</div>

		<div>
			<label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
			<textarea
				id="description"
				bind:value={description}
				placeholder="Category description (optional)"
				rows="3"
				class="w-full px-3 py-2 sm:px-4 sm:py-2.5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm sm:text-base"
				disabled={isSubmitting}
			></textarea>
		</div>

		<div>
			<label for="color" class="block text-sm font-medium text-gray-700 mb-1">Color</label>
			<div class="flex items-center gap-3">
				<input
					id="color"
					type="color"
					bind:value={color}
					class="h-10 w-16 border border-gray-300 rounded cursor-pointer"
					disabled={isSubmitting}
				/>
				<div class="flex gap-2 flex-wrap">
					{#each predefinedColors as presetColor (presetColor)}
						<button
							type="button"
							onclick={() => (color = presetColor)}
							class="w-8 h-8 rounded-full border-2 {color === presetColor ? 'border-gray-900' : 'border-gray-300'} hover:scale-110 transition-transform min-w-[44px] min-h-[44px]"
							style="background-color: {presetColor}"
							disabled={isSubmitting}
							aria-label={`Select color ${presetColor}`}
						></button>
					{/each}
				</div>
			</div>
		</div>

		<div>
			<label for="icon" class="block text-sm font-medium text-gray-700 mb-1">Icon (emoji or text)</label>
			<input
				id="icon"
				type="text"
				bind:value={icon}
				placeholder="📁 or 🎯"
				maxlength="50"
				class="w-full px-3 py-2 sm:px-4 sm:py-2.5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm sm:text-base"
				disabled={isSubmitting}
			/>
		</div>

		<div class="flex flex-col sm:flex-row gap-3 pt-4">
			<button
				type="button"
				onclick={() => handleSubmit()}
				disabled={isSubmitting}
				class="flex-1 px-4 py-3 sm:px-4 sm:py-2.5 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors min-h-[44px]"
			>
				{isSubmitting ? 'Saving...' : initialData ? 'Update' : 'Create'}
			</button>
			{#if onCancel}
				<button
					type="button"
					onclick={() => handleCancel()}
					disabled={isSubmitting}
					class="px-4 py-3 sm:px-4 sm:py-2.5 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed transition-colors min-h-[44px]"
				>
					Cancel
				</button>
			{/if}
		</div>
	</form>
</div>