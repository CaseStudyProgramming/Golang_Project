<script lang="ts">
	import { onMount } from 'svelte';
	import { Chart } from 'chart.js/auto';
	import type { TimeBasedData } from '../types/analytics.types';

	let { timeBasedData }: { timeBasedData: TimeBasedData[] } = $props();
	let canvasElement = $state<HTMLCanvasElement>();
	let chart: Chart | null = null;

	onMount(() => {
		if (canvasElement) {
			const ctx = canvasElement.getContext('2d');
			if (ctx) {
				chart = new Chart(ctx, {
					type: 'line',
					data: {
						labels: timeBasedData.map((d) => {
							const date = new Date(d.date);
							return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
						}),
						datasets: [
							{
								label: 'Completed',
								data: timeBasedData.map((d) => d.completed),
								borderColor: '#10B981',
								backgroundColor: 'rgba(16, 185, 129, 0.1)',
								fill: true,
								tension: 0.4
							},
							{
								label: 'Created',
								data: timeBasedData.map((d) => d.created),
								borderColor: '#3B82F6',
								backgroundColor: 'rgba(59, 130, 246, 0.1)',
								fill: true,
								tension: 0.4
							}
						]
					},
					options: {
						responsive: true,
						maintainAspectRatio: false,
						plugins: {
							legend: {
								position: 'top',
								labels: {
									padding: 20,
									usePointStyle: true
								}
							}
						},
						scales: {
							y: {
								beginAtZero: true,
								ticks: {
									precision: 0
								}
							}
						}
					}
				});
			}
		}

		return () => {
			if (chart) {
				chart.destroy();
			}
		};
	});

	$effect(() => {
		if (chart) {
			chart.data.labels = timeBasedData.map((d) => {
				const date = new Date(d.date);
				return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
			});
			chart.data.datasets[0].data = timeBasedData.map((d) => d.completed);
			chart.data.datasets[1].data = timeBasedData.map((d) => d.created);
			chart.update();
		}
	});
</script>

<div class="bg-white rounded-lg shadow-sm p-6 border border-gray-100">
	<h3 class="text-lg font-semibold text-gray-800 mb-4">Task Completion Trends</h3>
	<div class="h-64">
		<canvas bind:this={canvasElement}></canvas>
	</div>
</div>