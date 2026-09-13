<script lang="ts">
import { fly } from 'svelte/transition'
import { toastStore } from '../stores'
import type { Toast as ToastType } from '../stores/toast.store.svelte.ts'

let {
	onDismiss,
	toast,
}: {
	onDismiss?: (id: string) => void
	toast: ToastType
} = $props()

const _typeIcons = $derived(() => {
	switch (toast.type) {
		case 'error':
			return '<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>'
		case 'info':
			return '<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>'
		case 'success':
			return '<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>'
		case 'warning':
			return '<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>'
	}
})

const _typeColors = $derived(() => {
	switch (toast.type) {
		case 'error':
			return 'bg-red-50 border-red-200 text-red-800'
		case 'info':
			return 'bg-blue-50 border-blue-200 text-blue-800'
		case 'success':
			return 'bg-green-50 border-green-200 text-green-800'
		case 'warning':
			return 'bg-yellow-50 border-yellow-200 text-yellow-800'
	}
})

const _iconColors = $derived(() => {
	switch (toast.type) {
		case 'error':
			return 'text-red-500'
		case 'info':
			return 'text-blue-500'
		case 'success':
			return 'text-green-500'
		case 'warning':
			return 'text-yellow-500'
	}
})

function _handleDismiss() {
	if (onDismiss) {
		onDismiss(toast.id)
	} else {
		toastStore.removeToast(toast.id)
	}
}
</script>

<div
	class="pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 {_typeColors()} border"
	transition:fly={{ duration: 300, y: -20 }}
>
	<div class="p-4">
		<div class="flex items-start">
			<div class="flex-shrink-0 {_iconColors()}">
				<!-- eslint-disable svelte/no-at-html-tags -->
				<!-- Icons are hardcoded within the component, not user input -->
				{@html _typeIcons()}
			</div>
			<div class="ml-3 flex-1 w-0">
				<p class="text-sm font-medium">{toast.title}</p>
				{#if toast.message}
					<p class="mt-1 text-sm opacity-90">{toast.message}</p>
				{/if}
			</div>
			<div class="ml-4 flex flex-shrink-0">
				<button
					onclick={_handleDismiss}
					class="inline-flex rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-500"
					aria-label="Close"
				>
					<svg class="h-5 w-5 opacity-70 hover:opacity-100" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>
		</div>
	</div>
</div>