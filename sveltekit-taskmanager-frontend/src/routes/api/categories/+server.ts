import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'

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

export const GET: RequestHandler = async () => {
	return json(mockCategories)
}
