<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { CourseDetailView } from './view.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import NavigationButton from '$lib/components/button/NavigationButton.svelte';
	import { Star } from '@lucide/svelte';
	import type { CourseReview } from './model';

	let { data }: PageProps = $props();
	const View = new CourseDetailView(data.reviews, data.detail);

	const courseSchema = {
		'@context': 'https://schema.org',
		'@type': 'Course',
		name: data.detail.title,
		description: data.detail.description,
		provider: {
			'@type': 'Organization',
			name: data.detail.mentor_name,
			image: data.detail.mentor_profile_image
		},
		aggregateRating: {
			'@type': 'AggregateRating',
			ratingValue: data.reviews.entries.length > 0 
				? (data.reviews.entries.reduce((sum: number, r: CourseReview) => sum + r.rating, 0) / data.reviews.entries.length).toFixed(1)
				: '5',
			ratingCount: data.reviews.page_info.total_row
		},
		offers: {
			'@type': 'Offer',
			priceCurrency: 'IDR',
			price: data.detail.price.toString()
		}
	};
</script>

<svelte:head>
	<title>{data.detail.title} - Learn with {data.detail.mentor_name} | Privat Unmei</title>
	<meta name="description" content="{data.detail.description.substring(0, 160)}" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta property="og:title" content="{data.detail.title} - Privat Unmei" />
	<meta property="og:description" content="{data.detail.description.substring(0, 160)}" />
	<meta property="og:type" content="website" />
	<meta property="og:image" content="{data.detail.mentor_profile_image}" />
	<meta name="keywords" content="{data.detail.title}, {data.detail.mentor_name}, online course, learning, Privat Unmei" />
	<meta name="author" content="{data.detail.mentor_name}" />
	<script type="application/ld+json">
		{JSON.stringify(courseSchema)}
	</script>
</svelte:head>

