<script lang="ts">
import { browser } from '$app/environment'
import { toastStore } from '$lib/shared/stores'
import { tagStore } from '../stores/tag.store'
import type { Tag } from '../types/tag.types'

let {
	availableTags = $bindable(browser ? tagStore.state.tags : []),
	onCreateTag,
	placeholder = 'Add tags...',
	selectedTags = $bindable([]),
}: {
	availableTags?: Tag[]
	onCreateTag?: (name: string) => Promise<void>
	placeholder?: string
	selectedTags?: string[]
} = $props()

let input = $state('')
let isOpen = $state(false)
let highlightedIndex = $state(-1)
let isLoading = $state(false)

/**
 * Get filtered available tags
 */
const filteredTags = $derived(
	!input.trim()
		? availableTags
		: availableTags.filter((tag) => tag.name.toLowerCase().includes(input.toLowerCase()))
)

/**
 * Get selected tag objects
 */
const selectedTagObjects = $derived(availableTags.filter((tag) => selectedTags.includes(tag.id)))

/**
 * Handle blur with delay to allow click events
 */
function handleBlur() {
	setTimeout(() => {
		isOpen = false
		highlightedIndex = -1
	}, 200)
}

/**
 * Handle creating new tag
 */
async function handleCreateTag() {
	if (!input.trim() || !onCreateTag) return

	isLoading = true
	try {
		await onCreateTag(input.trim())
		toastStore.success('Tag created', `Tag "${input.trim()}" has been created successfully`)
		input = ''
		isOpen = false
	} catch (error) {
		console.error('Failed to create tag:', error)
		toastStore.error(
			'Failed to create tag',
			error instanceof Error ? error.message : 'An unexpected error occurred'
		)
	} finally {
		isLoading = false
	}
}

/**
 * Handle input change
 */
function handleInput(e: Event) {
	input = (e.target as HTMLInputElement).value
	isOpen = true
	highlightedIndex = -1
}

/**
 * Handle keyboard navigation
 */
function handleKeydown(e: KeyboardEvent) {
	if (!isOpen) {
		if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
			isOpen = true
		}
		return
	}

	switch (e.key) {
		case 'ArrowDown':
			e.preventDefault()
			highlightedIndex = Math.min(highlightedIndex + 1, filteredTags.length - 1)
			break
		case 'ArrowUp':
			e.preventDefault()
			highlightedIndex = Math.max(highlightedIndex - 1, 0)
			break
		case 'Backspace':
			if (!input && selectedTags.length > 0) {
				removeTag(selectedTags[selectedTags.length - 1])
			}
			break
		case 'Enter':
			e.preventDefault()
			if (highlightedIndex >= 0 && filteredTags[highlightedIndex]) {
				selectTag(filteredTags[highlightedIndex])
			} else if (input.trim()) {
				handleCreateTag()
			}
			break
		case 'Escape':
			isOpen = false
			highlightedIndex = -1
			break
		case 'Tab':
			isOpen = false
			highlightedIndex = -1
			break
	}
}

/**
 * Handle tag removal
 */
function removeTag(tagId: string) {
	selectedTags = selectedTags.filter((id) => id !== tagId)
}

/**
 * Handle tag selection
 */
function selectTag(tag: Tag) {
	if (!selectedTags.includes(tag.id)) {
		selectedTags = [...selectedTags, tag.id]
	}
	input = ''
	isOpen = false
}
</script>

<div class="relative">
	<!-- Selected tags display -->
	{#if selectedTagObjects.length > 0}
		<div class="flex flex-wrap gap-2 mb-2" role="list">
			{#each selectedTagObjects as tag (tag.id)}
				<div
					class="inline-flex items-center gap-1 px-2 py-1 rounded-full text-sm"
					style="background-color: {tag.color || '#3B82F6'}20; color: {tag.color || '#3B82F6'}; border: 1px solid {tag.color || '#3B82F6'}40"
					role="listitem"
				>
					<span>{tag.name}</span>
					<button
						onclick={() => removeTag(tag.id)}
						onkeydown={(e) => {
							if (e.key === 'Enter' || e.key === ' ') {
								e.preventDefault();
								removeTag(tag.id);
							}
						}}
						class="hover:opacity-70 transition-opacity"
						title="Remove tag"
						aria-label={`Remove ${tag.name} tag`}
					>
						<svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Input field -->
	<div class="relative">
		<input
			type="text"
			bind:value={input}
			oninput={handleInput}
			onkeydown={handleKeydown}
			onfocus={() => (isOpen = true)}
			onblur={handleBlur}
			placeholder={placeholder}
			class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			disabled={isLoading}
			role="combobox"
			aria-expanded={isOpen}
			aria-controls="tag-listbox"
			aria-haspopup="listbox"
			aria-autocomplete="list"
		/>
		{#if isLoading}
			<div class="absolute right-3 top-1/2 -translate-y-1/2">
				<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
			</div>
		{/if}
	</div>

	<!-- Dropdown suggestions -->
	{#if isOpen && (filteredTags.length > 0 || input.trim())}
		<div class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto" role="listbox" id="tag-listbox">
			{#if filteredTags.length > 0}
				{#each filteredTags as tag, index (tag.id)}
					<button
						type="button"
						onclick={() => selectTag(tag)}
						class="w-full text-left px-3 py-2 cursor-pointer hover:bg-gray-100 {index === highlightedIndex ? 'bg-blue-50' : ''}"
						role="option"
						aria-selected={selectedTags.includes(tag.id)}
						tabindex={-1}
					>
						<div class="flex items-center gap-2">
							<div
								class="w-3 h-3 rounded-full"
								style="background-color: {tag.color || '#3B82F6'}"
							></div>
							<span>{tag.name}</span>
							{#if selectedTags.includes(tag.id)}
								<svg class="h-4 w-4 ml-auto text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
								</svg>
							{/if}
						</div>
					</button>
				{/each}
			{:else if input.trim()}
				<button
					type="button"
					onclick={handleCreateTag}
					class="w-full text-left px-3 py-2 cursor-pointer hover:bg-gray-100 text-blue-600"
					role="option"
					aria-selected="false"
					tabindex={-1}
				>
					<div class="flex items-center gap-2">
						<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
						</svg>
						<span>Create "{input}"</span>
					</div>
				</button>
			{/if}
		</div>
	{/if}
</div>