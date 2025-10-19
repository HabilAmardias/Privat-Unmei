<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import { MentorManagerView } from './view.svelte';
	import Input from '$lib/components/form/Input.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { enhance } from '$app/forms';
	import Loading from '$lib/components/loader/Loading.svelte';
	import { ScrollArea } from 'bits-ui';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import { ArrowDownUp, UserPlus } from '@lucide/svelte';
	import Link from '$lib/components/button/Link.svelte';

	let { data }: PageProps = $props();
	const View = new MentorManagerView();

	onMount(() => {
		if (!data.isVerified) {
			goto('/manager/admin/verify', { replaceState: true });
			return;
		}
		View.setIsDesktop(window.innerWidth >= 768);
		function setIsDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', setIsDesktop);
		View.setMentors(data.mentorsList.entries);
		View.setPaginationData(
			data.mentorsList.page_info.page,
			data.mentorsList.page_info.limit,
			data.mentorsList.page_info.total_row
		);
		return () => {
			window.removeEventListener('resize', setIsDesktop);
		};
	});
</script>

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
		use:enhance={View.onUpdateMentors}
		class="grid grid-cols-4 gap-4"
		action="?/getMentors"
		method="POST"
	>
		<Input
			bind:value={View.search}
			width="full"
			placeholder="Search"
			id="search"
			name="search"
			type="text"
		/>
		<Button onClick={() => View.onSort()} full type="submit">
			<div class="flex w-full justify-between px-1">
				<p>Sort</p>
				<ArrowDownUp />
			</div>
		</Button>
		<Button onClick={() => View.resetFilterForm()} full type="submit">
			<div class="w-full">Reset</div>
		</Button>
		<Button disabled={View.mentorsIsLoading} type="submit" full formAction="?/getMentors"
			>Search</Button
		>
	</form>
	<div class="flex flex-1">
		{#if View.mentorsIsLoading}
			<Loading />
		{:else if !View.mentors || View.mentors.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No mentors found</b>
		{:else}
			<ScrollArea.Root class="h-full">
				<ScrollArea.Viewport class="h-full">
					{#each View.mentors as mentor (mentor.id)}
						<div>
							<p>{mentor.email}</p>
							<p>{mentor.name}</p>
							<p>{mentor.years_of_experience}</p>
							<form use:enhance={View.onDeleteMentor} method="POST" action="?/deleteMentor">
								<Button
									onClick={() => {
										View.setMentorToDelete(mentor.id);
									}}
									type="submit"
									formAction="?/deleteMentor">Delete</Button
								>
							</form>
						</div>
					{/each}
				</ScrollArea.Viewport>
			</ScrollArea.Root>
		{/if}
	</div>
	<form
		use:enhance={View.onUpdateMentors}
		action="?/getMentors"
		class="flex w-full items-center justify-center"
		method="POST"
	>
		<Pagination
			onPageChange={(num) => View.onPageChange(num)}
			pageNumber={View.page}
			perPage={View.limit}
			count={View.total_row}
			offset
		/>
	</form>
</div>
