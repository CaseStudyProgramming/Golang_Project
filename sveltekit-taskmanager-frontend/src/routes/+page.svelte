<script lang="ts">
	import { authStore } from '$lib/shared/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		if (authStore.state.isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
	<div class="max-w-4xl mx-auto text-center">
		<header class="mb-8">
			<h1 class="text-5xl font-bold text-gray-800 mb-4">Task Manager</h1>
			<p class="text-xl text-gray-600">Manage your tasks efficiently</p>
		</header>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="text-4xl mb-4">📋</div>
				<h3 class="text-xl font-semibold text-gray-800 mb-2">Organize Tasks</h3>
				<p class="text-gray-600">Create and manage your tasks with ease</p>
			</div>
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="text-4xl mb-4">📊</div>
				<h3 class="text-xl font-semibold text-gray-800 mb-2">Track Progress</h3>
				<p class="text-gray-600">Monitor your productivity and completion rates</p>
			</div>
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="text-4xl mb-4">🏷️</div>
				<h3 class="text-xl font-semibold text-gray-800 mb-2">Categories & Tags</h3>
				<p class="text-gray-600">Organize tasks with categories and tags</p>
			</div>
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="text-4xl mb-4">🔒</div>
				<h3 class="text-xl font-semibold text-gray-800 mb-2">Secure</h3>
				<p class="text-gray-600">Your data is safe and secure</p>
			</div>
		</div>

		<div class="flex flex-col sm:flex-row gap-4 justify-center">
			<a
				href="/auth/register"
				class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-lg font-medium"
			>
				Get Started
			</a>
			<a
				href="/auth/login"
				class="px-8 py-3 bg-white text-blue-600 border border-blue-600 rounded-lg hover:bg-blue-50 transition-colors text-lg font-medium"
			>
				Sign In
			</a>
		</div>
	</div>
</div>
