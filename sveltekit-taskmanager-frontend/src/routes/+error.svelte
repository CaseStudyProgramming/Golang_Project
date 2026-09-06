<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	function goHome() {
		goto('/');
	}

	function goBack() {
		window.history.back();
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-red-50 to-orange-100 flex items-center justify-center p-4">
	<div class="max-w-md w-full text-center">
		<div class="text-6xl mb-4">⚠️</div>
		<h1 class="text-4xl font-bold text-gray-800 mb-2">Oops!</h1>
		<p class="text-xl text-gray-600 mb-6">Something went wrong</p>

		{#if $page.error?.message}
			<div class="bg-white rounded-lg shadow-lg p-6 mb-6">
				<p class="text-gray-700">{$page.error.message}</p>
			</div>
		{/if}

		<div class="flex flex-col sm:flex-row gap-4 justify-center">
			<button
				onclick={goBack}
				class="px-6 py-2 bg-white text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
			>
				Go Back
			</button>
			<button
				onclick={goHome}
				class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
			>
				Go Home
			</button>
		</div>
	</div>
</div>
