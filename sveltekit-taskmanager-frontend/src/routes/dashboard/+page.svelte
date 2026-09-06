<script lang="ts">
	import { authStore } from '$lib/shared/stores';
	import { taskStore } from '$lib/shared/stores';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			await taskStore.fetchTasks();
		} catch (error) {
			console.error('Failed to fetch tasks:', error);
		}
	});
</script>

<div class="mb-8">
	<h1 class="text-3xl font-bold text-gray-800 mb-2">Welcome back!</h1>
	<p class="text-gray-600">
		{authStore.state.user?.name || authStore.state.user?.email || 'User'}
	</p>
</div>

<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
	<div class="bg-white rounded-lg shadow p-6">
		<div class="text-3xl font-bold text-blue-600 mb-2">{taskStore.state.tasks.length}</div>
		<div class="text-gray-600">Total Tasks</div>
	</div>
	<div class="bg-white rounded-lg shadow p-6">
		<div class="text-3xl font-bold text-green-600 mb-2">
			{taskStore.state.tasks.filter((t) => t.status === 'completed').length}
		</div>
		<div class="text-gray-600">Completed</div>
	</div>
	<div class="bg-white rounded-lg shadow p-6">
		<div class="text-3xl font-bold text-orange-600 mb-2">
			{taskStore.state.tasks.filter((t) => t.status === 'in_progress').length}
		</div>
		<div class="text-gray-600">In Progress</div>
	</div>
</div>

<div class="bg-white rounded-lg shadow p-6">
	<h2 class="text-xl font-semibold text-gray-800 mb-4">Recent Tasks</h2>
	{#if taskStore.state.isLoading}
		<div class="text-center py-8 text-gray-500">Loading tasks...</div>
	{:else if taskStore.state.tasks.length === 0}
		<div class="text-center py-8 text-gray-500">No tasks yet. Create your first task!</div>
	{:else}
		<div class="space-y-3">
			{#each taskStore.state.tasks.slice(0, 5) as task}
				<div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
					<div>
						<h3 class="font-medium text-gray-800">{task.title}</h3>
						<p class="text-sm text-gray-500">{task.status}</p>
					</div>
					<a href="/dashboard/tasks/{task.id}" class="text-blue-600 hover:text-blue-700 text-sm">
						View
					</a>
				</div>
			{/each}
		</div>
	{/if}
</div>
