<script lang="ts">
	import type { Category } from '$lib/features/categories';

	import { categoryStore } from '$lib/features/categories';
	import { CategoryForm, CategoryList } from '$lib/features/categories';
	import { onMount } from 'svelte';

	let categories = $derived(categoryStore.state.categories);
	let isLoading = $derived(categoryStore.state.isLoading);
	let error = $derived(categoryStore.state.error);

	let showForm = $state(false);
	let editingCategory = $state<Category | null>(null);

	/**
	 * Load categories on mount
	 */
	onMount(() => {
		categoryStore.fetchCategories();
	});

	/**
	 * Handle category creation
	 */
	async function handleCreateCategory(data: unknown): Promise<void> {
		try {
			await categoryStore.createCategory(data);
			showForm = false;
		} catch (error) {
			console.error('Failed to create category:', error);
		}
	}

	/**
	 * Handle create button click
	 */
	function handleCreateClick(): void {
		editingCategory = null;
		showForm = true;
	}

	/**
	 * Handle category deletion
	 */
	async function handleDeleteCategory(category: Category): Promise<void> {
		if (!confirm(`Are you sure you want to delete "${category.name}"?`)) return;

		try {
			await categoryStore.deleteCategory(category.id);
		} catch (error) {
			console.error('Failed to delete category:', error);
		}
	}

	/**
	 * Handle edit button click
	 */
	function handleEditCategory(category: Category): void {
		editingCategory = category;
		showForm = true;
	}

	/**
	 * Handle form cancel
	 */
	function handleFormCancel(): void {
		showForm = false;
		editingCategory = null;
	}

	/**
	 * Handle category update
	 */
	async function handleUpdateCategory(data: unknown): Promise<void> {
		if (!editingCategory) return;

		try {
			await categoryStore.updateCategory(editingCategory.id, data);
			editingCategory = null;
			showForm = false;
		} catch (error) {
			console.error('Failed to update category:', error);
		}
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Categories</h1>
			<p class="text-gray-600 mt-1">Manage your task categories</p>
		</div>
		<button
			onclick={handleCreateClick}
			class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			<svg class="h-5 w-5 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			Create Category
		</button>
	</div>

	{#if error}
		<div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
			<p class="text-red-600">{error}</p>
		</div>
	{/if}

	{#if showForm}
		<div class="mb-6">
			<CategoryForm
				initialData={editingCategory || undefined}
				onSubmit={editingCategory ? handleUpdateCategory : handleCreateCategory}
				onCancel={handleFormCancel}
			/>
		</div>
	{/if}

	<CategoryList
		bind:categories
		bind:isLoading
		onEditCategory={handleEditCategory}
		onDeleteCategory={handleDeleteCategory}
	/>
</div>