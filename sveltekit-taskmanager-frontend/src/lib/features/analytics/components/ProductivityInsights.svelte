<script lang="ts">
import type { ProductivityInsights } from '../types/analytics.types'

let { insights }: { insights: ProductivityInsights } = $props()

interface InsightCard {
	label: string
	value: string
	description: string
	icon: string
}

const _insightCards: InsightCard[] = $derived([
	{
		description: 'Tasks completed',
		icon: '🏆',
		label: 'Total Completed',
		value: insights.totalTasksCompleted.toString(),
	},
	{
		description: 'Time to complete tasks',
		icon: '⏱️',
		label: 'Avg. Completion Time',
		value: `${insights.averageCompletionTime}h`,
	},
	{
		description: 'Day with most completions',
		icon: '📅',
		label: 'Most Productive Day',
		value: insights.mostProductiveDay,
	},
	{
		description: 'Consecutive days of activity',
		icon: '🔥',
		label: 'Current Streak',
		value: `${insights.streakDays} days`,
	},
	{
		description: 'Average daily completions',
		icon: '📈',
		label: 'Tasks per Day',
		value: insights.tasksPerDay.toFixed(1),
	},
	{
		description: 'Tasks completed by due date',
		icon: '✅',
		label: 'On-Time Rate',
		value: `${insights.onTimeCompletionRate}%`,
	},
])
</script>

<div class="bg-white rounded-lg shadow-sm p-6 border border-gray-100">
	<h3 class="text-lg font-semibold text-gray-800 mb-4">Productivity Insights</h3>
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
		{#each _insightCards as card (card.label)}
			<div class="p-4 bg-gradient-to-br from-gray-50 to-white rounded-lg border border-gray-100">
				<div class="flex items-start gap-3">
					<div class="text-2xl">{card.icon}</div>
					<div class="flex-1">
						<p class="text-sm text-gray-500 mb-1">{card.label}</p>
						<p class="text-xl font-bold text-gray-800">{card.value}</p>
						<p class="text-xs text-gray-400 mt-1">{card.description}</p>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>