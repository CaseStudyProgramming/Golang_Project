/**
 * UI-related types and interfaces
 */

/**
 * Loading state interface
 */
export interface LoadingState {
	isLoading: boolean
	message?: string
}

/**
 * Modal interface
 */
export interface Modal {
	id: string
	type: ModalType
	title: string
	message?: string
	content?: string
	isOpen: boolean
	onConfirm?: () => Promise<void> | void
	onCancel?: () => void
	confirmText?: string
	cancelText?: string
	showCancel?: boolean
}

/**
 * Modal type
 */
export type ModalType = 'alert' | 'confirm' | 'custom' | 'prompt'

/**
 * Notification interface
 */
export interface Notification {
	id: string
	type: NotificationType
	title: string
	message?: string
	duration?: number
	isPersistent?: boolean
	createdAt: number
}

/**
 * Notification type
 */
export type NotificationType = 'error' | 'info' | 'success' | 'warning'
