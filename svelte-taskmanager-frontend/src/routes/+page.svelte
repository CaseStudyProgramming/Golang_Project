<script>
	import { onMount } from 'svelte';

	onMount(() => {
		// Redirect to dashboard if authenticated, otherwise to login
		const token = localStorage.getItem('auth_token');
		if (token) {
			window.location.href = '/dashboard';
		} else {
			window.location.href = '/auth/login';
		}
	});
</script>

<div class="loading-container">
	<p>Loading Task Manager...</p>
</div>

<style>
	.loading-container {
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100vh;
		font-size: 1.2rem;
		color: var(--secondary-color);
	}
</style>
