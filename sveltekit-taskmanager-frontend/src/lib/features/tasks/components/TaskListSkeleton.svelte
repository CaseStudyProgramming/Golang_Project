<script lang="ts">
import Skeleton from '$lib/shared/components/Skeleton.svelte'
</script>

<div class="space-y-3">
	{#each Array(3) as _, index (index)}
		<div class="bg-white rounded-lg shadow p-3 sm:p-4">
			<div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
				<div class="flex-1 min-w-0 space-y-3">
					<div class="flex items-center gap-2">
						<Skeleton variant="text" class="w-1/3 h-6" />
						<Skeleton variant="default" class="w-16 h-6 rounded-full" />
						<Skeleton variant="default" class="w-20 h-6 rounded-full" />
					</div>
					<Skeleton variant="text" class="w-full h-4" />
					<Skeleton variant="text" class="w-2/3 h-4" />
					<div class="flex items-center gap-4">
						<Skeleton variant="text" class="w-24 h-4" />
						<Skeleton variant="text" class="w-20 h-4" />
					</div>
				</div>
				<div class="flex items-center gap-2">
					<Skeleton variant="circular" class="w-10 h-10" />
					<Skeleton variant="circular" class="w-10 h-10" />
					<Skeleton variant="circular" class="w-10 h-10" />
				</div>
			</div>
		</div>
	{/each}
</div>