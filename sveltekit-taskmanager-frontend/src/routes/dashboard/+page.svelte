<script lang="ts">
import { onMount } from 'svelte'
import type { TimePeriod } from '$lib/features/analytics'
import { analyticsStore } from '$lib/features/analytics'
import type { TaskStatistics } from '$lib/features/analytics/types/analytics.types'
import TimePeriodSelector from '$lib/features/analytics/components/TimePeriodSelector.svelte'
import StatisticsCards from '$lib/features/analytics/components/StatisticsCards.svelte'
import ChartSkeleton from '$lib/features/analytics/components/ChartSkeleton.svelte'
import CompletionRateChart from '$lib/features/analytics/components/CompletionRateChart.svelte'
import PriorityChart from '$lib/features/analytics/components/PriorityChart.svelte'
import CategoryChart from '$lib/features/analytics/components/CategoryChart.svelte'
import OverdueTasksSummary from '$lib/features/analytics/components/OverdueTasksSummary.svelte'
import ProductivityInsights from '$lib/features/analytics/components/ProductivityInsights.svelte'
import { taskStore } from '$lib/features/tasks'
import type { Task } from '$lib/features/tasks/types/task.types'
import { authStore } from '$lib/features/auth'
import EmptyState from '$lib/shared/components/EmptyState.svelte'

onMount(async () => {
	try {
		await taskStore.fetchTasks()
		await analyticsStore.refreshAnalytics()
	} catch (error) {
		console.error('Failed to fetch data:', error)
	}
})

function handlePeriodChange(period: TimePeriod): void {
	analyticsStore.setPeriod(period)
}

const overdueTasks = $derived(
	taskStore.state.tasks.filter(
		(t: Task) => t.dueDate && new Date(t.dueDate) < new Date() && t.status !== 'completed'
	)
)

const defaultStatistics: TaskStatistics = {
	cancelled: 0,
	completed: 0,
	completionRate: 0,
	inProgress: 0,
	overdue: 0,
	todo: 0,
	total: 0,
}
</script>

<div class="mb-6 sm:mb-8">
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-2">
		<h1 class="text-2xl sm:text-3xl font-bold text-gray-800">Dashboard</h1>
		<TimePeriodSelector
			selectedPeriod={analyticsStore.state.selectedPeriod}
			onPeriodChange={(period: TimePeriod) => handlePeriodChange(period)}
		/>
	</div>
	<p class="text-gray-600 text-sm sm:text-base">
		Welcome back, {authStore.state.user?.name || authStore.state.user?.email || 'User'}!
	</p>
</div>

{#if analyticsStore.state.isLoading}
	<div class="space-y-4 sm:space-y-6">
		<!-- Statistics Cards Skeleton -->
		<StatisticsCards statistics={defaultStatistics} isLoading={true} />
		
		<!-- Charts Skeleton -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
			<ChartSkeleton title="Completion Rate" />
			<ChartSkeleton title="Priority Distribution" />
		</div>
		
		<!-- Additional Charts Skeleton -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
			<ChartSkeleton title="Category Distribution" />
			<ChartSkeleton title="Overdue Tasks" />
		</div>
	</div>
{:else if analyticsStore.state.data}
	<div class="space-y-4 sm:space-y-6">
		<!-- Statistics Cards -->
		<StatisticsCards statistics={analyticsStore.state.data.statistics} />

		<!-- Charts Row -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
			<CompletionRateChart timeBasedData={analyticsStore.state.data.timeBasedData} />
			<PriorityChart distribution={analyticsStore.state.data.priorityDistribution} />
		</div>

		<!-- Category Chart and Overdue Tasks -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
			<CategoryChart categoryDistribution={analyticsStore.state.data.categoryDistribution} />
			<OverdueTasksSummary overdueTasks={overdueTasks} />
		</div>

		<!-- Productivity Insights -->
		<ProductivityInsights insights={analyticsStore.state.data.productivityInsights} />
	</div>
{:else}
	<div class="bg-white rounded-lg shadow">
		<EmptyState 
			icon='<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>'
			title="No analytics data"
			description="Start creating tasks to see your productivity analytics and insights."
		/>
	</div>
{/if}
