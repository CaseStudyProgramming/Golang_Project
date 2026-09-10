/**
 * Global UI state store using Svelte 5 runes
 * Manages loading states, notifications, and modals
 */

import type { LoadingState, Modal, Notification, NotificationType } from '../types/ui.types';

/**
 * Options for adding notifications
 */
interface AddNotificationOptions {
	duration?: number;
	isPersistent?: boolean;
}

/**
 * Options for opening modals
 */
interface OpenModalOptions {
	confirmText?: string;
	cancelText?: string;
}

/**
 * UI store state interface
 */
interface UIState {
	loading: LoadingState;
	notifications: Notification[];
	modals: Modal[];
}

/**
 * Create UI store with Svelte 5 runes
 */
function createUIStore() {
	const state = $state<UIState>({
		loading: {
			isLoading: false,
			message: undefined
		},
		modals: [],
		notifications: []
	});

	/**
	 * Set loading state
	 */
	function setLoading(isLoading: boolean, message?: string): void {
		state.loading = { isLoading, message };
	}

	/**
	 * Clear loading state
	 */
	function clearLoading(): void {
		state.loading = { isLoading: false, message: undefined };
	}

	/**
	 * Add notification
	 */
	function addNotification(
		type: NotificationType,
		title: string,
		message?: string,
		options?: AddNotificationOptions
	): string {
		const id = Date.now().toString();
		const notification: Notification = {
			createdAt: Date.now(),
			duration: options?.duration || 5000,
			id,
			isPersistent: options?.isPersistent || false,
			message,
			title,
			type
		};

		state.notifications = [...state.notifications, notification];

		// Auto-remove notification if not persistent
		if (!notification.isPersistent && notification.duration) {
			setTimeout(() => {
				removeNotification(id);
			}, notification.duration);
		}

		return id;
	}

	/**
	 * Remove notification by ID
	 */
	function removeNotification(id: string): void {
		state.notifications = state.notifications.filter(n => n.id !== id);
	}

	/**
	 * Clear all notifications
	 */
	function clearNotifications(): void {
		state.notifications = [];
	}

	/**
	 * Add success notification
	 */
	function success(title: string, message?: string, options?: AddNotificationOptions): string {
		return addNotification('success', title, message, options);
	}

	/**
	 * Add error notification
	 */
	function error(title: string, message?: string, options?: AddNotificationOptions): string {
		return addNotification('error', title, message, { ...options, isPersistent: true });
	}

	/**
	 * Add warning notification
	 */
	function warning(title: string, message?: string, options?: AddNotificationOptions): string {
		return addNotification('warning', title, message, options);
	}

	/**
	 * Add info notification
	 */
	function info(title: string, message?: string, options?: AddNotificationOptions): string {
		return addNotification('info', title, message, options);
	}

	/**
	 * Open modal
	 */
	function openModal(modal: Omit<Modal, 'id' | 'isOpen'>): string {
		const id = Date.now().toString();
		const newModal: Modal = {
			...modal,
			id,
			isOpen: true
		};

		state.modals = [...state.modals, newModal];
		return id;
	}

	/**
	 * Close modal by ID
	 */
	function closeModal(id: string): void {
		state.modals = state.modals.map(modal =>
			modal.id === id ? { ...modal, isOpen: false } : modal
		);

		// Remove modal after animation
		setTimeout(() => {
			state.modals = state.modals.filter(modal => modal.id !== id);
		}, 300);
	}

	/**
	 * Close all modals
	 */
	function closeAllModals(): void {
		state.modals = state.modals.map(modal => ({ ...modal, isOpen: false }));

		// Remove all modals after animation
		setTimeout(() => {
			state.modals = [];
		}, 300);
	}

	/**
	 * Open confirm modal
	 */
	function confirm(
		title: string,
		message?: string,
		onConfirm?: () => Promise<void> | void,
		options?: OpenModalOptions
	): string {
		return openModal({
			cancelText: options?.cancelText || 'Cancel',
			confirmText: options?.confirmText || 'Confirm',
			message,
			onConfirm,
			showCancel: true,
			title,
			type: 'confirm'
		});
	}

	/**
	 * Open alert modal
	 */
	function alert(title: string, message?: string, options?: OpenModalOptions): string {
		return openModal({
			confirmText: options?.confirmText || 'OK',
			message,
			showCancel: false,
			title,
			type: 'alert'
		});
	}

	/**
	 * Reset store state
	 */
	function reset(): void {
		state.loading = { isLoading: false, message: undefined };
		state.notifications = [];
		state.modals = [];
	}

	return {
		addNotification,
		alert,
		clearLoading,
		clearNotifications,
		closeAllModals,
		closeModal,
		confirm,
		error,
		info,
		openModal,
		removeNotification,
		reset,
		setLoading,
		get state() {
			return state;
		},
		success,
		warning
	};
}

/**
 * Export UI store instance
 * Only create store instance on client side to avoid SSR issues
 */
let uiStoreInstance: null | ReturnType<typeof createUIStore> = null;

export const uiStore = new Proxy({} as ReturnType<typeof createUIStore>, {
	get(_target, prop) {
		if (!uiStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('uiStore can only be accessed on the client side');
			}
			uiStoreInstance = createUIStore();
		}
		return uiStoreInstance[prop as keyof ReturnType<typeof createUIStore>];
	}
});

/**
 * Export store creator for testing
 */
export { createUIStore };
