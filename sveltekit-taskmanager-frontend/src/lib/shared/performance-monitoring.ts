/**
 * Performance Monitoring Module
 *
 * This module provides performance monitoring capabilities including:
 * - Web Vitals tracking
 * - Custom performance metrics
 * - Error tracking integration
 * - Performance reporting
 */

import type { Metric } from 'web-vitals'

// Custom metrics
interface CustomMetrics {
	apiResponseTime: number
	renderTime: number
	componentLoadTime: number
}

// Performance data structure
interface PerformanceData {
	webVitals: WebVitalsData
	customMetrics: CustomMetrics
	timestamp: number
	url: string
	userAgent: string
}

// Web Vitals types
interface WebVitalsData {
	FCP?: number // First Contentful Paint
	LCP?: number // Largest Contentful Paint
	FID?: number // First Input Delay
	CLS?: number // Cumulative Layout Shift
	TTFB?: number // Time to First Byte
}

class PerformanceMonitor {
	private customMetrics: CustomMetrics = {
		apiResponseTime: 0,
		componentLoadTime: 0,
		renderTime: 0,
	}
	private isEnabled: boolean
	private metrics: WebVitalsData = {}
	private reportingEndpoint?: string

	constructor(options: { enabled?: boolean; reportingEndpoint?: string } = {}) {
		this.isEnabled = options.enabled ?? true
		this.reportingEndpoint = options.reportingEndpoint

		if (this.isEnabled && typeof window !== 'undefined') {
			this.initializeWebVitals()
		}
	}

	/**
	 * Clear all metrics
	 */
	clearMetrics(): void {
		this.metrics = {}
		this.customMetrics = {
			apiResponseTime: 0,
			componentLoadTime: 0,
			renderTime: 0,
		}
	}

	/**
	 * Get current performance data
	 */
	getPerformanceData(): PerformanceData {
		return {
			customMetrics: this.customMetrics,
			timestamp: Date.now(),
			url: typeof window !== 'undefined' ? window.location.href : '',
			userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
			webVitals: this.metrics,
		}
	}

	/**
	 * Enable/disable monitoring
	 */
	setEnabled(enabled: boolean): void {
		this.isEnabled = enabled
	}

	/**
	 * Set reporting endpoint
	 */
	setReportingEndpoint(endpoint: string): void {
		this.reportingEndpoint = endpoint
	}

	/**
	 * Track API response time
	 */
	trackApiCall(url: string, startTime: number): void {
		const duration = performance.now() - startTime
		this.customMetrics.apiResponseTime = duration

		if (this.isEnabled) {
			console.log(`[Performance] API call to ${url}:`, duration)

			if (this.reportingEndpoint) {
				this.sendToReportingService({
					duration,
					timestamp: Date.now(),
					type: 'api_call',
					url,
				})
			}
		}
	}

	/**
	 * Track component render time
	 */
	trackComponentRender(componentName: string, startTime: number): void {
		const duration = performance.now() - startTime
		this.customMetrics.componentLoadTime = duration

		if (this.isEnabled) {
			console.log(`[Performance] Component ${componentName} render:`, duration)

			if (this.reportingEndpoint) {
				this.sendToReportingService({
					componentName,
					duration,
					timestamp: Date.now(),
					type: 'component_render',
				})
			}
		}
	}

	/**
	 * Track page render time
	 */
	trackPageRender(pageName: string, startTime: number): void {
		const duration = performance.now() - startTime
		this.customMetrics.renderTime = duration

		if (this.isEnabled) {
			console.log(`[Performance] Page ${pageName} render:`, duration)

			if (this.reportingEndpoint) {
				this.sendToReportingService({
					duration,
					pageName,
					timestamp: Date.now(),
					type: 'page_render',
				})
			}
		}
	}

