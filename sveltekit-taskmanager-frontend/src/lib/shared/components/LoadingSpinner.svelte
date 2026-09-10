<script lang="ts">
	let { 
		size = 'md',
		color = 'blue',
		text = ''
	}: {
		size?: 'sm' | 'md' | 'lg';
		color?: 'blue' | 'white' | 'gray';
		text?: string;
	} = $props();

	const sizeClasses = $derived(() => {
		switch (size) {
			case 'sm':
				return 'h-4 w-4';
			case 'lg':
				return 'h-12 w-12';
			default:
				return 'h-8 w-8';
		}
	});

	const colorClasses = $derived(() => {
		switch (color) {
			case 'white':
				return 'border-white';
			case 'gray':
				return 'border-gray-400';
			default:
				return 'border-blue-600';
		}
	});
</script>

<div class="flex flex-col items-center justify-center">
	<div class="inline-block animate-spin rounded-full {sizeClasses()} border-2 {colorClasses()} border-t-transparent"></div>
	{#if text}
		<p class="mt-2 text-sm text-gray-500">{text}</p>
	{/if}
</div>