/**
 * Confirmation dialog store using Svelte 5 runes
 * Manages confirmation dialogs for destructive actions
 */

export interface ConfirmDialog {
	id: string
	isOpen: boolean
	title: string
	message: string
	confirmText: string
	cancelText: string
	type: ConfirmType
	onConfirm?: () => Promise<void> | void
	onCancel?: () => void
}

export type ConfirmType = 'danger' | 'info' | 'warning'

interface ConfirmState {
	dialog: ConfirmDialog | null
}

/**
 * Create confirmation dialog store with Svelte 5 runes
 */
function createConfirmStore() {
	const state = $state<ConfirmState>({
		dialog: null,
	})

	/**
	 * Show confirmation dialog
	 */
	function showConfirm(options: Omit<ConfirmDialog, 'id' | 'isOpen'>): Promise<boolean> {
		return new Promise((resolve) => {
			const id = `confirm-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`

			state.dialog = {
				...options,
				id,
				isOpen: true,
				onCancel: () => {
					options.onCancel?.()
					state.dialog = null
					resolve(false)
				},
				onConfirm: async () => {
					await options.onConfirm?.()
					state.dialog = null
					resolve(true)
				},
			}
		})
	}

	/**
	 * Close confirmation dialog
	 */
	function closeConfirm(): void {
		state.dialog = null
	}

	return {
		closeConfirm,
		showConfirm,
		get state() {
			return state
		},
	}
}

/**
 * Export confirm store instance
 * Only create store instance on client side to avoid SSR issues
 */
let confirmStoreInstance: null | ReturnType<typeof createConfirmStore> = null

// Create a safe SSR-compatible default store
function createSSRConfirmStore() {
	const state: ConfirmState = {
		dialog: null,
	}

	return {
		closeConfirm: () => {},
		showConfirm: async () => false,
		get state() {
			return state
		},
	}
}

export const confirmStore = new Proxy({} as ReturnType<typeof createConfirmStore>, {
	get(_target, prop) {
		if (!confirmStoreInstance) {
			if (typeof window === 'undefined') {
				// Return safe SSR-compatible store during server-side rendering
				confirmStoreInstance = createSSRConfirmStore()
			} else {
				confirmStoreInstance = createConfirmStore()
			}
		}
		return confirmStoreInstance[prop as keyof ReturnType<typeof createConfirmStore>]
	},
})

/**
 * Export store creator for testing
 */
export { createConfirmStore }
