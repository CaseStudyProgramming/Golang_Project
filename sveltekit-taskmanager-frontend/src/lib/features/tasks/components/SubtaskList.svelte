<script lang="ts">
	import { SvelteSet } from 'svelte/reactivity';

	import type { Subtask } from '../types/task.types';

	let {
		onAddSubtask,
		onBulkComplete,
		onBulkDelete,
		onDeleteSubtask,
		onToggleSubtask,
		subtasks = $bindable([])
	}: {
		onAddSubtask?: (title: string) => Promise<void>;
		onBulkComplete?: (subtaskIds: string[]) => Promise<void>;
		onBulkDelete?: (subtaskIds: string[]) => Promise<void>;
		onDeleteSubtask?: (subtaskId: string) => Promise<void>;
		onToggleSubtask?: (subtaskId: string) => Promise<void>;
		subtasks?: Subtask[];
	} = $props();

	let newSubtaskTitle = $state('');
	let isAdding = $state(false);
	let selectedSubtasks = new SvelteSet<string>();

	/**
	 * Calculate progress percentage
	 */
	const progress = $derived(() => {
		if (subtasks.length === 0) return 0;
		const completed = subtasks.filter((s) => s.isCompleted).length;
		return Math.round((completed / subtasks.length) * 100);
	});

	/**
	 * Handle adding new subtask
	 */
	async function handleAddSubtask() {
		if (!newSubtaskTitle.trim() || !onAddSubtask) return;

		try {
			await onAddSubtask(newSubtaskTitle.trim());
			newSubtaskTitle = '';
			isAdding = false;
		} catch (error) {
			console.error('Failed to add subtask:', error);
		}
	}

	/**
	 * Handle bulk complete
	 */
	async function handleBulkComplete() {
		if (!onBulkComplete || selectedSubtasks.size === 0) return;

		try {
			await onBulkComplete(Array.from(selectedSubtasks));
			selectedSubtasks.clear();
		} catch (error) {
			console.error('Failed to bulk complete subtasks:', error);
		}
	}

	/**
	 * Handle bulk delete
	 */
	async function handleBulkDelete() {
		if (!onBulkDelete || selectedSubtasks.size === 0) return;

		try {
			await onBulkDelete(Array.from(selectedSubtasks));
			selectedSubtasks.clear();
		} catch (error) {
			console.error('Failed to bulk delete subtasks:', error);
		}
	}

	/**
	 * Handle deleting subtask
	 */
	async function handleDeleteSubtask(subtaskId: string) {
		if (!onDeleteSubtask) return;

		try {
			await onDeleteSubtask(subtaskId);
		} catch (error) {
			console.error('Failed to delete subtask:', error);
		}
	}

	/**
	 * Handle toggling subtask completion
	 */
	async function handleToggleSubtask(subtaskId: string) {
		if (!onToggleSubtask) return;

		try {
			await onToggleSubtask(subtaskId);
		} catch (error) {
			console.error('Failed to toggle subtask:', error);
		}
	}

	/**
	 * Handle subtask selection
	 */
	function toggleSelection(subtaskId: string) {
		if (selectedSubtasks.has(subtaskId)) {
			selectedSubtasks.delete(subtaskId);
		} else {
			selectedSubtasks.add(subtaskId);
		}
	}
</script>

<div class="bg-gray-50 rounded-lg p-4 mt-4">
	<div class="flex items-center justify-between mb-4">
		<h4 class="text-sm font-medium text-gray-900">Subtasks</h4>
		<div class="flex items-center gap-3">
			<div class="text-sm text-gray-600">
				{progress}% complete ({subtasks.filter((s) => s.isCompleted).length}/{subtasks.length})
			</div>
			<button
				onclick={() => (isAdding = true)}
				class="inline-flex items-center px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-700"
			>
				<svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
				</svg>
				Add Subtask
			</button>
		</div>
	</div>

	<!-- Progress bar -->
	<div class="w-full bg-gray-200 rounded-full h-2 mb-4">
		<div class="bg-blue-600 h-2 rounded-full transition-all duration-300" style="width: {progress}%"></div>
	</div>

	{#if isAdding}
		<div class="flex gap-2 mb-4">
			<input
				type="text"
				bind:value={newSubtaskTitle}
				placeholder="Enter subtask title..."
				class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				onkeydown={(e) => e.key === 'Enter' && handleAddSubtask()}
			/>
			<button
				onclick={handleAddSubtask}
				class="px-3 py-2 bg-blue-600 text-white rounded-md text-sm hover:bg-blue-700"
			>
				Add
			</button>
			<button
				onclick={() => (isAdding = false)}
				class="px-3 py-2 bg-gray-200 text-gray-700 rounded-md text-sm hover:bg-gray-300"
			>
				Cancel
			</button>
		</div>
	{/if}

	{#if subtasks.length === 0}
		<p class="text-sm text-gray-500 text-center py-4">No subtasks yet. Add one to get started.</p>
	{:else}
		<div class="space-y-2">
			{#each subtasks as subtask (subtask.id)}
				<div class="flex items-center gap-3 p-2 bg-white rounded-md hover:bg-gray-50 transition-colors">
					<input
						type="checkbox"
						checked={subtask.isCompleted}
						onchange={() => handleToggleSubtask(subtask.id)}
						class="h-4 w-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
					/>
					<input
						type="checkbox"
						checked={selectedSubtasks.has(subtask.id)}
						onchange={() => toggleSelection(subtask.id)}
						class="h-4 w-4 text-gray-600 rounded border-gray-300 focus:ring-gray-500"
					/>
					<span class="flex-1 text-sm {subtask.isCompleted ? 'line-through text-gray-400' : 'text-gray-900'}">
						{subtask.title}
					</span>
					<button
						onclick={() => handleDeleteSubtask(subtask.id)}
						class="p-1 text-gray-400 hover:text-red-600 rounded hover:bg-red-50 transition-colors"
						title="Delete subtask"
					>
						<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>
			{/each}
		</div>

		{#if selectedSubtasks.size > 0}
			<div class="flex gap-2 mt-4 pt-4 border-t border-gray-200">
				<button
					onclick={handleBulkComplete}
					class="px-3 py-1.5 bg-green-600 text-white rounded-md text-sm hover:bg-green-700"
				>
					Complete Selected ({selectedSubtasks.size})
				</button>
				<button
					onclick={handleBulkDelete}
					class="px-3 py-1.5 bg-red-600 text-white rounded-md text-sm hover:bg-red-700"
				>
					Delete Selected ({selectedSubtasks.size})
				</button>
				<button
					onclick={() => selectedSubtasks.clear()}
					class="px-3 py-1.5 bg-gray-200 text-gray-700 rounded-md text-sm hover:bg-gray-300"
				>
					Clear Selection
				</button>
			</div>
		{/if}
	{/if}
</div>