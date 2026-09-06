<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore } from '$lib/features/auth';
	import { onMount, onDestroy } from 'svelte';

	let { children } = $props();

	onMount(() => {
		// Initialize auth store on client side
		authStore.initialize();
	});

	onDestroy(() => {
		// Cleanup auth store when component unmounts
		authStore.cleanup();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>Task Manager</title>
</svelte:head>

{@render children()}
