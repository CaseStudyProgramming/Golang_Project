import { expect, test } from '@playwright/test'
import { AuthHelper } from './helpers/auth-helper'

test.describe('Task Management Critical Flows', () => {
	test.beforeAll(async ({ request }) => {
		// Register test user before running tests
		await AuthHelper.registerTestUser(request)
	})

	test.beforeEach(async ({ page }) => {
		// Login with test user before each test
		await AuthHelper.loginTestUser(page)
	})

	test.afterAll(async ({ request }) => {
		// Cleanup test user data after all tests
		await AuthHelper.cleanupTestUser(request)
	})

	test.describe('Create Task Flow', () => {
		test('should create a new task and view it', async ({ page }) => {
			// Already on tasks page from beforeEach

			// Click create task button
			await page.click('[data-testid="create-task-button"]')

			// Wait for modal to appear
			await page.waitForSelector('.fixed.inset-0.bg-black.bg-opacity-50', { state: 'visible' })

			// Fill in task details
			await page.fill('[data-testid="task-title-input"]', 'E2E Test Task')
			await page.fill(
				'[data-testid="task-description-input"]',
				'This is a test task created by E2E tests'
			)
			await page.selectOption('[data-testid="task-priority-select"]', 'high')

			// Submit the form
			await page.click('[data-testid="save-task-button"]')

			// Wait for modal to close
			await page
				.waitForSelector('.fixed.inset-0.bg-black.bg-opacity-50', { state: 'hidden' })
				.catch(() => {
					// Modal might be removed instead of hidden
					return page.waitForTimeout(1000)
				})

			// Wait for task to appear in the list (without navigation)
			await page.waitForTimeout(2000)

			// Verify task was created - wait for task list to be visible
			await page.waitForSelector('[data-testid="task-list"]')

			// Check if the task appears in the list
			const taskExists = await page
				.locator('[data-testid="task-list"]')
				.getByText('E2E Test Task')
				.count()
			expect(taskExists).toBeGreaterThan(0)
		})
	})

	test.describe('Update Task Flow', () => {
		test('should update an existing task and verify changes', async ({ page }) => {
			// Already on tasks page from beforeEach

			// Wait for task list to load
			await page.waitForSelector('[data-testid="task-list"]')

			// Click on first task's view button (eye icon)
			const firstTask = page.locator('[data-testid="task-item"]').first()
			await firstTask.locator('button[aria-label="View task"]').click()

			// Wait for task detail page to load
			await page.waitForURL(/\/tasks\/.+/)
			await page.waitForLoadState('networkidle')

			// Click edit button
			await page.click('text=Edit Task')

			// Wait for edit form to appear
			await page.waitForSelector('[data-testid="task-title-input"]', { state: 'visible' })

			// Update task title
			await page.fill('[data-testid="task-title-input"]', 'Updated E2E Test Task')

			// Submit the form
			await page.click('[data-testid="save-task-button"]')

			// Wait for edit mode to close
			await page.waitForTimeout(2000)

			// Verify the updated title is visible
			await expect(page.locator('text=Updated E2E Test Task')).toBeVisible()
		})
	})

	test.describe('Delete Task Flow', () => {
		test('should delete a task and verify it is removed', async ({ page }) => {
			// Already on tasks page from beforeEach

			// Wait for task list to load
			await page.waitForSelector('[data-testid="task-list"]')

			// Get initial task count
			const initialTaskCount = await page.locator('[data-testid="task-item"]').count()

			// Click on first task's view button
			const firstTask = page.locator('[data-testid="task-item"]').first()
			await firstTask.locator('button[aria-label="View task"]').click()

			// Wait for task detail page to load
			await page.waitForURL(/\/tasks\/.+/)
			await page.waitForLoadState('networkidle')

			// Click delete button
			await page.click('text=Delete Task')

			// Confirm deletion
			await page.click('text=Delete')

			// Wait for redirect back to tasks page
			await page.waitForURL('/tasks')
			await page.waitForLoadState('networkidle')

			// Verify task was removed
			const finalTaskCount = await page.locator('[data-testid="task-item"]').count()
			expect(finalTaskCount).toBeLessThan(initialTaskCount)
		})
	})
})
