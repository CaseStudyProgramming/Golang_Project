<script lang="ts">
	import { fly } from 'svelte/transition';

	let {
		cancelText = 'Cancel',
		confirmText = 'Confirm',
		isOpen = false,
		message = 'Are you sure you want to proceed?',
		onCancel,
		onConfirm,
		title = 'Confirm Action',
		type = 'danger'
	}: {
		cancelText?: string;
		confirmText?: string;
		isOpen?: boolean;
		message?: string;
		onCancel?: () => void;
		onConfirm?: () => Promise<void> | void;
		title?: string;
		type?: 'danger' | 'info' | 'warning';
	} = $props();

	let isConfirming = $state(false);

	const typeClasses = $derived(() => {
		switch (type) {
			case 'danger':
				return 'bg-red-50 border-red-200 text-red-800';
			case 'info':
				return 'bg-blue-50 border-blue-200 text-blue-800';
			case 'warning':
				return 'bg-yellow-50 border-yellow-200 text-yellow-800';
		}
	});

	const buttonClasses = $derived(() => {
		switch (type) {
			case 'danger':
				return 'bg-red-600 hover:bg-red-700 focus:ring-red-500';
			case 'info':
				return 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500';
			case 'warning':
				return 'bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500';
		}
	});

	function handleCancel() {
		onCancel?.();
	}

	async function handleConfirm() {
		isConfirming = true;
		try {
			await onConfirm?.();
		} finally {
			isConfirming = false;
		}
	}
</script>

{#if isOpen}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<button 
			class="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
			onclick={handleCancel}
			onkeydown={(e) => { if (e.key === 'Escape') handleCancel(); }}
			aria-label="Close dialog"
		></button>
		
		<div 
			class="relative bg-white rounded-lg shadow-xl max-w-md w-full p-6 border {typeClasses()}"
			transition:fly={{ duration: 200, y: 20 }}
			role="dialog"
			aria-modal="true"
			aria-labelledby="dialog-title"
		>
			<h3 id="dialog-title" class="text-lg font-medium text-gray-900 mb-2">{title}</h3>
			<p class="text-sm text-gray-600 mb-6">{message}</p>
			
			<div class="flex flex-col sm:flex-row justify-end gap-3">
				<button
					onclick={handleCancel}
					disabled={isConfirming}
					class="w-full sm:w-auto px-4 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 active:bg-gray-100 disabled:opacity-50 transition-colors min-h-[44px]"
				>
					{cancelText}
				</button>
				<button
					onclick={handleConfirm}
					disabled={isConfirming}
					class="w-full sm:w-auto px-4 py-2.5 text-sm font-medium text-white {buttonClasses()} rounded-lg disabled:opacity-50 transition-colors min-h-[44px] flex items-center justify-center gap-2"
				>
					{#if isConfirming}
						<div class="inline-block animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
					{/if}
					{isConfirming ? 'Processing...' : confirmText}
				</button>
			</div>
		</div>
	</div>
{/if}
