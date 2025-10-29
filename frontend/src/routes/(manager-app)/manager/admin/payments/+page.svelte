<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { PaymentManagementView } from './view.svelte';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/button/Button.svelte';
	import { Pencil, Trash } from '@lucide/svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Link from '$lib/components/button/Link.svelte';

	let { data }: PageProps = $props();
	const View = new PaymentManagementView();
	onMount(() => {
		if (!data.isVerified) {
			goto('/managet/admin/verify', { replaceState: true });
		}
		View.setLastID(data.payments.page_info.last_id);
		View.setPayments(data.payments.entries);
		View.setTotalRow(data.payments.page_info.total_row);
	});
</script>

<div class="h-full p-4">
	<div class="mb-4 flex items-center justify-between">
		<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Payment Methods</h3>
		<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<Link theme="light" href="/manager/admin/payments/create">Create New Payment</Link>
		</div>
	</div>
	<form action="?/deletePayment"></form>
	<form action="?/getPayments" class="mb-4 flex gap-4">
		<Input type="text" name="search" id="search" placeholder="Search" />
		<Button>Search</Button>
	</form>
	<table class="w-full text-center">
		<thead>
			<tr>
				<th class="text-[var(--tertiary-color)]">Name</th>
			</tr>
		</thead>
		<tbody>
			{#each View.payments as py (py.payment_method_id)}
				<tr>
					<td>
						{py.payment_method_name}
					</td>
					<td>
						<div class="ml-auto w-fit">
							<Button><Trash /></Button>
						</div>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
