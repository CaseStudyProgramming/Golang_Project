import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'

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

export const GET: RequestHandler = async () => {
	return json(mockTags)
}

export const POST: RequestHandler = async ({ request }) => {
	const newTag = await request.json()
	const createdTag = {
		...newTag,
		id: 'new-tag-id',
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
	}
	return json(createdTag, { status: 201 })
}
