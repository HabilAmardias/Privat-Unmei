<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import { View } from './view.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import { CircleCheck } from '@lucide/svelte';
	import Button from '$lib/components/button/Button.svelte';

	let { data }: PageProps = $props();
	onMount(() => {
		if (data.success) {
			View.setOpenDialog(true);
		} else {
			goto('/', { replaceState: true });
		}
		return () => {
			View.setOpenDialog(false);
		};
	});
	function navigateToLogin() {
		goto('/', { replaceState: true });
	}
</script>

{#snippet dialogTitle()}
	Successfully Verified
{/snippet}

{#snippet dialogContent()}
	<CircleCheck size={128} />
	<p>Your account is now successfully verified</p>
	<p>You can now login with your account</p>
	<Button full={true} onClick={() => navigateToLogin()}>Login</Button>
{/snippet}

<Dialog bind:open={View.openDialog} {dialogTitle} {dialogContent} />
