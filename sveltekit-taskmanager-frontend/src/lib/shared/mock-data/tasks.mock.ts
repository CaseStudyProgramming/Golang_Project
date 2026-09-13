// Mock data for e2e testing - shared across API routes
let mockTasks = [
	{
		id: '1',
		title: 'Test Task 1',
		description: 'This is a test task',
		status: 'todo',
		priority: 'medium',
		categoryId: 'cat1',
		tags: ['tag1'],
		dueDate: '2024-12-31',
		createdAt: '2024-01-01T00:00:00Z',
		updatedAt: '2024-01-01T00:00:00Z',
	},
	{
		id: '2',
		title: 'Test Task 2',
		description: 'Another test task',
		status: 'in_progress',
		priority: 'high',
		categoryId: 'cat2',
		tags: ['tag2'],
		dueDate: '2024-12-25',
		createdAt: '2024-01-02T00:00:00Z',
		updatedAt: '2024-01-02T00:00:00Z',
	},
]

export function getMockTasks() {
	return mockTasks
}

export function updateMockTasks(newTasks: typeof mockTasks) {
	mockTasks = newTasks
}

export function resetMockData() {
	mockTasks = [
		{
			id: '1',
			title: 'Test Task 1',
			description: 'This is a test task',
			status: 'todo',
			priority: 'medium',
			categoryId: 'cat1',
			tags: ['tag1'],
			dueDate: '2024-12-31',
			createdAt: '2024-01-01T00:00:00Z',
			updatedAt: '2024-01-01T00:00:00Z',
		},
		{
			id: '2',
			title: 'Test Task 2',
			description: 'Another test task',
			status: 'in_progress',
			priority: 'high',
			categoryId: 'cat2',
			tags: ['tag2'],
			dueDate: '2024-12-25',
			createdAt: '2024-01-02T00:00:00Z',
			updatedAt: '2024-01-02T00:00:00Z',
		},
	]
}
