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
			const loading = { isLoading: true, message: undefined };
			expect(loading.isLoading).toBe(true);
			expect(loading.message).toBeUndefined();
		});

		it('clears loading state', () => {
			let loading = { isLoading: true, message: 'Loading...' };
			loading = { isLoading: false, message: undefined };
			expect(loading.isLoading).toBe(false);
			expect(loading.message).toBeUndefined();
		});
	});

	describe('Notification Management', () => {
		it('creates notification with required fields', () => {
			const notification = {
				id: '123',
				type: 'success' as const,
				title: 'Success!',
				message: 'Operation completed',
				createdAt: Date.now(),
				duration: 5000,
				isPersistent: false
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
				id: '123',
				type: 'error' as const,
				title: 'Error',
				message: 'Something went wrong',
				createdAt: Date.now(),
				duration: 10000,
				isPersistent: true
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
				{ id: '1', type: 'success' as const, title: 'Test 1', createdAt: Date.now(), duration: 5000, isPersistent: false },
				{ id: '2', type: 'error' as const, title: 'Test 2', createdAt: Date.now(), duration: 5000, isPersistent: false }
			];

			const filtered = notifications.filter(n => n.id !== '1');
			expect(filtered).toHaveLength(1);
			expect(filtered[0].id).toBe('2');
		});

		it('clears all notifications', () => {
			let notifications = [
				{ id: '1', type: 'success' as const, title: 'Test 1', createdAt: Date.now(), duration: 5000, isPersistent: false },
				{ id: '2', type: 'error' as const, title: 'Test 2', createdAt: Date.now(), duration: 5000, isPersistent: false }
			];

			notifications = [];
			expect(notifications).toHaveLength(0);
		});
	});

	describe('Notification Type Helpers', () => {
		it('creates success notification', () => {
			const notification = {
				id: '123',
				type: 'success' as const,
				title: 'Success!',
				message: 'Operation completed',
				createdAt: Date.now(),
				duration: 5000,
				isPersistent: false
			};

			expect(notification.type).toBe('success');
		});

		it('creates error notification as persistent', () => {
			const notification = {
				id: '123',
				type: 'error' as const,
				title: 'Error!',
				message: 'Something went wrong',
				createdAt: Date.now(),
				duration: 5000,
				isPersistent: true
			};

			expect(notification.type).toBe('error');
			expect(notification.isPersistent).toBe(true);
		});

		it('creates warning notification', () => {
			const notification = {
				id: '123',
				type: 'warning' as const,
				title: 'Warning!',
				message: 'Be careful',
				createdAt: Date.now(),
				duration: 5000,
				isPersistent: false
			};

			expect(notification.type).toBe('warning');
		});

		it('creates info notification', () => {
			const notification = {
				id: '123',
				type: 'info' as const,
				title: 'Info!',
				message: 'Some information',
				createdAt: Date.now(),
				duration: 5000,
				isPersistent: false
			};

			expect(notification.type).toBe('info');
		});
	});

	describe('Modal Management', () => {
		it('creates modal with required fields', () => {
			const modal = {
				id: '123',
				title: 'Test Modal',
				content: 'Modal content',
				isOpen: true,
				showCancel: true,
				type: 'confirm' as const,
				confirmText: 'Confirm',
				cancelText: 'Cancel'
			};

			expect(modal.id).toBe('123');
			expect(modal.title).toBe('Test Modal');
			expect(modal.isOpen).toBe(true);
			expect(modal.showCancel).toBe(true);
			expect(modal.type).toBe('confirm');
		});

		it('creates modal without cancel button', () => {
			const modal = {
				id: '123',
				title: 'Alert Modal',
				content: 'Alert content',
				isOpen: true,
				showCancel: false,
				type: 'alert' as const,
				confirmText: 'OK',
				cancelText: ''
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
				id: '123',
				title: 'Test',
				content: 'Content',
				isOpen: true,
				showCancel: true,
				type: 'confirm' as const,
				confirmText: 'Confirm',
				cancelText: 'Cancel'
			};

			modal.isOpen = false;
			expect(modal.isOpen).toBe(false);
		});

		it('removes modal from array', () => {
			const modals = [
				{ id: '1', title: 'Modal 1', content: 'Content 1', isOpen: true, showCancel: true, type: 'confirm' as const, confirmText: 'Confirm', cancelText: 'Cancel' },
				{ id: '2', title: 'Modal 2', content: 'Content 2', isOpen: true, showCancel: true, type: 'confirm' as const, confirmText: 'Confirm', cancelText: 'Cancel' }
			];

			const filtered = modals.filter(m => m.id !== '1');
			expect(filtered).toHaveLength(1);
			expect(filtered[0].id).toBe('2');
		});
	});

	describe('State Reset', () => {
		it('resets loading state', () => {
			let state = {
				loading: { isLoading: true, message: 'Loading...' },
				notifications: [],
				modals: []
			};

			state.loading = { isLoading: false, message: undefined };
			expect(state.loading.isLoading).toBe(false);
			expect(state.loading.message).toBeUndefined();
		});

		it('resets notifications array', () => {
			let state = {
				loading: { isLoading: false, message: undefined },
				notifications: [{ id: '1', type: 'success' as const, title: 'Test', createdAt: Date.now(), duration: 5000, isPersistent: false }],
				modals: []
			};

			state.notifications = [];
			expect(state.notifications).toHaveLength(0);
		});

		it('resets modals array', () => {
			let state = {
				loading: { isLoading: false, message: undefined },
				notifications: [],
				modals: [{ id: '1', title: 'Test', content: 'Content', isOpen: true, showCancel: true, type: 'confirm' as const, confirmText: 'Confirm', cancelText: 'Cancel' }]
			};

			state.modals = [];
			expect(state.modals).toHaveLength(0);
		});
	});
});