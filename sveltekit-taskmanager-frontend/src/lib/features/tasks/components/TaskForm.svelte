<script lang="ts">
	import { taskStore } from '../stores/task.store';
	import { createTaskSchema, type CreateTaskPayload, type UpdateTaskPayload } from '../schemas/task.schemas';
	import type { TaskPriority } from '../types/task.types';
	import { categoryStore } from '$lib/features/categories';
	import { tagStore } from '$lib/features/tags';
	import { TagInput } from '$lib/features/tags';
	import type { Category } from '$lib/features/categories';

	let {
		mode = 'create',
		initialData,
		onSubmit,
		onCancel
	}: {
		mode?: 'create' | 'edit';
		initialData?: Partial<CreateTaskPayload>;
		onSubmit?: (data: CreateTaskPayload) => Promise<void>;
		onCancel?: () => void;
	} = $props();

	let title = $state('');
	let description = $state('');
	let priority = $state<TaskPriority>('medium');
	let dueDate = $state('');
	let categoryId = $state('');
	let tags = $state<string[]>([]);
	let errors = $state<Record<string, string>>({});
	let isSubmitting = $state(false);

	/**
	 * Get derived values from initialData
	 */
	const initialTitle = $derived(initialData?.title || '');
	const initialDescription = $derived(initialData?.description || '');
	const initialPriority = $derived<TaskPriority>(initialData?.priority || 'medium');
	const initialDueDate = $derived(initialData?.dueDate || '');
	const initialCategoryId = $derived(initialData?.categoryId || '');
	const initialTags = $derived<string[]>(initialData?.tags || []);

	/**
	 * Sync form state when initialData changes
	 */
	$effect(() => {
		title = initialTitle;
		description = initialDescription;
		priority = initialPriority;
		dueDate = initialDueDate;
		categoryId = initialCategoryId;
		tags = initialTags;
	});

	let categories = $derived(categoryStore.state.categories);
	let availableTags = $derived(tagStore.state.tags);

	/**
	 * Initialize categories and tags on mount
	 */
	$effect(() => {
		categoryStore.fetchCategories();
		tagStore.fetchTags();
	});

	/**
	 * Handle creating new tag
	 */
	async function handleCreateTag(name: string): Promise<void> {
		try {
			const newTag = await tagStore.createTag({ name });
			tags = [...tags, newTag.id];
		} catch (error) {
			console.error('Failed to create tag:', error);
		}
	}

	/**
	 * Validate form
	 */
	function validateForm(): boolean {
		try {
			createTaskSchema.parse({
				title,
				description,
				priority,
				dueDate: dueDate || undefined,
				categoryId: categoryId || undefined,
				tags
			});
			errors = {};
			return true;
		} catch (error: any) {
			if (error.name === 'ZodError') {
				const newErrors: Record<string, string> = {};
				error.errors.forEach((err: any) => {
					const path = err.path.join('.');
					newErrors[path] = err.message;
				});
				errors = newErrors;
			}
			return false;
		}
	}

	/**
	 * Handle form submission
	 */
	async function handleSubmit(): Promise<void> {
		if (!validateForm()) {
			return;
		}

		isSubmitting = true;
		try {
			const payload: CreateTaskPayload = {
				title,
				description,
				priority,
				dueDate: dueDate || undefined,
				categoryId: categoryId || undefined,
				tags
			};

			if (onSubmit) {
				await onSubmit(payload);
			} else {
				await taskStore.createTask(payload);
			}

			// Reset form on success
			if (mode === 'create') {
				title = '';
				description = '';
				priority = 'medium';
				dueDate = '';
				categoryId = '';
				tags = [];
			}
		} catch (error) {
			console.error('Failed to submit task:', error);
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Handle cancel
	 */
	function handleCancel(): void {
		onCancel?.();
	}
</script>

<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-6">
	<div>
		<label for="title" class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
		<input
			id="title"
			type="text"
			bind:value={title}
			class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			placeholder="Enter task title"
			required
		/>
		{#if errors.title}
			<p class="mt-1 text-sm text-red-600">{errors.title}</p>
		{/if}
	</div>

	<div>
		<label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
		<textarea
			id="description"
			bind:value={description}
			rows="3"
			class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
			placeholder="Enter task description"
		></textarea>
		{#if errors.description}
			<p class="mt-1 text-sm text-red-600">{errors.description}</p>
		{/if}
	</div>

	<div>
		<label for="priority" class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
		<select
			id="priority"
			bind:value={priority}
			class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
		>
			<option value="low">Low</option>
			<option value="medium">Medium</option>
			<option value="high">High</option>
			<option value="urgent">Urgent</option>
		</select>
	</div>

	<div>
		<label for="dueDate" class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
		<input
			id="dueDate"
			type="date"
			bind:value={dueDate}
			class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
		/>
		{#if errors.dueDate}
			<p class="mt-1 text-sm text-red-600">{errors.dueDate}</p>
		{/if}
	</div>

	<div>
		<label for="categoryId" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
		<select
			id="categoryId"
			bind:value={categoryId}
			class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
		>
			<option value="">No category</option>
			{#each categories as category}
				<option value={category.id}>{category.icon ? category.icon + ' ' : ''}{category.name}</option>
			{/each}
		</select>
	</div>

	<div>
		<label for="tags" class="block text-sm font-medium text-gray-700 mb-1">Tags</label>
		<TagInput
			bind:selectedTags={tags}
			bind:availableTags={availableTags}
			onCreateTag={handleCreateTag}
			placeholder="Add tags..."
		/>
	</div>

	<div class="flex justify-end gap-3">
		{#if onCancel}
			<button
				type="button"
				onclick={handleCancel}
				disabled={isSubmitting}
				class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 transition-colors"
			>
				Cancel
			</button>
		{/if}
		<button
			type="submit"
			disabled={isSubmitting}
			class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
		>
			{isSubmitting ? 'Saving...' : mode === 'create' ? 'Create Task' : 'Update Task'}
		</button>
	</div>
</form>