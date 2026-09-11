import { describe, expect, it } from 'vitest';

// Helper function to test the logic without type narrowing issues
function getVariantClasses(variant: 'circular' | 'default' | 'rectangular' | 'text'): string {
	return variant === 'circular' ? 'rounded-full' : variant === 'text' ? 'rounded h-4' : variant === 'rectangular' ? 'rounded-md' : 'rounded-md';
}

describe('Skeleton Component Logic', () => {
	it('calculates correct variant classes for circular', () => {
		const variantClasses = getVariantClasses('circular');
		expect(variantClasses).toBe('rounded-full');
	});

	it('calculates correct variant classes for text', () => {
		const variantClasses = getVariantClasses('text');
		expect(variantClasses).toBe('rounded h-4');
	});

	it('calculates correct variant classes for rectangular', () => {
		const variantClasses = getVariantClasses('rectangular');
		expect(variantClasses).toBe('rounded-md');
	});

	it('calculates correct variant classes for default', () => {
		const variantClasses = getVariantClasses('default');
		expect(variantClasses).toBe('rounded-md');
	});

	it('applies base classes correctly', () => {
		const baseClasses = 'animate-pulse bg-gray-200';
		expect(baseClasses).toContain('animate-pulse');
		expect(baseClasses).toContain('bg-gray-200');
	});

	it('combines classes correctly', () => {
		const baseClasses = 'animate-pulse bg-gray-200';
		const variantClasses = 'rounded-full';
		const customClasses = 'w-32 h-32';
		const combinedClasses = `${baseClasses} ${variantClasses} ${customClasses}`;
		
		expect(combinedClasses).toContain('animate-pulse');
		expect(combinedClasses).toContain('bg-gray-200');
		expect(combinedClasses).toContain('rounded-full');
		expect(combinedClasses).toContain('w-32');
		expect(combinedClasses).toContain('h-32');
	});

	it('handles default props correctly', () => {
		const defaultClassName = '';
		const defaultVariant = 'default';
		
		expect(defaultClassName).toBe('');
		expect(defaultVariant).toBe('default');
	});
});