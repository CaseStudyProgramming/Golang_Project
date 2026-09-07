<script lang="ts">
	let {
		currentPage = $bindable(1),
		totalPages,
		onPageChange
	}: {
		currentPage?: number;
		totalPages: number;
		onPageChange?: (page: number) => void;
	} = $props();

	/**
	 * Generate page numbers to display
	 */
	function getPageNumbers(): (number | string)[] {
		const pages: (number | string)[] = [];
		const maxVisible = 5;

		if (totalPages <= maxVisible) {
			for (let i = 1; i <= totalPages; i++) {
				pages.push(i);
			}
		} else {
			if (currentPage <= 3) {
				for (let i = 1; i <= 4; i++) {
					pages.push(i);
				}
				pages.push('...');
				pages.push(totalPages);
			} else if (currentPage >= totalPages - 2) {
				pages.push(1);
				pages.push('...');
				for (let i = totalPages - 3; i <= totalPages; i++) {
					pages.push(i);
				}
			} else {
				pages.push(1);
				pages.push('...');
				for (let i = currentPage - 1; i <= currentPage + 1; i++) {
					pages.push(i);
				}
				pages.push('...');
				pages.push(totalPages);
			}
		}

		return pages;
	}

	/**
	 * Handle page change
	 */
	function handlePageChange(page: number): void {
		if (page >= 1 && page <= totalPages && page !== currentPage) {
			currentPage = page;
			onPageChange?.(page);
		}
	}
</script>

<div class="flex items-center justify-between">
	<div class="text-sm text-gray-700">
		Page <span class="font-medium">{currentPage}</span> of <span class="font-medium">{totalPages}</span>
	</div>
	<div class="flex items-center gap-2">
		<button
			onclick={() => handlePageChange(currentPage - 1)}
			disabled={currentPage === 1}
			class="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
		>
			Previous
		</button>
		<div class="flex items-center gap-1">
			{#each getPageNumbers() as page}
				{#if page === '...'}
					<span class="px-3 py-2 text-sm text-gray-500">...</span>
				{:else}
					<button
						onclick={() => handlePageChange(page as number)}
						disabled={page === currentPage}
						class="px-3 py-2 text-sm font-medium rounded-lg transition-colors {page === currentPage
							? 'bg-blue-600 text-white'
							: 'text-gray-700 bg-white border border-gray-300 hover:bg-gray-50'}"
					>
						{page}
					</button>
				{/if}
			{/each}
		</div>
		<button
			onclick={() => handlePageChange(currentPage + 1)}
			disabled={currentPage === totalPages}
			class="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
		>
			Next
		</button>
	</div>
</div>