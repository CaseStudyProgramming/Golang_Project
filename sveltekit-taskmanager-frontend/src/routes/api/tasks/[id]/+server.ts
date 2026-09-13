import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'
import { getMockTasks, updateMockTasks } from '$lib/shared/mock-data/tasks.mock'

export const GET: RequestHandler = async ({ params }) => {
	const mockTasks = getMockTasks()
	const task = mockTasks.find((t) => t.id === params.id)
	if (!task) {
		return json({ error: 'Task not found' }, { status: 404 })
	}
	return json(task)
}

export const PUT: RequestHandler = async ({ params, request }) => {
	const updatedData = await request.json()
	const mockTasks = getMockTasks()
	const taskIndex = mockTasks.findIndex((t) => t.id === params.id)
	if (taskIndex === -1) {
		return json({ error: 'Task not found' }, { status: 404 })
	}
	const updatedTask = {
		...mockTasks[taskIndex],
		...updatedData,
		updatedAt: new Date().toISOString(),
	}
	const updatedTasks = [...mockTasks]
	updatedTasks[taskIndex] = updatedTask
	updateMockTasks(updatedTasks)
	return json(updatedTask)
}

export const DELETE: RequestHandler = async ({ params }) => {
	const mockTasks = getMockTasks()
	const taskIndex = mockTasks.findIndex((t) => t.id === params.id)
	if (taskIndex === -1) {
		return json({ error: 'Task not found' }, { status: 404 })
	}
	const updatedTasks = [...mockTasks]
	updatedTasks.splice(taskIndex, 1)
	updateMockTasks(updatedTasks)
	return json({ message: 'Task deleted successfully' })
}
