<script lang="ts">
	let { password }: { password: string } = $props();

	const strengthLevels = [
		{ label: 'Weak', color: 'bg-red-500', minScore: 0 },
		{ label: 'Fair', color: 'bg-orange-500', minScore: 2 },
		{ label: 'Good', color: 'bg-yellow-500', minScore: 3 },
		{ label: 'Strong', color: 'bg-green-500', minScore: 4 }
	];

	const requirements = [
		{ label: 'At least 8 characters', test: (pwd: string) => pwd.length >= 8 },
		{ label: 'Uppercase letter', test: (pwd: string) => /[A-Z]/.test(pwd) },
		{ label: 'Lowercase letter', test: (pwd: string) => /[a-z]/.test(pwd) },
		{ label: 'Number', test: (pwd: string) => /[0-9]/.test(pwd) },
		{ label: 'Special character', test: (pwd: string) => /[^A-Za-z0-9]/.test(pwd) }
	];

	const getStrength = () => {
		let score = 0;
		for (const req of requirements) {
			if (req.test(password)) score++;
		}
		return score;
	};

	const getStrengthLevel = () => {
		const strength = getStrength();
		for (let i = strengthLevels.length - 1; i >= 0; i--) {
			if (strength >= strengthLevels[i].minScore) {
				return strengthLevels[i];
			}
		}
		return strengthLevels[0];
	};
</script>

<div class="space-y-3">
	<div class="flex gap-1">
		{#each requirements as req}
			<div
				class="h-1 flex-1 rounded-full {req.test(password) ? getStrengthLevel().color : 'bg-gray-200'}"
				title={req.label}
			></div>
		{/each}
	</div>

	<div class="flex justify-between items-center">
		<span class="text-sm font-medium {getStrengthLevel().color.replace('bg-', 'text-')}">
			{getStrengthLevel().label}
		</span>
		<span class="text-xs text-gray-500">{getStrength()}/5 requirements met</span>
	</div>

	<ul class="space-y-1 text-xs">
		{#each requirements as req}
			<li class="flex items-center gap-2 {req.test(password) ? 'text-green-600' : 'text-gray-400'}">
				<span class="w-4 h-4 flex items-center justify-center">
					{req.test(password) ? '✓' : '○'}
				</span>
				{req.label}
			</li>
		{/each}
	</ul>
</div>
