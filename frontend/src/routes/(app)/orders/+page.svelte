<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import toast from 'svelte-french-toast';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import { enhance } from '$app/forms';

	async function onRefreshSubmit(args: EnhancementArgs) {
		const loadID = toast.loading('loading....', { position: 'top-right' });
		return async ({ result, update }: EnhancementReturn) => {
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

<form use:enhance={onRefreshSubmit} action="?/refresh" method="POST">
	<Button type="submit">Refresh</Button>
</form>
