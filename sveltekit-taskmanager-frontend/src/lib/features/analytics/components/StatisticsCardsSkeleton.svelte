<script lang="ts">
	import { Skeleton } from '$lib/shared/components';
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
	{#each Array(6) as _}
		<div class="bg-white rounded-lg shadow p-4 sm:p-6">
			<Skeleton variant="text" class="w-1/2 h-4 mb-2" />
			<Skeleton variant="text" class="w-1/3 h-8" />
			<Skeleton variant="text" class="w-1/4 h-4 mt-2" />
		</div>
	{/each}
</div>