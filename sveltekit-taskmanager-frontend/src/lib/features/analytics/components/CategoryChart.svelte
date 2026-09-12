<script lang="ts">
import { Chart } from 'chart.js/auto'
import { onMount } from 'svelte'

import type { CategoryDistribution } from '../types/analytics.types'

let { categoryDistribution }: { categoryDistribution: CategoryDistribution[] } = $props()
let canvasElement = $state<HTMLCanvasElement>()
let chart: Chart | null = null

onMount(() => {
	if (canvasElement) {
		const ctx = canvasElement.getContext('2d')
		if (ctx) {
			const colors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#06B6D4']

			chart = new Chart(ctx, {
				data: {
					datasets: [
						{
							backgroundColor: colors,
							borderRadius: 4,
							data: categoryDistribution.map((c) => c.count),
							label: 'Total',
						},
						{
							backgroundColor: colors.map((c) => `${c}80`),
							borderRadius: 4,
							data: categoryDistribution.map((c) => c.completed),
							label: 'Completed',
						},
					],
					labels: categoryDistribution.map((c) => c.categoryName),
				},
				options: {
					maintainAspectRatio: false,
					plugins: {
						legend: {
							labels: {
								padding: 20,
								usePointStyle: true,
							},
							position: 'top',
						},
					},
					responsive: true,
					scales: {
						y: {
							beginAtZero: true,
							ticks: {
								precision: 0,
							},
						},
					},
				},
				type: 'bar',
			})
		}
	}

	return () => {
		if (chart) {
			chart.destroy()
		}
	}
})

$effect(() => {
	if (chart) {
		chart.data.labels = categoryDistribution.map((c) => c.categoryName)
		chart.data.datasets[0].data = categoryDistribution.map((c) => c.count)
		chart.data.datasets[1].data = categoryDistribution.map((c) => c.completed)
		chart.update()
	}
})
</script>

<div class="bg-white rounded-lg shadow-sm p-6 border border-gray-100">
	<h3 class="text-lg font-semibold text-gray-800 mb-4">Category Distribution</h3>
	{#if categoryDistribution.length === 0}
		<div class="text-center py-8 text-gray-500">No category data available</div>
	{:else}
		<div class="h-64">
			<canvas bind:this={canvasElement}></canvas>
		</div>
	{/if}
</div>