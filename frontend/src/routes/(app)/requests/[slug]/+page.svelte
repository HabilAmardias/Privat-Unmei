<script lang="ts">
	import type { PageProps } from './$types';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';
	import { RequestDetailView } from './view.svelte';

	let { data }: PageProps = $props();
	const View = new RequestDetailView();
</script>

<svelte:head>
	<title>Request Detail - Privat Unmei</title>
	<meta name="description" content="Request Detail - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<h1 class="text-xl font-bold text-[var(--tertiary-color)]">{data.detail.course_name}</h1>
	<div class="flex flex-col gap-2">
		<div class="flex gap-2">
			<p class="font-bold text-[var(--tertiary-color)]">Mentor Name:</p>
			<p>{data.detail.mentor_name}</p>
		</div>
		<div class="flex gap-2">
			<p class="font-bold text-[var(--tertiary-color)]">Mentor Email:</p>
			<p>{data.detail.mentor_email}</p>
		</div>
		<div class="flex gap-2">
			<p class="font-bold text-[var(--tertiary-color)]">Participant:</p>
			<p>{data.detail.number_of_participant}</p>
		</div>
		<div class="flex gap-2">
			<p class="font-bold text-[var(--tertiary-color)]">Session:</p>
			<p>{data.detail.number_of_sessions}</p>
		</div>
		<div class="flex gap-2">
			<p class="font-bold text-[var(--tertiary-color)]">Status:</p>
			<p>{View.capitalizeFirstLetter(data.detail.status)}</p>
		</div>
		{#if data.detail.expired_at}
			<div class="flex gap-2">
				<p class="font-bold text-[var(--tertiary-color)]">Expired At:</p>
				<p>{View.convertToDatetime(data.detail.expired_at)}</p>
			</div>
		{/if}
	</div>
	<h2 class="text-lg font-bold text-[var(--tertiary-color)]">Payment Info</h2>
	<div class="grid grid-cols-2 gap-4 text-center md:flex md:justify-between">
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Payment Method:</p>
			<p>{data.detail.payment_method}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Account Number:</p>
			<p>{data.detail.account_number}</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Subtotal:</p>
			<p>
				{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
					data.detail.subtotal
				)}
			</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Operational Cost:</p>
			<p>
				{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
					data.detail.operational_cost
				)}
			</p>
		</div>
		<div>
			<p class="font-bold text-[var(--tertiary-color)]">Total Price:</p>
			<p>
				{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
					data.detail.total_price
				)}
			</p>
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<p class="font-bold text-[var(--tertiary-color)]">Schedules:</p>
		{#if data.detail.schedules}
			<ScrollArea orientation="vertical" viewportClasses="max-h-[300px]">
				<ul class="flex flex-col gap-4">
					{#each data.detail.schedules as sch (sch.date)}
						<li class="rounded-lg bg-[var(--tertiary-color)] p-4 text-[var(--secondary-color)]">
							<p class="font-bold text-[var(--primary-color)]">{View.convertToDate(sch.date)}</p>
							<p>{sch.start_time} - {sch.end_time}</p>
						</li>
					{/each}
				</ul>
			</ScrollArea>
		{:else}
			<b class="text-[var(--tertiary-color)]">No schedules found</b>
		{/if}
	</div>
</div>
