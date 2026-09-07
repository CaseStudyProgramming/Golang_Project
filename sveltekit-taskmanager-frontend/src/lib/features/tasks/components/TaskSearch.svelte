<script lang="ts">
	let {
		searchQuery = $bindable(''),
		onSearch,
		placeholder = 'Search tasks...',
		debounceMs = 300
	}: {
		searchQuery?: string;
		onSearch?: (query: string) => void;
		placeholder?: string;
		debounceMs?: number;
	} = $props();

	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	/**
	 * Handle search input with debouncing
	 */
	function handleSearchInput(event: Event): void {
		const target = event.target as HTMLInputElement;
		searchQuery = target.value;

		if (debounceTimer) {
			clearTimeout(debounceTimer);
		}

		debounceTimer = setTimeout(() => {
			onSearch?.(searchQuery);
		}, debounceMs);
	}

	/**
	 * Clear search
	 */
	function clearSearch(): void {
		searchQuery = '';
		onSearch?.('');
	}
</script>

<div class="relative">
	<div class="relative">
		<input
			type="text"
			value={searchQuery}
			oninput={handleSearchInput}
			placeholder={placeholder}
			class="w-full pl-10 pr-10 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
		/>
		<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
			<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
			</svg>
		</div>
		{#if searchQuery}
			<button
				type="button"
				onclick={clearSearch}
				class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
				title="Clear search"
			>
				<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		{/if}
	</div>
</div>