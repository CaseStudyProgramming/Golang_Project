import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'
import { resetMockData } from '$lib/shared/mock-data/tasks.mock'

export const POST: RequestHandler = async () => {
	resetMockData()
	return json({ message: 'Mock data reset successfully' })
}
