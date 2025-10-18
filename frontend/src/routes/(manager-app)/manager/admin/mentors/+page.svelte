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

<div class="flex flex-1 flex-col gap-4">
	<h3 class="text-xl font-bold text-[var(--tertiary-color)]">Mentors</h3>
	<form use:enhance class="grid grid-cols-3 gap-4" action="?/getMentors" method="POST">
		<Input width="full" placeholder="Search" id="search" name="search" type="text" />
		<Button disabled={View.mentorsIsLoading} type="submit" full formAction="?/myOrders"
			>Search</Button
		>
	</form>
	<div class="flex flex-1">
		{#if View.mentorsIsLoading}
			<Loading />
		{:else if !View.mentors || View.mentors.length === 0}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No orders found</b>
		{:else}
			<ScrollArea.Root class="h-full">
				<ScrollArea.Viewport class="h-full">
					{#each View.mentors as mentor (mentor.id)}
						<div>
							<p>{mentor.email}</p>
							<p>{mentor.name}</p>
							<p>{mentor.years_of_experience}</p>
							<form use:enhance method="POST" action="?/deleteMentor">
								<Button type="submit" formAction="?/deleteMentor">Delete</Button>
							</form>
						</div>
					{/each}
				</ScrollArea.Viewport>
			</ScrollArea.Root>
		{/if}
	</div>
</div>
