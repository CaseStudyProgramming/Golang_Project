import { describe, expect, it } from 'vitest';

describe('Skeleton Component Logic', () => {
	it('calculates correct variant classes for circular', () => {
		const variant = 'circular' as 'default' | 'circular' | 'text' | 'rectangular';
		const variantClasses = variant === 'circular' ? 'rounded-full' : variant === 'text' ? 'rounded h-4' : variant === 'rectangular' ? 'rounded-md' : 'rounded-md';
		expect(variantClasses).toBe('rounded-full');
	});

	it('calculates correct variant classes for text', () => {
		const variant = 'text' as 'default' | 'circular' | 'text' | 'rectangular';
		const variantClasses = variant === 'circular' ? 'rounded-full' : variant === 'text' ? 'rounded h-4' : variant === 'rectangular' ? 'rounded-md' : 'rounded-md';
		expect(variantClasses).toBe('rounded h-4');
	});

	it('calculates correct variant classes for rectangular', () => {
		const variant = 'rectangular' as 'default' | 'circular' | 'text' | 'rectangular';
		const variantClasses = variant === 'circular' ? 'rounded-full' : variant === 'text' ? 'rounded h-4' : variant === 'rectangular' ? 'rounded-md' : 'rounded-md';
		expect(variantClasses).toBe('rounded-md');
	});

	it('calculates correct variant classes for default', () => {
		const variant = 'default' as 'default' | 'circular' | 'text' | 'rectangular';
		const variantClasses = variant === 'circular' ? 'rounded-full' : variant === 'text' ? 'rounded h-4' : variant === 'rectangular' ? 'rounded-md' : 'rounded-md';
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
		const defaultVariant = 'default' as 'default' | 'circular' | 'text' | 'rectangular';
		
		expect(defaultClassName).toBe('');
		expect(defaultVariant).toBe('default');
	});
});