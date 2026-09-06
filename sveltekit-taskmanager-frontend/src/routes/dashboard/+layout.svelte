<script lang="ts">
	import { authStore } from '$lib/shared/stores';
	import { goto } from '$app/navigation';

	let { children } = $props();

	function handleLogout() {
		authStore.logout();
		goto('/auth/login');
	}
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm border-b border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/dashboard" class="text-xl font-bold text-gray-800">Task Manager</a>
				</div>
				<div class="flex items-center space-x-4">
					<a
						href="/dashboard"
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Dashboard
					</a>
					<a
						href="/dashboard/tasks"
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Tasks
					</a>
					<button
						onclick={handleLogout}
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Logout
					</button>
				</div>
			</div>
		</div>
	</nav>

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		{@render children()}
	</main>
</div>
