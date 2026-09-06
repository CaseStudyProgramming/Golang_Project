/**
 * UI-related types and interfaces
 */

/**
 * Notification type
 */
export type NotificationType = 'success' | 'error' | 'warning' | 'info';

/**
 * Notification interface
 */
export interface Notification {
	id: string;
	type: NotificationType;
	title: string;
	message?: string;
	duration?: number;
	isPersistent?: boolean;
	createdAt: number;
}

/**
 * Modal type
 */
export type ModalType = 'confirm' | 'alert' | 'prompt' | 'custom';

/**
 * Modal interface
 */
export interface Modal {
	id: string;
	type: ModalType;
	title: string;
	message?: string;
	content?: string;
	isOpen: boolean;
	onConfirm?: () => void | Promise<void>;
	onCancel?: () => void;
	confirmText?: string;
	cancelText?: string;
	showCancel?: boolean;
}

/**
 * Loading state interface
 */
export interface LoadingState {
	isLoading: boolean;
	message?: string;
}
