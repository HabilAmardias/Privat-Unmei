<script lang="ts">
	import { enhance } from '$app/forms';
	import Button from '$lib/components/button/Button.svelte';
	import Card from '$lib/components/card/Card.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import toast from 'svelte-french-toast';
	import { View } from './view.svelte';

	function onSendSubmit(args: EnhancementArgs) {
		View.setIsLoading(true);
		const loadID = toast.loading('loading....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
			View.setIsLoading(false);
			toast.dismiss(loadID);
			if (result.type === 'success') {
				toast.success(result.data?.message, { position: 'top-right' });
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
		<form use:enhance={onSendSubmit} action="?/send" method="post" class="flex flex-col gap-4">
			<Input bind:value={View.email} type="email" name="email" placeholder="Email" id="email" />
			<Button disabled={View.isLoading} type="submit" full={true}>Submit</Button>
		</form>
	</Card>
</div>
