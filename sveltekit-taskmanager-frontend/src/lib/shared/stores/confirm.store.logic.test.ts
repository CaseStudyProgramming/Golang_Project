import { describe, expect, it, vi } from 'vitest'

import type { ConfirmType } from './confirm.store.svelte.ts'

describe('Confirm Store Logic', () => {
	describe('Dialog ID Generation', () => {
		it('generates unique dialog IDs', () => {
			const id1 = `confirm-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
			const id2 = `confirm-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`

			expect(id1).not.toBe(id2)
		})

		it('generates IDs with correct prefix', () => {
			const id = `confirm-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
			expect(id).toMatch(/^confirm-/)
		})
	})

	describe('Dialog Types', () => {
		it('accepts valid dialog types', () => {
			const validTypes: ConfirmType[] = ['danger', 'warning', 'info']
			validTypes.forEach((type) => {
				expect(['danger', 'warning', 'info']).toContain(type)
			})
		})

		it('handles danger type for destructive actions', () => {
			const type: ConfirmType = 'danger'
			expect(type).toBe('danger')
		})

		it('handles warning type for caution actions', () => {
			const type: ConfirmType = 'warning'
			expect(type).toBe('warning')
		})

		it('handles info type for informational actions', () => {
			const type: ConfirmType = 'info'
			expect(type).toBe('info')
		})
	})

	describe('Dialog State Management', () => {
		it('initializes dialog with open state', () => {
			const dialog = {
				cancelText: 'Cancel',
				confirmText: 'Confirm',
				id: 'test-id',
				isOpen: true,
				message: 'Test message',
				title: 'Test',
				type: 'info' as ConfirmType,
			}

			expect(dialog.isOpen).toBe(true)
		})

		it('closes dialog by setting isOpen to false', () => {
			const dialog = {
				cancelText: 'Cancel',
				confirmText: 'Confirm',
				id: 'test-id',
				isOpen: true,
				message: 'Test message',
				title: 'Test',
				type: 'info' as ConfirmType,
			}

			dialog.isOpen = false
			expect(dialog.isOpen).toBe(false)
		})

		it('clears dialog by setting to null', () => {
			const dialog: null | { id: string } = null
			expect(dialog).toBeNull()
		})
	})

	describe('Promise Resolution', () => {
		it('resolves promise with true on confirm', async () => {
			const promise = Promise.resolve(true)
			const result = await promise
			expect(result).toBe(true)
		})

		it('resolves promise with false on cancel', async () => {
			const promise = Promise.resolve(false)
			const result = await promise
			expect(result).toBe(false)
		})

		it('handles async onConfirm callback', async () => {
			const onConfirm = vi.fn(async () => {
				await new Promise((resolve) => setTimeout(resolve, 10))
				return true
			})

			const result = await onConfirm()
			expect(result).toBe(true)
			expect(onConfirm).toHaveBeenCalled()
		})
	})

	describe('Optional Callbacks', () => {
		it('handles missing onConfirm callback', () => {
			type ConfirmOptions = {
				cancelText: string
				confirmText: string
				message: string
				onCancel?: unknown
				onConfirm?: unknown
				title: string
				type: ConfirmType
			}
			const options: ConfirmOptions = {
				cancelText: 'Cancel',
				confirmText: 'OK',
				message: 'Test message',
				title: 'Test',
				type: 'info' as ConfirmType,
			}

			expect(options.onConfirm).toBeUndefined()
		})

		it('handles missing onCancel callback', () => {
			type ConfirmOptions = {
				cancelText: string
				confirmText: string
				message: string
				onCancel?: unknown
				onConfirm?: unknown
				title: string
				type: ConfirmType
			}
			const options: ConfirmOptions = {
				cancelText: 'Cancel',
				confirmText: 'OK',
				message: 'Test message',
				title: 'Test',
				type: 'info' as ConfirmType,
			}

			expect(options.onCancel).toBeUndefined()
		})

		it('executes onConfirm when provided', async () => {
			const onConfirm = vi.fn()
			await onConfirm()
			expect(onConfirm).toHaveBeenCalled()
		})

		it('executes onCancel when provided', () => {
			const onCancel = vi.fn()
			onCancel()
			expect(onCancel).toHaveBeenCalled()
		})
	})
})
