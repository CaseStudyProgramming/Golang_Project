<script lang="ts">
	import { confirmStore } from '$lib/shared/stores';
	import { ConfirmDialogComponent } from '$lib/shared/components';
</script>

{#if confirmStore.state.dialog}
	<ConfirmDialogComponent
		isOpen={confirmStore.state.dialog.isOpen}
		title={confirmStore.state.dialog.title}
		message={confirmStore.state.dialog.message}
		confirmText={confirmStore.state.dialog.confirmText}
		cancelText={confirmStore.state.dialog.cancelText}
		type={confirmStore.state.dialog.type}
		onConfirm={confirmStore.state.dialog.onConfirm}
		onCancel={confirmStore.state.dialog.onCancel}
	/>
{/if}
