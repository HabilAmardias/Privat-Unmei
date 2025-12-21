<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import { VerifyView } from './view.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import { CircleCheck } from '@lucide/svelte';
	import Button from '$lib/components/button/Button.svelte';

	const View = new VerifyView();

	let { data }: PageProps = $props();
	onMount(() => {
		if (data.success) {
			View.setOpenDialog(true);
		}
		return () => {
			View.setOpenDialog(false);
		};
	});
	function navigateToLogin() {
		goto('/logout', { replaceState: true });
	}
</script>

{#snippet dialogTitle()}
	Successfully Verified
{/snippet}

{#snippet dialogContent()}{/snippet}

<Dialog buttonText="" bind:open={View.openDialog} title={dialogTitle} description={dialogContent}>
	<div class="flex w-full flex-col items-center gap-4">
		<CircleCheck size={128} class="text-[var(--tertiary-color)]" />
		<p>Your account is now successfully verified</p>
		<p>You can now login with your account</p>
		<Button full={true} onClick={() => navigateToLogin()}>Login</Button>
	</div>
</Dialog>
