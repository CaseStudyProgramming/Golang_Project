/**
 * Confirmation dialog store using Svelte 5 runes
 * Manages confirmation dialogs for destructive actions
 */

export type ConfirmType = 'danger' | 'warning' | 'info';

export interface ConfirmDialog {
	id: string;
	isOpen: boolean;
	title: string;
	message: string;
	confirmText: string;
	cancelText: string;
	type: ConfirmType;
	onConfirm?: () => void | Promise<void>;
	onCancel?: () => void;
}

interface ConfirmState {
	dialog: ConfirmDialog | null;
}

/**
 * Create confirmation dialog store with Svelte 5 runes
 */
function createConfirmStore() {
	const state = $state<ConfirmState>({
		dialog: null
	});

	/**
	 * Show confirmation dialog
	 */
	function showConfirm(options: Omit<ConfirmDialog, 'id' | 'isOpen'>): Promise<boolean> {
		return new Promise((resolve) => {
			const id = `confirm-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
			
			state.dialog = {
				...options,
				id,
				isOpen: true,
				onConfirm: async () => {
					await options.onConfirm?.();
					state.dialog = null;
					resolve(true);
				},
				onCancel: () => {
					options.onCancel?.();
					state.dialog = null;
					resolve(false);
				}
			};
		});
	}

	/**
	 * Close confirmation dialog
	 */
	function closeConfirm(): void {
		state.dialog = null;
	}

	return {
		showConfirm,
		closeConfirm,
		get state() {
			return state;
		}
	};
}

/**
 * Export confirm store instance
 * Only create store instance on client side to avoid SSR issues
 */
let confirmStoreInstance: null | ReturnType<typeof createConfirmStore> = null;

export const confirmStore = new Proxy({} as ReturnType<typeof createConfirmStore>, {
	get(_target, prop) {
		if (!confirmStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('confirmStore can only be accessed on the client side');
			}
			confirmStoreInstance = createConfirmStore();
		}
		return confirmStoreInstance[prop as keyof ReturnType<typeof createConfirmStore>];
	}
});

/**
 * Export store creator for testing
 */
export { createConfirmStore };
