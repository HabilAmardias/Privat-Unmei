<script lang="ts">
	import { enhance } from '$app/forms';
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import { ResetView } from './view.svelte';
	import type { PageProps } from './$types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { CreateToast, DismissToast } from '$lib/utils/helper';

	const View = new ResetView();
	let { data }: PageProps = $props();

	onMount(() => {
		if (data.returnHome) {
			goto('/', { replaceState: true });
		}
	});

	function onSendSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result, update }: EnhancementReturn) => {
			View.setIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			update();
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
		<form use:enhance={onSendSubmit} action="?/send" method="post" class="flex flex-col gap-4">
			<Input
				err={View.emailError}
				onBlur={() => View.emailOnBlur()}
				bind:value={View.email}
				type="email"
				name="email"
				placeholder="Email"
				id="email"
			/>
			<Button disabled={View.sendDisabled} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>