<main class="space-y-6 px-4 py-6 sm:px-6 md:space-y-8 md:py-8">
	<!-- Header Section -->
	<section class="space-y-4">
		<div class="space-y-3">
			<h1 class="text-2xl font-bold text-(--tertiary-color) sm:text-3xl md:text-4xl">
				{data.detail.title}
			</h1>
			<p class="line-clamp-2 text-sm text-(--tertiary-color) sm:text-base">
				Learn from {data.detail.mentor_name}
			</p>
		</div>
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
			<div class="inline-flex rounded-lg bg-(--tertiary-color) px-3 py-2 sm:px-4 sm:py-3">
				<p class="font-bold text-(--secondary-color) text-sm sm:text-base">
					{new Intl.NumberFormat('id-ID', { currency: 'IDR', style: 'currency' }).format(
						data.detail.price
					)} / session
				</p>
			</div>
			{#if data.profile}
				<Link href={`/requests/create/${data.detail.id}`}>
					<button class="w-full rounded-lg bg-(--tertiary-color) px-4 py-2 font-semibold text-(--secondary-color) transition-transform hover:scale-105 sm:w-auto">
						Buy Course
					</button>
				</Link>
			{/if}
		</div>
	</section>

	<!-- Course Information Section -->
	<section class="space-y-3">
		<NavigationButton
			menus={[
				{
					header: 'Description',
					onClick: () => (View.detailState = 'description')
				},
				{
					header: 'Details',
					onClick: () => (View.detailState = 'detail')
				}
			]}
		/>
		<div class="space-y-3 rounded-lg rounded-tl-none bg-(--tertiary-color) p-4 sm:p-5">
			{#if View.detailState === 'detail'}
				<div class="space-y-3">
					<div class="flex flex-col gap-1 sm:flex-row sm:gap-3">
						<p class="font-bold text-(--primary-color)">Method:</p>
						<p class="text-(--secondary-color) text-sm">{View.capitalizeFirstLetter(data.detail.method)}</p>
					</div>
					<div class="flex flex-col gap-1 sm:flex-row sm:gap-3">
						<p class="font-bold text-(--primary-color)">Domicile:</p>
						<p class="text-(--secondary-color) text-sm">{data.detail.domicile}</p>
					</div>
					<div class="flex flex-col gap-1 sm:flex-row sm:gap-3">
						<p class="font-bold text-(--primary-color)">Duration:</p>
						<p class="text-(--secondary-color) text-sm">{data.detail.session_duration_minutes} minutes</p>
					</div>
					<div class="flex flex-col gap-2">
						<p class="font-bold text-(--primary-color)">Categories:</p>
						<ScrollArea orientation="horizontal" viewportClasses="max-w-full">
							<ul class="flex gap-2 pr-4">
								{#each data.courseCategories as cc (cc.id)}
									<li class="shrink-0 rounded-md bg-[rgba(255,255,255,0.1)] px-3 py-1 text-(--secondary-color) text-xs sm:text-sm">
										{cc.name}
									</li>
								{/each}
							</ul>
						</ScrollArea>
					</div>
				</div>
			{:else}
				<ScrollArea orientation="vertical" viewportClasses="h-56 max-h-56 sm:h-64">
					<p class="whitespace-pre-wrap text-justify text-(--secondary-color) text-sm leading-relaxed pr-4">
						{data.detail.description}
					</p>
				</ScrollArea>
			{/if}
		</div>
	</section>

	<!-- Mentor Section -->
	<section class="space-y-3">
		<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">Expert Mentor</h2>
		<article class="flex flex-col gap-3 rounded-lg bg-(--tertiary-color) p-4 sm:flex-row sm:items-start sm:gap-4 sm:p-5">
			<div class="shrink-0">
				<CldImage
					src={data.detail.mentor_profile_image}
					width={70}
					height={70}
					className="rounded-full object-cover"
				/>
			</div>	
			<div class="min-w-0 flex-1 space-y-2">
				<Link href={`/mentors/${data.detail.mentor_id}`}>
					<h3 class="font-bold text-(--primary-color) hover:text-(--secondary-color) transition-colors">
						{data.detail.mentor_name}
					</h3>
				</Link>
				<p class="text-(--secondary-color) text-xs sm:text-sm">{data.detail.mentor_public_id}</p>
				{#if data.profile}
					<form method="POST" action="?/messageMentor" use:enhance={View.onMessageMentor} class="pt-2">
						<Button withPadding={false} withBg={false} type="submit">
							💬 Send Message
						</Button>
					</form>
				{/if}
			</div>
		</article>
	</section>

	<!-- Topics Section -->
	{#if data.topics && data.topics.length > 0}
		<section class="space-y-3">
			<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">Course Topics</h2>
			<ScrollArea orientation="vertical" viewportClasses="max-h-96">
				<ul class="space-y-3 pr-4">
					{#each data.topics as t (t.title)}
						<li class="rounded-lg bg-(--tertiary-color) p-4 space-y-2 hover:shadow-md transition-shadow">
							<h3 class="font-bold text-(--primary-color) text-sm sm:text-base">{t.title}</h3>
							<p class="text-(--secondary-color) text-xs sm:text-sm leading-relaxed">{t.description}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		</section>
	{/if}

	<!-- Reviews Section -->
	<section class="space-y-4">
		<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">Student Reviews</h2>
		<div>
			{#if View.isLoading}
				<Loading />
			{:else if !View.reviews || View.reviews.length === 0}
				<div class="flex min-h-32 items-center justify-center rounded-lg bg-[rgba(54,84,134,0.5)]">
					<p class="font-semibold text-(--tertiary-color)">No reviews yet</p>
				</div>
			{:else}
				<ScrollArea orientation="vertical" viewportClasses="h-96 max-h-96 sm:h-[500px] sm:max-h-[500px]">
					<ul class="space-y-3 pr-4 md:grid md:grid-cols-2 lg:grid-cols-3 md:gap-3 md:space-y-0">
						{#each View.reviews as r (r.id)}
							<li class="space-y-2 rounded-lg bg-(--tertiary-color) p-3 sm:p-4">
								<div class="flex items-start justify-between gap-2">
									<p class="font-bold text-(--primary-color) text-sm sm:text-base truncate">{r.name}</p>
									<div class="flex shrink-0 items-center gap-1">
										<Star class="h-4 w-4 sm:h-5 sm:w-5 fill-current text-(--primary-color)" />
										<p class="font-bold text-(--primary-color) text-sm">{r.rating}</p>
									</div>
								</div>
								<p class="line-clamp-3 text-(--secondary-color) text-xs sm:text-sm leading-relaxed">
									{r.feedback}
								</p>
								<p class="text-end text-(--secondary-color) text-xs">{View.getDate(r.created_at)}</p>
							</li>
						{/each}
					</ul>
				</ScrollArea>
			{/if}
		</div>

		<!-- Pagination -->
		<form
			use:enhance={View.onPageChangeEnhance}
			action="?/getReviews"
			class="flex w-full items-center justify-center pt-2"
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
	</section>
</main>
