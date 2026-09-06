<script lang="ts">
	let taskCount = $state(0);
	let newTaskTitle = $state('');

	function addTask() {
		if (newTaskTitle.trim()) {
			taskCount++;
			newTaskTitle = '';
		}
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-8">
	<div class="max-w-4xl mx-auto">
		<header class="mb-8">
			<h1 class="text-4xl font-bold text-gray-800 mb-2">Task Manager</h1>
			<p class="text-gray-600">Manage your tasks efficiently</p>
		</header>

		<div class="bg-white rounded-lg shadow-lg p-6 mb-6">
			<form onsubmit={(e) => { e.preventDefault(); addTask(); }} class="flex gap-4">
				<input
					type="text"
					bind:value={newTaskTitle}
					placeholder="Add a new task..."
					class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				<button
					type="submit"
					class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				>
					Add Task
				</button>
			</form>
		</div>

		<div class="bg-white rounded-lg shadow-lg p-6">
			<h2 class="text-xl font-semibold text-gray-800 mb-4">Tasks ({taskCount})</h2>
			{#if taskCount === 0}
				<p class="text-gray-500 text-center py-8">No tasks yet. Add one above!</p>
			{:else}
				<div class="space-y-2">
					{#each Array(taskCount) as _, i}
						<div class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
							<input type="checkbox" class="w-5 h-5 text-blue-600 rounded" />
							<span class="flex-1 text-gray-700">Task {i + 1}</span>
							<button class="text-red-500 hover:text-red-700">Delete</button>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
