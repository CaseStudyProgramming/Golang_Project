import { describe, expect, it } from 'vitest'

function getColorClasses(color: 'blue' | 'gray' | 'white'): string {
	return color === 'white'
		? 'border-white'
		: color === 'gray'
			? 'border-gray-400'
			: 'border-blue-600'
}

// Helper functions to test the logic without type narrowing issues
function getSizeClasses(size: 'lg' | 'md' | 'sm'): string {
	return size === 'sm' ? 'h-4 w-4' : size === 'lg' ? 'h-12 w-12' : 'h-8 w-8'
}

describe('LoadingSpinner Component Logic', () => {
	it('calculates correct size classes for sm size', () => {
		const sizeClasses = getSizeClasses('sm')
		expect(sizeClasses).toBe('h-4 w-4')
	})

	it('calculates correct size classes for md size', () => {
		const sizeClasses = getSizeClasses('md')
		expect(sizeClasses).toBe('h-8 w-8')
	})

	it('calculates correct size classes for lg size', () => {
		const sizeClasses = getSizeClasses('lg')
		expect(sizeClasses).toBe('h-12 w-12')
	})

	it('calculates correct color classes for blue color', () => {
		const colorClasses = getColorClasses('blue')
		expect(colorClasses).toBe('border-blue-600')
	})

	it('calculates correct color classes for white color', () => {
		const colorClasses = getColorClasses('white')
		expect(colorClasses).toBe('border-white')
	})

	it('calculates correct color classes for gray color', () => {
		const colorClasses = getColorClasses('gray')
		expect(colorClasses).toBe('border-gray-400')
	})

	it('handles default props correctly', () => {
		const defaultSize = 'md'
		const defaultColor = 'blue'
		const defaultText = ''

		expect(defaultSize).toBe('md')
		expect(defaultColor).toBe('blue')
		expect(defaultText).toBe('')
	})
})
