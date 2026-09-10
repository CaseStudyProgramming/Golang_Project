<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/features/auth';

	let { children } = $props();
	let isLoggingOut = $state(false);
	let isMobileMenuOpen = $state(false);

	async function handleLogout() {
		isLoggingOut = true;
		try {
			await authStore.logout();
			goto('/auth/login');
		} catch (error) {
			console.error('Logout failed:', error);
		} finally {
			isLoggingOut = false;
		}
	}

	function toggleMobileMenu() {
		isMobileMenuOpen = !isMobileMenuOpen;
	}

	function closeMobileMenu() {
		isMobileMenuOpen = false;
	}
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm border-b border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/dashboard" class="text-lg sm:text-xl font-bold text-gray-800">Task Manager</a>
				</div>
				
				<!-- Desktop Navigation -->
				<div class="hidden md:flex items-center space-x-4">
					<a
						href="/dashboard"
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Dashboard
					</a>
					<a
						href="/dashboard/tasks"
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
					>
						Tasks
					</a>
					<button
						onclick={handleLogout}
						disabled={isLoggingOut}
						class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{isLoggingOut ? 'Logging out...' : 'Logout'}
					</button>
				</div>

				<!-- Mobile menu button -->
				<div class="md:hidden flex items-center">
					<button
						onclick={toggleMobileMenu}
						class="text-gray-600 hover:text-gray-900 p-2 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500"
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
						onclick={closeMobileMenu}
						class="block text-gray-600 hover:text-gray-900 hover:bg-gray-50 px-3 py-2 rounded-md text-base font-medium"
					>
						Dashboard
					</a>
					<a
						href="/dashboard/tasks"
						onclick={closeMobileMenu}
						class="block text-gray-600 hover:text-gray-900 hover:bg-gray-50 px-3 py-2 rounded-md text-base font-medium"
					>
						Tasks
					</a>
					<button
						onclick={() => {
							handleLogout();
							closeMobileMenu();
						}}
						disabled={isLoggingOut}
						class="w-full text-left text-gray-600 hover:text-gray-900 hover:bg-gray-50 px-3 py-2 rounded-md text-base font-medium disabled:opacity-50 disabled:cursor-not-allowed"
					>
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
