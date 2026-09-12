<script lang="ts">
let {
	class: className = '',
	variant = 'default',
}: {
	class?: string
	variant?: 'circular' | 'default' | 'rectangular' | 'text'
} = $props()

const _baseClasses = 'animate-pulse bg-gray-200'

const _variantClasses = $derived(() => {
	switch (variant) {
		case 'circular':
			return 'rounded-full'
		case 'rectangular':
			return 'rounded-md'
		case 'text':
			return 'rounded h-4'
		default:
			return 'rounded-md'
	}
})
</script>

<div class={`${_baseClasses} ${_variantClasses()} ${className}`}></div>