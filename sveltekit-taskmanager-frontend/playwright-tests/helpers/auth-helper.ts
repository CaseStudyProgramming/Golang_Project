import { Page, APIRequestContext } from '@playwright/test'

const TEST_USER = {
  email: 'e2e-test@example.com',
  password: 'testpassword123',
  name: 'E2E Test User'
}

export class AuthHelper {
  static async registerTestUser(request: APIRequestContext): Promise<void> {
    try {
      const response = await request.post('http://localhost:8080/auth/register', {
        data: TEST_USER
      })
      // Ignore if user already exists
      if (response.status() !== 201 && response.status() !== 400) {
        console.log('Registration status:', response.status())
      }
    } catch (error) {
      // Ignore registration errors (user might already exist)
      console.log('Registration note:', error)
    }
  }

  static async loginTestUser(page: Page): Promise<void> {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    
    await page.fill('[data-testid="email-input"]', TEST_USER.email)
    await page.fill('[data-testid="password-input"]', TEST_USER.password)
    await page.click('[data-testid="login-button"]')
    
    await page.waitForURL('/tasks')
    await page.waitForLoadState('networkidle')
  }

  static async cleanupTestUser(request: APIRequestContext): Promise<void> {
    try {
      // Login to get token
      const loginResponse = await request.post('http://localhost:8080/auth/login', {
        data: {
          email: TEST_USER.email,
          password: TEST_USER.password
        }
      })
      
      if (loginResponse.ok()) {
        const data = await loginResponse.json()
        const token = data.data.token
        
        // Delete all tasks for this user
        await request.get('http://localhost:8080/tasks', {
          headers: { Authorization: `Bearer ${token}` }
        })
        
        // Logout
        await request.post('http://localhost:8080/auth/logout', {
          headers: { Authorization: `Bearer ${token}` }
        })
      }
    } catch (error) {
      console.log('Cleanup note:', error)
    }
  }
}
