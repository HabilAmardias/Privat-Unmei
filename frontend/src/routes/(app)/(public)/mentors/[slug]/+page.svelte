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
</script>

<svelte:head>
	<title>{data.profile.name} - Privat Unmei</title>
	<meta name="description" content="Mentor Detail - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<div class="flex items-center gap-4">
		<CldImage
			className="rounded-full"
			height={View.size}
			width={View.size}
			src={data.profile.profile_image}
		/>
		<div class="flex flex-col gap-1">
			<p class="font-bold text-[var(--tertiary-color)]">{data.profile.name}</p>
			<p>{data.profile.email}</p>
			{#if data.studentProfile}
				<form method="POST" action="?/messageMentor" use:enhance={View.onMessageMentor}>
					<Button type="submit">Message</Button>
				</form>
			{/if}
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Experience:</p>
			<p>{data.profile.years_of_experience} Year</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Campus:</p>
			<p>{data.profile.campus}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Degree:</p>
			<p>{View.capitalizeFirstLetter(data.profile.degree)}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Major:</p>
			<p>{data.profile.major}</p>
		</div>
	</div>
	<div class="flex flex-col">
		<p class="font-bold text-[var(--tertiary-color)]">Bio:</p>
		<p>{data.profile.bio}</p>
	</div>
	<div class="flex flex-col gap-4">
		<p class="font-bold text-[var(--tertiary-color)]">Schedules:</p>
		{#if data.schedules}
			<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
				<ul class="flex flex-col gap-4">
					{#each data.schedules as sch, i (i)}
						<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p>{dowMap.get(sch.day_of_week)}, {sch.start_time} - {sch.end_time}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{:else}
			<b class="mx-auto self-center text-[var(--tertiary-color)]">No schedules found</b>
		{/if}
	</div>

	<h3 class="font-bold text-[var(--tertiary-color)]">Courses:</h3>
	<ScrollArea orientation="vertical" viewportClasses="h-[300px] max-h-[300px]">
		{#if View.isLoading}
			<Loading />
		{:else if View.courses.length === 0}
			<b class="text-[var(--tertiary-color)]">No courses found</b>
		{:else}
			<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
				{#each View.courses as c (c.id)}
					<li class="transition-transform hover:-translate-y-1">
						<Link href={`/courses/${c.id}`}>
							<div
								class="flex h-[100px] flex-col justify-between rounded-lg bg-[var(--tertiary-color)] p-4"
							>
								<p class="font-bold text-[var(--primary-color)]">{c.title}</p>
								<div>
									<p class="text-[var(--secondary-color)]">
										{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
											c.price
										)}
									</p>
								</div>
							</div>
						</Link>
					</li>
				{/each}
			</ul>
		{/if}
	</ScrollArea>
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
