<script lang="ts">
	import { authApi } from '$lib/features/auth';
	import { forgotPasswordSchema } from '$lib/features/auth/schemas/auth.schemas';
	import { getErrorMessage, isValidationError } from '$lib/shared/utils/error.utils';
	import ErrorToast from '$lib/shared/components/ErrorToast.svelte';
	import { goto } from '$app/navigation';

	let email = $state('');
	let isLoading = $state(false);
	let isSuccess = $state(false);
	let error = $state('');
	let fieldErrors = $state<Record<string, string>>({});
	let toastError = $state('');

	async function handleForgotPassword(e: Event) {
		e.preventDefault();
		isLoading = true;
		error = '';
		fieldErrors = {};
		toastError = '';

		try {
			const validatedData = forgotPasswordSchema.parse({ email });
			await authApi.forgotPassword(validatedData);
			isSuccess = true;
		} catch (err) {
			if (err instanceof Error && err.name === 'ZodError') {
				const zodError = err as any;
				if (zodError.errors && zodError.errors[0]) {
					fieldErrors.email = zodError.errors[0].message;
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
		<h1 class="text-3xl font-bold text-gray-800 mb-2">Forgot Password</h1>
		<p class="text-gray-600">Enter your email to receive a password reset link</p>
	</div>

	{#if isSuccess}
		<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-4">
			<p class="font-medium">Check your email</p>
			<p class="text-sm mt-1">We've sent a password reset link to your email address.</p>
		</div>
		<div class="mt-6 text-center">
			<a href="/auth/login" class="text-blue-600 hover:text-blue-700 font-medium">Back to login</a>
		</div>
	{:else}
		{#if error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
				{error}
			</div>
		{/if}

		<form onsubmit={handleForgotPassword} class="space-y-4">
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

			<button
				type="submit"
				disabled={isLoading}
				class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
			>
				{isLoading ? 'Sending...' : 'Send Reset Link'}
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
