<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { CourseDetailView } from './view.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';

	let { data }: PageProps = $props();
	const View = new CourseDetailView(data.reviews);
</script>

<svelte:head>
	<title>{data.detail.title} - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<h1 class="text-xl font-bold text-[var(--tertiary-color)]">{data.detail.title}</h1>
	<div class="grid grid-cols-2 gap-4 text-center md:flex md:justify-between">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Method:</p>
			<p>{data.detail.method}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Domicile:</p>
			<p>{data.detail.domicile}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Max Session:</p>
			<p>{data.detail.max_total_session}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Session Duration:</p>
			<p>{data.detail.session_duration_minutes}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Price:</p>
			<p>{data.detail.price}</p>
		</div>
	</div>
	<div class="flex flex-col">
		<p class="font-bold text-[var(--tertiary-color)]">Description:</p>
		<p>{data.detail.description}</p>
	</div>
	<div class="flex flex-col gap-4 md:grid md:grid-cols-2">
		<div class="flex flex-col gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Categories:</p>
			{#if data.courseCategories}
				<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
					<ul class="flex flex-col gap-4">
						{#each data.courseCategories as cc, i (cc.id)}
							<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
								<p>{cc.name}</p>
							</li>
						{/each}
					</ul>
				</ScrollArea>
			{:else}
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No categories found</b>
			{/if}
		</div>
		<div class="flex flex-col gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Topics:</p>
			{#if data.topics}
				<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
					<ul class="flex flex-col gap-4">
						{#each data.topics as t, i (i)}
							<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
								<p class="font-bold text-[var(--primary-color)]">{t.title}</p>
								<p>{t.description}</p>
							</li>
						{/each}
					</ul>
				</ScrollArea>
			{:else}
				<b class="text-[var(--tertiary-color)]">No topic found</b>
			{/if}
		</div>
	</div>
	<h2 class="text-xl font-bold text-[var(--tertiary-color)]">Reviews</h2>
	<div class="flex-1">
		{#if View.isLoading}
			<Loading />
		{:else if !View.reviews || View.reviews.length === 0}
			<div class="flex h-full items-center justify-center">
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No reviews found</b>
			</div>
		{:else}
			<ScrollArea orientation="vertical" viewportClasses="h-[400px] max-h-[400px]">
				<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
					{#each View.reviews as r (r.id)}
						<li>
							<div class="flex w-full justify-between">
								<p>{r.name}</p>
								<p>{r.rating}</p>
							</div>
							<p>{r.feedback}</p>
							<p class="text-end">{r.created_at}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{/if}
	</div>
	<form
		use:enhance={View.onPageChangeEnhance}
		action="?/getReviews"
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
