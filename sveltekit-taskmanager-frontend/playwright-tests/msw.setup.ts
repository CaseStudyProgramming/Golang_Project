import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

// Mock data with persistent storage across requests
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

const mockCategories = [
	{
		id: 'cat1',
		name: 'Work',
		icon: '💼',
		color: '#3B82F6',
		description: 'Work related tasks',
		createdAt: '2024-01-01T00:00:00Z',
		updatedAt: '2024-01-01T00:00:00Z',
	},
	{
		id: 'cat2',
		name: 'Personal',
		icon: '🏠',
		color: '#10B981',
		description: 'Personal tasks',
		createdAt: '2024-01-01T00:00:00Z',
		updatedAt: '2024-01-01T00:00:00Z',
	},
]

const mockTags = [
	{
		id: 'tag1',
		name: 'Important',
		color: '#EF4444',
		createdAt: '2024-01-01T00:00:00Z',
		updatedAt: '2024-01-01T00:00:00Z',
	},
	{
		id: 'tag2',
		name: 'Urgent',
		color: '#F59E0B',
		createdAt: '2024-01-01T00:00:00Z',
		updatedAt: '2024-01-01T00:00:00Z',
	},
]

// Create MSW server
export const server = setupServer(
	// Mock GET /api/tasks
	http.get('/api/tasks', () => {
		return HttpResponse.json({
			tasks: mockTasks,
			pagination: {
				page: 1,
				limit: 10,
				total: mockTasks.length,
				totalPages: 1,
			},
		})
	}),

	// Mock GET /api/tasks/:id
	http.get('/api/tasks/:id', ({ params }) => {
		const task = mockTasks.find((t) => t.id === params.id)
		if (!task) {
			return HttpResponse.json({ error: 'Task not found' }, { status: 404 })
		}
		return HttpResponse.json(task)
	}),

	// Mock POST /api/tasks
	http.post('/api/tasks', async ({ request }) => {
		const newTask = (await request.json()) as Partial<(typeof mockTasks)[0]>
		const createdTask = {
			...newTask,
			id: `task-${Date.now()}`,
			status: 'todo',
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
		} as (typeof mockTasks)[0]
		mockTasks.push(createdTask)
		return HttpResponse.json(createdTask, { status: 201 })
	}),

	// Mock PUT /api/tasks/:id
	http.put('/api/tasks/:id', async ({ params, request }) => {
		const updatedData = (await request.json()) as Partial<(typeof mockTasks)[0]>
		const taskIndex = mockTasks.findIndex((t) => t.id === params.id)
		if (taskIndex === -1) {
			return HttpResponse.json({ error: 'Task not found' }, { status: 404 })
		}
		const updatedTask = {
			...mockTasks[taskIndex],
			...updatedData,
			updatedAt: new Date().toISOString(),
		} as (typeof mockTasks)[0]
		mockTasks[taskIndex] = updatedTask
		return HttpResponse.json(updatedTask)
	}),

	// Mock DELETE /api/tasks/:id
	http.delete('/api/tasks/:id', ({ params }) => {
		const taskIndex = mockTasks.findIndex((t) => t.id === params.id)
		if (taskIndex === -1) {
			return HttpResponse.json({ error: 'Task not found' }, { status: 404 })
		}
		mockTasks.splice(taskIndex, 1)
		return HttpResponse.json({ message: 'Task deleted successfully' })
	}),

	// Mock GET /api/categories
	http.get('/api/categories', () => {
		return HttpResponse.json(mockCategories)
	}),

	// Mock GET /api/tags
	http.get('/api/tags', () => {
		return HttpResponse.json(mockTags)
	}),

	// Mock POST /api/tags
	http.post('/api/tags', async ({ request }) => {
		const newTag = (await request.json()) as Partial<(typeof mockTags)[0]>
		const createdTag = {
			...newTag,
			id: 'new-tag-id',
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
		} as (typeof mockTags)[0]
		return HttpResponse.json(createdTag, { status: 201 })
	}),

	// Mock authentication endpoints
	http.post('/api/auth/login', async ({ request }) => {
		const credentials = (await request.json()) as { email: string }
		// Accept any credentials for testing
		return HttpResponse.json({
			user: {
				id: 'user-1',
				email: credentials.email,
				name: 'Test User',
			},
			token: 'mock-jwt-token',
		})
	}),

	http.post('/api/auth/register', async ({ request }) => {
		const userData = (await request.json()) as { email: string; name: string }
		return HttpResponse.json(
			{
				user: {
					id: 'user-1',
					email: userData.email,
					name: userData.name,
				},
				token: 'mock-jwt-token',
			},
			{ status: 201 }
		)
	}),

	// Mock health check
	http.get('/api/health', () => {
		return HttpResponse.json({ status: 'ok' })
	})
)

// Server lifecycle
export const startServer = () => server.listen({ onUnhandledRequest: 'bypass' })
export const stopServer = () => server.close()
export const resetHandlers = () => server.resetHandlers()

// Reset mock data to initial state
export const resetMockData = () => {
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
