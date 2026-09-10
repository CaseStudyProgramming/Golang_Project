<script lang="ts">
	import type { Task } from '$lib/features/tasks/types/task.types';

	let { overdueTasks }: { overdueTasks: Task[] } = $props();
</script>

<div class="bg-white rounded-lg shadow-sm p-6 border border-gray-100">
	<h3 class="text-lg font-semibold text-gray-800 mb-4">Overdue Tasks</h3>
	{#if overdueTasks.length === 0}
		<div class="text-center py-8 text-gray-500">No overdue tasks 🎉</div>
	{:else}
		<div class="space-y-3">
			{#each overdueTasks.slice(0, 5) as task}
				<div class="flex items-center justify-between p-4 bg-red-50 rounded-lg border border-red-100">
					<div class="flex-1">
						<h4 class="font-medium text-gray-800">{task.title}</h4>
						<p class="text-sm text-red-600 mt-1">
							Due: {task.dueDate ? new Date(task.dueDate).toLocaleDateString() : 'N/A'}
						</p>
					</div>
					<div class="flex items-center gap-2">
						<span class="px-2 py-1 text-xs font-medium bg-red-100 text-red-700 rounded">
							{task.priority}
						</span>
						<a href="/tasks/{task.id}" class="text-blue-600 hover:text-blue-700 text-sm font-medium">
							View
						</a>
					</div>
				</div>
			{/each}
			{#if overdueTasks.length > 5}
				<div class="text-center text-sm text-gray-500 pt-2">
					+{overdueTasks.length - 5} more overdue tasks
				</div>
			{/if}
		</div>
	{/if}
</div>