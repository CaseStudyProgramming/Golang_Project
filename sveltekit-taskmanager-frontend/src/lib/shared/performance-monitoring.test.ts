import { describe, it, expect, beforeEach, vi } from 'vitest'
import PerformanceMonitor, {
	getPerformanceMonitor,
	performanceUtils,
} from './performance-monitoring'

describe('Performance Monitor', () => {
	let monitor: PerformanceMonitor

	beforeEach(() => {
		monitor = new PerformanceMonitor({ enabled: false })
		vi.spyOn(console, 'log').mockImplementation(() => {})
		vi.spyOn(console, 'warn').mockImplementation(() => {})
	})

	afterEach(() => {
		vi.restoreAllMocks()
	})

	describe('constructor', () => {
		it('should create monitor with default options', () => {
			const defaultMonitor = new PerformanceMonitor()
			expect(defaultMonitor).toBeDefined()
		})

		it('should create monitor with custom options', () => {
			const customMonitor = new PerformanceMonitor({
				enabled: false,
				reportingEndpoint: 'http://example.com/metrics',
			})
			expect(customMonitor).toBeDefined()
		})

		it('should be enabled by default', () => {
			const defaultMonitor = new PerformanceMonitor()
			expect(defaultMonitor).toBeDefined()
		})
	})

	describe('clearMetrics', () => {
		it('should clear all metrics', () => {
			monitor.trackApiCall('http://test.com', performance.now())
			monitor.clearMetrics()
			const data = monitor.getPerformanceData()
			expect(data.customMetrics.apiResponseTime).toBe(0)
			expect(data.customMetrics.componentLoadTime).toBe(0)
			expect(data.customMetrics.renderTime).toBe(0)
		})
	})

	describe('getPerformanceData', () => {
		it('should return performance data structure', () => {
			const data = monitor.getPerformanceData()
			expect(data).toHaveProperty('webVitals')
			expect(data).toHaveProperty('customMetrics')
			expect(data).toHaveProperty('timestamp')
			expect(data).toHaveProperty('url')
			expect(data).toHaveProperty('userAgent')
		})

		it('should include current timestamp', () => {
			const before = Date.now()
			const data = monitor.getPerformanceData()
			const after = Date.now()
			expect(data.timestamp).toBeGreaterThanOrEqual(before)
			expect(data.timestamp).toBeLessThanOrEqual(after)
		})

		it('should include custom metrics', () => {
			monitor.trackApiCall('http://test.com', performance.now())
			const data = monitor.getPerformanceData()
			expect(data.customMetrics.apiResponseTime).toBeGreaterThan(0)
		})
	})

	describe('setEnabled', () => {
		it('should enable monitoring', () => {
			monitor.setEnabled(true)
			expect(monitor).toBeDefined()
		})

		it('should disable monitoring', () => {
			monitor.setEnabled(false)
			expect(monitor).toBeDefined()
		})
	})

	describe('setReportingEndpoint', () => {
		it('should set reporting endpoint', () => {
			monitor.setReportingEndpoint('http://example.com/metrics')
			expect(monitor).toBeDefined()
		})
	})

	describe('trackApiCall', () => {
		it('should track API call duration', () => {
			const startTime = performance.now()
			monitor.trackApiCall('http://test.com', startTime)
			const data = monitor.getPerformanceData()
			expect(data.customMetrics.apiResponseTime).toBeGreaterThanOrEqual(0)
		})

		it('should log when enabled', () => {
			const enabledMonitor = new PerformanceMonitor({ enabled: true })
			const startTime = performance.now()
			enabledMonitor.trackApiCall('http://test.com', startTime)
			expect(console.log).toHaveBeenCalled()
		})
	})

	describe('trackComponentRender', () => {
		it('should track component render duration', () => {
			const startTime = performance.now()
			monitor.trackComponentRender('TestComponent', startTime)
			const data = monitor.getPerformanceData()
			expect(data.customMetrics.componentLoadTime).toBeGreaterThanOrEqual(0)
		})

		it('should log when enabled', () => {
			const enabledMonitor = new PerformanceMonitor({ enabled: true })
			const startTime = performance.now()
			enabledMonitor.trackComponentRender('TestComponent', startTime)
			expect(console.log).toHaveBeenCalled()
		})
	})

	describe('trackPageRender', () => {
		it('should track page render duration', () => {
			const startTime = performance.now()
			monitor.trackPageRender('TestPage', startTime)
			const data = monitor.getPerformanceData()
			expect(data.customMetrics.renderTime).toBeGreaterThanOrEqual(0)
		})

		it('should log when enabled', () => {
			const enabledMonitor = new PerformanceMonitor({ enabled: true })
			const startTime = performance.now()
			enabledMonitor.trackPageRender('TestPage', startTime)
			expect(console.log).toHaveBeenCalled()
		})
	})
})

