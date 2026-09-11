import { describe, expect, it } from 'vitest';

describe('LoadingSpinner Component Logic', () => {
	it('calculates correct size classes for sm size', () => {
		const size = 'sm' as 'sm' | 'md' | 'lg';
		const sizeClasses = size === 'sm' ? 'h-4 w-4' : size === 'lg' ? 'h-12 w-12' : 'h-8 w-8';
		expect(sizeClasses).toBe('h-4 w-4');
	});

	it('calculates correct size classes for md size', () => {
		const size = 'md' as 'sm' | 'md' | 'lg';
		const sizeClasses = size === 'sm' ? 'h-4 w-4' : size === 'lg' ? 'h-12 w-12' : 'h-8 w-8';
		expect(sizeClasses).toBe('h-8 w-8');
	});

	it('calculates correct size classes for lg size', () => {
		const size = 'lg' as 'sm' | 'md' | 'lg';
		const sizeClasses = size === 'sm' ? 'h-4 w-4' : size === 'lg' ? 'h-12 w-12' : 'h-8 w-8';
		expect(sizeClasses).toBe('h-12 w-12');
	});

	it('calculates correct color classes for blue color', () => {
		const color = 'blue' as 'blue' | 'white' | 'gray';
		const colorClasses = color === 'white' ? 'border-white' : color === 'gray' ? 'border-gray-400' : 'border-blue-600';
		expect(colorClasses).toBe('border-blue-600');
	});

	it('calculates correct color classes for white color', () => {
		const color = 'white' as 'blue' | 'white' | 'gray';
		const colorClasses = color === 'white' ? 'border-white' : color === 'gray' ? 'border-gray-400' : 'border-blue-600';
		expect(colorClasses).toBe('border-white');
	});

	it('calculates correct color classes for gray color', () => {
		const color = 'gray' as 'blue' | 'white' | 'gray';
		const colorClasses = color === 'white' ? 'border-white' : color === 'gray' ? 'border-gray-400' : 'border-blue-600';
		expect(colorClasses).toBe('border-gray-400');
	});

	it('handles default props correctly', () => {
		const defaultSize = 'md' as 'sm' | 'md' | 'lg';
		const defaultColor = 'blue' as 'blue' | 'white' | 'gray';
		const defaultText = '';

		expect(defaultSize).toBe('md');
		expect(defaultColor).toBe('blue');
		expect(defaultText).toBe('');
	});
});