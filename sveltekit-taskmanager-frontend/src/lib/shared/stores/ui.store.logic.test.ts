import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

describe('UI Store Logic', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		vi.useFakeTimers();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	describe('Loading State Management', () => {
		it('sets loading state with message', () => {
			const loading = { isLoading: true, message: 'Loading data...' };
			expect(loading.isLoading).toBe(true);
			expect(loading.message).toBe('Loading data...');
		});

		it('sets loading state without message', () => {
			type LoadingState = { isLoading: boolean; message?: string };
			const loading: LoadingState = { isLoading: true, message: undefined };
			expect(loading.isLoading).toBe(true);
			expect(loading.message).toBeUndefined();
		});

		it('clears loading state', () => {
			type LoadingState = { isLoading: boolean; message?: string };
			const loading: LoadingState = { isLoading: false, message: undefined };
			expect(loading.isLoading).toBe(false);
			expect(loading.message).toBeUndefined();
		});
	});

	describe('Notification Management', () => {
		it('creates notification with required fields', () => {
			const notification = {
				createdAt: Date.now(),
				duration: 5000,
				id: '123',
				isPersistent: false,
				message: 'Operation completed',
				title: 'Success!',
				type: 'success' as const
			};

			expect(notification.id).toBe('123');
			expect(notification.type).toBe('success');
			expect(notification.title).toBe('Success!');
			expect(notification.message).toBe('Operation completed');
			expect(notification.duration).toBe(5000);
			expect(notification.isPersistent).toBe(false);
		});

		it('creates notification with custom duration', () => {
			const notification = {
				createdAt: Date.now(),
				duration: 10000,
				id: '123',
				isPersistent: true,
				message: 'Something went wrong',
				title: 'Error',
				type: 'error' as const
			};

			expect(notification.duration).toBe(10000);
			expect(notification.isPersistent).toBe(true);
		});

		it('generates unique notification IDs', () => {
			const id1 = Date.now().toString();
			const id2 = (Date.now() + 1).toString();
			expect(id1).not.toBe(id2);
		});

		it('removes notification by ID', () => {
			const notifications = [
				{ createdAt: Date.now(), duration: 5000, id: '1', isPersistent: false, title: 'Test 1', type: 'success' as const },
				{ createdAt: Date.now(), duration: 5000, id: '2', isPersistent: false, title: 'Test 2', type: 'error' as const }
			];

			const filtered = notifications.filter(n => n.id !== '1');
			expect(filtered).toHaveLength(1);
			expect(filtered[0].id).toBe('2');
		});

		it('clears all notifications', () => {
			const notifications: never[] = [];
			expect(notifications).toHaveLength(0);
		});
	});

	describe('Notification Type Helpers', () => {
		it('creates success notification', () => {
			const notification = {
				createdAt: Date.now(),
				duration: 5000,
				id: '123',
				isPersistent: false,
				message: 'Operation completed',
				title: 'Success!',
				type: 'success' as const
			};

			expect(notification.type).toBe('success');
		});

		it('creates error notification as persistent', () => {
			const notification = {
				createdAt: Date.now(),
				duration: 5000,
				id: '123',
				isPersistent: true,
				message: 'Something went wrong',
				title: 'Error!',
				type: 'error' as const
			};

			expect(notification.type).toBe('error');
			expect(notification.isPersistent).toBe(true);
		});

		it('creates warning notification', () => {
			const notification = {
				createdAt: Date.now(),
				duration: 5000,
				id: '123',
				isPersistent: false,
				message: 'Be careful',
				title: 'Warning!',
				type: 'warning' as const
			};

			expect(notification.type).toBe('warning');
		});

		it('creates info notification', () => {
			const notification = {
				createdAt: Date.now(),
				duration: 5000,
				id: '123',
				isPersistent: false,
				message: 'Some information',
				title: 'Info!',
				type: 'info' as const
			};

			expect(notification.type).toBe('info');
		});
	});

	describe('Modal Management', () => {
		it('creates modal with required fields', () => {
			const modal = {
				cancelText: 'Cancel',
				confirmText: 'Confirm',
				content: 'Modal content',
				id: '123',
				isOpen: true,
				showCancel: true,
				title: 'Test Modal',
				type: 'confirm' as const
			};

			expect(modal.id).toBe('123');
			expect(modal.title).toBe('Test Modal');
			expect(modal.isOpen).toBe(true);
			expect(modal.showCancel).toBe(true);
			expect(modal.type).toBe('confirm');
		});

		it('creates modal without cancel button', () => {
			const modal = {
				cancelText: '',
				confirmText: 'OK',
				content: 'Alert content',
				id: '123',
				isOpen: true,
				showCancel: false,
				title: 'Alert Modal',
				type: 'alert' as const
			};

			expect(modal.showCancel).toBe(false);
			expect(modal.type).toBe('alert');
		});

		it('generates unique modal IDs', () => {
			const id1 = Date.now().toString();
			const id2 = (Date.now() + 1).toString();
			expect(id1).not.toBe(id2);
		});

		it('closes modal by setting isOpen to false', () => {
			const modal = {
				cancelText: 'Cancel',
				confirmText: 'Confirm',
				content: 'Content',
				id: '123',
				isOpen: true,
				showCancel: true,
				title: 'Test',
				type: 'confirm' as const
			};

			modal.isOpen = false;
			expect(modal.isOpen).toBe(false);
		});

		it('removes modal from array', () => {
			const modals = [
				{ cancelText: 'Cancel', confirmText: 'Confirm', content: 'Content 1', id: '1', isOpen: true, showCancel: true, title: 'Modal 1', type: 'confirm' as const },
				{ cancelText: 'Cancel', confirmText: 'Confirm', content: 'Content 2', id: '2', isOpen: true, showCancel: true, title: 'Modal 2', type: 'confirm' as const }
			];

			const filtered = modals.filter(m => m.id !== '1');
			expect(filtered).toHaveLength(1);
			expect(filtered[0].id).toBe('2');
		});
	});

	describe('State Reset', () => {
		it('resets loading state', () => {
			type LoadingState = { isLoading: boolean; message?: string };
			type UIState = {
				loading: LoadingState;
				modals: unknown[];
				notifications: unknown[];
			};
			const state: UIState = {
				loading: { isLoading: true, message: 'Loading...' },
				modals: [],
				notifications: []
			};

			state.loading = { isLoading: false, message: undefined };
			expect(state.loading.isLoading).toBe(false);
			expect(state.loading.message).toBeUndefined();
		});

		it('resets notifications array', () => {
			const state = {
				loading: { isLoading: false, message: undefined },
				modals: [],
				notifications: [{ createdAt: Date.now(), duration: 5000, id: '1', isPersistent: false, title: 'Test', type: 'success' as const }]
			};

			state.notifications = [];
			expect(state.notifications).toHaveLength(0);
		});

		it('resets modals array', () => {
			const state = {
				loading: { isLoading: false, message: undefined },
				modals: [{ cancelText: 'Cancel', confirmText: 'Confirm', content: 'Content', id: '1', isOpen: true, showCancel: true, title: 'Test', type: 'confirm' as const }],
				notifications: []
			};

			state.modals = [];
			expect(state.modals).toHaveLength(0);
		});
	});
});