<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { MentorDetailView } from './view.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { dowMap } from './constants';
	import Link from '$lib/components/button/Link.svelte';
	import { enhance } from '$app/forms';
	import Pagination from '$lib/components/pagination/Pagination.svelte';
	import Loading from '$lib/components/loader/Loading.svelte';
	import Button from '$lib/components/button/Button.svelte';
	import { Star } from '@lucide/svelte';

	let { data }: PageProps = $props();
	const View = new MentorDetailView(data.courses);

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function handleDisplay() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', handleDisplay);
		return () => {
			window.removeEventListener('resize', handleDisplay);
		};
	});

	const mentorSchema = {
		'@context': 'https://schema.org',
		'@type': 'Person',
		name: data.profile.name,
		description: data.profile.bio,
		image: data.profile.profile_image,
		jobTitle: 'Online Mentor',
		aggregateRating: data.profile.rating > 0 ? {
			'@type': 'AggregateRating',
			ratingValue: data.profile.rating.toString(),
			ratingCount: '1'
		} : undefined,
		qualifications: {
			'@type': 'EducationalOccupationalCredential',
			educationalLevel: data.profile.degree,
			fieldOfStudy: data.profile.major,
			organization: {
				'@type': 'Organization',
				name: data.profile.campus
			}
		}
	};
</script>

<svelte:head>
	<title>{data.profile.name} - Expert Mentor | Privat Unmei</title>
	<meta name="description" content="Learn from {data.profile.name}, an experienced mentor with {data.profile.years_of_experience} years of experience in {data.profile.major}. Rating: {data.profile.rating}/5 on Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta property="og:title" content="{data.profile.name} - Expert Mentor | Privat Unmei" />
	<meta property="og:description" content="Learn from {data.profile.name}, an experienced {data.profile.degree} mentor specializing in {data.profile.major}." />
	<meta property="og:type" content="profile" />
	<meta property="og:image" content="{data.profile.profile_image}" />
	<meta name="keywords" content="{data.profile.name}, mentor, {data.profile.major}, online teaching, Privat Unmei" />
	<meta name="author" content="{data.profile.name}" />
	<script type="application/ld+json">
		{JSON.stringify(mentorSchema)}
	</script>
</svelte:head>

