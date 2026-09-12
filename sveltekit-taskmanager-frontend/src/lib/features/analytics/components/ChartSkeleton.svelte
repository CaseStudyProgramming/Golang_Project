<script lang="ts">
let { title = '' }: { title?: string } = $props()
</script>

<div class="bg-white rounded-lg shadow p-4 sm:p-6">
	<div class="animate-pulse">
		{#if title}
			<div class="h-4 bg-gray-200 rounded w-1/3 mb-4"></div>
		{/if}
		<div class="h-48 bg-gray-200 rounded"></div>
	</div>
</div>