	/**
	 * Initialize Web Vitals monitoring
	 */
	private async initializeWebVitals(): Promise<void> {
		try {
			const { onCLS, onFCP, onFID, onLCP, onTTFB } = await import('web-vitals')

			onCLS((metric: Metric) => {
				this.metrics.CLS = metric.value
				this.reportMetric('CLS', metric)
			})

			onFID((metric: Metric) => {
				this.metrics.FID = metric.value
				this.reportMetric('FID', metric)
			})

			onFCP((metric: Metric) => {
				this.metrics.FCP = metric.value
				this.reportMetric('FCP', metric)
			})

			onLCP((metric: Metric) => {
				this.metrics.LCP = metric.value
				this.reportMetric('LCP', metric)
			})

			onTTFB((metric: Metric) => {
				this.metrics.TTFB = metric.value
				this.reportMetric('TTFB', metric)
			})
		} catch (error) {
			console.warn('Failed to initialize Web Vitals:', error)
		}
	}

	/**
	 * Report individual metric
	 */
	private reportMetric(name: string, metric: Metric): void {
		if (!this.isEnabled) return

		console.log(`[Performance] ${name}:`, metric.value, metric)

		// Send to reporting endpoint if configured
		if (this.reportingEndpoint) {
			this.sendToReportingService({
				name,
				navigationType: metric.navigationType,
				rating: metric.rating,
				timestamp: Date.now(),
				value: metric.value,
			})
		}
	}

	/**
	 * Send metrics to reporting service
	 */
	private sendToReportingService(data: Record<string, unknown>): void {
		if (!this.reportingEndpoint) return

		// Use navigator.sendBeacon for non-blocking requests
		if (navigator.sendBeacon) {
			const blob = new Blob([JSON.stringify(data)], { type: 'application/json' })
			navigator.sendBeacon(this.reportingEndpoint, blob)
		} else {
			// Fallback to fetch
			fetch(this.reportingEndpoint, {
				body: JSON.stringify(data),
				keepalive: true,
				method: 'POST',
			}).catch((err) => console.warn('Failed to send metrics:', err))
		}
	}
}

// Singleton instance
let performanceMonitor: PerformanceMonitor

/**
 * Get performance monitor instance
 */
export function getPerformanceMonitor(): PerformanceMonitor {
	if (!performanceMonitor) {
		const enabled = import.meta.env.PUBLIC_ENABLE_PERFORMANCE_MONITORING === 'true'
		const endpoint = import.meta.env.PUBLIC_PERFORMANCE_ENDPOINT

		performanceMonitor = new PerformanceMonitor({
			enabled,
			reportingEndpoint: endpoint,
		})
	}

	return performanceMonitor
}

/**
 * Performance monitoring utility functions
 */
export const performanceUtils = {
	/**
	 * Get navigation timing
	 */
	getNavigationTiming(): null | Record<string, number> {
		if (typeof performance === 'undefined' || !performance.timing) {
			return null
		}

		const timing = performance.timing
		return {
			dns: timing.domainLookupEnd - timing.domainLookupStart,
			domProcessing: timing.domComplete - timing.domLoading,
			download: timing.responseEnd - timing.responseStart,
			tcp: timing.connectEnd - timing.connectStart,
			total: timing.loadEventEnd - timing.navigationStart,
			ttfb: timing.responseStart - timing.requestStart,
		}
	},

	/**
	 * Create performance mark
	 */
	mark(name: string): void {
		if (typeof performance !== 'undefined') {
			performance.mark(name)
		}
	},

	/**
	 * Create performance measure
	 */
	measure(name: string, startMark: string, endMark: string): number {
		if (typeof performance !== 'undefined') {
			try {
				performance.measure(name, startMark, endMark)
				const entries = performance.getEntriesByName(name, 'measure')
				return entries[entries.length - 1]?.duration || 0
			} catch (error) {
				console.warn('Failed to create performance measure:', error)
				return 0
			}
		}
		return 0
	},

	/**
	 * Measure function execution time
	 */
	async measureFunction<T>(name: string, fn: () => Promise<T>): Promise<T> {
		const start = performance.now()
		try {
			const result = await fn()
			const duration = performance.now() - start

			if (performanceMonitor) {
				console.log(`[Performance] ${name}:`, duration)
			}

			return result
		} catch (error) {
			const duration = performance.now() - start
			console.error(`[Performance] ${name} failed after:`, duration)
			throw error
		}
	},
}

export default PerformanceMonitor
