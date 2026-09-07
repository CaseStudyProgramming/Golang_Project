<script lang="ts">
	import { categoryStore } from '../stores/category.store';
	import type { Category } from '../types/category.types';
	import type { CreateCategoryPayload, UpdateCategoryPayload } from '../schemas/category.schemas';

	let {
		initialData,
		onSubmit,
		onCancel
	}: {
		initialData?: Category;
		onSubmit?: (data: CreateCategoryPayload | UpdateCategoryPayload) => Promise<void>;
		onCancel?: () => void;
	} = $props();

	let name = $state(initialData?.name || '');
	let description = $state(initialData?.description || '');
	let color = $state(initialData?.color || '#3B82F6');
	let icon = $state(initialData?.icon || '');
	let isSubmitting = $state(false);
	let error = $state('');

	const predefinedColors = [
		'#3B82F6', // blue
		'#10B981', // green
		'#F59E0B', // yellow
		'#EF4444', // red
		'#8B5CF6', // purple
		'#EC4899', // pink
		'#6366F1', // indigo
		'#14B8A6'  // teal
	];

	/**
	 * Handle form submission
	 */
	async function handleSubmit() {
		if (!name.trim()) {
			error = 'Name is required';
			return;
		}

		isSubmitting = true;
		error = '';

		try {
			const payload = {
				name: name.trim(),
				description: description.trim() || undefined,
				color: color || undefined,
				icon: icon.trim() || undefined
			};

			if (onSubmit) {
				await onSubmit(payload);
			} else {
				if (initialData) {
					await categoryStore.updateCategory(initialData.id, payload);
				} else {
					await categoryStore.createCategory(payload);
				}
			}

			// Reset form if creating new category
			if (!initialData) {
				name = '';
				description = '';
				color = '#3B82F6';
				icon = '';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save category';
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Handle cancel
	 */
	function handleCancel() {
		if (onCancel) {
			onCancel();
		}
	}
</script>

<div class="bg-white rounded-lg shadow p-6">
	<h2 class="text-lg font-medium text-gray-900 mb-4">
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
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
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
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
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
					{#each predefinedColors as presetColor}
						<button
							type="button"
							onclick={() => (color = presetColor)}
							class="w-8 h-8 rounded-full border-2 {color === presetColor ? 'border-gray-900' : 'border-gray-300'} hover:scale-110 transition-transform"
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
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				disabled={isSubmitting}
			/>
		</div>

		<div class="flex gap-3 pt-4">
			<button
				type="button"
				onclick={handleSubmit}
				disabled={isSubmitting}
				class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
			>
				{isSubmitting ? 'Saving...' : initialData ? 'Update' : 'Create'}
			</button>
			{#if onCancel}
				<button
					type="button"
					onclick={handleCancel}
					disabled={isSubmitting}
					class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed transition-colors"
				>
					Cancel
				</button>
			{/if}
		</div>
	</form>
</div>