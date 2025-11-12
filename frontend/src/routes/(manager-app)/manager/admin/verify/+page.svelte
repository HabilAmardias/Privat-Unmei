<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import Card from '$lib/components/card/Card.svelte';
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import Input from '$lib/components/form/Input.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { VerifyAdminView } from './view.svelte';
	import InputSecret from '$lib/components/form/InputSecret.svelte';
	import { CreateToast, DismissToast } from '$lib/utils/helper';

	let { data }: PageProps = $props();
	const View = new VerifyAdminView();

	onMount(() => {
		if (data.isVerified) {
			goto('/manager/admin', { replaceState: true });
		}
	});

	function onVerifySubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			View.setIsLoading(false);
			if (result.type === 'success') {
				loadingStore.setLogOutLoadID(loadID);
				await goto('/manager/login', { replaceState: true });
			}
			if (result.type === 'failure') {
				DismissToast(loadID);
				CreateToast('error', result.data?.message);
			}
		};
	}
</script>

<svelte:head>
	<title>Verify Account - Privat Unmei</title>
	<meta name="description" content="Verify Account - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-dvh w-full items-center justify-center">
	<Card>
		<h2 class="mb-3 text-2xl font-bold text-[var(--tertiary-color)]">Verify</h2>
		<form
			use:enhance={onVerifySubmit}
			action="?/verifyAdmin"
			method="post"
			class="flex flex-col gap-4"
		>
			<Input
				err={View.emailError}
				onBlur={() => View.emailOnBlur()}
				bind:value={View.email}
				type="email"
				name="email"
				placeholder="Email"
				id="email"
			/>
			<InputSecret
				err={View.passwordError}
				onBlur={() => View.passwordOnBlur()}
				id="password"
				placeholder="Password"
				name="password"
				bind:value={View.password}
			/>
			<Button disabled={View.verifyDisabled} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>
