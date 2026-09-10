<script lang="ts">
	let { 
		class: className = '',
		variant = 'default'
	}: {
		class?: string;
		variant?: 'default' | 'circular' | 'text' | 'rectangular';
	} = $props();

	const baseClasses = 'animate-pulse bg-gray-200';
	
	const variantClasses = $derived(() => {
		switch (variant) {
			case 'circular':
				return 'rounded-full';
			case 'text':
				return 'rounded h-4';
			case 'rectangular':
				return 'rounded-md';
			default:
				return 'rounded-md';
		}
	});
</script>

<div class={`${baseClasses} ${variantClasses()} ${className}`}></div>