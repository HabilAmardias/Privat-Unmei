<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { MentorDetailView } from './view.svelte';
	import AlertDialog from '$lib/components/dialog/AlertDialog.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { dowMap } from './constants';

	const View = new MentorDetailView();
	let { data }: PageProps = $props();

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
	});
</script>

{#snippet dialogTitle()}
	Delete Mentor Confirmation
{/snippet}

{#snippet dialogDescription()}
	Irreversible action, are you sure want to proceed?
{/snippet}

<svelte:head>
	<title>Mentor Detail - Privat Unmei</title>
	<meta name="description" content="Mentor Detail - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<div class="flex gap-4">
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
			<AlertDialog
				action="?/deleteMentor"
				bind:open={View.alertOpen}
				enhancement={View.onDeleteMentor}
				title={dialogTitle}
				description={dialogDescription}>Delete Account</AlertDialog
			>
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
				<ScrollArea class="h-[100px]" orientation="vertical" viewportClasses="max-h-[100px]">
					{#each data.schedules as sch, i (i)}
						<div class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p>{dowMap.get(sch.day_of_week)}, {sch.start_time} - {sch.end_time}</p>
						</div>
					{/each}
				</ScrollArea>
			{:else}
				<b class="mx-auto self-center text-[var(--tertiary-color)]">No schedules found</b>
			{/if}
		</div>
		<div class="flex flex-col gap-4">
			<p class="font-bold text-[var(--tertiary-color)]">Payments:</p>
			{#if data.payments}
				<ScrollArea class="h-[100px]" orientation="vertical" viewportClasses="max-h-[100px]">
					{#each data.payments as py, i (i)}
						<div class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p>{py.payment_method_name} - {py.account_number}</p>
						</div>
					{/each}
				</ScrollArea>
			{:else}
				<b class="text-[var(--tertiary-color)]">No payments found</b>
			{/if}
		</div>
	</div>
</div>
