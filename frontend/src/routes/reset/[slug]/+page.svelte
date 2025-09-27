<script lang="ts">
	import { View } from '../view.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import toast from 'svelte-french-toast';
	import InputSecret from '$lib/components/form/InputSecret.svelte';

	function onChangePasswordSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('loading....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			toast.dismiss(loadID);
			View.setIsLoading(false);
			if (result.type === 'redirect') {
				toast.success('successfully reset password', { position: 'top-right' });
			}
			if (result.type === 'failure') {
				toast.error(result.data?.message, { position: 'top-right' });
			}
			update();
		};
	}
</script>

<div class="flex h-full w-full items-center justify-center">
	<Card>
		<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Reset Password</h2>
		<form
			use:enhance={onChangePasswordSubmit}
			action="?/reset"
			method="post"
			class="flex flex-col gap-4"
		>
			<InputSecret
				id="password"
				placeholder="New Password"
				bind:value={View.password}
				name="password"
			/>
			<InputSecret
				id="password"
				placeholder="Repeat New Password"
				bind:value={View.repeatPassword}
				name="repeat-password"
			/>
			<Button disabled={View.isLoading} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>
