import { test as base } from '@playwright/test'

// Extend base test with setup/teardown for e2e tests
export const test = base.extend({
  // Setup before all tests
  context: async ({ context }, use) => {
    // Additional context setup if needed
    await use(context)
  }
})

export const expect = base.expect
