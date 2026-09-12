<script lang="ts">
import Skeleton from '$lib/shared/components/Skeleton.svelte'
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
	{#each Array(6) as _, index (index)}
		<div class="bg-white rounded-lg shadow p-3 sm:p-4 border-l-4 border-gray-200">
			<div class="flex items-start justify-between mb-2">
				<div class="flex items-center gap-2">
					<Skeleton variant="circular" class="w-8 h-8" />
					<Skeleton variant="text" class="w-24 h-6" />
				</div>
				<div class="flex items-center gap-1">
					<Skeleton variant="circular" class="w-8 h-8" />
					<Skeleton variant="circular" class="w-8 h-8" />
				</div>
			</div>
			<Skeleton variant="text" class="w-full h-4 mb-2" />
			<Skeleton variant="text" class="w-2/3 h-4" />
		</div>
	{/each}
</div>