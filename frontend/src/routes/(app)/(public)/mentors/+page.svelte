<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import Input from '$lib/components/form/Input.svelte';
	import { MentorListView } from './view.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';

	let { data }: PageProps = $props();
	const View = new MentorListView(data.mentors);
</script>

<svelte:head>
	<title>Mentors - Privat Unmei</title>
	<meta name="description" content="Mentors - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>
<div class="flex h-full flex-col gap-4 p-4">
	<h1 class="text-2xl font-bold text-[var(--tertiary-color)]">Mentors</h1>
	<form
		bind:this={View.searchForm}
		use:enhance={View.onSearchMentorEnhance}
		class="flex flex-col gap-4 md:flex-row"
		action="?/getMentors"
		method="POST"
	>
		<Input
			bind:value={View.search}
			onInput={View.onSearchMentor}
			width="full"
			placeholder="Search"
			id="search"
			name="search"
			type="text"
		/>
	</form>
	<div class="flex-1">
		{#if View.isLoading}
			<Loading />
		{:else if !View.mentors || View.mentors.length === 0}
			<div class="flex h-full items-center justify-center">
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No mentors found</b>
			</div>
		{:else}
			<ScrollArea orientation="vertical" viewportClasses="h-[400px] max-h-[300px]">
				<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
					{#each data.mentors.entries as m (m.id)}
						<li class=" transition-transform hover:-translate-y-1">
							<Link href={`/mentors/${m.id}`}>
								<div class="flex h-[100px] gap-4 rounded-lg bg-[var(--tertiary-color)] p-4">
									<CldImage src={m.profile_image} width={70} height={70} className="rounded-full" />
									<div>
										<p class="font-bold text-[var(--primary-color)]">{m.name}</p>
										<p class="text-[var(--secondary-color)]">{m.email}</p>
									</div>
								</div>
							</Link>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{/if}
	</div>
	<form
		use:enhance={View.onPageChangeEnhance}
		action="?/getMentors"
		class="flex w-full items-center justify-center"
		method="POST"
		bind:this={View.paginationForm}
	>
		<Pagination
			onPageChange={View.onPageChange}
			bind:pageNumber={View.page}
			perPage={View.limit}
			count={View.totalRow}
		/>
	</form>
</div>
