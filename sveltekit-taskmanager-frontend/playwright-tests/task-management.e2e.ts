import { expect, test } from '@playwright/test'

test.describe('Task Management Critical Flows', () => {
	test.describe('Create Task Flow', () => {
		test('should create a new task and view it', async ({ page }) => {
			// Navigate to tasks page
			await page.goto('/tasks')

			// Wait for page to load
			await page.waitForLoadState('networkidle')

			// Click create task button
			await page.click('[data-testid="create-task-button"]')

			// Fill in task details
			await page.fill('[data-testid="task-title-input"]', 'E2E Test Task')
			await page.fill(
				'[data-testid="task-description-input"]',
				'This is a test task created by E2E tests'
			)
			await page.selectOption('[data-testid="task-priority-select"]', 'high')

			// Submit the form
			await page.click('[data-testid="save-task-button"]')

			// Wait for modal to close and task to appear
			await page.waitForTimeout(2000)

			// Navigate to tasks page to verify
			await page.goto('/tasks')
			await page.waitForLoadState('networkidle')

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
			// Navigate to tasks page
			await page.goto('/tasks')
			await page.waitForLoadState('networkidle')

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
			// Navigate to tasks page
			await page.goto('/tasks')
			await page.waitForLoadState('networkidle')

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
