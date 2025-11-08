<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { PaymentManagementView } from './view.svelte';
	import { goto } from '$app/navigation';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';

	let { data }: PageProps = $props();
	const View = new PaymentManagementView(data.payments);
	onMount(() => {
		if (!data.isVerified) {
			goto('/managet/admin/verify', { replaceState: true });
		}
	});
</script>

{#snippet deleteDialogTitle()}
	Delete Payment Confirmation
{/snippet}

{#snippet deleteDialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

{#snippet createDialogTitle()}
	Create Payment
{/snippet}

{#snippet createDialogDescription()}
	Add New Payment Method
{/snippet}

<div class="flex h-full flex-col p-4">
	<div class="mb-4 flex items-center justify-between">
		<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Payment Methods</h3>
		<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<Dialog
				bind:open={View.createDialogOpen}
				buttonText="Create Payment"
				title={createDialogTitle}
				description={createDialogDescription}
				buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
			>
				<form
					class="flex flex-col gap-4"
					action="?/createPayment"
					method="POST"
					use:enhance={View.onCreatePayment}
				>
					<Input type="text" id="name" name="name" placeholder="Insert Payment Method Name" />
					<Button full type="submit">Create</Button>
				</form>
			</Dialog>
		</div>
	</div>
	<form
		bind:this={View.searchForm}
		use:enhance={View.onSearchPayments}
		action="?/getPayments"
		class="mb-4 flex gap-4"
		method="POST"
	>
		<Input
			onInput={View.onSearchInput}
			type="text"
			name="search"
			id="search"
			placeholder="Search"
		/>
	</form>
	<div class="flex flex-1">
		{#if View.isLoading}
			<Loading />
		{:else if !View.payments || View.payments.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No payments found</b>
		{:else}
			<ScrollArea orientation="vertical" class="flex-1" viewportClasses="h-full w-full max-h-full">
				<table class="w-full border-separate border-spacing-4">
					<thead>
						<tr>
							<th class="text-[var(--tertiary-color)]">Name</th>
						</tr>
					</thead>
					<tbody>
						{#each View.payments as py (py.payment_method_id)}
							<tr>
								<td class="text-center">
									{py.payment_method_name}
								</td>
								<td>
									<div class="flex items-center justify-center">
										<AlertDialog
											action="?/deletePayment"
											bind:open={View.deleteDialogOpen}
											enhancement={View.onDeletePayment}
											title={deleteDialogTitle}
											onClick={() => {
												View.setPaymentToDelete(py.payment_method_id);
											}}
											description={deleteDialogDescription}>Delete</AlertDialog
										>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</ScrollArea>
		{/if}
	</div>
	<form
		use:enhance={View.onSearchPayments}
		action="?/getPayments"
		class="flex w-full items-center justify-center"
		method="POST"
	>
		<Pagination
			onPageChange={(num) => View.onPageChange(num)}
			pageNumber={View.page}
			perPage={View.limit}
			count={View.totalRow}
		/>
	</form>
</div>
