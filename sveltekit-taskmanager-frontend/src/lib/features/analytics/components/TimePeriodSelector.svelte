<script lang="ts">
	import type { TimePeriod } from '../types/analytics.types';

	let { selectedPeriod, onPeriodChange }: {
		selectedPeriod: TimePeriod;
		onPeriodChange: (period: TimePeriod) => void;
	} = $props();

	const periods: { value: TimePeriod; label: string }[] = [
		{ value: 'daily', label: 'Last 7 days' },
		{ value: 'weekly', label: 'Last 30 days' },
		{ value: 'monthly', label: 'Last 90 days' }
	];
</script>

<div class="flex items-center gap-2">
	{#each periods as period}
		<button
			onclick={() => onPeriodChange(period.value)}
			class="px-4 py-2 rounded-lg text-sm font-medium transition-colors {selectedPeriod === period.value
				? 'bg-blue-600 text-white'
				: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
		>
			{period.label}
		</button>
	{/each}
</div>