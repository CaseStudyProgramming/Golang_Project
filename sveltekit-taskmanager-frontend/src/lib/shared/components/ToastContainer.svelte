<script lang="ts">
	import { toastStore } from '$lib/shared/stores';
	import { ToastComponent } from '$lib/shared/components';
</script>

<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
	{#each toastStore.state.toasts as toast (toast.id)}
		<ToastComponent {toast} />
	{/each}
</div>