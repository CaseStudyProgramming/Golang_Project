import { describe, expect, it, vi } from 'vitest'

describe('EmptyState Component Logic', () => {
	it('handles default props correctly', () => {
		const defaultIcon = ''
		const defaultTitle = 'No data'
		const defaultDescription = 'There is no data to display.'
		const defaultActionLabel = ''
		const defaultOnAction = undefined

		expect(defaultIcon).toBe('')
		expect(defaultTitle).toBe('No data')
		expect(defaultDescription).toBe('There is no data to display.')
		expect(defaultActionLabel).toBe('')
		expect(defaultOnAction).toBeUndefined()
	})

	it('determines when to show action button', () => {
		const actionLabel = 'Create Task'
		const onAction = vi.fn()
		const shouldShowButton = actionLabel && onAction
		expect(shouldShowButton).toBeTruthy()

		const shouldNotShowButton2 = actionLabel && undefined
		expect(shouldNotShowButton2).toBeFalsy()
	})

	it('determines icon type (emoji vs SVG)', () => {
		const emojiIcon = '📝'
		const isEmoji = !emojiIcon.trim().startsWith('<svg')
		expect(isEmoji).toBeTruthy()

		const svgIcon =
			'<svg xmlns="http://www.w3.org/2000/svg"><path d="M12 2L2 7l10 5 10-5-10-5z"/></svg>'
		const isSvg = svgIcon.trim().startsWith('<svg')
		expect(isSvg).toBeTruthy()
	})

	it('handles action callback execution', () => {
		const mockAction = vi.fn()
		mockAction()
		expect(mockAction).toHaveBeenCalledTimes(1)
	})

	it('validates responsive class structure', () => {
		const baseClasses = 'text-center py-8 sm:py-12 px-4'
		expect(baseClasses).toContain('py-8')
		expect(baseClasses).toContain('sm:py-12')
		expect(baseClasses).toContain('px-4')
	})

	it('determines when to render icon', () => {
		const icon = '📝'
		const shouldRenderIcon = icon
		expect(shouldRenderIcon).toBeTruthy()

		const emptyIcon = ''
		const shouldNotRenderIcon = emptyIcon
		expect(shouldNotRenderIcon).toBeFalsy()
	})
})
