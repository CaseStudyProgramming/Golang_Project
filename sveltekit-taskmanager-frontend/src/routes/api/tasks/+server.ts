import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'
import { getMockTasks, updateMockTasks } from '$lib/shared/mock-data/tasks.mock'

export const GET: RequestHandler = async ({ url }) => {
	const mockTasks = getMockTasks()
	const page = parseInt(url.searchParams.get('page') || '1', 10)
	const limit = parseInt(url.searchParams.get('limit') || '10', 10)

	return json({
		tasks: mockTasks,
		pagination: {
			page,
			limit,
			total: mockTasks.length,
			totalPages: Math.ceil(mockTasks.length / limit),
		},
	})
}

export const POST: RequestHandler = async ({ request }) => {
	const newTask = await request.json()
	const createdTask = {
		...newTask,
		id: `task-${Date.now()}`,
		status: 'todo',
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
	}
	const mockTasks = getMockTasks()
	mockTasks.push(createdTask)
	updateMockTasks(mockTasks)
	return json(createdTask, { status: 201 })
}
