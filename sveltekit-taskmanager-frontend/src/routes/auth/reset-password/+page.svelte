<script lang="ts">
	import { authApi } from '$lib/features/auth';
	import { resetPasswordSchema } from '$lib/features/auth/schemas/auth.schemas';
	import { getErrorMessage, isValidationError } from '$lib/shared/utils/error.utils';
	import ErrorToast from '$lib/shared/components/ErrorToast.svelte';
	import PasswordStrength from '$lib/shared/components/PasswordStrength.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let password = $state('');
	let confirmPassword = $state('');
	let isLoading = $state(false);
	let isSuccess = $state(false);
	let error = $state('');
	let fieldErrors = $state<Record<string, string>>({});
	let toastError = $state('');

	// Get reset token from URL
	const getToken = () => {
		const urlParams = new URLSearchParams($page.url.search);
		return urlParams.get('token');
	};

	// Redirect if no token
	if (!getToken()) {
		goto('/auth/forgot-password');
	}

	async function handleResetPassword(e: Event) {
		e.preventDefault();
		isLoading = true;
		error = '';
		fieldErrors = {};
		toastError = '';

		if (password !== confirmPassword) {
			fieldErrors.confirmPassword = 'Passwords do not match';
			error = 'Please fix the errors below.';
			isLoading = false;
			return;
		}

		try {
			const validatedData = resetPasswordSchema.parse({
				token: getToken(),
				password,
				confirmPassword
			});
			await authApi.resetPassword(validatedData);
			isSuccess = true;
		} catch (err) {
			if (err instanceof Error && err.name === 'ZodError') {
				const zodError = err as any;
				if (zodError.errors && zodError.errors[0]) {
					const field = zodError.errors[0].path[0] as string;
					fieldErrors[field] = zodError.errors[0].message;
					error = 'Please fix the errors below.';
				}
			} else if (isValidationError(err)) {
				fieldErrors[err.field] = err.message;
				error = 'Please fix the errors below.';
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
		<h1 class="text-3xl font-bold text-gray-800 mb-2">Reset Password</h1>
		<p class="text-gray-600">Enter your new password below</p>
	</div>

	{#if isSuccess}
		<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-4">
			<p class="font-medium">Password reset successful</p>
			<p class="text-sm mt-1">You can now log in with your new password.</p>
		</div>
		<div class="mt-6 text-center">
			<a href="/auth/login" class="text-blue-600 hover:text-blue-700 font-medium">Go to login</a>
		</div>
	{:else}
		{#if error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
				{error}
			</div>
		{/if}

		<form onsubmit={handleResetPassword} class="space-y-4">
			<div>
				<label for="password" class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
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
				{#if password}
					<div class="mt-2">
						<PasswordStrength password={password} />
					</div>
				{/if}
			</div>

			<div>
				<label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
				<input
					id="confirmPassword"
					type="password"
					bind:value={confirmPassword}
					required
					disabled={isLoading}
					class="w-full px-4 py-2 border {fieldErrors.confirmPassword ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
					placeholder="••••••••"
				/>
				{#if fieldErrors.confirmPassword}
					<p class="mt-1 text-sm text-red-600">{fieldErrors.confirmPassword}</p>
				{/if}
			</div>

			<button
				type="submit"
				disabled={isLoading}
				class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
			>
				{isLoading ? 'Resetting...' : 'Reset Password'}
			</button>
		</form>

		<div class="mt-6 text-center">
			<p class="text-gray-600">
				Remember your password?
				<a href="/auth/login" class="text-blue-600 hover:text-blue-700 font-medium">Sign in</a>
			</p>
		</div>
	{/if}
</div>

{#if toastError}
	<ErrorToast message={toastError} onDismiss={dismissToast} />
{/if}
