/**
 * Toast notification store using Svelte 5 runes
 * Manages toast notifications with different types and auto-dismissal
 */

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
	id: string;
	type: ToastType;
	title: string;
	message?: string;
	duration?: number;
}

interface ToastState {
	toasts: Toast[];
}

/**
 * Create toast store with Svelte 5 runes
 */
function createToastStore() {
	const state = $state<ToastState>({
		toasts: []
	});

	/**
	 * Add a new toast notification
	 */
	function addToast(toast: Omit<Toast, 'id'>): string {
		const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
		const newToast: Toast = {
			...toast,
			id,
			duration: toast.duration ?? 5000
		};

		state.toasts = [...state.toasts, newToast];

		// Auto-dismiss after duration
		if (newToast.duration && newToast.duration > 0) {
			setTimeout(() => {
				removeToast(id);
			}, newToast.duration);
		}

		return id;
	}

	/**
	 * Remove a toast notification
	 */
	function removeToast(id: string): void {
		state.toasts = state.toasts.filter((toast) => toast.id !== id);
	}

	/**
	 * Clear all toast notifications
	 */
	function clearToasts(): void {
		state.toasts = [];
	}

	/**
	 * Show success toast
	 */
	function success(title: string, message?: string, duration?: number): string {
		return addToast({ type: 'success', title, message, duration });
	}

	/**
	 * Show error toast
	 */
	function error(title: string, message?: string, duration?: number): string {
		return addToast({ type: 'error', title, message, duration });
	}

	/**
	 * Show warning toast
	 */
	function warning(title: string, message?: string, duration?: number): string {
		return addToast({ type: 'warning', title, message, duration });
	}

	/**
	 * Show info toast
	 */
	function info(title: string, message?: string, duration?: number): string {
		return addToast({ type: 'info', title, message, duration });
	}

	return {
		addToast,
		clearToasts,
		error,
		info,
		removeToast,
		success,
		warning,
		get state() {
			return state;
		}
	};
}

/**
 * Export toast store instance
 * Only create store instance on client side to avoid SSR issues
 */
let toastStoreInstance: null | ReturnType<typeof createToastStore> = null;

// Create a safe SSR-compatible default store
function createSSRToastStore() {
	const state: ToastState = {
		toasts: []
	};

	return {
		addToast: () => 'ssr-toast-id',
		clearToasts: () => {},
		error: () => 'ssr-toast-id',
		info: () => 'ssr-toast-id',
		removeToast: () => {},
		success: () => 'ssr-toast-id',
		warning: () => 'ssr-toast-id',
		get state() {
			return state;
		}
	};
}

export const toastStore = new Proxy({} as ReturnType<typeof createToastStore>, {
	get(_target, prop) {
		if (!toastStoreInstance) {
			if (typeof window === 'undefined') {
				// Return safe SSR-compatible store during server-side rendering
				toastStoreInstance = createSSRToastStore();
			} else {
				toastStoreInstance = createToastStore();
			}
		}
		return toastStoreInstance[prop as keyof ReturnType<typeof createToastStore>];
	}
});