<main class="space-y-6 px-4 py-6 sm:px-6 md:space-y-8 md:py-8">
	<!-- Mentor Header Section -->
	<section class="space-y-4">
		<article class="flex flex-col gap-4 sm:flex-row sm:items-start sm:gap-6">
			<div class="shrink-0">
				<CldImage
					className="h-24 w-24 rounded-full object-cover sm:h-32 sm:w-32"
					height={128}
					width={128}
					src={data.profile.profile_image}
				/>
			</div>
			<div class="min-w-0 flex-1 space-y-3">
				<div>
					<h1 class="text-2xl font-bold text-(--tertiary-color) sm:text-3xl">
						{data.profile.name}
					</h1>
					<p class="text-sm text-(--tertiary-color) sm:text-base">{data.profile.public_id}</p>
				</div>
				<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-4">
					{#if data.profile.rating > 0}
						<div class="flex items-center gap-2">
							<div class="flex items-center gap-1">
								<p class="text-2xl font-bold text-(--tertiary-color) sm:text-3xl">{data.profile.rating}</p>
								<Star class="h-5 w-5 fill-current text-(--tertiary-color) sm:h-6 sm:w-6" />
							</div>
							<span class="text-xs text-(--tertiary-color) sm:text-sm">({data.profile.years_of_experience}y exp)</span>
						</div>
					{/if}
					{#if data.studentProfile}
						<form method="POST" action="?/messageMentor" use:enhance={View.onMessageMentor}>
							<Button type="submit">💬 Message</Button>
						</form>
					{/if}
				</div>
			</div>
		</article>
	</section>

	<!-- Education & Experience Section -->
	<section class="space-y-3">
		<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">Qualifications</h2>
		<div class="grid grid-cols-2 gap-3 rounded-lg bg-(--tertiary-color) p-4 sm:grid-cols-2 sm:gap-4 sm:p-5 md:grid-cols-4">
			<div class="space-y-1">
				<p class="text-xs font-bold text-(--primary-color) sm:text-sm">Experience</p>
				<p class="text-sm font-semibold text-(--secondary-color) sm:text-base">
					{data.profile.years_of_experience} {data.profile.years_of_experience === 1 ? 'year' : 'years'}
				</p>
			</div>
			<div class="space-y-1">
				<p class="text-xs font-bold text-(--primary-color) sm:text-sm">Degree</p>
				<p class="truncate text-sm font-semibold text-(--secondary-color) sm:text-base">
					{View.capitalizeFirstLetter(data.profile.degree)}
				</p>
			</div>
			<div class="space-y-1">
				<p class="text-xs font-bold text-(--primary-color) sm:text-sm">Major</p>
				<p class="truncate text-sm font-semibold text-(--secondary-color) sm:text-base">
					{data.profile.major}
				</p>
			</div>
			<div class="space-y-1">
				<p class="text-xs font-bold text-(--primary-color) sm:text-sm">Campus</p>
				<p class="truncate text-sm font-semibold text-(--secondary-color) sm:text-base">
					{data.profile.campus}
				</p>
			</div>
		</div>
	</section>

	<!-- Bio Section -->
	<section class="space-y-3">
		<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">About</h2>
		<div class="rounded-lg bg-(--tertiary-color) p-4 sm:p-5">
			<ScrollArea orientation="vertical" viewportClasses="h-48 max-h-48 sm:h-56">
				<p class="whitespace-pre-wrap text-justify text-(--secondary-color) text-sm leading-relaxed pr-4 sm:text-base">
					{data.profile.bio}
				</p>
			</ScrollArea>
		</div>
	</section>

	<!-- Schedules Section -->
	{#if data.schedules && data.schedules.length > 0}
		<section class="space-y-3">
			<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">Availability</h2>
			<ScrollArea orientation="vertical" viewportClasses="max-h-80">
				<ul class="space-y-2 pr-4 sm:space-y-3">
					{#each data.schedules as sch (sch.day_of_week + sch.start_time)}
						<li class="rounded-lg bg-(--tertiary-color) px-4 py-3 transition-shadow hover:shadow-md sm:px-5">
							<p class="font-medium text-(--secondary-color) text-sm sm:text-base">
								{dowMap.get(sch.day_of_week)}: {sch.start_time} - {sch.end_time}
							</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		</section>
	{/if}

	<!-- Courses Section -->
	<section class="space-y-3">
		<h2 class="text-lg font-bold text-(--tertiary-color) sm:text-xl">Courses</h2>
		<div class="space-y-4">
			{#if View.isLoading}
				<Loading />
			{:else if View.courses.length === 0}
				<div class="flex min-h-32 items-center justify-center rounded-lg bg-[rgba(54,84,134,0.5)]">
					<p class="font-semibold text-(--tertiary-color)">No courses available</p>
				</div>
			{:else}
				<ScrollArea orientation="vertical" viewportClasses="h-96 max-h-96 sm:h-[450px] sm:max-h-[450px]">
					<ul class="space-y-3 pr-4 md:grid md:grid-cols-2 lg:grid-cols-3 md:gap-3 md:space-y-0">
						{#each View.courses as c (c.id)}
							<li>
								<Link href={`/courses/${c.id}`}>
									<article class="group flex h-24 flex-col justify-between rounded-lg bg-(--tertiary-color) p-3 transition-all duration-200 hover:-translate-y-1 hover:shadow-md sm:p-4">
										<h3 class="truncate font-bold text-(--primary-color) text-sm sm:text-base">
											{c.title}
										</h3>
										<p class="text-xs font-semibold text-(--secondary-color) sm:text-sm">
											{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
												c.price
											)}
										</p>
									</article>
								</Link>
							</li>
						{/each}
					</ul>
				</ScrollArea>
			{/if}

			<!-- Pagination -->
			<form
				use:enhance={View.onPageChangeEnhance}
				action="?/getCourses"
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
	</section>
</main>
