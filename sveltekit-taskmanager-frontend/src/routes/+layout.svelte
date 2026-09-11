<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore } from '$lib/features/auth';
	import { ConfirmDialogContainer, ToastContainer } from '$lib/shared/components';
	import { onDestroy, onMount } from 'svelte';

	let { children } = $props();

	onMount(() => {
		// Initialize auth store on client side
		try {
			authStore.initialize();
		} catch (e) {
			// Store will initialize on first access
		}
	});

	onDestroy(() => {
		// Cleanup auth store when component unmounts
		try {
			authStore.cleanup();
		} catch (e) {
			// Ignore cleanup errors
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>Task Manager</title>
</svelte:head>

{@render children()}

<ToastContainer />
<ConfirmDialogContainer />
