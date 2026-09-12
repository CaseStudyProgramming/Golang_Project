import { stopServer } from './msw.setup'

export default async function globalTeardown() {
	console.log('Stopping MSW server...')
	stopServer()
}