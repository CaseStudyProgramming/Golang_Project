<script lang="ts">
import { Chart } from 'chart.js/auto'
import { onMount } from 'svelte'

import type { TimeBasedData } from '../types/analytics.types'

let { timeBasedData }: { timeBasedData: TimeBasedData[] } = $props()
let canvasElement = $state<HTMLCanvasElement>()
let chart: Chart | null = null

onMount(() => {
	if (canvasElement) {
		const ctx = canvasElement.getContext('2d')
		if (ctx) {
			chart = new Chart(ctx, {
				data: {
					datasets: [
						{
							backgroundColor: 'rgba(16, 185, 129, 0.1)',
							borderColor: '#10B981',
							data: timeBasedData.map((d) => d.completed),
							fill: true,
							label: 'Completed',
							tension: 0.4,
						},
						{
							backgroundColor: 'rgba(59, 130, 246, 0.1)',
							borderColor: '#3B82F6',
							data: timeBasedData.map((d) => d.created),
							fill: true,
							label: 'Created',
							tension: 0.4,
						},
					],
					labels: timeBasedData.map((d) => {
						const date = new Date(d.date)
						return date.toLocaleDateString('en-US', { day: 'numeric', month: 'short' })
					}),
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
				type: 'line',
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
		chart.data.labels = timeBasedData.map((d) => {
			const date = new Date(d.date)
			return date.toLocaleDateString('en-US', { day: 'numeric', month: 'short' })
		})
		chart.data.datasets[0].data = timeBasedData.map((d) => d.completed)
		chart.data.datasets[1].data = timeBasedData.map((d) => d.created)
		chart.update()
	}
})
</script>

<div class="bg-white rounded-lg shadow-sm p-6 border border-gray-100">
	<h3 class="text-lg font-semibold text-gray-800 mb-4">Task Completion Trends</h3>
	<div class="h-64">
		<canvas bind:this={canvasElement}></canvas>
	</div>
</div>