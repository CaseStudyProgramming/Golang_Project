<script lang="ts">
	import { authStore } from '$lib/features/auth';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getErrorMessage, isValidationError, isAuthenticationError } from '$lib/shared/utils/error.utils';
	import ErrorToast from '$lib/shared/components/ErrorToast.svelte';

	let email = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let error = $state('');
	let fieldErrors = $state<Record<string, string>>({});
	let toastError = $state('');

	// Get redirect URL from query params
	const getRedirectTo = () => {
		const urlParams = new URLSearchParams($page.url.search);
		return urlParams.get('redirectTo');
	};

	async function handleLogin(e: Event) {
		e.preventDefault();
		isLoading = true;
		error = '';
		fieldErrors = {};
		toastError = '';

		// Save redirect URL before login
		const redirectTo = getRedirectTo();
		if (redirectTo) {
			authStore.saveRedirectUrl(redirectTo);
		}

		try {
			await authStore.login({ email, password });
			await authStore.redirectAfterAuth();
		} catch (err) {
			if (isValidationError(err)) {
				fieldErrors[err.field] = err.message;
				error = 'Please fix the errors below.';
			} else if (isAuthenticationError(err)) {
				error = err.message;
				toastError = err.message;
			} else {
				error = getErrorMessage(err);
				toastError = getErrorMessage(err);
			}
		} finally {
			isLoading = false;
		}
	}

	function dismissToast() {
		toastError = '';
	}
</script>

<div class="bg-white rounded-lg shadow-lg p-8">
	<div class="text-center mb-6">
		<h1 class="text-3xl font-bold text-gray-800 mb-2">Welcome Back</h1>
		<p class="text-gray-600">Sign in to your account</p>
	</div>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
			{error}
		</div>
	{/if}

	<form onsubmit={handleLogin} class="space-y-4">
		<div>
			<label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				required
				disabled={isLoading}
				class="w-full px-4 py-2 border {fieldErrors.email ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
				placeholder="you@example.com"
			/>
			{#if fieldErrors.email}
				<p class="mt-1 text-sm text-red-600">{fieldErrors.email}</p>
			{/if}
		</div>

		<div>
			<label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				required
				disabled={isLoading}
				class="w-full px-4 py-2 border {fieldErrors.password ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
				placeholder="••••••••"
			/>
			{#if fieldErrors.password}
				<p class="mt-1 text-sm text-red-600">{fieldErrors.password}</p>
			{/if}
		</div>

		<button
			type="submit"
			disabled={isLoading}
			class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
		>
			{isLoading ? 'Signing in...' : 'Sign In'}
		</button>
	</form>

	<div class="mt-6 text-center">
		<p class="text-gray-600">
			Don't have an account?
			<a href="/auth/register" class="text-blue-600 hover:text-blue-700 font-medium">Sign up</a>
		</p>
	</div>

	<div class="mt-4 text-center">
		<a href="/auth/forgot-password" class="text-sm text-blue-600 hover:text-blue-700">Forgot password?</a>
	</div>
</div>

{#if toastError}
	<ErrorToast message={toastError} onDismiss={dismissToast} />
{/if}
