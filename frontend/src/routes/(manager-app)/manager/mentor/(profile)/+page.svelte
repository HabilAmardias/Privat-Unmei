<script lang="ts">
	import type { PageProps } from './$types';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { MentorDetailView } from './view.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { dowMap } from './constants';
	import Link from '$lib/components/button/Link.svelte';

	const View = new MentorDetailView();
	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Profile - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
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
		<div class="flex w-full items-center justify-between">
			<div class="flex flex-col gap-1">
				<p class="font-bold text-[var(--tertiary-color)]">{data.profile.name}</p>
				<p>{data.profile.email}</p>
			</div>
		</div>
		<div class="h-fit w-fit rounded-lg bg-[var(--tertiary-color)] p-2">
			<Link href="/manager/mentor/update">Update</Link>
		</div>
	</div>
	<div class="flex justify-between text-center">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">YoE:</p>
			<p>{data.profile.years_of_experience}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Campus:</p>
			<p>{data.profile.campus}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Degree:</p>
			<p>{data.profile.degree}</p>
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
	<div class="flex flex-col gap-4 md:grid md:grid-cols-2">
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
		<div class="flex flex-col gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Payments:</p>
			{#if data.payments}
				<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
					<ul class="flex flex-col gap-4">
						{#each data.payments as py, i (i)}
							<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
								<p>{py.payment_method_name} - {py.account_number}</p>
							</li>
						{/each}
					</ul>
				</ScrollArea>
			{:else}
				<b class="text-[var(--tertiary-color)]">No payments found</b>
			{/if}
		</div>
	</div>
</div>
