<script lang="ts">
	let { message, onDismiss }: { message: string; onDismiss: () => void } = $props();
	let isVisible = $state(true);

	function handleDismiss() {
		isVisible = false;
		setTimeout(() => {
			onDismiss();
		}, 300);
	}

	// Auto-dismiss after 5 seconds
	setTimeout(() => {
		if (isVisible) {
			handleDismiss();
		}
	}, 5000);
</script>

{#if isVisible}
	<div class="fixed top-4 right-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 transition-opacity duration-300">
		<span class="flex-1">{message}</span>
		<button
			onclick={handleDismiss}
			class="text-red-500 hover:text-red-700"
			aria-label="Dismiss"
		>
			✕
		</button>
	</div>
{/if}
