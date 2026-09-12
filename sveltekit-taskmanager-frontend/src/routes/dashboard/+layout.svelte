<script lang="ts">
import { goto } from '$app/navigation'
import { authStore } from '$lib/features/auth'
import { confirmStore, toastStore } from '$lib/shared/stores'
import LoadingSpinner from '$lib/shared/components/LoadingSpinner.svelte'

let { children } = $props()
let isLoggingOut = $state(false)
let isMobileMenuOpen = $state(false)

function closeMobileMenu() {
	isMobileMenuOpen = false
}

async function handleLogout() {
	const confirmed = await confirmStore.showConfirm({
		cancelText: 'Cancel',
		confirmText: 'Logout',
		message: 'Are you sure you want to logout?',
		title: 'Logout',
		type: 'info',
	})

	if (!confirmed) return

	isLoggingOut = true
	try {
		await authStore.logout()
		toastStore.success('Logged out', 'You have been logged out successfully')
		await goto('/auth/login')
	} catch (error) {
		console.error('Logout failed:', error)
		toastStore.error(
			'Logout failed',
			error instanceof Error ? error.message : 'An unexpected error occurred'
		)
	} finally {
		isLoggingOut = false
	}
}

function toggleMobileMenu() {
	isMobileMenuOpen = !isMobileMenuOpen
}
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm border-b border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a
						href="/dashboard"
						onclick={(e) => { e.preventDefault(); goto('/dashboard'); }}
						class="text-lg sm:text-xl font-bold text-gray-800"
					>Task Manager</a>
				</div>

				<!-- Desktop Navigation -->
				<div class="hidden md:flex items-center space-x-4">
					<a
						href="/dashboard"
						onclick={(e) => { e.preventDefault(); goto('/dashboard'); }}
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Dashboard
					</a>
					<a
						href="/dashboard/tasks"
						onclick={(e) => { e.preventDefault(); goto('/dashboard/tasks'); }}
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Tasks
					</a>
					<button
						onclick={() => handleLogout()}
						disabled={isLoggingOut}
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					>
						{#if isLoggingOut}
							<LoadingSpinner size="sm" color="gray" />
						{/if}
						{isLoggingOut ? 'Logging out...' : 'Logout'}
					</button>
				</div>

				<!-- Mobile menu button -->
				<div class="md:hidden flex items-center">
					<button
						onclick={() => toggleMobileMenu()}
						class="text-gray-600 hover:text-gray-900 p-3 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 active:bg-gray-100 min-w-[44px] min-h-[44px]"
						aria-label="Toggle menu"
						aria-expanded={isMobileMenuOpen}
					>
						<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							{#if isMobileMenuOpen}
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							{:else}
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
							{/if}
						</svg>
					</button>
				</div>
			</div>
		</div>

		<!-- Mobile Navigation -->
		{#if isMobileMenuOpen}
			<div class="md:hidden bg-white border-t border-gray-200">
				<div class="px-4 pt-2 pb-4 space-y-1">
					<a
						href="/dashboard"
						onclick={(e) => { e.preventDefault(); closeMobileMenu(); goto('/dashboard'); }}
						class="block text-gray-600 hover:text-gray-900 hover:bg-gray-50 active:bg-gray-100 px-3 py-3 rounded-md text-base font-medium min-h-[44px] flex items-center"
					>
						Dashboard
					</a>
					<a
						href="/dashboard/tasks"
						onclick={(e) => { e.preventDefault(); closeMobileMenu(); goto('/dashboard/tasks'); }}
						class="block text-gray-600 hover:text-gray-900 hover:bg-gray-50 active:bg-gray-100 px-3 py-3 rounded-md text-base font-medium min-h-[44px] flex items-center"
					>
						Tasks
					</a>
					<button
						onclick={() => handleLogout()}
						disabled={isLoggingOut}
						class="w-full text-left text-gray-600 hover:text-gray-900 hover:bg-gray-50 active:bg-gray-100 px-3 py-3 rounded-md text-base font-medium disabled:opacity-50 disabled:cursor-not-allowed min-h-[44px] flex items-center gap-2"
					>
						{#if isLoggingOut}
							<LoadingSpinner size="sm" color="gray" />
						{/if}
						{isLoggingOut ? 'Logging out...' : 'Logout'}
					</button>
				</div>
			</div>
		{/if}
	</nav>

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
		{@render children()}
	</main>
</div>
