<script lang="ts">
import { onMount } from 'svelte'
import type { Tag } from '$lib/features/tags'
import { tagStore } from '$lib/features/tags'

let tags = $derived(tagStore.state.tags)
let isLoading = $derived(tagStore.state.isLoading)
let error = $derived(tagStore.state.error)

let showForm = $state(false)
let editingTag = $state<null | Tag>(null)
let name = $state('')
let color = $state('#3B82F6')
let isSubmitting = $state(false)

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
 * Load tags on mount
 */
onMount(() => {
	tagStore.fetchTags()
})

/**
 * Handle create button click
 */
function handleCreateClick(): void {
	editingTag = null
	name = ''
	color = '#3B82F6'
	showForm = true
}

/**
 * Handle tag creation
 */
async function handleCreateTag(): Promise<void> {
	if (!name.trim()) return

	isSubmitting = true
	try {
		await tagStore.createTag({ color, name: name.trim() })
		name = ''
		color = '#3B82F6'
		showForm = false
	} catch (error) {
		console.error('Failed to create tag:', error)
	} finally {
		isSubmitting = false
	}
}

/**
 * Handle tag deletion
 */
async function handleDeleteTag(tag: Tag): Promise<void> {
	if (!confirm(`Are you sure you want to delete "${tag.name}"?`)) return

	try {
		await tagStore.deleteTag(tag.id)
	} catch (error) {
		console.error('Failed to delete tag:', error)
	}
}

/**
 * Handle edit button click
 */
function handleEditTag(tag: Tag): void {
	editingTag = tag
	name = tag.name
	color = tag.color || '#3B82F6'
	showForm = true
}

/**
 * Handle form cancel
 */
function handleFormCancel(): void {
	showForm = false
	editingTag = null
	name = ''
	color = '#3B82F6'
}

/**
 * Handle tag update
 */
async function handleUpdateTag(): Promise<void> {
	if (!editingTag || !name.trim()) return

	isSubmitting = true
	try {
		await tagStore.updateTag(editingTag.id, { color, name: name.trim() })
		editingTag = null
		showForm = false
		name = ''
		color = '#3B82F6'
	} catch (error) {
		console.error('Failed to update tag:', error)
	} finally {
		isSubmitting = false
	}
}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Tags</h1>
			<p class="text-gray-600 mt-1">Manage your task tags</p>
		</div>
		<button
			onclick={() => handleCreateClick()}
			class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
		>
			<svg class="h-5 w-5 inline mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			Create Tag
		</button>
	</div>

	{#if error}
		<div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
			<p class="text-red-600">{error}</p>
		</div>
	{/if}

	{#if showForm}
		<div class="mb-6 bg-white rounded-lg shadow p-6">
			<h2 class="text-lg font-medium text-gray-900 mb-4">
				{editingTag ? 'Edit Tag' : 'Create Tag'}
			</h2>
			<div class="space-y-4">
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name *</label>
					<input
						id="name"
						type="text"
						bind:value={name}
						placeholder="Tag name"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						disabled={isSubmitting}
					/>
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
									class="w-8 h-8 rounded-full border-2 {color === presetColor ? 'border-gray-900' : 'border-gray-300'} hover:scale-110 transition-transform"
									style="background-color: {presetColor}"
									disabled={isSubmitting}
									aria-label={`Select color ${presetColor}`}
								></button>
							{/each}
						</div>
					</div>
				</div>

				<div class="flex gap-3 pt-4">
					<button
						type="button"
						onclick={() => editingTag ? handleUpdateTag() : handleCreateTag()}
						disabled={isSubmitting}
						class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
					>
						{isSubmitting ? 'Saving...' : editingTag ? 'Update' : 'Create'}
					</button>
					<button
						type="button"
						onclick={() => handleFormCancel()}
						disabled={isSubmitting}
						class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed transition-colors"
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	{/if}

	{#if isLoading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-500">Loading tags...</p>
		</div>
	{:else if tags.length === 0}
		<div class="text-center py-12 bg-white rounded-lg shadow">
			<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
			</svg>
			<h3 class="mt-2 text-sm font-medium text-gray-900">No tags</h3>
			<p class="mt-1 text-sm text-gray-500">Get started by creating a new tag.</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each tags as tag (tag.id)}
				<div class="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-4 border-l-4" style:border-left-color={tag.color || '#3B82F6'}>
					<div class="flex items-start justify-between mb-2">
						<div class="flex items-center gap-2">
							<div
								class="w-4 h-4 rounded-full"
								style="background-color: {tag.color || '#3B82F6'}"
							></div>
							<h3 class="text-lg font-medium text-gray-900">{tag.name}</h3>
						</div>
						<div class="flex items-center gap-1">
							<button
								onclick={() => handleEditTag(tag)}
								class="p-1.5 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
								title="Edit tag"
							>
								<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
								</svg>
							</button>
							<button
								onclick={() => handleDeleteTag(tag)}
								class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
								title="Delete tag"
							>
								<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
								</svg>
							</button>
						</div>
					</div>
					<div class="mt-3 text-xs text-gray-500">
						Created {new Date(tag.createdAt).toLocaleDateString()}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>