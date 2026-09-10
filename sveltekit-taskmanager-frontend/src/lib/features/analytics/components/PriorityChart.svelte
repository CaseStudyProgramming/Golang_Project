<script lang="ts">
	import { onMount } from 'svelte';
	import { Chart } from 'chart.js/auto';
	import type { PriorityDistribution } from '../types/analytics.types';

	let { distribution }: { distribution: PriorityDistribution } = $props();
	let canvasElement = $state<HTMLCanvasElement>();
	let chart: Chart | null = null;

	onMount(() => {
		if (canvasElement) {
			const ctx = canvasElement.getContext('2d');
			if (ctx) {
				chart = new Chart(ctx, {
					type: 'doughnut',
					data: {
						labels: ['Low', 'Medium', 'High', 'Urgent'],
						datasets: [
							{
								data: [distribution.low, distribution.medium, distribution.high, distribution.urgent],
								backgroundColor: ['#10B981', '#F59E0B', '#EF4444', '#7C3AED'],
								borderWidth: 0
							}
						]
					},
					options: {
						responsive: true,
						maintainAspectRatio: false,
						plugins: {
							legend: {
								position: 'bottom',
								labels: {
									padding: 20,
									usePointStyle: true
								}
							}
						},
						cutout: '60%'
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
			chart.data.datasets[0].data = [
				distribution.low,
				distribution.medium,
				distribution.high,
				distribution.urgent
			];
			chart.update();
		}
	});
</script>

<div class="bg-white rounded-lg shadow-sm p-6 border border-gray-100">
	<h3 class="text-lg font-semibold text-gray-800 mb-4">Priority Distribution</h3>
	<div class="h-64">
		<canvas bind:this={canvasElement}></canvas>
	</div>
</div>