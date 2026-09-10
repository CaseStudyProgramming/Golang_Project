import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';

import LoadingSpinner from './LoadingSpinner.svelte';

describe('LoadingSpinner Component', () => {
	describe('Rendering', () => {
		it('renders spinner with default props', () => {
			render(LoadingSpinner);
			const spinner = screen.getByRole('status');
			expect(spinner).toBeInTheDocument();
		});

		it('renders with custom text', () => {
			render(LoadingSpinner, { text: 'Loading data...' });
			expect(screen.getByText('Loading data...')).toBeInTheDocument();
		});

		it('renders without text when not provided', () => {
			render(LoadingSpinner);
			const textElement = screen.queryByText(/.+/);
			expect(textElement).not.toBeInTheDocument();
		});
	});

	describe('Size Variations', () => {
		it('renders small spinner', () => {
			const { container } = render(LoadingSpinner, { size: 'sm' });
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('h-4', 'w-4');
		});

		it('renders medium spinner by default', () => {
			const { container } = render(LoadingSpinner);
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('h-8', 'w-8');
		});

		it('renders large spinner', () => {
			const { container } = render(LoadingSpinner, { size: 'lg' });
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('h-12', 'w-12');
		});
	});

	describe('Color Variations', () => {
		it('renders blue spinner by default', () => {
			const { container } = render(LoadingSpinner);
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('border-blue-600');
		});

		it('renders white spinner', () => {
			const { container } = render(LoadingSpinner, { color: 'white' });
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('border-white');
		});

		it('renders gray spinner', () => {
			const { container } = render(LoadingSpinner, { color: 'gray' });
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('border-gray-400');
		});
	});

	describe('Accessibility', () => {
		it('has aria-live region for screen readers', () => {
			render(LoadingSpinner);
			const spinner = screen.getByRole('status');
			expect(spinner).toHaveAttribute('aria-live', 'polite');
		});

		it('includes text in aria-label when provided', () => {
			render(LoadingSpinner, { text: 'Loading tasks...' });
			const spinner = screen.getByRole('status');
			expect(spinner).toHaveTextContent('Loading tasks...');
		});
	});

	describe('Layout', () => {
		it('centers content vertically and horizontally', () => {
			const { container } = render(LoadingSpinner);
			const wrapper = container.querySelector('.flex');
			expect(wrapper).toHaveClass('items-center', 'justify-center');
		});

		it('uses flex column layout', () => {
			const { container } = render(LoadingSpinner);
			const wrapper = container.querySelector('.flex');
			expect(wrapper).toHaveClass('flex-col');
		});
	});

	describe('Animation', () => {
		it('has spin animation', () => {
			const { container } = render(LoadingSpinner);
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('animate-spin');
		});

		it('has rounded border', () => {
			const { container } = render(LoadingSpinner);
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('rounded-full');
		});

		it('has transparent top border for animation effect', () => {
			const { container } = render(LoadingSpinner);
			const spinner = container.querySelector('.animate-spin');
			expect(spinner).toHaveClass('border-t-transparent');
		});
	});
});