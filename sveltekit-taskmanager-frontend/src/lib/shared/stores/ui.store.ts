/**
 * Global UI state store using Svelte 5 runes
 * Manages loading states, notifications, and modals
 */

import type { Notification, Modal, LoadingState, NotificationType } from '../types/ui.types';

/**
 * UI store state interface
 */
interface UIState {
	loading: LoadingState;
	notifications: Notification[];
	modals: Modal[];
}

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
 * Create UI store with Svelte 5 runes
 */
function createUIStore() {
	const state = $state<UIState>({
		loading: {
			isLoading: false,
			message: undefined
		},
		notifications: [],
		modals: []
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
			id,
			type,
			title,
			message,
			duration: options?.duration || 5000,
			isPersistent: options?.isPersistent || false,
			createdAt: Date.now()
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
		state.notifications = state.notifications.filter((n) => n.id !== id);
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
		state.modals = state.modals.map((modal) =>
			modal.id === id ? { ...modal, isOpen: false } : modal
		);

		// Remove modal after animation
		setTimeout(() => {
			state.modals = state.modals.filter((modal) => modal.id !== id);
		}, 300);
	}

	/**
	 * Close all modals
	 */
	function closeAllModals(): void {
		state.modals = state.modals.map((modal) => ({ ...modal, isOpen: false }));

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
		onConfirm?: () => void | Promise<void>,
		options?: OpenModalOptions
	): string {
		return openModal({
			type: 'confirm',
			title,
			message,
			onConfirm,
			confirmText: options?.confirmText || 'Confirm',
			cancelText: options?.cancelText || 'Cancel',
			showCancel: true
		});
	}

	/**
	 * Open alert modal
	 */
	function alert(title: string, message?: string, options?: OpenModalOptions): string {
		return openModal({
			type: 'alert',
			title,
			message,
			confirmText: options?.confirmText || 'OK',
			showCancel: false
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
		get state() {
			return state;
		},
		setLoading,
		clearLoading,
		addNotification,
		removeNotification,
		clearNotifications,
		success,
		error,
		warning,
		info,
		openModal,
		closeModal,
		closeAllModals,
		confirm,
		alert,
		reset
	};
}

/**
 * Export UI store instance
 * Only create store instance on client side to avoid SSR issues
 */
let uiStoreInstance: ReturnType<typeof createUIStore> | null = null;

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
