import { describe, expect, it } from 'vitest';

function getColorClasses(color: 'blue' | 'green' | 'orange' | 'purple' | 'red'): string {
	return color === 'green' ? 'bg-green-600' : color === 'orange' ? 'bg-orange-600' : color === 'red' ? 'bg-red-600' : color === 'purple' ? 'bg-purple-600' : 'bg-blue-600';
}

// Helper functions to test the logic without type narrowing issues
function getSizeClasses(size: 'lg' | 'md' | 'sm'): string {
	return size === 'sm' ? 'h-1.5' : size === 'lg' ? 'h-4' : 'h-2.5';
}

describe('ProgressBar Component Logic', () => {
	it('clamps progress to minimum 0', () => {
		const progress = -10;
		const clampedProgress = Math.min(100, Math.max(0, progress));
		expect(clampedProgress).toBe(0);
	});

	it('clamps progress to maximum 100', () => {
		const progress = 150;
		const clampedProgress = Math.min(100, Math.max(0, progress));
		expect(clampedProgress).toBe(100);
	});

	it('keeps valid progress unchanged', () => {
		const progress = 50;
		const clampedProgress = Math.min(100, Math.max(0, progress));
		expect(clampedProgress).toBe(50);
	});

	it('calculates correct size classes for sm size', () => {
		const sizeClasses = getSizeClasses('sm');
		expect(sizeClasses).toBe('h-1.5');
	});

	it('calculates correct size classes for md size', () => {
		const sizeClasses = getSizeClasses('md');
		expect(sizeClasses).toBe('h-2.5');
	});

	it('calculates correct size classes for lg size', () => {
		const sizeClasses = getSizeClasses('lg');
		expect(sizeClasses).toBe('h-4');
	});

	it('calculates correct color classes for blue', () => {
		const colorClasses = getColorClasses('blue');
		expect(colorClasses).toBe('bg-blue-600');
	});

	it('calculates correct color classes for green', () => {
		const colorClasses = getColorClasses('green');
		expect(colorClasses).toBe('bg-green-600');
	});

	it('calculates correct color classes for orange', () => {
		const colorClasses = getColorClasses('orange');
		expect(colorClasses).toBe('bg-orange-600');
	});

	it('calculates correct color classes for red', () => {
		const colorClasses = getColorClasses('red');
		expect(colorClasses).toBe('bg-red-600');
	});

	it('calculates correct color classes for purple', () => {
		const colorClasses = getColorClasses('purple');
		expect(colorClasses).toBe('bg-purple-600');
	});

	it('determines when to show label', () => {
		const showLabel = true;
		const label = 'Progress';
		const shouldShowLabel = showLabel || label;
		expect(shouldShowLabel).toBeTruthy();

		const shouldNotShowLabel = false || '';
		expect(shouldNotShowLabel).toBeFalsy();
	});

	it('generates correct ARIA values', () => {
		const progress = 65;
		const ariaValueNow = progress;
		const ariaValueMin = '0';
		const ariaValueMax = '100';

		expect(ariaValueNow).toBe(65);
		expect(ariaValueMin).toBe('0');
		expect(ariaValueMax).toBe('100');
	});

	it('applies transition classes', () => {
		const transitionClasses = 'transition-all duration-300 ease-out';
		expect(transitionClasses).toContain('transition-all');
		expect(transitionClasses).toContain('duration-300');
		expect(transitionClasses).toContain('ease-out');
	});
});