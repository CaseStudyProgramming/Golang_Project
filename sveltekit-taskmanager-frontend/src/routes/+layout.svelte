<script lang="ts">
import './layout.css'
import { onDestroy, onMount } from 'svelte'
import { authStore } from '$lib/features/auth'
import ToastContainer from '$lib/shared/components/ToastContainer.svelte'
import ConfirmDialogContainer from '$lib/shared/components/ConfirmDialogContainer.svelte'

let { children } = $props()

const favicon = $derived('/favicon.ico')

onMount(() => {
	// Initialize auth store on client side
	try {
		authStore.initialize()
	} catch {
		// Store will initialize on first access
	}
})

onDestroy(() => {
	// Cleanup auth store when component unmounts
	try {
		authStore.cleanup()
	} catch {
		// Ignore cleanup errors
	}
})
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>Task Manager</title>
</svelte:head>

{@render children()}

<ToastContainer />
<ConfirmDialogContainer />
