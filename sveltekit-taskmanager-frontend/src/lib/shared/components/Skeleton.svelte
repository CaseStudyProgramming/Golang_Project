<script lang="ts">
	let { 
		class: className = '',
		variant = 'default'
	}: {
		class?: string;
		variant?: 'circular' | 'default' | 'rectangular' | 'text';
	} = $props();

	const baseClasses = 'animate-pulse bg-gray-200';
	
	const variantClasses = $derived(() => {
		switch (variant) {
			case 'circular':
				return 'rounded-full';
			case 'rectangular':
				return 'rounded-md';
			case 'text':
				return 'rounded h-4';
			default:
				return 'rounded-md';
		}
	});
</script>

<div class={`${baseClasses} ${variantClasses()} ${className}`}></div>