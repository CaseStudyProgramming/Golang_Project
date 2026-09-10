<script lang="ts">
	import { fly } from 'svelte/transition';

	let {
		isOpen = false,
		title = 'Confirm Action',
		message = 'Are you sure you want to proceed?',
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		type = 'danger',
		onConfirm,
		onCancel
	}: {
		isOpen?: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		cancelText?: string;
		type?: 'danger' | 'warning' | 'info';
		onConfirm?: () => void | Promise<void>;
		onCancel?: () => void;
	} = $props();

	let isConfirming = $state(false);

	const typeClasses = $derived(() => {
		switch (type) {
			case 'danger':
				return 'bg-red-50 border-red-200 text-red-800';
			case 'warning':
				return 'bg-yellow-50 border-yellow-200 text-yellow-800';
			case 'info':
				return 'bg-blue-50 border-blue-200 text-blue-800';
		}
	});

	const buttonClasses = $derived(() => {
		switch (type) {
			case 'danger':
				return 'bg-red-600 hover:bg-red-700 focus:ring-red-500';
			case 'warning':
				return 'bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500';
			case 'info':
				return 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500';
		}
	});

	async function handleConfirm() {
		isConfirming = true;
		try {
			await onConfirm?.();
		} finally {
			isConfirming = false;
		}
	}

	function handleCancel() {
		onCancel?.();
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
			transition:fly={{ y: 20, duration: 200 }}
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
