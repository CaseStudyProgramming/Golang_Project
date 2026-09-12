import { startServer } from './msw.setup'

export default async function globalSetup() {
	console.log('Starting MSW server for E2E tests...')
	startServer()
}
