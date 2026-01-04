<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { CourseDetailView } from './view.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import RatingGroup from '$lib/components/rating/RatingGroup.svelte';
	import Textarea from '$lib/components/form/Textarea.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';
	import { Star } from '@lucide/svelte';

	let { data }: PageProps = $props();
	const View = new CourseDetailView(data.reviews, data.detail);
</script>

<svelte:head>
	<title>{data.detail.title} - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<div class="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
		<h1 class="text-2xl font-bold text-[var(--tertiary-color)]">{data.detail.title}</h1>
		<div class="flex flex-col gap-2">
			<div class="w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
				<p class="font-bold text-[var(--secondary-color)]">
					{new Intl.NumberFormat('id-ID', { currency: 'IDR', style: 'currency' }).format(
						data.detail.price
					)} / session
				</p>
			</div>
		</div>
	</div>
	{#if data.profile}
		<div class="w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<Link href={`/requests/create/${data.detail.id}`}>Buy Course</Link>
		</div>
	{/if}

	<div>
		<NavigationButton
			menus={[
				{
					header: 'Description',
					onClick: () => (View.detailState = 'description')
				},
				{
					header: 'Detail',
					onClick: () => (View.detailState = 'detail')
				}
			]}
		/>
		<div
			class="flex h-[150px] flex-col gap-2 rounded-lg rounded-tl-none bg-[var(--tertiary-color)] p-4"
		>
			{#if View.detailState === 'detail'}
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Method:</p>
					<p class="text-[var(--secondary-color)]">
						{View.capitalizeFirstLetter(data.detail.method)}
					</p>
				</div>
				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Domicile:</p>
					<p class="text-[var(--secondary-color)]">{data.detail.domicile}</p>
				</div>
				<div class="flex items-center gap-2">
					<p class="font-bold text-[var(--primary-color)]">Per Session Duration:</p>
					<p class="text-[var(--secondary-color)]">
						{data.detail.session_duration_minutes} Minutes
					</p>
				</div>

				<div class="flex gap-2">
					<p class="font-bold text-[var(--primary-color)]">Categories:</p>
					<ScrollArea orientation="horizontal" viewportClasses="max-w-[200px]">
						<ul class="flex gap-4">
							{#each data.courseCategories as cc, i (cc.id)}
								<li
									class="w-fit rounded-lg bg-[var(--tertiary-color)] text-[var(--secondary-color)]"
								>
									<p>{cc.name}</p>
								</li>
							{/each}
						</ul>
					</ScrollArea>
				</div>
			{:else}
				<ScrollArea orientation="vertical" viewportClasses="h-[150px] max-h-[150px]">
					<p class="text-justify text-[var(--secondary-color)]">{data.detail.description}</p>
				</ScrollArea>
			{/if}
		</div>
	</div>

	<p class="font-bold text-[var(--tertiary-color)]">Mentor</p>
	<div class="flex items-center gap-4 rounded-md bg-[var(--tertiary-color)] p-2">
		<CldImage
			src={data.detail.mentor_profile_image}
			width={50}
			height={50}
			className="rounded-full"
		/>
		<div>
			<Link href={`/mentors/${data.detail.mentor_id}`}>
				<p class="font-bold text-[var(--primary-color)] hover:text-[var(--secondary-color)]">
					{data.detail.mentor_name}
				</p>
			</Link>
			<p class="text-[var(--secondary-color)]">{data.detail.mentor_public_id}</p>
			{#if data.profile}
				<form method="POST" action="?/messageMentor" use:enhance={View.onMessageMentor}>
					<Button withPadding={false} withBg={false} type="submit">Message</Button>
				</form>
			{/if}
		</div>
	</div>

	<div class="flex flex-col gap-4">
		<div class="flex flex-col gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Topics</p>
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
	{#if data.profile}
		<form
			use:enhance={View.onCreateReview}
			class="flex flex-col gap-4"
			action="?/createReview"
			method="post"
		>
			<RatingGroup bind:value={View.star} name="rating" />
			<Textarea
				err={View.feedbackErr}
				bind:value={View.feedback}
				name="feedback"
				id="feedback"
				placeholder="please insert feedback"
			>
				<p class="font-bold text-[var(--tertiary-color)]">Feedback:</p>
			</Textarea>
			<Button full disabled={View.reviewDisabled} type="submit">Submit</Button>
		</form>
	{/if}
	<div>
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
						<li class="flex flex-col gap-2 rounded-lg bg-[var(--tertiary-color)] p-2">
							<div class="flex w-full justify-between">
								<p class="font-bold text-[var(--primary-color)]">{r.name}</p>
								<div class="flex items-center gap-2">
									<Star class="fill-current text-[var(--primary-color)]" />
									<p class="font-bold text-[var(--primary-color)]">{r.rating}</p>
								</div>
							</div>
							<p class="text-[var(--secondary-color)]">{r.feedback}</p>
							<p class="text-end text-[var(--secondary-color)]">{View.getDate(r.created_at)}</p>
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
