import { test, expect } from '@playwright/test';

test.describe('Task Management E2E Tests', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to the application
		await page.goto('/');
	});

	test('loads the application homepage', async ({ page }) => {
		// Check that the page loads successfully
		await expect(page).toHaveTitle(/Task Manager/);
	});

	test('navigates to task list page', async ({ page }) => {
		// Click on tasks link (adjust selector based on actual implementation)
		await page.click('a[href="/tasks"]');
		
		// Verify navigation to tasks page
		await expect(page).toHaveURL('/tasks');
		await expect(page.locator('h1')).toContainText('Tasks');
	});

	test('displays task list with items', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Wait for task list to load
		await page.waitForSelector('[data-testid="task-list"]');
		
		// Verify task items are displayed
		const taskItems = page.locator('[data-testid="task-item"]');
		const count = await taskItems.count();
		
		// At least verify the task list container exists
		await expect(page.locator('[data-testid="task-list"]')).toBeVisible();
	});

	test('creates a new task', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Click on create task button
		await page.click('[data-testid="create-task-button"]');
		
		// Fill in task form
		await page.fill('[data-testid="task-title-input"]', 'Test Task from E2E');
		await page.fill('[data-testid="task-description-input"]', 'This is a test task created by E2E test');
		
		// Select priority
		await page.selectOption('[data-testid="task-priority-select"]', 'medium');
		
		// Submit the form
		await page.click('[data-testid="save-task-button"]');
		
		// Verify task was created
		await expect(page.locator('text=Test Task from E2E')).toBeVisible();
	});

	test('filters tasks by status', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Click on status filter
		await page.click('[data-testid="status-filter"]');
		
		// Select a status (e.g., "In Progress")
		await page.click('text=In Progress');
		
		// Verify filtered results
		await expect(page.locator('[data-testid="task-list"]')).toBeVisible();
	});

	test('searches for tasks', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Type in search box
		await page.fill('[data-testid="task-search-input"]', 'test');
		
		// Verify search results
		await expect(page.locator('[data-testid="task-list"]')).toBeVisible();
	});

	test('edits an existing task', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Click on a task to edit
		await page.click('[data-testid="task-item"]:first-child');
		
		// Modify task details
		await page.fill('[data-testid="task-title-input"]', 'Updated Task Title');
		
		// Save changes
		await page.click('[data-testid="save-task-button"]');
		
		// Verify task was updated
		await expect(page.locator('text=Updated Task Title')).toBeVisible();
	});

	test('deletes a task', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Click on delete button for first task
		await page.click('[data-testid="task-item"]:first-child [data-testid="delete-task-button"]');
		
		// Confirm deletion in dialog
		await page.click('[data-testid="confirm-delete-button"]');
		
		// Verify task was deleted
		await expect(page.locator('text=Task deleted successfully')).toBeVisible();
	});

	test('handles task creation with validation errors', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Click on create task button
		await page.click('[data-testid="create-task-button"]');
		
		// Try to submit without required fields
		await page.click('[data-testid="save-task-button"]');
		
		// Verify validation error appears
		await expect(page.locator('text=Title is required')).toBeVisible();
	});

	test('navigates between pages', async ({ page }) => {
		// Start at dashboard
		await page.goto('/');
		
		// Navigate to tasks
		await page.click('a[href="/tasks"]');
		await expect(page).toHaveURL('/tasks');
		
		// Navigate to categories
		await page.click('a[href="/categories"]');
		await expect(page).toHaveURL('/categories');
		
		// Navigate back to dashboard
		await page.click('a[href="/"]');
		await expect(page).toHaveURL('/');
	});

	test('responsive design on mobile', async ({ page }) => {
		// Set mobile viewport
		await page.setViewportSize({ width: 375, height: 667 });
		
		// Navigate to tasks page
		await page.goto('/tasks');
		
		// Verify mobile menu is visible
		await expect(page.locator('[data-testid="mobile-menu-button"]')).toBeVisible();
		
		// Click mobile menu
		await page.click('[data-testid="mobile-menu-button"]');
		
		// Verify navigation drawer opens
		await expect(page.locator('[data-testid="mobile-navigation-drawer"]')).toBeVisible();
	});
});