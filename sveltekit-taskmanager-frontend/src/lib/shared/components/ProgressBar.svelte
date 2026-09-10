<script lang="ts">
	let { 
		progress = 0,
		color = 'blue',
		size = 'md',
		showLabel = false,
		label = ''
	}: {
		progress?: number;
		color?: 'blue' | 'green' | 'orange' | 'red' | 'purple';
		size?: 'sm' | 'md' | 'lg';
		showLabel?: boolean;
		label?: string;
	} = $props();

	const clampedProgress = $derived(Math.min(100, Math.max(0, progress)));

	const sizeClasses = $derived(() => {
		switch (size) {
			case 'sm':
				return 'h-1.5';
			case 'lg':
				return 'h-4';
			default:
				return 'h-2.5';
		}
	});

	const colorClasses = $derived(() => {
		switch (color) {
			case 'green':
				return 'bg-green-600';
			case 'orange':
				return 'bg-orange-600';
			case 'red':
				return 'bg-red-600';
			case 'purple':
				return 'bg-purple-600';
			default:
				return 'bg-blue-600';
		}
	});
</script>

<div class="w-full">
	{#if showLabel || label}
		<div class="flex justify-between items-center mb-1">
			{#if label}
				<span class="text-sm font-medium text-gray-700">{label}</span>
			{/if}
			<span class="text-sm font-medium text-gray-700">{clampedProgress}%</span>
		</div>
	{/if}
	<div class="w-full bg-gray-200 rounded-full {sizeClasses()}">
		<div 
			class="{colorClasses()} {sizeClasses()} rounded-full transition-all duration-300 ease-out"
			style="width: {clampedProgress}%"
			role="progressbar"
			aria-valuenow={clampedProgress}
			aria-valuemin="0"
			aria-valuemax="100"
		></div>
	</div>
</div>