<script lang="ts">
	import type { TaskStatistics } from '../types/analytics.types';
	import StatisticsCardsSkeleton from './StatisticsCardsSkeleton.svelte';

	let { 
		statistics,
		isLoading = false 
	}: { 
		statistics: TaskStatistics;
		isLoading?: boolean;
	} = $props();

	type CardColor = 'blue' | 'green' | 'orange' | 'purple' | 'red' | 'emerald';

	interface StatCard {
		label: string;
		value: string | number;
		color: CardColor;
		icon: string;
	}

	const cards: StatCard[] = $derived([
		{
			label: 'Total Tasks',
			value: statistics.total,
			color: 'blue',
			icon: '📋'
		},
		{
			label: 'Completed',
			value: statistics.completed,
			color: 'green',
			icon: '✅'
		},
		{
			label: 'In Progress',
			value: statistics.inProgress,
			color: 'orange',
			icon: '🔄'
		},
		{
			label: 'To Do',
			value: statistics.todo,
			color: 'purple',
			icon: '📝'
		},
		{
			label: 'Overdue',
			value: statistics.overdue,
			color: 'red',
			icon: '⚠️'
		},
		{
			label: 'Completion Rate',
			value: `${statistics.completionRate}%`,
			color: 'emerald',
			icon: '📊'
		}
	]);

	const colorClasses: Record<CardColor, string> = {
		blue: 'text-blue-600 bg-blue-50',
		green: 'text-green-600 bg-green-50',
		orange: 'text-orange-600 bg-orange-50',
		purple: 'text-purple-600 bg-purple-50',
		red: 'text-red-600 bg-red-50',
		emerald: 'text-emerald-600 bg-emerald-50'
	};
</script>

{#if isLoading}
	<StatisticsCardsSkeleton />
{:else}
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
		{#each cards as card}
			<div class="bg-white rounded-lg shadow-sm p-4 sm:p-6 border border-gray-100 hover:shadow-md transition-shadow">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-medium text-gray-500 mb-1">{card.label}</p>
						<p class="text-2xl sm:text-3xl font-bold {colorClasses[card.color].split(' ')[0]}">{card.value}</p>
					</div>
					<div class="text-2xl sm:text-3xl">{card.icon}</div>
				</div>
			</div>
		{/each}
	</div>
{/if}