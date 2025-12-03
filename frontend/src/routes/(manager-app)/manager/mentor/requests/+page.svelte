<script lang="ts">
	import type { PageProps } from './$types';
	import { RequestManagementView } from './view.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import Select from '$lib/components/select/Select.svelte';
	import { statusOpts } from './constants';

	let { data }: PageProps = $props();
	const View = new RequestManagementView(data.requests);
</script>

<svelte:head>
	<title>Request Management - Privat Unmei</title>
	<meta name="description" content="Request Management - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full flex-col p-4">
	<h3 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">Requests</h3>
	<form
		use:enhance={View.onSearchRequest}
		action="?/getRequests"
		class="mb-4 flex gap-4"
		method="POST"
	>
		<Select
			bind:value={View.status}
			options={statusOpts}
			defaultLable="Select Status"
			name="status"
		/>
		<Button type="submit">Search</Button>
	</form>
	<div class="flex flex-1">
		{#if View.isLoading}
			<Loading />
		{:else if !View.requests || View.requests.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No requests found</b>
		{:else}
			<ScrollArea orientation="vertical" class="flex-1" viewportClasses="max-h-[500px]">
				<table class="w-full table-fixed border-separate border-spacing-4">
					<thead>
						<tr>
							<th class="text-[var(--tertiary-color)]">Title</th>
							<th class="text-[var(--tertiary-color)]">Status</th>
							<th class="text-[var(--tertiary-color)]">Student Name</th>
							<th class="text-[var(--tertiary-color)]">Student Email</th>
							<th class="text-[var(--tertiary-color)]">Total Price</th>
						</tr>
					</thead>
					<tbody>
						{#each View.requests as r (r.id)}
							<tr>
								<td class="text-center">
									{r.course_name}
								</td>
								<td class="text-center">
									{View.capitalizeFirstLetter(r.status)}
								</td>
								<td class="text-center">
									{r.name}
								</td>
								<td class="text-center">
									{r.email}
								</td>
								<td class="text-center">
									{r.total_price}
								</td>
								<td>
									<div class="flex items-center justify-center">
										<div class="h-fit w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
											<Link href={`/manager/mentor/requests/${r.id}`}>Detail</Link>
										</div>
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
		use:enhance={View.onPageChange}
		action="?/getRequests"
		class="flex w-full items-center justify-center"
		method="POST"
		bind:this={View.paginationForm}
	>
		<Pagination
			onPageChange={View.setPageNumber}
			bind:pageNumber={View.page}
			perPage={View.limit}
			count={View.totalRow}
		/>
	</form>
</div>
