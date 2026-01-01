<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import Loading from '$lib/components/loader/Loading.svelte';
	import { GoogleCallbackView } from './view.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import TAC from './TAC.svelte';
	import { enhance } from '$app/forms';

	let { data }: PageProps = $props();
	let View = new GoogleCallbackView();

	onMount(() => {
		if (data.success && data.userStatus === 'verified') {
			View.verifiedLogin();
		}
		if (data.success && data.userStatus === 'unverified') {
			View.openDialog = true;
		}
	});
</script>

{#snippet TACDialogTitle()}
	Terms and Condition
{/snippet}

{#snippet TACDialogDescription()}
	<p class="text-justify text-[var(--secondary-color)]"></p>
{/snippet}

<div class="h-dvh">
	{#if data.userStatus === 'verified'}
		<Loading />
	{:else}
		<form
			bind:this={View.verifyForm}
			method="POST"
			action="?/verify"
			use:enhance={View.onVerify}
		></form>
		<Dialog
			bind:open={View.openDialog}
			buttonText=""
			buttonClass="text-[var(--tertiary-color)] hover:text-[var(--primary-color)] cursor-pointer"
			title={TACDialogTitle}
			description={TACDialogDescription}
		>
			<TAC onClick={() => View.verifyForm?.requestSubmit()} />
		</Dialog>
	{/if}
</div>