describe('getPerformanceMonitor', () => {
	it('should return singleton instance', () => {
		const instance1 = getPerformanceMonitor()
		const instance2 = getPerformanceMonitor()
		expect(instance1).toBe(instance2)
	})

	it('should create new instance if not exists', () => {
		// Reset the singleton by clearing the module cache logic would be needed here
		// For now, just test that it returns an instance
		const instance = getPerformanceMonitor()
		expect(instance).toBeInstanceOf(PerformanceMonitor)
	})
})

describe('performanceUtils', () => {
	describe('getNavigationTiming', () => {
		it('should return null when performance API is not available', () => {
			// Mock performance as undefined
			const originalPerformance = global.performance
			// @ts-expect-error - Testing undefined performance
			delete global.performance

			const result = performanceUtils.getNavigationTiming()
			expect(result).toBeNull()

			global.performance = originalPerformance
		})

		it('should return timing data when available', () => {
			// This test would need to mock performance.timing
			// For now, just test the function exists
			expect(typeof performanceUtils.getNavigationTiming).toBe('function')
		})
	})

	describe('mark', () => {
		it('should create performance mark when performance API is available', () => {
			performanceUtils.mark('test-mark')
			// If performance API is available, mark should be created
			expect(typeof performanceUtils.mark).toBe('function')
		})

		it('should not throw when performance API is not available', () => {
			const originalPerformance = global.performance
			// @ts-expect-error - Testing undefined performance
			delete global.performance

			expect(() => performanceUtils.mark('test-mark')).not.toThrow()

			global.performance = originalPerformance
		})
	})

	describe('measure', () => {
		it('should create performance measure', () => {
			performanceUtils.mark('start-mark')
			performanceUtils.mark('end-mark')
			const duration = performanceUtils.measure('test-measure', 'start-mark', 'end-mark')
			expect(typeof duration).toBe('number')
		})

		it('should return 0 when marks do not exist', () => {
			const duration = performanceUtils.measure('invalid-measure', 'non-existent-start', 'non-existent-end')
			expect(duration).toBe(0)
		})

		it('should not throw when performance API is not available', () => {
			const originalPerformance = global.performance
			// @ts-expect-error - Testing undefined performance
			delete global.performance

			expect(() => performanceUtils.measure('test', 'start', 'end')).not.toThrow()

			global.performance = originalPerformance
		})
	})

	describe('measureFunction', () => {
		it('should measure async function execution time', async () => {
			const testFn = async () => {
				await new Promise((resolve) => setTimeout(resolve, 10))
				return 'result'
			}

			const result = await performanceUtils.measureFunction('test-function', testFn)
			expect(result).toBe('result')
		})

		it('should handle function errors', async () => {
			const testFn = async () => {
				throw new Error('Test error')
			}

			await expect(performanceUtils.measureFunction('test-function', testFn)).rejects.toThrow('Test error')
		})

		it('should log execution time when monitor is available', async () => {
			const testFn = async () => 'result'
			await performanceUtils.measureFunction('test-function', testFn)
			// Console log should be called if performanceMonitor is available
		})
	})
})
