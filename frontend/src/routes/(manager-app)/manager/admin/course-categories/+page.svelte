<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { CourseCategoryManagementView } from './view.svelte';
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
	const View = new CourseCategoryManagementView(data.categories);
	onMount(() => {
		if (!data.isVerified) {
			goto('/managet/admin/verify', { replaceState: true });
		}
	});
</script>

{#snippet deleteDialogTitle()}
	Delete Category Confirmation
{/snippet}

{#snippet deleteDialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

{#snippet createDialogTitle()}
	Create Course Category
{/snippet}

{#snippet createDialogDescription()}
	Add New Course Category
{/snippet}

{#snippet updateDialogTitle()}
	Update Course Category
{/snippet}

{#snippet updateDialogDescription()}
	Update Course Category Name
{/snippet}

<div class="flex h-full flex-col p-4">
	<div class="mb-4 flex items-center justify-between">
		<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Course Categories</h3>
		<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<Dialog
				bind:open={View.createDialogOpen}
				buttonText="Create Category"
				title={createDialogTitle}
				description={createDialogDescription}
				buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
			>
				<form
					class="flex flex-col gap-4"
					action="?/createCategory"
					method="POST"
					use:enhance={View.onCreateCategory}
				>
					<Input type="text" id="name" name="name" placeholder="Insert Course Category Name" />
					<Button full type="submit">Create</Button>
				</form>
			</Dialog>
		</div>
	</div>
	<form
		bind:this={View.searchForm}
		use:enhance={View.onSearchCategory}
		action="?/getCategories"
		class="mb-4 flex gap-4"
		method="POST"
	>
		<Input
			bind:value={View.search}
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
		{:else if !View.categories || View.categories.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No categories found</b>
		{:else}
			<ScrollArea orientation="vertical" class="flex-1" viewportClasses="max-h-[500px]">
				<table class="w-full table-fixed border-separate border-spacing-4">
					<thead>
						<tr>
							<th class="text-[var(--tertiary-color)]">Name</th>
						</tr>
					</thead>
					<tbody>
						{#each View.categories as c (c.id)}
							<tr>
								<td class="text-center">
									{c.name}
								</td>
								<td>
									<div class="flex items-center justify-center gap-4">
										<div class="h-fit rounded-lg bg-[var(--tertiary-color)] p-2">
											<Dialog
												bind:open={View.updateDialogOpen}
												buttonText="Update"
												title={updateDialogTitle}
												buttonOnClick={() => {
													View.categoryToUpdate = c.id;
												}}
												description={updateDialogDescription}
												buttonClass="text-[var(--secondary-color)] cursor-pointer hover:text-[var(--primary-color)]"
											>
												<form
													class="flex flex-col gap-4"
													action="?/updateCategory"
													method="POST"
													use:enhance={View.onUpdateCategory}
												>
													<Input
														type="text"
														id="name"
														name="name"
														placeholder="Insert Course Category Name"
													/>
													<Button full type="submit">Update</Button>
												</form>
											</Dialog>
										</div>
										<AlertDialog
											action="?/deleteCategory"
											bind:open={View.deleteDialogOpen}
											enhancement={View.onDeleteCategory}
											title={deleteDialogTitle}
											onClick={() => {
												View.categoryToDelete = c.id;
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
		use:enhance={View.onPageChangeForm}
		action="?/getCategories"
		class="flex w-full items-center justify-center"
		method="POST"
		bind:this={View.paginationForm}
	>
		<Pagination
			onPageChange={(num) => View.setPageNumber(num)}
			pageNumber={View.pageNumber}
			perPage={View.limit}
			count={View.totalRow}
			offset
		/>
	</form>
</div>
