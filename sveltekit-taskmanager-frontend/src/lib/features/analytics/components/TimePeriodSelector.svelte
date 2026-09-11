<script lang="ts">
	import type { TimePeriod } from '../types/analytics.types';

	let { onPeriodChange, selectedPeriod }: {
		onPeriodChange: (period: TimePeriod) => void;
		selectedPeriod: TimePeriod;
	} = $props();

	const periods: { label: string; value: TimePeriod; }[] = [
		{ label: 'Last 7 days', value: 'daily' },
		{ label: 'Last 30 days', value: 'weekly' },
		{ label: 'Last 90 days', value: 'monthly' }
	];
</script>

<div class="flex items-center gap-2">
	{#each periods as period (period.value)}
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