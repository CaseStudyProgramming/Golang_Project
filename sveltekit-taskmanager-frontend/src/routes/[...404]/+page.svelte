<script lang="ts">
	import { goto } from '$app/navigation';

	function goHome() {
		goto('/');
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
	<div class="max-w-md w-full text-center">
		<div class="text-8xl mb-4">🔍</div>
		<h1 class="text-6xl font-bold text-gray-800 mb-2">404</h1>
		<p class="text-2xl text-gray-600 mb-6">Page Not Found</p>
		<p class="text-gray-500 mb-8">
			The page you're looking for doesn't exist or has been moved.
		</p>

		<button
			onclick={goHome}
			class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-lg font-medium"
		>
			Go Home
		</button>
	</div>
</div>
