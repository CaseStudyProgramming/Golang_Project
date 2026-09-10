<script lang="ts">
	import { categoryStore } from '../stores/category.store';
	import type { Category } from '../types/category.types';
	import CategoryListSkeleton from './CategoryListSkeleton.svelte';

	let {
		categories = $bindable(categoryStore.state.categories),
		isLoading = $bindable(categoryStore.state.isLoading),
		onEditCategory,
		onDeleteCategory,
		onSelectCategory
	}: {
		categories?: Category[];
		isLoading?: boolean;
		onEditCategory?: (category: Category) => void;
		onDeleteCategory?: (category: Category) => void;
		onSelectCategory?: (category: Category) => void;
	} = $props();
</script>

<div class="space-y-3 sm:space-y-4">
	{#if isLoading}
		<CategoryListSkeleton />
	{:else if categories.length === 0}
		<div class="text-center py-8 sm:py-12 bg-white rounded-lg shadow">
			<svg class="mx-auto h-10 w-10 sm:h-12 sm:w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
			</svg>
			<h3 class="mt-2 text-sm font-medium text-gray-900">No categories</h3>
			<p class="mt-1 text-sm text-gray-500">Get started by creating a new category.</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
			{#each categories as category}
				<div
					class="bg-white rounded-lg shadow hover:shadow-md active:shadow-lg transition-shadow p-3 sm:p-4 border-l-4 cursor-pointer"
					style:border-left-color={category.color || '#3B82F6'}
					style:background-color={`${category.color || '#3B82F6'}10`}
					onclick={() => onSelectCategory?.(category)}
					onkeydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') {
							e.preventDefault();
							onSelectCategory?.(category);
						}
					}}
					role="button"
					tabindex="0"
				>
					<div class="flex items-start justify-between mb-2">
						<div class="flex items-center gap-2">
							{#if category.icon}
								<span class="text-xl sm:text-2xl">{category.icon}</span>
							{/if}
							<h3 class="text-base sm:text-lg font-medium text-gray-900">{category.name}</h3>
						</div>
						<div class="flex items-center gap-1">
							{#if onEditCategory}
								<button
									onclick={(e) => {
										e.stopPropagation();
										onEditCategory(category);
									}}
									class="p-2 sm:p-1.5 text-gray-400 hover:text-green-600 hover:bg-green-50 active:bg-green-100 rounded-lg transition-colors min-w-[36px] min-h-[36px] sm:min-w-0 sm:min-h-0"
									title="Edit category"
									aria-label="Edit category"
								>
									<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
									</svg>
								</button>
							{/if}
							{#if onDeleteCategory}
								<button
									onclick={(e) => {
										e.stopPropagation();
										onDeleteCategory(category);
									}}
									class="p-2 sm:p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 active:bg-red-100 rounded-lg transition-colors min-w-[36px] min-h-[36px] sm:min-w-0 sm:min-h-0"
									title="Delete category"
									aria-label="Delete category"
								>
									<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
									</svg>
								</button>
							{/if}
						</div>
					</div>
					{#if category.description}
						<p class="text-sm text-gray-600 line-clamp-2">{category.description}</p>
					{/if}
					<div class="mt-3 text-xs text-gray-500">
						Created {new Date(category.createdAt).toLocaleDateString()}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>