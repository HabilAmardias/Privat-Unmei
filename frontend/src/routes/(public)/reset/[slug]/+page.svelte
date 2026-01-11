<script lang="ts">
	import { ResetView } from '../view.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { goto } from '$app/navigation';
	import type { PageProps } from './$types';
	import { onMount } from 'svelte';
	import { CreateToast, DismissToast } from '$lib/utils/helper';

	const View = new ResetView();
	let { data }: PageProps = $props();

	onMount(() => {
		if (data.returnHome) {
			goto('/', { replaceState: true });
		}
	});

	function onChangePasswordSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			View.setIsLoading(false);
			if (result.type === 'success') {
				CreateToast('success', 'successfully reset password');
				await goto('/logout', { replaceState: true });
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	}
</script>

<svelte:head>
	<title>Reset Password - Privat Unmei</title>
	<meta name="description" content="Reset Password - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-dvh w-full items-center justify-center">
	<Card>
		<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Reset Password</h2>
		<form
			use:enhance={onChangePasswordSubmit}
			action="?/reset"
			method="post"
			class="flex flex-col gap-4"
		>
			<InputSecret
				err={View.passwordError}
				id="password"
				placeholder="New Password"
				bind:value={View.password}
				name="password"
			/>
			<InputSecret
				err={View.repeatPasswordError}
				id="password"
				placeholder="Repeat New Password"
				bind:value={View.repeatPassword}
				name="repeat-password"
			/>
			<Button disabled={View.resetDisabled} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>
