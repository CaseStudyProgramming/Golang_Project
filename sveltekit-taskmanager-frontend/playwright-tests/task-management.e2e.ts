import { expect, test } from '@playwright/test'

test.describe('Task Management E2E Tests', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to the application
		await page.goto('/')
	})

	test('loads the application homepage', async ({ page }) => {
		// Check that the page loads successfully
		await expect(page).toHaveTitle(/Task Manager/)
		await expect(page.locator('h1')).toContainText('Task Manager')
	})

	test('navigates to task list page', async ({ page }) => {
		// Click on tasks link
		await page.click('a[href="/tasks"]')

		// Verify navigation to tasks page
		await expect(page).toHaveURL('/tasks')
		await expect(page.locator('h1')).toContainText('Tasks')
	})

	test('displays task list with items', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks')

		// Wait for task list to load
		await page.waitForSelector('[data-testid="task-list"]')

		// Verify task list container exists
		await expect(page.locator('[data-testid="task-list"]')).toBeVisible()
	})

	test('navigates between pages', async ({ page }) => {
		// Start at homepage
		await page.goto('/')

		// Navigate to tasks
		await page.click('a[href="/tasks"]')
		await expect(page).toHaveURL('/tasks')

		// Navigate to categories
		await page.click('a[href="/categories"]')
		await expect(page).toHaveURL('/categories')

		// Navigate back to homepage
		await page.click('a[href="/"]')
		await expect(page).toHaveURL('/')
	})

	test('basic task form rendering', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks')

		// Click on create task button
		await page.click('[data-testid="create-task-button"]')

		// Verify form elements are present
		await expect(page.locator('[data-testid="task-title-input"]')).toBeVisible()
		await expect(page.locator('[data-testid="task-description-input"]')).toBeVisible()
		await expect(page.locator('[data-testid="task-priority-select"]')).toBeVisible()
		await expect(page.locator('[data-testid="save-task-button"]')).toBeVisible()
	})

	test('filter components are present', async ({ page }) => {
		// Navigate to tasks page
		await page.goto('/tasks')

		// Verify filter components exist
		await expect(page.locator('[data-testid="status-filter"]')).toBeVisible()
		await expect(page.locator('[data-testid="task-search-input"]')).toBeVisible()
	})
})
