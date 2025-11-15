<script lang="ts">
	import type { PageProps } from './$types';
	import { CostManagementView } from './view.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import Dialog from '$lib/components/dialog/Dialog.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';

	let { data }: PageProps = $props();
	const View = new CostManagementView(data.costs, data.discounts);
</script>

<svelte:head>
	<title>Cost Management - Privat Unmei</title>
	<meta name="description" content="Cost Management - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

{#snippet deleteCostDialogTitle()}
	Delete Cost Confirmation
{/snippet}

{#snippet deleteDiscountDialogTitle()}
	Delete Discount Confirmation
{/snippet}

{#snippet deleteDialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

{#snippet createCostDialogTitle()}
	Create Additional Cost
{/snippet}

{#snippet createDiscountDialogTitle()}
	Create Discount
{/snippet}

{#snippet createCostDialogDescription()}
	Add New Cost
{/snippet}

{#snippet createDiscountDialogDescription()}
	Add New Discount
{/snippet}

{#snippet updateCostDialogTitle()}
	Update Additional Cost
{/snippet}

{#snippet updateDiscountDialogTitle()}
	Update Discount
{/snippet}

{#snippet updateCostDialogDescription()}
	Update Cost Amount
{/snippet}

{#snippet updateDiscountDialogDescription()}
	Update Discount Amount
{/snippet}

<div class="flex h-full flex-col p-4">
	<div class="mb-4">
		<NavigationButton
			menus={[
				{
					header: 'Costs',
					onClick: () => (View.menuState = 'costs')
				},
				{
					header: 'Discounts',
					onClick: () => (View.menuState = 'discounts')
				}
			]}
		/>
	</div>
	{#if View.menuState === 'costs'}
		<div class="mb-4 flex items-center justify-between">
			<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Additional Costs</h3>
			<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
				<Dialog
					bind:open={View.createCostDialogOpen}
					buttonText="Create Cost"
					title={createCostDialogTitle}
					description={createCostDialogDescription}
					buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
				>
					<form
						class="flex flex-col gap-4"
						action="?/addCost"
						method="POST"
						use:enhance={View.onCreateCost}
					>
						<Input type="text" id="name" name="name" placeholder="Insert Cost Name" />
						<Input type="number" id="amount" name="amount" placeholder="Insert Cost Amount" />
						<Button full type="submit">Create</Button>
					</form>
				</Dialog>
			</div>
		</div>
		<div class="flex flex-1">
			{#if View.isLoading}
				<Loading />
			{:else if !View.costs || View.costs.length === 0}
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No Cost found</b>
			{:else}
				<ScrollArea orientation="vertical" class="flex-1" viewportClasses="max-h-[500px]">
					<table class="w-full table-fixed border-separate border-spacing-4">
						<thead>
							<tr>
								<th class="text-[var(--tertiary-color)]">Name</th>
								<th class="text-[var(--tertiary-color)]">Amount</th>
							</tr>
						</thead>
						<tbody>
							{#each View.costs as c (c.id)}
								<tr>
									<td class="text-center">
										{c.name}
									</td>
									<td class="text-center">
										{c.amount}
									</td>
									<td>
										<div class="flex items-center justify-center gap-4">
											<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
												<Dialog
													bind:open={View.updateCostDialogOpen}
													buttonText="Update"
													title={updateCostDialogTitle}
													buttonOnClick={() => {
														View.costToUpdate = c.id;
													}}
													description={updateCostDialogDescription}
													buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
												>
													<form
														class="flex flex-col gap-4"
														action="?/updateCostAmount"
														method="POST"
														use:enhance={View.onUpdateCost}
													>
														<Input
															type="number"
															id="amount"
															name="amount"
															placeholder="Insert Cost Amount"
														/>
														<Button full type="submit">Update</Button>
													</form>
												</Dialog>
											</div>
											<AlertDialog
												action="?/deleteCost"
												bind:open={View.deleteCostDialogOpen}
												enhancement={View.onDeleteCost}
												title={deleteCostDialogTitle}
												onClick={() => {
													View.costToDelete = c.id;
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
			use:enhance={View.onCostPageChangeForm}
			action="?/getCosts"
			class="flex w-full items-center justify-center"
			method="POST"
			bind:this={View.costPaginationForm}
		>
			<Pagination
				onPageChange={(num) => View.setPageChange(num)}
				pageNumber={View.costPageNumber}
				perPage={View.costLimit}
				count={View.costTotalRow}
				offset
			/>
		</form>
	{:else}
		<div class="mb-4 flex items-center justify-between">
			<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Discounts</h3>
			<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
				<Dialog
					bind:open={View.createDiscountDialogOpen}
					buttonText="Create Discount"
					title={createDiscountDialogTitle}
					description={createDiscountDialogDescription}
					buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
				>
					<form
						class="flex flex-col gap-4"
						action="?/addDiscount"
						method="POST"
						use:enhance={View.onCreateDiscount}
					>
						<Input
							type="number"
							id="number_of_participant"
							name="number_of_participant"
							placeholder="Number Of Participant"
						/>
						<Input type="number" id="amount" name="amount" placeholder="Insert Cost Amount" />
						<Button full type="submit">Create</Button>
					</form>
				</Dialog>
			</div>
		</div>
		<div class="flex flex-1">
			{#if View.isLoading}
				<Loading />
			{:else if !View.discounts || View.discounts.length === 0}
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No Discount found</b>
			{:else}
				<ScrollArea orientation="vertical" class="flex-1" viewportClasses="max-h-[500px]">
					<table class="w-full table-fixed border-separate border-spacing-4">
						<thead>
							<tr>
								<th class="text-[var(--tertiary-color)]">Number Of Participant</th>
								<th class="text-[var(--tertiary-color)]">Amount</th>
							</tr>
						</thead>
						<tbody>
							{#each View.discounts as d (d.id)}
								<tr>
									<td class="text-center">
										{d.number_of_participant}
									</td>
									<td class="text-center">
										{d.amount}
									</td>
									<td>
										<div class="flex items-center justify-center gap-4">
											<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
												<Dialog
													bind:open={View.updateDiscountDialogOpen}
													buttonText="Update"
													title={updateDiscountDialogTitle}
													buttonOnClick={() => {
														View.discountToUpdate = d.id;
													}}
													description={updateDiscountDialogDescription}
													buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
												>
													<form
														class="flex flex-col gap-4"
														action="?/updateDiscountAmount"
														method="POST"
														use:enhance={View.onUpdateDiscount}
													>
														<Input
															type="number"
															id="amount"
															name="amount"
															placeholder="Insert Discount Amount"
														/>
														<Button full type="submit">Update</Button>
													</form>
												</Dialog>
											</div>
											<AlertDialog
												action="?/deleteDiscount"
												bind:open={View.deleteDiscountDialogOpen}
												enhancement={View.onDeleteDiscount}
												title={deleteDiscountDialogTitle}
												onClick={() => {
													View.discountToDelete = d.id;
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
			use:enhance={View.onCostPageChangeForm}
			action="?/getCosts"
			class="flex w-full items-center justify-center"
			method="POST"
			bind:this={View.costPaginationForm}
		>
			<Pagination
				onPageChange={(num) => View.setPageChange(num)}
				pageNumber={View.costPageNumber}
				perPage={View.costLimit}
				count={View.costTotalRow}
				offset
			/>
		</form>
	{/if}
</div>
