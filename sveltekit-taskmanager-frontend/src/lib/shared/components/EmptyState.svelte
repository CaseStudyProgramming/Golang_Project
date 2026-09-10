<script lang="ts">
	let { 
		icon = '',
		title = 'No data',
		description = 'There is no data to display.',
		actionLabel = '',
		onAction
	}: {
		icon?: string;
		title?: string;
		description?: string;
		actionLabel?: string;
		onAction?: () => void;
	} = $props();
</script>

<div class="text-center py-8 sm:py-12 px-4">
	{#if icon}
		<div class="mx-auto h-16 w-16 sm:h-20 sm:w-20 text-gray-400 mb-4 flex items-center justify-center">
			{#if icon.trim().startsWith('<svg')}
				<div class="h-full w-full">
					{@html icon}
				</div>
			{:else}
				<span class="text-4xl sm:text-5xl">{icon}</span>
			{/if}
		</div>
	{/if}
	<h3 class="text-base sm:text-lg font-medium text-gray-900 mb-2">{title}</h3>
	<p class="text-sm sm:text-base text-gray-500 mb-6 max-w-md mx-auto">{description}</p>
	{#if actionLabel && onAction}
		<button
			onclick={onAction}
			class="inline-flex items-center px-4 py-2 sm:px-6 sm:py-2.5 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 min-h-[44px]"
		>
			{actionLabel}
		</button>
	{/if}
</div>