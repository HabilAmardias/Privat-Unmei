<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { MentorManagerView } from './view.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { enhance } from '$app/forms';
	import Loading from '$lib/components/loader/Loading.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import { ArrowDown, ArrowUp, UserPlus } from '@lucide/svelte';
	import Link from '$lib/components/button/Link.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';

	let { data }: PageProps = $props();
	const View = new MentorManagerView(data.mentorsList);

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function setIsDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', setIsDesktop);

		return () => {
			window.removeEventListener('resize', setIsDesktop);
		};
	});
</script>

{#snippet dialogTitle()}
	Delete Mentor Confirmation
{/snippet}

{#snippet dialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

<svelte:head>
	<title>Mentors - Privat Unmei</title>
	<meta name="description" content="Mentors - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full flex-col gap-4 p-4">
	<div class="flex justify-between gap-4">
		<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Mentors</h3>
		<Link href="/manager/admin/mentors/create" theme="light">
			<div class="flex gap-3 rounded-lg bg-[var(--tertiary-color)] p-2">
				<p>Create Mentor</p>
				<UserPlus />
			</div>
		</Link>
	</div>
	<form
		use:enhance={View.onSearchMentors}
		bind:this={View.searchForm}
		class="grid grid-cols-2 gap-4"
		action="?/getMentors"
		method="POST"
	>
		<Input
			onInput={View.onSearchInput}
			width="full"
			placeholder="Search"
			id="search"
			name="search"
			type="text"
		/>
		<Button onClick={() => View.onSort()} full type="submit">
			<div class="flex w-full justify-between px-1">
				<p>Sort</p>
				{#if View.sortByYears === true}
					<ArrowUp />
				{:else if View.sortByYears === false}
					<ArrowDown />
				{/if}
			</div>
		</Button>
	</form>
	<div class="flex flex-1">
		{#if View.mentorsIsLoading}
			<Loading />
		{:else if !View.mentors || View.mentors.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No mentors found</b>
		{:else}
			<ScrollArea orientation="vertical" class="flex-1" viewportClasses="w-full max-h-[500px]">
				<table class="w-full border-separate border-spacing-4">
					<thead>
						<tr>
							<th class="text-[var(--tertiary-color)]">Email</th>
							<th class="text-[var(--tertiary-color)]">Name</th>
							<th class="text-[var(--tertiary-color)]">YoE</th>
						</tr>
					</thead>
					<tbody>
						{#each View.mentors as mentor (mentor.id)}
							<tr>
								<td class="overflow-x-auto text-center">{mentor.email}</td>
								<td class="overflow-x-auto text-center">{mentor.name}</td>
								<td class="overflow-x-auto text-center">{mentor.years_of_experience}</td>
								<td class="text-center">
									<Link theme="dark" href={`/manager/admin/mentors/${mentor.id}`}>Detail</Link>
								</td>
								<td>
									<AlertDialog
										action="?/deleteMentor"
										bind:open={View.alertOpen}
										enhancement={View.onDeleteMentor}
										title={dialogTitle}
										onClick={() => {
											View.setMentorToDelete(mentor.id);
										}}
										description={dialogDescription}>Delete</AlertDialog
									>
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
		action="?/getMentors"
		class="flex w-full items-center justify-center"
		method="POST"
	>
		<Pagination bind:pageNumber={View.page} perPage={View.limit} count={View.total_row} />
	</form>
</div>
