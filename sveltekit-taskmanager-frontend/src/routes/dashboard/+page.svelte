<script lang="ts">
	import { authStore } from '$lib/features/auth';
	import { taskStore } from '$lib/features/tasks';
	import { analyticsStore } from '$lib/features/analytics';
	import { onMount } from 'svelte';
	import StatisticsCards from '$lib/features/analytics/components/StatisticsCards.svelte';
	import PriorityChart from '$lib/features/analytics/components/PriorityChart.svelte';
	import CompletionRateChart from '$lib/features/analytics/components/CompletionRateChart.svelte';
	import CategoryChart from '$lib/features/analytics/components/CategoryChart.svelte';
	import OverdueTasksSummary from '$lib/features/analytics/components/OverdueTasksSummary.svelte';
	import ProductivityInsights from '$lib/features/analytics/components/ProductivityInsights.svelte';
	import TimePeriodSelector from '$lib/features/analytics/components/TimePeriodSelector.svelte';
	import ChartSkeleton from '$lib/features/analytics/components/ChartSkeleton.svelte';
	import type { TimePeriod } from '$lib/features/analytics';
	import type { TaskStatistics } from '$lib/features/analytics/types/analytics.types';
	import type { Task } from '$lib/features/tasks/types/task.types';

	onMount(async () => {
		try {
			await taskStore.fetchTasks();
			await analyticsStore.refreshAnalytics();
		} catch (error) {
			console.error('Failed to fetch data:', error);
		}
	});

	function handlePeriodChange(period: TimePeriod) {
		analyticsStore.setPeriod(period);
	}

	const overdueTasks = $derived(
		taskStore.state.tasks.filter(
			(t: Task) => t.dueDate && new Date(t.dueDate) < new Date() && t.status !== 'completed'
		)
	);

	const defaultStatistics: TaskStatistics = {
		total: 0,
		completed: 0,
		inProgress: 0,
		todo: 0,
		overdue: 0,
		cancelled: 0,
		completionRate: 0
	};
</script>

<div class="mb-6 sm:mb-8">
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-2">
		<h1 class="text-2xl sm:text-3xl font-bold text-gray-800">Dashboard</h1>
		<TimePeriodSelector
			selectedPeriod={analyticsStore.state.selectedPeriod}
			onPeriodChange={handlePeriodChange}
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
	<div class="text-center py-8 sm:py-12">
		<p class="text-gray-500 text-sm sm:text-base">No analytics data available</p>
	</div>
{/if}
